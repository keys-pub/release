package main

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func cmdBrew() *cli.Command {
	return &cli.Command{
		Name: "brew",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "version",
				Usage: "version",
			},
		},
		Action: func(c *cli.Context) error {
			return brew(c.String("version"))
		},
	}
}

func brew(version string) error {
	if version == "" {
		return errors.Errorf("no version specified")
	}

	url := fmt.Sprintf("https://github.com/keys-pub/keys-ext/releases/download/v%s/keys_%s_darwin_x86_64.tar.gz", version, version)
	sha256, err := downloadCalculateHash(url)
	if err != nil {
		return err
	}

	brew := fmt.Sprintf(`class Keys < Formula
  desc "Key management"
  homepage "https://keys.pub"
  version "%s"
  bottle :unneeded
  
  if OS.mac?
    url "%s"
    sha256 "%s"
  elsif OS.linux?
  end
  
  def install
    bin.install "keys"
    bin.install "keysd"
  end
end
`, version, url, sha256)
	log.Printf("%s:\n", brew)

	if err := updateRepo("keys-pub", "homebrew-tap", "keys.rb", []byte(brew), version); err != nil {
		return err
	}
	return nil
}
