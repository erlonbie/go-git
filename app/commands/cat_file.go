package commands

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"strings"
)

func CatFile(sha string) {
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

	zlibReader, err := zlib.NewReader(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating zlib reader: %s\n", err)
		os.Exit(1)
	}
	defer zlibReader.Close()

	stream, err := io.ReadAll(zlibReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading zlib stream: %s\n", err)
		os.Exit(1)
	}

	parts := strings.Split(string(stream), "\x00")
	if len(parts) > 1 {
		fmt.Print(parts[1])
	}
}
