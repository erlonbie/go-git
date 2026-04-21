package commands

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestHashObject(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)

	Init()

	testFile := "hello.txt"
	content := "hello world"
	os.WriteFile(testFile, []byte(content), 0644)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	HashObject(testFile)

	w.Close()
	os.Stdout = old
	var out bytes.Buffer
	io.Copy(&out, r)

	expectedSha := "95d09f2b10159347eece71399a7e2e907ea3df4f"
	if !strings.Contains(out.String(), expectedSha) {
		t.Errorf("Esperava SHA %s, recebeu %s", expectedSha, out.String())
	}

	objectPath := ".git/objects/95/d09f2b10159347eece71399a7e2e907ea3df4f"
	if _, err := os.Stat(objectPath); os.IsNotExist(err) {
		t.Errorf("Objeto não foi salvo em %s", objectPath)
	}
}
