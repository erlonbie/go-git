package commands

import (
	"fmt"
	"os"
func Clone(repoURL, targetDir string) {
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating target directory: %v\n", err)
		os.Exit(1)
	}

	if err := os.Chdir(targetDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error entering the target directory: %v\n", err)
		os.Exit(1)
	}

	Init()
}
