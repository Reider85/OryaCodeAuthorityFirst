package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSafeParseYAMLNilOutput(t *testing.T) {
	if err := SafeParseYAML([]byte("a: 1"), nil); err == nil {
		t.Fatal("SafeParseYAML should reject nil output")
	}
}

func TestWriteYAMLMissingDir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "no", "such", "dir", "x.yml")
	if err := WriteYAML(path, map[string]string{"a": "b"}); err == nil {
		t.Fatal("WriteYAML should fail for missing directory")
	}
}

func TestEnsureCAFDir(t *testing.T) {
	tempDir := t.TempDir()
	cafDir, err := EnsureCAFDir(tempDir)
	if err != nil {
		t.Fatalf("EnsureCAFDir: %v", err)
	}
	if cafDir != filepath.Join(tempDir, ".caf") {
		t.Errorf("cafDir = %q", cafDir)
	}
	for _, sub := range []string{"sessions", "decisions", "mental-model", "cache", "metrics", "reviews"} {
		info, err := os.Stat(filepath.Join(cafDir, sub))
		if err != nil || !info.IsDir() {
			t.Errorf("expected subdir %s to exist", sub)
		}
	}
}
