package commands

import (
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)

	Init()

	dirs := []string{".git", ".git/objects", ".git/refs"}
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Diretório %s não foi criado", dir)
		}
	}

	if _, err := os.Stat(".git/HEAD"); os.IsNotExist(err) {
		t.Error("Arquivo .git/HEAD não foi criado")
	}
}
