package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

const updaterVersion = "0.1.3"

func cmdDownloadExtra() *cli.Command {
	return &cli.Command{
		Name: "download-extra",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "version",
				Usage: "version",
			},
			&cli.StringFlag{
				Name:  "out",
				Usage: "out",
			},
			&cli.StringFlag{
				Name:  "platform",
				Usage: "platform",
				Value: runtime.GOOS,
			},
		},
		Action: func(c *cli.Context) error {
			return downloadExtra(c.String("version"), c.String("platform"), c.String("out"))
		},
	}
}

func downloadExtra(version string, platform string, out string) error {
	if version == "" {
		return errors.Errorf("no version specified")
	}

	skip := []string{"README.md", "LICENSE"}

	keysFile := fmt.Sprintf("keys_%s_%s_x86_64.tar.gz", version, platform)
	keysURLString := fmt.Sprintf("https://github.com/keys-pub/keysd/releases/download/v%s/%s", version, keysFile)

	log.Printf("Extracting %s\n", keysURLString)
	if err := extractURL(keysURLString, out, skip); err != nil {
		return err
	}

	updaterFile := fmt.Sprintf("updater_%s_%s_x86_64.tar.gz", updaterVersion, platform)
	updaterURLString := fmt.Sprintf("https://github.com/keys-pub/updater/releases/download/v%s/%s", updaterVersion, updaterFile)
	log.Printf("Extracting %s\n", updaterURLString)
	if err := extractURL(updaterURLString, out, skip); err != nil {
		return err
	}

	return nil
}
