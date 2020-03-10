package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func cmdPublish() cli.Command {
	return cli.Command{
		Name: "publish",
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
		},
		Action: func(c *cli.Context) error {
			return publish(c.String("version"), c.String("in"))
		},
	}
}

func publish(version string, in string) error {
	if version == "" {
		return errors.Errorf("no version specified")
	}

	owner := "keys-pub"
	repo := "app"
	var upload []string

	switch runtime.GOOS {
	case "darwin":
		upload = []string{
			fmt.Sprintf("Keys-%s-mac.zip", version),
			fmt.Sprintf("Keys-%s.dmg", version),
			"latest-mac.yml",
		}
	case "windows":
		upload = []string{
			fmt.Sprintf("Keys %s.msi", version),
			"latest-windows.yml",
		}
	}

	ctx := context.Background()
	client, err := newGithubClient(ctx)
	if err != nil {
		return err
	}

	releases, _, err := client.Repositories.ListReleases(ctx, owner, repo, nil)
	if err != nil {
		return err
	}
	var release *github.RepositoryRelease
	for _, r := range releases {
		if *r.Name == version {
			release = r
		}
	}

	if release != nil {
		log.Printf("Found release: %s\n", *release.Name)
	} else {
		commits, _, err := client.Repositories.ListCommits(ctx, owner, repo, nil)
		if len(commits) == 0 {
			return errors.Errorf("no commits")
		}
		commit := commits[0].SHA
		log.Printf("Commit: %s", *commit)

		log.Printf("Creating release: %s\n", version)
		tag := fmt.Sprintf("v%s", version)
		draft := true

		release = &github.RepositoryRelease{
			Name:            &version,
			TagName:         &tag,
			TargetCommitish: commit,
			Draft:           &draft,
		}
		r, _, err := client.Repositories.CreateRelease(ctx, owner, repo, release)
		if err != nil {
			return err
		}
		release = r
	}

	id := *release.ID
	log.Printf("Release ID: %d\n", id)
	assets, _, err := client.Repositories.ListReleaseAssets(ctx, owner, repo, id, nil)
	if err != nil {
		return err
	}
	existing := []string{}
	for _, a := range assets {
		existing = append(existing, *a.Name)
	}

	for _, u := range upload {
		if contains(u, existing) {
			log.Printf("Asset already exists %s", u)
			continue
		}

		p := filepath.Join(in, u)
		log.Printf("Uploading asset %s", p)
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()
		opts := &github.UploadOptions{
			Name: u,
		}
		asset, _, err := client.Repositories.UploadReleaseAsset(ctx, owner, repo, id, opts, f)
		if err != nil {
			return err
		}
		log.Printf("Uploaded (%d)\n", *asset.ID)
		f.Close()
	}

	return nil
}
