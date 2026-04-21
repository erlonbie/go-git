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
		
		content, err := io.ReadAll(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
			os.Exit(1)
		}
		
		header := fmt.Sprintf("blob %d\x00", len(content))
		hasher := sha1.New()
		hasher.Write([]byte(header))
		hasher.Write(content)
		
		sha := hex.EncodeToString(hasher.Sum(nil))
		fmt.Println(sha)
		
		dir := ".git/objects/" + sha[:2]
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
		}

		var stringBuilder strings.Builder
		zlibWriter := zlib.NewWriter(&stringBuilder)

		objectContents := header + string(content)
		if _, err := zlibWriter.Write([]byte(objectContents)); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to zlib writer: %s\n", err)
			os.Exit(1)
		}

		if err := zlibWriter.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing zlib writer: %s\n", err)
			os.Exit(1)
		}

		objectPath := fmt.Sprintf(".git/objects/%s/%s", sha[:2], sha[2:])
		if err := os.WriteFile(objectPath, []byte(stringBuilder.String()), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing object file: %s\n", err)
			os.Exit(1)
		}
	case "ls-tree":
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

		nullIndex := strings.Index(string(stream), "\x00")
		data := stream[nullIndex+1:]
		i := 0

		for i < len(data) {
			spaceIndex := -1
			for j := i; j < len(data); j++ {
				if data[j] == ' ' {
					spaceIndex = j
					break
				}
			}

			mode := string(data[i:spaceIndex])
			i = spaceIndex + 1

			nullIndex := -1
			for j := i; j < len(data); j++ {
				if data[j] == 0 {
					nullIndex = j
					break
				}
			}

			name := string(data[i:nullIndex])
			i = nullIndex + 1

			shaBytes := data[i : i+20]
			i += 20

			sha := fmt.Sprintf("%x", shaBytes)

			if os.Args[2] == "--name-only" {
				fmt.Println(name)
			} else {
				fmt.Printf("%s %s %s\n", mode, sha, name)
			}
		}

		if err := zlibReader.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing zlib reader: %s\n", err)
			os.Exit(1)
		}
	case "write-tree":
		entries := os.Args[2:]

		var stringBuilder strings.Builder
		for _, entry := range entries {
			parts := strings.SplitN(entry, " ", 3)
			mode := parts[0]
			name := parts[1]
			sha := parts[2]

			shaBytes, err := hex.DecodeString(sha)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error decoding SHA: %s\n", err)
				os.Exit(1)
			}

			fmt.Fprintf(&stringBuilder, "%s %s\x00", mode, name)
			stringBuilder.Write(shaBytes)
		}

		content := stringBuilder.String()
		header := fmt.Sprintf("tree %d\x00", len(content))
		hasher := sha1.New()
		hasher.Write([]byte(header))
		hasher.Write([]byte(content))

		sha := hex.EncodeToString(hasher.Sum(nil))
		fmt.Println(sha)

		dir := ".git/objects/" + sha[:2]
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
		}

		var zlibBuilder strings.Builder
		zlibWriter := zlib.NewWriter(&zlibBuilder)

		objectContents := header + content
		if _, err := zlibWriter.Write([]byte(objectContents)); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to zlib writer: %s\n", err)
			os.Exit(1)
		}

		if err := zlibWriter.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing zlib writer: %s\n", err)
			os.Exit(1)
		}

		objectPath := fmt.Sprintf(".git/objects/%s/%s", sha[:2], sha[2:])
		if err := os.WriteFile(objectPath, []byte(zlibBuilder.String()), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing object file: %s\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
