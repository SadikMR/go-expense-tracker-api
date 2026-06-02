package services

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/SadikMR/go-expense-tracker-api/models"
	"github.com/SadikMR/go-expense-tracker-api/utils"
)

func serviceRepoRoot(t *testing.T) string {
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

func usersCSVPath(t *testing.T) string {
	return filepath.Join(dataDir(t), "users.csv")
}

func backupFile(t *testing.T, path string) ([]byte, bool) {
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

func restoreFile(t *testing.T, path string, backup []byte, existed bool) {
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

func writeUsersCSV(t *testing.T, path string, rows [][]string) {
	t.Helper()
	head := []string{"id", "name", "email", "password", "created_at"}
	if err := utils.RewriteCSV(path, head, rows); err != nil {
		t.Fatalf("failed to write users.csv: %v", err)
	}
}

func TestRegisterAndLogin(t *testing.T) {
	path := usersCSVPath(t)
	backup, existed := backupFile(t, path)
	t.Cleanup(func() {
		restoreFile(t, path, backup, existed)
	})

	writeUsersCSV(t, path, [][]string{
		{"1", "Existing", "existing@example.com", "secret", "2025-01-01T00:00:00Z"},
	})

	if err := Register("New User", "new@example.com", "password123"); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	user, err := models.GetUserByEmail("new@example.com")
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}
	if user == nil || user.Name != "New User" {
		t.Fatalf("expected New User, got %+v", user)
	}

	if err := Register("New User", "new@example.com", "password123"); err != ErrEmailExists {
		t.Fatalf("expected ErrEmailExists, got %v", err)
	}

	authUser, err := Login("new@example.com", "password123")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if authUser == nil || authUser.Email != "new@example.com" {
		t.Fatalf("expected login user new@example.com, got %+v", authUser)
	}

	_, err = Login("new@example.com", "wrong")
	if err != ErrInvalidCreds {
		t.Fatalf("expected ErrInvalidCreds, got %v", err)
	}
}
