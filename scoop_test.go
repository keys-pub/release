package main

import "testing"

func TestScoop(t *testing.T) {
	url32 := "https://github.com/keys-pub/keys-ext/releases/download/v0.0.48/keys_0.0.48_windows_i386.tar.gz"
	hash, err := downloadCalculateHash(url32)
	if err != nil {
		t.Fatal(err)
	}
	if hash != "d3794365f598bb0c21432e2195fb7c49e8a47402ea5089bcee2a3ec9e7aeb2c9" {
		t.Fatalf("invalid hash %s", hash)
	}
}
