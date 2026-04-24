package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/git-starter-go/app/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		commands.Init()

	case "cat-file":
		if len(os.Args) < 4 {
			fmt.Fprintf(os.Stderr, "usage: mygit cat-file -p <sha>\n")
			os.Exit(1)
		}
		sha := os.Args[3]
		commands.CatFile(sha)

	case "hash-object":
		if len(os.Args) < 4 {
			fmt.Fprintf(os.Stderr, "usage: mygit hash-object -w <file>\n")
			os.Exit(1)
		}
		fileName := os.Args[3]
		commands.HashObject(fileName)

	case "ls-tree":
		if len(os.Args) < 4 {
			fmt.Fprintf(os.Stderr, "usage: mygit ls-tree [--name-only] <sha>\n")
			os.Exit(1)
		}
		nameOnly := os.Args[2] == "--name-only"
		sha := os.Args[len(os.Args)-1]
		commands.LsTree(sha, nameOnly)

	case "write-tree":
		commands.WriteTreeCmd()

	case "commit-tree":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "usage: mygit commit-tree <tree_sha> [-p <parent_sha>] -m <message>\n")
			os.Exit(1)
		}
		treeSha := os.Args[2]
		var parentSha, message string

		for i := 3; i < len(os.Args); i++ {
			if os.Args[i] == "-p" && i+1 < len(os.Args) {
				parentSha = os.Args[i+1]
				i++
			} else if os.Args[i] == "-m" && i+1 < len(os.Args) {
				message = os.Args[i+1]
				i++
			}
		}
		commands.CommitTree(treeSha, parentSha, message)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
