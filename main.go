package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

const version = "0.1.1"

func cmdVersion() cli.Command {
	return cli.Command{
		Name: "version",
		Action: func(c *cli.Context) error {
			fmt.Printf("%s\n", version)
			return nil
		},
	}
}

func main() {
	app := &cli.App{
		Commands: []cli.Command{
			cmdVersion(),
			cmdDownloadExtra(),
			cmdFixBuild(),
			cmdLatestYAML(),
			cmdPublish(),
			cmdCask(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
