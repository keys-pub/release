package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/v29/github"
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
		},
		Action: func(c *cli.Context) error {
			return cask(c.String("version"))
		},
	}
}

func download(url string, file string) error {
	log.Printf("Downloading %s to %s\n", url, file)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(file)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func cask(version string) error {
	if version == "" {
		return errors.Errorf("no version specified")
	}

	file := fmt.Sprintf("Keys-%s-mac.zip", version)
	url := fmt.Sprintf("https://github.com/keys-pub/app/releases/download/v%s/Keys-%s-mac.zip", version, version)
	dl := filepath.Join(os.TempDir(), file)
	if err := download(url, dl); err != nil {
		return err
	}

	sha256, err := sha256FileToHex(dl)
	if err != nil {
		return err
	}

	cask := fmt.Sprintf(`cask 'keys' do
    version '%s'
    sha256 '%s'

    url "%s"
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
`, version, sha256, url)
	log.Printf("%s:\n", cask)

	ctx := context.Background()
	client, err := newGithubClient(ctx)
	if err != nil {
		return err
	}

	log.Printf("Updating repo...\n")
	owner := "keys-pub"
	repo := "homebrew-tap"

	content, _, _, err := client.Repositories.GetContents(ctx, owner, repo, "Casks/keys.rb", &github.RepositoryContentGetOptions{})
	if err != nil {
		return err
	}
	b, err := base64.StdEncoding.DecodeString(*content.Content)
	if err != nil {
		return err
	}
	if bytes.Equal(b, []byte(cask)) {
		log.Printf("Content already exists")
		return nil
	}

	msg := fmt.Sprintf("Update Casks/keys.rb (%s)", version)
	opts := &github.RepositoryContentFileOptions{
		Message: &msg,
		Content: []byte(cask),
		SHA:     content.SHA,
	}
	if _, _, err := client.Repositories.UpdateFile(ctx, owner, repo, "Casks/keys.rb", opts); err != nil {
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
