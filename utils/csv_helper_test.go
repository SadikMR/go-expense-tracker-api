package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func testTempDir(t *testing.T) (string, func()) {
	t.Helper()
	tmp, err := os.MkdirTemp("", "csv-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	return tmp, func() {
		os.RemoveAll(tmp)
	}
}

func TestReadCSV(t *testing.T) {
	tmp, cleanup := testTempDir(t)
	defer cleanup()

	tests := []struct {
		name      string
		setup     func(string) error
		wantRows  int
		wantError bool
	}{
		{
			name: "read valid csv with multiple rows",
			setup: func(p string) error {
				return RewriteCSV(p, []string{"id", "name"}, [][]string{
					{"1", "Alice"},
					{"2", "Bob"},
				})
			},
			wantRows:  2,
			wantError: false,
		},
		{
			name: "read csv with header only",
			setup: func(p string) error {
				return RewriteCSV(p, []string{"id", "name"}, [][]string{})
			},
			wantRows:  0,
			wantError: false,
		},
		{
			name: "read non-existent file",
			setup: func(p string) error {
				return nil
			},
			wantRows:  0,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testPath := filepath.Join(tmp, "test-"+tt.name+".csv")
			if err := tt.setup(testPath); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			rows, err := ReadCSV(testPath)
			if tt.wantError && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(rows) != tt.wantRows {
				t.Fatalf("expected %d rows, got %d", tt.wantRows, len(rows))
			}
		})
	}
}

func TestAppendCSV(t *testing.T) {
	tmp, cleanup := testTempDir(t)
	defer cleanup()

	tests := []struct {
		name      string
		initial   [][]string
		rows      [][]string
		wantError bool
		validate  func(*testing.T, string)
	}{
		{
			name:    "append to empty csv",
			initial: [][]string{},
			rows:    [][]string{{"1", "Alice"}},
			validate: func(t *testing.T, p string) {
				rows, err := ReadCSV(p)
				if err != nil {
					t.Fatalf("read failed: %v", err)
				}
				if len(rows) != 1 || rows[0][0] != "1" {
					t.Fatalf("expected 1 row with id 1, got %v", rows)
				}
			},
		},
		{
			name:    "append multiple rows",
			initial: [][]string{{"1", "Alice"}},
			rows:    [][]string{{"2", "Bob"}, {"3", "Charlie"}},
			validate: func(t *testing.T, p string) {
				rows, err := ReadCSV(p)
				if err != nil {
					t.Fatalf("read failed: %v", err)
				}
				if len(rows) != 3 {
					t.Fatalf("expected 3 rows, got %d", len(rows))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testPath := filepath.Join(tmp, "append-"+tt.name+".csv")
			header := []string{"id", "name"}

			if len(tt.initial) > 0 || tt.name == "append to empty csv" {
				if err := RewriteCSV(testPath, header, tt.initial); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			for _, row := range tt.rows {
				if err := AppendCSV(testPath, row); err != nil && !tt.wantError {
					t.Fatalf("AppendCSV failed: %v", err)
				}
			}

			if tt.validate != nil {
				tt.validate(t, testPath)
			}
		})
	}
}

func TestRewriteCSV(t *testing.T) {
	tmp, cleanup := testTempDir(t)
	defer cleanup()

	tests := []struct {
		name      string
		path      string
		header    []string
		rows      [][]string
		wantError bool
	}{
		{
			name:      "rewrite with valid header and rows",
			path:      filepath.Join(tmp, "rewrite1.csv"),
			header:    []string{"id", "name", "email"},
			rows:      [][]string{{"1", "Alice", "alice@example.com"}},
			wantError: false,
		},
		{
			name:      "rewrite with empty rows",
			path:      filepath.Join(tmp, "rewrite2.csv"),
			header:    []string{"id", "name"},
			rows:      [][]string{},
			wantError: false,
		},
		{
			name:      "rewrite with multiple rows",
			path:      filepath.Join(tmp, "rewrite3.csv"),
			header:    []string{"id"},
			rows:      [][]string{{"1"}, {"2"}, {"3"}},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RewriteCSV(tt.path, tt.header, tt.rows)
			if tt.wantError && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tt.wantError {
				rows, err := ReadCSV(tt.path)
				if err != nil {
					t.Fatalf("verify read failed: %v", err)
				}
				if len(rows) != len(tt.rows) {
					t.Fatalf("expected %d rows, got %d", len(tt.rows), len(rows))
				}
			}
		})
	}
}

func TestEnsureCSVExists(t *testing.T) {
	tmp, cleanup := testTempDir(t)
	defer cleanup()

	tests := []struct {
		name       string
		priorExist bool
		wantError  bool
	}{
		{
			name:       "create csv when not exists",
			priorExist: false,
			wantError:  false,
		},
		{
			name:       "skip when csv already exists",
			priorExist: true,
			wantError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmp, tt.name+".csv")

			if tt.priorExist {
				if err := RewriteCSV(path, []string{"id"}, [][]string{}); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			err := EnsureCSVExists(path, []string{"id", "name"})
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
