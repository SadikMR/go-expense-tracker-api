package models

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/SadikMR/go-expense-tracker-api/utils"
)

func testRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatal("failed to locate repo root")
	return ""
}

func testBackupFile(t *testing.T, path string) ([]byte, bool) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false
		}
		t.Fatalf("failed to back up file %q: %v", path, err)
	}
	return content, true
}

func testRestoreFile(t *testing.T, path string, backup []byte, existed bool) {
	t.Helper()
	if existed {
		if err := os.WriteFile(path, backup, 0644); err != nil {
			t.Fatalf("failed to restore file %q: %v", path, err)
		}
		return
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed to remove file %q: %v", path, err)
	}
}

func testWriteCSV(t *testing.T, path string, header []string, rows [][]string) {
	t.Helper()
	if err := utils.RewriteCSV(path, header, rows); err != nil {
		t.Fatalf("failed to write CSV %q: %v", path, err)
	}
}

func RewriteCSV(path string, header []string, rows [][]string) error {
	return utils.RewriteCSV(path, header, rows)
}

func testEnterTempDir(t *testing.T) func() {
	t.Helper()
	tmp, err := os.MkdirTemp("", "go-expense-tracker-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	if err := os.Mkdir(filepath.Join(tmp, "data"), 0755); err != nil {
		os.RemoveAll(tmp)
		t.Fatalf("failed to create temp data dir: %v", err)
	}
	orig, err := os.Getwd()
	if err != nil {
		os.RemoveAll(tmp)
		t.Fatalf("failed to get current working directory: %v", err)
	}
	if err := os.Chdir(tmp); err != nil {
		os.RemoveAll(tmp)
		t.Fatalf("failed to change directory to temp dir: %v", err)
	}
	return func() {
		if err := os.Chdir(orig); err != nil {
			t.Fatalf("failed to restore working directory: %v", err)
		}
		if err := os.RemoveAll(tmp); err != nil {
			t.Fatalf("failed to remove temp dir: %v", err)
		}
	}
}
