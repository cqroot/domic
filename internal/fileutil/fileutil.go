package fileutil

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func MD5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func SameContents(a, b string) (bool, error) {
	ha, err := MD5(a)
	if err != nil {
		return false, fmt.Errorf("hash %s: %w", a, err)
	}
	hb, err := MD5(b)
	if err != nil {
		return false, fmt.Errorf("hash %s: %w", b, err)
	}
	return ha == hb, nil
}
