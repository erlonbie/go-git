package commands

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

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

	refsURL := fmt.Sprintf("%s/info/refs?service=git-upload-pack", repoURL)
	resp, err := http.Get(refsURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching references: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "HTTP error discovering references: %s\n", resp.Status)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading refs response: %v\n", err)
		os.Exit(1)
	}

	refs := parsePktLines(body)
	
	var headSha string
	for _, ref := range refs {
		parts := strings.Split(ref, "\x00")
		refInfo := strings.Fields(parts[0])
		
		if len(refInfo) >= 2 {
			sha := refInfo[0]
			name := refInfo[1]
			
			if name == "HEAD" {
				headSha = sha
				break
			}
		}
	}

	if headSha == "" {
		fmt.Fprintf(os.Stderr, "Could not find the HEAD reference on the server.\n")
		os.Exit(1)
	}
}

func parsePktLines(data []byte) []string {
	var lines []string
	buf := bytes.NewReader(data)

	for {
		lengthBytes := make([]byte, 4)
		_, err := io.ReadFull(buf, lengthBytes)
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		length, err := strconv.ParseInt(string(lengthBytes), 16, 32)
		if err != nil || length == 0 {
			continue
		}

		payload := make([]byte, length-4)
		_, err = io.ReadFull(buf, payload)
		if err != nil {
			break
		}

		line := strings.TrimSuffix(string(payload), "\n")
		
		if !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}

	return lines
}
