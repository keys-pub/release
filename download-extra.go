package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

const updaterVersion = "0.2.2"

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

	if out == "" {
		out = "."
	}

	log.Printf("Extracting %s\n", keysURLString)
	if _, err := extractURL(keysURLString, out, skip); err != nil {
		return err
	}

	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		if err := makeExecutable([]string{filepath.Join(out, "keys"), filepath.Join(out, "keysd")}); err != nil {
			return err
		}
	}

	fido2File := fmt.Sprintf("fido2_%s_%s_x86_64.tar.gz", version, platform)
	fido2URLString := fmt.Sprintf("https://github.com/keys-pub/keysd/releases/download/v%s/%s", version, fido2File)

	if out == "" {
		out = "."
	}

	log.Printf("Extracting %s\n", fido2URLString)
	if _, err := extractURL(fido2URLString, out, skip); err != nil {
		return err
	}

	updaterFile := fmt.Sprintf("updater_%s_%s_x86_64.tar.gz", updaterVersion, platform)
	updaterURLString := fmt.Sprintf("https://github.com/keys-pub/updater/releases/download/v%s/%s", updaterVersion, updaterFile)
	log.Printf("Extracting %s\n", updaterURLString)
	if _, err := extractURL(updaterURLString, out, skip); err != nil {
		return err
	}

	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		if err := makeExecutable([]string{filepath.Join(out, "updater")}); err != nil {
			return err
		}
	}

	return nil
}

func makeExecutable(paths []string) error {
	for _, p := range paths {
		log.Printf("chmod 0755 %s\n", p)
		if err := os.Chmod(p, 0755); err != nil {
			return err
		}
	}
	return nil
}
