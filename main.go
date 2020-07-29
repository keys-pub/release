package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const toolVersion = "0.1.2"

func cmdVersion() *cli.Command {
	return &cli.Command{
		Name: "version",
		Action: func(c *cli.Context) error {
			fmt.Printf("%s\n", toolVersion)
			return nil
		},
	}
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			cmdVersion(),
			cmdDownloadExtra(),
			cmdFixBuild(),
			cmdLatestYAML(),
			cmdPublish(),
			cmdCask(),
			cmdBrew(),
			cmdScoop(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
