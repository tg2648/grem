package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tg2648/grem/internal/models"
	"github.com/urfave/cli/v2"
)

const appDirName = ".grem"
const dueDateLayout = "2006-01-02"

func directoryExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		// If the error is not "file does not exist," return false
		if os.IsNotExist(err) {
			return false
		}
	}
	return info.IsDir()
}

func setup(db *sql.DB) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to determine user's home directory: %s", err.Error())
	}

	appDir := homeDir + "/" + appDirName
	if !directoryExists(appDir) {
		if err = os.Mkdir(appDir, 0755); err != nil {
			log.Fatalf("Unable to create the application directory directory: %s", err.Error())
		}
		log.Printf("Created the application directory in %s", appDir)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS reminders (
		id INTEGER NOT NULL PRIMARY KEY,
		title TEXT NOT NULL,
		due_at DATE NOT NULL,
		dismissed_at DATETIME DEFAULT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

// The openDb() function wraps [sql.Open] and returns a sql.DB connection pool
// for a given DSN.
func openDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// type application struct {
// 	errorLog *log.Logger
// 	reminders *models.ReminderModel
// 	cliApp *cli.App
// }

func main() {
	dsn := "file:grem.db"
	db, err := openDb(dsn)

	if err != nil {
		log.Fatal("db error: ", err.Error())
	}
	defer db.Close()

	setup(db)
	reminders := &models.ReminderModel{DB: db}

	app := &cli.App{
		Name:            "grem",
		Usage:           "A program for managing terminal-based reminders",
		Version:         "v0.0.1",
		HideHelpCommand: true,
		Commands: []*cli.Command{
			{
				Name:  "add",
				Usage: "Schedule a new reminder",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "title",
						Aliases: []string{"t"},
						Usage:   "Reminder's `title`",
					},
					&cli.TimestampFlag{
						Name:    "due",
						Aliases: []string{"d"},
						Usage:   "Reminder's due `date` in the YYYY-MM-DD format",
						Layout:  dueDateLayout,
					},
				},
				Action: func(ctx *cli.Context) error {
					title := ctx.String("title")
					due := ctx.Timestamp("due")

					// fmt.Printf("title: %q\n", title)
					// fmt.Printf("due: %q\n", due)

					// TODO: if no flags provided, show a TUI interface
					if title == "" || due == nil {
						return cli.Exit("Error: title and due date are required", 1)
					}

					_, err := reminders.Insert(title, due)
					if err != nil {
						return err
					}
					fmt.Println("Reminder added: ", title)

					return nil
				},
				OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
					fmt.Println("Usage error: ", err.Error())
					return nil
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			// When ran without arguments, show all reminders due today
			if ctx.Args().Present() {
				return cli.Exit(
					fmt.Sprintf(
						"Error: command %q not found. Run %s --help for a list of available commands.\n",
						ctx.Args().First(),
						ctx.App.Name,
					), 1)
			}

			due, err := reminders.GetDueToday()
			if err != nil {
				return err
			}

			fmt.Println("grem reminders (oldest first):")
			for _, reminder := range due {
				fmt.Printf("%s (id: %03d): %s \n", reminder.DueAt.Format(time.DateOnly), reminder.ID, reminder.Title)
			}

			return nil
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
