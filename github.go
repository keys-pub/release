package main

import (
	"context"
	"os"

	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func newGithubClient(ctx context.Context) (*github.Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, errors.Errorf("no GITHUB_TOKEN set")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return client, nil
}
