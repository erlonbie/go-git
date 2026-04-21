package commands

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

type TreeEntry struct {
	Mode string
	Name string
	Sha  string
}

func WriteTree(dirPath string) (string, error) {
	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		return "", fmt.Errorf("error reading directory %s: %w", dirPath, err)
	}

	var entries []TreeEntry

	for _, entry := range dirEntries {
		if entry.Name() == ".git" {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return "", err
		}

		fullPath := filepath.Join(dirPath, entry.Name())
		var mode, sha string

		if info.IsDir() {
			mode = "40000"
			sha, err = WriteTree(fullPath)
			if err != nil {
				return "", err
			}
		} else {
			mode = "100644"
			sha, err = WriteBlob(fullPath)
			if err != nil {
				return "", err
			}
		}

		entries = append(entries, TreeEntry{
			Mode: mode,
			Name: entry.Name(),
			Sha:  sha,
		})
	}

	var treeContent bytes.Buffer
	for _, entry := range entries {
		shaBytes, err := hex.DecodeString(entry.Sha)
		if err != nil {
			return "", err
		}
		treeContent.WriteString(fmt.Sprintf("%s %s\x00", entry.Mode, entry.Name))
		treeContent.Write(shaBytes)
	}

	return WriteObject("tree", treeContent.Bytes())
}

func WriteBlob(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %w", filePath, err)
	}
	return WriteObject("blob", content)
}

func WriteObject(objType string, content []byte) (string, error) {
	header := fmt.Sprintf("%s %d\x00", objType, len(content))
	fullContent := append([]byte(header), content...)

	hasher := sha1.New()
	hasher.Write(fullContent)
	sha := hex.EncodeToString(hasher.Sum(nil))

	dir := filepath.Join(".git", "objects", sha[:2])
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	objectPath := filepath.Join(dir, sha[2:])
	if _, err := os.Stat(objectPath); err == nil {
		return sha, nil
	}

	var zlibBuffer bytes.Buffer
	zlibWriter := zlib.NewWriter(&zlibBuffer)
	if _, err := zlibWriter.Write(fullContent); err != nil {
		return "", err
	}
	zlibWriter.Close()

	if err := os.WriteFile(objectPath, zlibBuffer.Bytes(), 0644); err != nil {
		return "", err
	}

	return sha, nil
}
