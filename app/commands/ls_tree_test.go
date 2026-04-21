package commands

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestLsTree(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)

	Init()

	os.WriteFile("a.txt", []byte("a"), 0644)
	os.WriteFile("b.txt", []byte("b"), 0644)
	sha, _ := WriteTree(".")

	// Teste com nameOnly = true
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	LsTree(sha, true)

	w.Close()
	os.Stdout = old
	var out bytes.Buffer
	io.Copy(&out, r)

	output := out.String()
	if !strings.Contains(output, "a.txt") || !strings.Contains(output, "b.txt") {
		t.Errorf("Listagem --name-only falhou. Output: %s", output)
	}

	// Teste com nameOnly = false
	r2, w2, _ := os.Pipe()
	os.Stdout = w2

	LsTree(sha, false)

	w2.Close()
	os.Stdout = old
	var out2 bytes.Buffer
	io.Copy(&out2, r2)

	output2 := out2.String()
	if !strings.Contains(output2, "100644") || !strings.Contains(output2, "a.txt") {
		t.Errorf("Listagem detalhada falhou. Output: %s", output2)
	}
}
