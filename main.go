package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Commands: []cli.Command{
			cmdDownloadExtra(),
			cmdFixBuild(),
			cmdLatestYAML(),
			cmdPublish(),
			cmdCask(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
