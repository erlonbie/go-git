package commands

import (
	"fmt"
	"os"
)

func CatFile(sha string) {
	_, content, err := ReadObject(sha)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading object: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(string(content))
}
