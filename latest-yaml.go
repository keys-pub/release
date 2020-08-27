package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func cmdLatestYAML() *cli.Command {
	return &cli.Command{
		Name: "latest-yaml",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "version",
				Usage: "version",
			},
			&cli.StringFlag{
				Name:  "in",
				Usage: "in",
				Value: "release",
			},
			&cli.StringFlag{
				Name:  "out",
				Usage: "out",
				Value: "release",
			},
			&cli.StringFlag{
				Name:  "platform",
				Usage: "platform",
				Value: runtime.GOOS,
			},
		},
		Action: func(c *cli.Context) error {
			return latestYAML(c.String("platform"), c.String("version"), c.String("in"), c.String("out"))
		},
	}
}

type pkg struct {
	In  string
	Out string
}

func latestYAML(platform string, version string, in string, out string) error {
	if version == "" {
		return errors.Errorf("no version specified")
	}

	var pkgs []pkg
	switch platform {
	case "darwin":
		pkgs = []pkg{
			pkg{
				In:  fmt.Sprintf("Keys-%s-mac.zip", version),
				Out: "latest-mac.yml",
			},
		}
	case "windows":
		pkgs = []pkg{
			pkg{
				In:  fmt.Sprintf("Keys-%s.msi", version),
				Out: "latest-windows.yml",
			},
		}
	case "linux":
		pkgs = []pkg{
			pkg{
				In:  fmt.Sprintf("Keys-%s.AppImage", version),
				Out: "latest-linux.yml",
			},
		}
	}

	for _, pkg := range pkgs {
		inPath := filepath.Join(in, pkg.In)
		outPath := filepath.Join(out, pkg.Out)

		f, err := os.Open(inPath)
		if err != nil {
			return err
		}
		defer f.Close()
		hasher := sha512.New()
		if _, err := io.Copy(hasher, f); err != nil {
			log.Fatal(err)
		}
		encoded := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

		sha512 := encoded
		releaseDate := time.Now().Format(time.RFC3339)

		s := fmt.Sprintf(`version: %s
path: %s
sha512: %s
releaseDate: '%s'
`, version, pkg.In, sha512, releaseDate)

		log.Printf("Writing %s:\n%s\n", outPath, s)
		if err := ioutil.WriteFile(outPath, []byte(s), 0644); err != nil {
			return err
		}
	}

	return nil
}
