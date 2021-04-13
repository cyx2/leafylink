package main

import (
	"testing"
)

func TestUrlHash(t *testing.T) {
	testUrl := "https://www.mongodb.com/2"
	key := urlHash(testUrl)

	expKey := "15c8bc"

	if key != expKey {
		t.Errorf("URL Hash does not match, got %s, expected %s", key, expKey)
	}
}
