package main

import (
	"context"
	"fmt"
	"github.com/test/warno-utils/pkg/update"
	"github.com/test/warno-utils/pkg/utils"
	"os"
	"path/filepath"

	"github.com/test/warno-utils/pkg/switcher"
	urfavealt "github.com/urfave/cli-altsrc/v3"
	urfave "github.com/urfave/cli/v3"
)

var version = "dev"

func main() {
	docs, err := utils.GetUserDocsDir()
	if err != nil {
		fmt.Println("Failed to get user documents directory:", err)
		os.Exit(1)
	}

	dir, err := utils.GetBinaryDir()
	if err != nil {
		fmt.Println("Failed to get binary directory:", err)
		os.Exit(1)
	}
	configFiles := []string{filepath.Join(docs, "config.yaml"), filepath.Join(dir, "config.yaml"), filepath.Join(docs, "config.yml"), filepath.Join(dir, "config.yml")}
	app := &urfave.Command{
		Name:                  "warno-util",
		Usage:                 "Utility for WARNO configurations",
		EnableShellCompletion: true,
		Version:               version,
		Commands: []*urfave.Command{
			{
				Name:                  "switch",
				Usage:                 "Switch WARNO configurations between VIP beta and standard. Only works for people with access to the VIP beta. The first time you do a version switch you will need to install whichever version you didn't already have installed",
				EnableShellCompletion: true,
				Flags: []urfave.Flag{
					&urfave.StringFlag{
						Name:    "steamAppsPath",
						Usage:   "Path to Steam steamapps directory",
						Sources: urfavealt.YAML("switch.steamAppsPath", configFiles...),
						Aliases: []string{"a"},
					},
					&urfave.StringFlag{
						Name:    "steamUserDataPath",
						Usage:   "Path to Steam userdata directory",
						Sources: urfavealt.YAML("switch.steamUserDataPath", configFiles...),
						Aliases: []string{"u"},
					},
					&urfave.StringFlag{
						Name:    "steamExecutablePath",
						Usage:   "Path to Steam executable steam.exe",
						Sources: urfavealt.YAML("switch.steamExecutablePath", configFiles...),
						Aliases: []string{"e"},
					},
					&urfave.StringFlag{
						Name:    "config",
						Usage:   "Path to YAML config file. Will default to looking in the user's Documents directory and the directory of the binary for config.yaml or config.yml",
						Value:   "config.yaml",
						Aliases: []string{"c"},
					},
				},
				Action: func(c *urfave.Context) error {
					fmt.Println(configFiles)

					if err := c.Err(); err != nil {
						return err
					}
					cfg := switcher.Config{
						SteamAppsPath:       c.String("steamAppsPath"),
						SteamUserDataPath:   c.String("steamUserDataPath"),
						SteamExecutablePath: c.String("steamExecutablePath"),
					}
					fmt.Println("Steam Steamapps Path:", cfg.SteamAppsPath)
					fmt.Println("Steam Userdata Path:", cfg.SteamUserDataPath)
					fmt.Println("Steam Executable Path:", cfg.SteamExecutablePath)

					if err := switcher.Switcher(cfg); err != nil {
						return fmt.Errorf("failed to execute switcher: %w", err)
					}

					return nil
				},
			},
			{
				Name:                  "update",
				Usage:                 "update warnoutil",
				EnableShellCompletion: true,
				Flags: []urfave.Flag{
					&urfave.StringFlag{
						Name:    "auto-approve",
						Aliases: []string{"y"},
					},
				},
				Action: func(c *urfave.Context) error {
					if err := c.Err(); err != nil {
						return err
					}

					u := update.Updater{Version: version}
					return u.RunUpdate()
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
