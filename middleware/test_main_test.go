package middleware

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "middleware-test-*")
	if err != nil {
		panic("TestMain: failed to create temp dir: " + err.Error())
	}
	defer os.RemoveAll(tmp)

	dataDir := filepath.Join(tmp, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		panic("TestMain: failed to create data dir: " + err.Error())
	}
	if err := os.Setenv("DATA_DIR", dataDir); err != nil {
		panic("TestMain: failed to set DATA_DIR: " + err.Error())
	}

	code := m.Run()
	os.RemoveAll(tmp)
	os.Exit(code)
}
