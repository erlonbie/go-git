package commands

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestCommitTree(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	Init()

	treeSha := "d8329fc1cc938780ffdd9f94e0d364e0ea74f579"
	message := "Test commit message"

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	CommitTree(treeSha, "", message)

	w.Close()
	os.Stdout = old
	var out bytes.Buffer
	io.Copy(&out, r)

	sha := strings.TrimSpace(out.String())
	if len(sha) != 40 {
		t.Errorf("Esperava SHA de 40 chars, recebeu %d", len(sha))
	}

	r2, w2, _ := os.Pipe()
	os.Stdout = w2

	CatFile(sha)

	w2.Close()
	os.Stdout = old
	var out2 bytes.Buffer
	io.Copy(&out2, r2)

	if !strings.Contains(out2.String(), message) {
		t.Errorf("O commit não contém a mensagem esperada. Conteúdo: %s", out2.String())
	}
	if !strings.Contains(out2.String(), "tree "+treeSha) {
		t.Errorf("O commit não contém o tree SHA correto")
	}
}
