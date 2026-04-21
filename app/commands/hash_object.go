package commands

import (
	"fmt"
	"os"
)

func HashObject(filePath string) {
	sha, err := WriteBlob(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error hashing object: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(sha)
}
