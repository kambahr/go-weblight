package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// getDirectories returns sub-directories from a root directory.
func getDirectories(root string) ([]string, error) {
	var dirs []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})
	return dirs, err
}

// fileOrDirExists checks to see if a file or directory exists.
func fileOrDirExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

// getMd5String computes the MD5 and returns string format.
func getMd5String(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

// getRandCharacter returns a string of random characters.
func getRandCharacter() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)

	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x%x%x", b[0:4], b[4:6], b[6:8])
}
