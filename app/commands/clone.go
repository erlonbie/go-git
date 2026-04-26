package commands

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/codecrafters-io/git-starter-go/app/plumbing"
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

	packURL := fmt.Sprintf("%s/git-upload-pack", repoURL)

	requestBody := fmt.Sprintf("0032want %s\n00000009done\n", headSha)

	req, err := http.NewRequest("POST", packURL, bytes.NewBufferString(requestBody))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating request for packfile: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/x-git-upload-pack-request")

	client := &http.Client{}
	packResp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error requesting packfile: %v\n", err)
		os.Exit(1)
	}
	defer packResp.Body.Close()

	if packResp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "HTTP error downloading packfile: %s\n", packResp.Status)
		os.Exit(1)
	}

	packData, err := io.ReadAll(packResp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading packfile body: %v\n", err)
		os.Exit(1)
	}

	packIndex := bytes.Index(packData, []byte("PACK"))
	if packIndex == -1 {
		fmt.Fprintf(os.Stderr, "'PACK' signature not found in response\n")
		os.Exit(1)
	}

	packfileContent := packData[packIndex:]

	objects, err := plumbing.ParsePackfile(packfileContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing packfile: %v\n", err)
		os.Exit(1)
	}

	for _, obj := range objects {
		_, err := WriteObject(obj.Type, obj.Content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving extracted object: %v\n", err)
			os.Exit(1)
		}
	}

	_, commitContent, err := ReadObject(headSha)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading HEAD commit: %v\n", err)
		os.Exit(1)
	}

	commitParts := strings.SplitN(string(commitContent), "\n", 2)
	treeSha := strings.TrimPrefix(commitParts[0], "tree ")

	if err := checkoutTree(treeSha, "."); err != nil {
		fmt.Fprintf(os.Stderr, "Error checking out tree: %v\n", err)
		os.Exit(1)
	}
}

func checkoutTree(treeSha, currentPath string) error {
	_, treeContent, err := ReadObject(treeSha)
	if err != nil {
		return fmt.Errorf("error reading tree %s: %w", treeSha, err)
	}

	i := 0
	for i < len(treeContent) {
		spaceIndex := bytes.IndexByte(treeContent[i:], ' ')
		if spaceIndex == -1 {
			break
		}
		spaceIndex += i

		mode := string(treeContent[i:spaceIndex])
		i = spaceIndex + 1

		nullIndex := bytes.IndexByte(treeContent[i:], 0)
		if nullIndex == -1 {
			break
		}
		nullIndex += i

		name := string(treeContent[i:nullIndex])
		i = nullIndex + 1

		shaBytes := treeContent[i : i+20]
		i += 20

		shaHex := fmt.Sprintf("%x", shaBytes)
		fullPath := filepath.Join(currentPath, name)

		if mode == "40000" {
			if err := os.MkdirAll(fullPath, 0755); err != nil {
				return err
			}
			if err := checkoutTree(shaHex, fullPath); err != nil {
				return err
			}
		} else {
			_, fileContent, err := ReadObject(shaHex)
			if err != nil {
				return err
			}
			
			perm := os.FileMode(0644)
			if mode == "100755" {
				perm = 0755
			}
			
			if err := os.WriteFile(fullPath, fileContent, perm); err != nil {
				return err
			}
		}
	}

	return nil
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
