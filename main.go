package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tg2648/grem/internal/models"
	"github.com/urfave/cli/v2"
)

const configDirName = ".grem"
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

	configDir := homeDir + "/" + configDirName
	if !directoryExists(configDir) {
		if err = os.Mkdir(configDir, 0755); err != nil {
			log.Fatalf("Unable to create the config directory: %s", err.Error())
		}
		log.Printf("Created config directory in %s", configDir)
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
				Name:            "reminders",
				Usage:           "Manage reminders",
				HideHelpCommand: true,
				Subcommands: []*cli.Command{
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

							fmt.Printf("title: %q\n", title)
							fmt.Printf("due: %q\n", due)

							// TODO: if no flags provided, show a TUI interface
							if title == "" || due == nil {
								return cli.Exit("Error: title and date are required", 1)
							}

							fmt.Println("Adding reminder")
							id, err := reminders.Insert(title, due)
							if err != nil {
								return err
							}
							fmt.Println("Reminder added: ", id)

							return nil
						},
						OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
							fmt.Println("Usage error: ", err.Error())
							return nil
						},
					},
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Present() {
				return cli.Exit(
					fmt.Sprintf(
						"Error: command %q not found. Run %s --help for a list of available commands.\n",
						ctx.Args().First(),
						ctx.App.Name,
					), 1)
			}

			r, err := reminders.Get(4)
			if err != nil {
				if errors.Is(err, models.ErrNoRecord) {
					log.Println("Reminder not found")
				} else {
					log.Println("Error getting reminder: ", err.Error())
				}
			} else {
				log.Printf("Reminder (%d) %q due at %q\n", r.ID, r.Title, r.DueAt.Format("2006-01-02"))
			}

			fmt.Println("All Reminders:")
			date, _ := time.Parse("2006-01-02", "2024-12-15")
			due, err := reminders.GetDue(&date)
			if err != nil {
				return err
			}

			for _, reminder := range due {
				fmt.Printf("%+v\n", reminder)
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
