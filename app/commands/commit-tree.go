package commands

import (
	"fmt"
	"os"
	"time"
)

func CommitTree(treeSha, parentSha, message string) {
	now := time.Now().Unix()
	author := fmt.Sprintf("John Doe <john@example.com> %d +0000", now)

	content := fmt.Sprintf("tree %s\n", treeSha)
	if parentSha != "" {
		content += fmt.Sprintf("parent %s\n", parentSha)
	}
	content += fmt.Sprintf("author %s\n", author)
	content += fmt.Sprintf("committer %s\n", author)
	content += fmt.Sprintf("\n%s\n", message)

	sha, err := WriteObject("commit", []byte(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing commit: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(sha)
}
