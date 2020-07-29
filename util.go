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
	"net/url"
	"os"
	"path/filepath"

	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

func download(urs string) (string, error) {
	ur, err := url.Parse(urs)
	if err != nil {
		return "", err
	}
	path := filepath.Join(os.TempDir(), filepath.Base(ur.Path))

	log.Printf("Downloading %s to %s\n", urs, path)
	resp, err := http.Get(urs)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if resp.StatusCode != 200 {
		return "", errors.Errorf("http status %d", resp.StatusCode)
	}

	_, err = io.Copy(out, resp.Body)
	return path, err
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

func downloadCalculateHash(url string) (string, error) {
	path, err := download(url)
	if err != nil {
		return "", err
	}
	defer func() { _ = os.Remove(path) }()

	log.Printf("Calculating sha256...\n")
	sha256, err := sha256FileToHex(path)
	if err != nil {
		return "", err
	}
	log.Printf("Calculated sha256: %s\n", sha256)
	return sha256, nil
}

func updateRepo(owner string, repo string, file string, data []byte, version string) error {
	ctx := context.Background()
	client, err := newGithubClient(ctx)
	if err != nil {
		return err
	}

	log.Printf("Updating repo...\n")

	content, _, _, err := client.Repositories.GetContents(ctx, owner, repo, file, &github.RepositoryContentGetOptions{})
	if err != nil {
		return err
	}
	b, err := base64.StdEncoding.DecodeString(*content.Content)
	if err != nil {
		return err
	}
	if bytes.Equal(b, data) {
		log.Printf("Content already exists")
		return nil
	}

	msg := fmt.Sprintf("Update %s (%s)", file, version)
	opts := &github.RepositoryContentFileOptions{
		Message: &msg,
		Content: data,
		SHA:     content.SHA,
	}
	if _, _, err := client.Repositories.UpdateFile(ctx, owner, repo, file, opts); err != nil {
		return err
	}

	return nil
}
