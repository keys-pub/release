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
	"github.com/urfave/cli"
)

func cmdLatestYAML() cli.Command {
	return cli.Command{
		Name: "latest-yaml",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "version, v",
				Usage: "version",
			},
			cli.StringFlag{
				Name:  "in",
				Usage: "in",
				Value: ".",
			},
			cli.StringFlag{
				Name:  "out",
				Usage: "out",
				Value: ".",
			},
		},
		Action: func(c *cli.Context) error {
			return latestYAML(c.String("version"), c.String("in"), c.String("out"))
		},
	}
}

func latestYAML(version string, in string, out string) error {
	if version == "" {
		return errors.Errorf("no version specified")
	}

	var inFile string
	var outFile string
	switch runtime.GOOS {
	case "darwin":
		inFile = fmt.Sprintf("Keys-%s-mac.zip", version)
		outFile = "latest-mac.yml"
	case "windows":
		inFile = fmt.Sprintf("Keys %s.msi", version)
		outFile = "latest-windows.yml"

	}

	inPath := filepath.Join(in, inFile)
	outPath := filepath.Join(out, outFile)

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
`, version, inFile, sha512, releaseDate)

	log.Printf("Writing %s:\n%s\n", outPath, s)
	if err := ioutil.WriteFile(outPath, []byte(s), 0644); err != nil {
		return err
	}

	return nil
}
