package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")
	content := []byte("TEST_KEY_FOO=bar_value\n# comment\nTEST_QUOTED=\"hello world\"\n")
	if err := os.WriteFile(envFile, content, 0644); err != nil {
		t.Fatalf("error creando archivo temporal: %v", err)
	}

	LoadEnv(envFile)

	if val := os.Getenv("TEST_KEY_FOO"); val != "bar_value" {
		t.Fatalf("esperaba bar_value, obtuve %s", val)
	}
	if val := os.Getenv("TEST_QUOTED"); val != "hello world" {
		t.Fatalf("esperaba hello world, obtuve %s", val)
	}
}
