package commands

import (
	"fmt"
	"os"
)

func WriteTreeCmd() {
	sha, err := WriteTree(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing tree: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(sha)
}
