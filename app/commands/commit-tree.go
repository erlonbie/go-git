package commands

import (
	"crypto/sha1"
	"fmt"
)

func CommitTree(treeSha, parentSha, message string) {
	commitContent := fmt.Sprintf(
		"tree %s\nparent %s\nauthor John Doe <john@example.com> 1700000000 +0000\ncommitter John Doe <john@example.com> 1700000000 +0000\n\n%s",
		treeSha,
		parentSha,
		message,
	)

	header := fmt.Sprintf("commit %d\x00", len(commitContent))
	fullContent := append([]byte(header), []byte(commitContent)...)

	hasher := sha1.New()
	hasher.Write(fullContent)
	sha := fmt.Sprintf("%x", hasher.Sum(nil))
	fmt.Println(sha)
}
