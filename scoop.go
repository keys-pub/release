package main

import (
	"fmt"
	"log"

	"github.com/urfave/cli/v2"
)

func cmdScoop() *cli.Command {
	return &cli.Command{
		Name: "scoop",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "version",
				Usage: "version",
			},
		},
		Action: func(c *cli.Context) error {
			return scoop(c.String("version"))
		},
	}
}

func scoop(version string) error {

	url32 := fmt.Sprintf("https://github.com/keys-pub/keys-ext/releases/download/v%s/keys_%s_windows_i386.tar.gz", version, version)
	hash32, err := downloadCalculateHash(url32)
	if err != nil {
		return err
	}
	url64 := fmt.Sprintf("https://github.com/keys-pub/keys-ext/releases/download/v%s/keys_%s_windows_x86_64.tar.gz", version, version)
	hash64, err := downloadCalculateHash(url64)
	if err != nil {
		return err
	}

	scoop := fmt.Sprintf(`
{
    "version": "%s",
    "architecture": {
        "32bit": {
            "url": "%s",
            "bin": [
                "keys.exe",
                "keysd.exe"
            ],
            "hash": "%s"
        },
        "64bit": {
            "url": "%s",
            "bin": [
                "keys.exe",
                "keysd.exe"
            ],
            "hash": "%s"
        }
    },
    "homepage": "https://keys.pub",
    "license": "MIT"
}
`, version, url32, hash32, url64, hash64)

	log.Printf("%s:\n", scoop)

	if err := updateRepo("keys-pub", "scoop-bucket", "keys.json", []byte(scoop), version); err != nil {
		return err
	}
	return nil
}
