package models

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/SadikMR/go-expense-tracker-api/utils"
)

func userFilePath(t *testing.T) string {
	return filepath.Join("data", "users.csv")
}

func writeUsersCSV(t *testing.T, path string, rows [][]string) {
	t.Helper()
	head := []string{"id", "name", "email", "password", "created_at"}
	testWriteCSV(t, path, head, rows)
}

func TestNextUserID(t *testing.T) {
	cleanup := testEnterTempDir(t)
	t.Cleanup(cleanup)

	path := userFilePath(t)
	backup, existed := testBackupFile(t, path)
	t.Cleanup(func() {
		testRestoreFile(t, path, backup, existed)
	})

	writeUsersCSV(t, path, [][]string{
		{"5", "Existing", "existing@example.com", "secret", "2025-01-01T00:00:00Z"},
	})

	id, err := NextUserID()
	if err != nil {
		t.Fatalf("NextUserID failed: %v", err)
	}
	if id != 6 {
		t.Fatalf("expected next user ID 6, got %d", id)
	}
}

func TestGetUserByEmailAndID(t *testing.T) {
	cleanup := testEnterTempDir(t)
	t.Cleanup(cleanup)

	path := userFilePath(t)
	backup, existed := testBackupFile(t, path)
	t.Cleanup(func() {
		testRestoreFile(t, path, backup, existed)
	})

	writeUsersCSV(t, path, [][]string{
		{"1", "Alice", "alice@example.com", "secret", "2025-01-01T00:00:00Z"},
	})

	user, err := GetUserByEmail("alice@example.com")
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}
	if user == nil || user.Name != "Alice" {
		t.Fatalf("expected Alice, got %+v", user)
	}

	missing, err := GetUserByEmail("missing@example.com")
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}
	if missing != nil {
		t.Fatal("expected nil for missing email")
	}

	userByID, err := GetUserByID(1)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if userByID == nil || userByID.Email != "alice@example.com" {
		t.Fatalf("expected alice@example.com, got %+v", userByID)
	}
}

func TestCreateUserAppends(t *testing.T) {
	cleanup := testEnterTempDir(t)
	t.Cleanup(cleanup)

	path := userFilePath(t)
	backup, existed := testBackupFile(t, path)
	t.Cleanup(func() {
		testRestoreFile(t, path, backup, existed)
	})

	writeUsersCSV(t, path, [][]string{
		{"1", "Alice", "alice@example.com", "secret", "2025-01-01T00:00:00Z"},
	})

	user := &User{ID: 2, Name: "Bob", Email: "bob@example.com", Password: "hunter2"}
	if err := CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	created, err := GetUserByID(2)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if created == nil || created.Email != "bob@example.com" {
		t.Fatalf("expected bob@example.com, got %+v", created)
	}
}

func TestEnsureUsersCSV(t *testing.T) {
	tests := []struct {
		name       string
		priorExist bool
		wantError  bool
	}{
		{
			name:       "create when not exists",
			priorExist: false,
			wantError:  false,
		},
		{
			name:       "skip when exists",
			priorExist: true,
			wantError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := testEnterTempDir(t)
			t.Cleanup(cleanup)

			path := "data/users.csv"
			if tt.priorExist {
				if err := os.Mkdir("data", 0755); err != nil && !os.IsExist(err) {
					t.Fatalf("mkdir failed: %v", err)
				}
				if err := utils.RewriteCSV(path, []string{"id"}, [][]string{}); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			err := EnsureUsersCSV()
			if tt.wantError && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if _, err := os.Stat(path); err != nil {
				t.Fatalf("csv file not created: %v", err)
			}
		})
	}
}
