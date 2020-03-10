package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func cmdCask() cli.Command {
	return cli.Command{
		Name: "cask",
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
				Value: filepath.Join(".", "keys.rb"),
			},
		},
		Action: func(c *cli.Context) error {
			return cask(c.String("version"), c.String("in"), c.String("out"))
		},
	}
}

func cask(version string, in string, out string) error {
	if version == "" {
		return errors.Errorf("no version specified")
	}

	inFile := fmt.Sprintf("Keys-%s-mac.zip", version)
	inPath := filepath.Join(in, inFile)

	sha256, err := sha256FileToHex(inPath)
	if err != nil {
		return err
	}

	cask := fmt.Sprintf(`cask 'keys' do
    version '%s'
    sha256 '%s'

    url "https://github.com/keys-pub/app/releases/download/v#{version}/Keys-#{version}-mac.zip"
    name 'Keys'
    homepage 'https://keys.pub'

    depends_on macos: '>= :sierra'

    app 'Keys.app'

    uninstall delete: [
        '/usr/local/bin/keys'
    ]

    zap trash: [
        '~/Library/Application Support/Keys',
        '~/Library/Caches/Keys',
        '~/Library/Logs/Keys',
        '~/Library/Preferences/pub.Keys.plist',
    ]
end
`, version, sha256)

	log.Printf("Writing %s:\n%s\n", out, cask)
	if err := ioutil.WriteFile(out, []byte(cask), 0644); err != nil {
		return err
	}

	return nil
}

func downloadURLString(url string, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func sha256FileToHex(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
