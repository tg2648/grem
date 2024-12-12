package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
)

const configDirName = ".grem"

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

func setup() {
	var homeDir, err = os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to determine user's home directory: %s", err.Error())
	}

	var configDir = homeDir + "/" + configDirName
	if !directoryExists(configDir) {
		if err = os.Mkdir(configDir, 0755); err != nil {
			log.Fatalf("Unable to create the config directory: %s", err.Error())
		}
		log.Printf("Created config directory in %s", configDir)
	}
}

func main() {
	setup()

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
								Name:     "title",
								Aliases:  []string{"t"},
								Usage:    "Reminder's `title`",
								Required: true,
							},
						},
						Action: func(ctx *cli.Context) error {
							fmt.Println("Adding reminder", ctx.String("title"))
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
			fmt.Println("Notifications:")
			return nil
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
