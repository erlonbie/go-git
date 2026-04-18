package main

import (
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

// Usage: your_program.sh <command> <arg1> <arg2> ...
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintf(os.Stderr, "Logs from your program will appear here!\n")

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		// TODO: Uncomment the code below to pass the first stage!
		//
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")
	case "cat-file":
		sha := os.Args[3]
		path := fmt.Sprintf(".git/objects/%s/%s", sha[:2], sha[2:])

		file, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %s\n", err)
			os.Exit(1)
		}

		defer func() {
			if err := file.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Error closing file: %s\n", err)
				os.Exit(1)
			}
		}()

		reader := io.Reader(file)

		zlibReader, err := zlib.NewReader(reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating zlib reader: %s\n", err)
			os.Exit(1)
		}

		stream, err := io.ReadAll(zlibReader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading zlib stream: %s\n", err)
			os.Exit(1)
		}

		parts := strings.Split(string(stream), "\x00")
		fmt.Print(parts[1])

		if err := zlibReader.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing zlib reader: %s\n", err)
			os.Exit(1)
		}
	case "hash-object":
		fileName := os.Args[3]
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %s\n", err)
			os.Exit(1)
		}
		defer func() {
			if err := file.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Error closing file: %s\n", err)
				os.Exit(1)
			}
		}()
		
		hasher := sha1.New()
		
		if _, err := io.Copy(hasher, file); err != nil {
			fmt.Fprintf(os.Stderr, "Error while copying file to hasher: %s\n", err)
			os.Exit(1)
		}

		sha := hex.EncodeToString(hasher.Sum(nil))
		for _, dir := range []string{".git/objects/"+sha[:2]} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		objectContents := fmt.Sprintf("blob %d\x00%s", hasher.Size(), string(hasher.Sum(nil)))
		var b strings.Builder
		w, err := zlib.NewWriterLevel(&b, zlib.BestCompression)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating zlib writer: %s\n", err)
			os.Exit(1)
		}
	
		if _, err := w.Write([]byte(objectContents)); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to zlib writer: %s\n", err)
			os.Exit(1)
		}

		if err := w.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing zlib writer: %s\n", err)
			os.Exit(1)
		}
	
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
