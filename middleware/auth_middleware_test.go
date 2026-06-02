package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/SadikMR/go-expense-tracker-api/utils"
	"github.com/beego/beego/v2/server/web/context"
)

func repoRoot(t *testing.T) string {
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

func dataDir(t *testing.T) string {
	t.Helper()
	if dir := os.Getenv("DATA_DIR"); dir != "" {
		return dir
	}
	return filepath.Join(repoRoot(t), "data")
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

func enterRepoRoot(t *testing.T) func() {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	if err := os.Chdir(repoRoot(t)); err != nil {
		t.Fatalf("failed to change directory to repo root: %v", err)
	}
	return func() {
		if err := os.Chdir(orig); err != nil {
			t.Fatalf("failed to restore working directory: %v", err)
		}
	}
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

func newTestContext(t *testing.T, req *http.Request) *context.Context {
	t.Helper()
	ctx := context.NewContext()
	ctx.Request = req
	rec := httptest.NewRecorder()
	ctx.ResponseWriter = &context.Response{ResponseWriter: rec}
	ctx.Input.Reset(ctx)
	ctx.Output.Reset(ctx)
	return ctx
}

func TestAuthMiddlewareMissingHeader(t *testing.T) {
	cleanup := enterRepoRoot(t)
	t.Cleanup(cleanup)

	path := usersCSVPath(t)
	backup, existed := backupFile(t, path)
	t.Cleanup(func() {
		restoreFile(t, path, backup, existed)
	})

	writeUsersCSV(t, path, [][]string{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	ctx := newTestContext(t, req)

	AuthMiddleware(ctx)
	if !ctx.ResponseWriter.Started {
		t.Fatal("expected response to be started")
	}
	if ctx.ResponseWriter.Status != 401 {
		t.Fatalf("expected status 401, got %d", ctx.ResponseWriter.Status)
	}
}

func TestAuthMiddlewareInvalidHeader(t *testing.T) {
	cleanup := enterRepoRoot(t)
	t.Cleanup(cleanup)

	path := usersCSVPath(t)
	backup, existed := backupFile(t, path)
	t.Cleanup(func() {
		restoreFile(t, path, backup, existed)
	})

	writeUsersCSV(t, path, [][]string{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	req.Header.Set("X-User-ID", "abc")
	ctx := newTestContext(t, req)

	AuthMiddleware(ctx)
	if ctx.ResponseWriter.Status != 401 {
		t.Fatalf("expected status 401, got %d", ctx.ResponseWriter.Status)
	}
}

func TestAuthMiddlewareNonexistentUser(t *testing.T) {
	cleanup := enterRepoRoot(t)
	t.Cleanup(cleanup)

	path := usersCSVPath(t)
	backup, existed := backupFile(t, path)
	t.Cleanup(func() {
		restoreFile(t, path, backup, existed)
	})

	writeUsersCSV(t, path, [][]string{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	req.Header.Set("X-User-ID", "999")
	ctx := newTestContext(t, req)

	AuthMiddleware(ctx)
	if ctx.ResponseWriter.Status != 401 {
		t.Fatalf("expected status 401, got %d", ctx.ResponseWriter.Status)
	}
}

func TestAuthMiddlewareValidHeader(t *testing.T) {
	cleanup := enterRepoRoot(t)
	t.Cleanup(cleanup)

	path := usersCSVPath(t)
	backup, existed := backupFile(t, path)
	t.Cleanup(func() {
		restoreFile(t, path, backup, existed)
	})

	writeUsersCSV(t, path, [][]string{
		{"1", "Alice", "alice@example.com", "secret", "2025-01-01T00:00:00Z"},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	req.Header.Set("X-User-ID", "1")
	ctx := newTestContext(t, req)

	AuthMiddleware(ctx)
	if ctx.ResponseWriter.Started {
		t.Fatal("expected middleware to allow request and not write a response")
	}
	if ctx.ResponseWriter.Status != 0 {
		t.Fatalf("expected status 0, got %d", ctx.ResponseWriter.Status)
	}
}
