package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

// UrlHash computes the MD5 hash of a longUrl
// and returns the first 6 characters
func urlHash(longUrl string) string {
	h := md5.New()
	io.WriteString(h, longUrl)
	return hex.EncodeToString(h.Sum(nil))[:6]
}
