package main

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
)

func cmdPublish() cli.Command {
	return cli.Command{
		Name: "publish",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "version, v",
				Usage: "version",
			},
		},
		Action: func(c *cli.Context) error {
			return publish(c.String("version"))
		},
	}
}

func publish(version string) error {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return errors.Errorf("no GITHUB_TOKEN set")
	}
	if version == "" {
		return errors.Errorf("no version specified")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, "", nil)
	if err != nil {
		return err
	}
	log.Printf("Repos: %+v\n", repos)
	return nil
}
