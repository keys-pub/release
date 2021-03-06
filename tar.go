package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func contains(s string, from []string) bool {
	for _, i := range from {
		if s == i {
			return true
		}
	}
	return false
}

func extractURL(url string, out string, skip []string) ([]string, error) {
	// f, err := os.Open(srcFile)
	// if err != nil {
	// 	return err
	// }
	// defer f.Close()
	if out == "" {
		out = "."
	}
	if err := os.MkdirAll(out, 0755); err != nil {
		return nil, err
	}

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	reader, err := gzip.NewReader(response.Body)
	if err != nil {
		return nil, err
	}

	tarReader := tar.NewReader(reader)

	files := []string{}
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if contains(header.Name, skip) {
				continue
			}
			if err := os.Mkdir(header.Name, 0755); err != nil {
				return nil, err
			}
		case tar.TypeReg:
			p := filepath.Join(out, header.Name)
			if contains(header.Name, skip) {
				continue
			}
			files = append(files, p)
			log.Printf("%s\n", p)
			outFile, err := os.Create(p)
			if err != nil {
				return nil, err
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return nil, err
			}
			outFile.Close()
		default:
			return nil, errors.Errorf("unknown type: %b in %s", header.Typeflag, header.Name)
		}
	}
	return files, nil
}
