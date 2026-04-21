package commands

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestCatFile(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)

	Init()

	testFile := "test.txt"
	content := "cat content"
	os.WriteFile(testFile, []byte(content), 0644)
	sha, _ := WriteBlob(testFile)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	CatFile(sha)

	w.Close()
	os.Stdout = old
	var out bytes.Buffer
	io.Copy(&out, r)

	if out.String() != content {
		t.Errorf("Esperava '%s', recebeu '%s'", content, out.String())
	}
}
