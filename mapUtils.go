package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func urlHash(longUrl string) string {
	h := md5.New()
	io.WriteString(h, longUrl)
	return hex.EncodeToString(h.Sum(nil))[:6]
}

// TODO: De-Dup Keys
// TODO: De-Dup longURLs
