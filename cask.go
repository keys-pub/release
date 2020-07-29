package main

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func cmdCask() *cli.Command {
	return &cli.Command{
		Name: "cask",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "version",
				Usage: "version",
			},
		},
		Action: func(c *cli.Context) error {
			return cask(c.String("version"))
		},
	}
}

func cask(version string) error {
	if version == "" {
		return errors.Errorf("no version specified")
	}

	url := fmt.Sprintf("https://github.com/keys-pub/app/releases/download/v%s/Keys-%s-mac.zip", version, version)
	sha256, err := downloadCalculateHash(url)
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

	if err := updateRepo("keys-pub", "homebrew-tap", "Casks/keys.rb", []byte(cask), version); err != nil {
		return err
	}
	return nil

}
