package commands

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"strings"
)

func LsTree(sha string, nameOnly bool) {
	path := fmt.Sprintf(".git/objects/%s/%s", sha[:2], sha[2:])

	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

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

		shaHex := fmt.Sprintf("%x", shaBytes)

		if nameOnly {
			fmt.Println(name)
		} else {
			fmt.Printf("%s %s %s\n", mode, shaHex, name)
		}
	}
}
