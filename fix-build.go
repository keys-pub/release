package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func cmdFixBuild() *cli.Command {
	return &cli.Command{
		Name: "fix-build",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "version",
				Usage: "version",
			},
			&cli.StringFlag{
				Name:  "in",
				Usage: "in",
			},
			&cli.StringFlag{
				Name:  "out",
				Usage: "out",
			},
		},
		Action: func(c *cli.Context) error {
			switch runtime.GOOS {
			case "darwin":
				return fixBuildDarwin(c.String("version"), c.String("in"), c.String("out"))
			default:
				log.Printf("No fixes required for %s", runtime.GOOS)
				return nil
			}
		},
	}
}

func fixBuildDarwin(version string, in string, out string) error {
	if version == "" {
		return errors.Errorf("no version specified")
	}
	if in == "" {
		return errors.Errorf("no in specified")
	}
	if out == "" {
		return errors.Errorf("no out specified")
	}

	zipBase := fmt.Sprintf("Keys-%s-mac.zip", version)
	appDir := filepath.Dir(in)
	app := filepath.Base(in)

	// Chdir to zip dir
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(appDir); err != nil {
		return err
	}

	log.Printf("Re-zipping %s (from %s)\n", zipBase, in)
	cmd := exec.Command("ditto", "-c", "-k", "--sequesterRsrc", "--keepParent", app, zipBase)
	if err := cmd.Run(); err != nil {
		return err
	}

	// Chdir back
	if err := os.Chdir(cwd); err != nil {
		return err
	}

	zipOut := filepath.Join(out, zipBase)
	if _, err := os.Stat(zipOut); err == nil {
		if err := os.Remove(zipOut); err != nil {
			return err
		}
	}
	zipIn := filepath.Join(appDir, zipBase)
	log.Printf("Moving %s to %s\n", zipIn, zipOut)
	return os.Rename(zipIn, zipOut)
}
