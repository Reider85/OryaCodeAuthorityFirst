package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestManifestCreateLoadRoundTrip(t *testing.T) {
	tempDir := t.TempDir()
	m := NewManifestManager(filepath.Join(tempDir, "manifest.yml"))

	if err := m.CreateManifest("OryaSecurity"); err != nil {
		t.Fatalf("CreateManifest: %v", err)
	}

	loaded, err := m.LoadManifest()
	if err != nil {
		t.Fatalf("LoadManifest: %v", err)
	}
	if loaded.Version != "0.4.0-triz2" {
		t.Errorf("Version = %q, want 0.4.0-triz2", loaded.Version)
	}
	if loaded.Repo != "OryaSecurity" {
		t.Errorf("Repo = %q, want OryaSecurity", loaded.Repo)
	}
	if loaded.Initialized.IsZero() {
		t.Error("Initialized should be set")
	}
	if !loaded.LastSync.IsZero() {
		t.Error("LastSync should be zero before first sync")
	}
}

func TestManifestSaveLoadWithCrossRepo(t *testing.T) {
	tempDir := t.TempDir()
	m := NewManifestManager(filepath.Join(tempDir, "manifest.yml"))

	manifest := Manifest{
		Version:     "0.4.0-triz2",
		Repo:        "OryaSecurity",
		Initialized: time.Date(2026, 10, 24, 11, 0, 0, 0, time.UTC),
		CrossRepo: CrossRepo{
			Enabled: true,
			Repos:   []string{"OryaSecurity", "OtherRepo"},
		},
	}
	if err := m.SaveManifest(manifest); err != nil {
		t.Fatalf("SaveManifest: %v", err)
	}

	loaded, err := m.LoadManifest()
	if err != nil {
		t.Fatalf("LoadManifest: %v", err)
	}
	if !loaded.CrossRepo.Enabled {
		t.Error("CrossRepo.Enabled should be true")
	}
	if len(loaded.CrossRepo.Repos) != 2 {
		t.Errorf("CrossRepo.Repos = %v, want 2 entries", loaded.CrossRepo.Repos)
	}
	if !loaded.Initialized.Equal(manifest.Initialized) {
		t.Errorf("Initialized = %v, want %v", loaded.Initialized, manifest.Initialized)
	}
}

func TestLoadManifestMissingFile(t *testing.T) {
	m := NewManifestManager(filepath.Join(t.TempDir(), "missing.yml"))
	if _, err := m.LoadManifest(); err == nil {
		t.Fatal("LoadManifest should fail for missing file")
	}
}

func TestUpdateLastSync(t *testing.T) {
	tempDir := t.TempDir()
	m := NewManifestManager(filepath.Join(tempDir, "manifest.yml"))
	if err := m.CreateManifest("OryaSecurity"); err != nil {
		t.Fatalf("CreateManifest: %v", err)
	}

	if err := m.UpdateLastSync(); err != nil {
		t.Fatalf("UpdateLastSync: %v", err)
	}

	loaded, err := m.LoadManifest()
	if err != nil {
		t.Fatalf("LoadManifest: %v", err)
	}
	if loaded.LastSync.IsZero() {
		t.Error("LastSync should be set after UpdateLastSync")
	}
	if loaded.Repo != "OryaSecurity" {
		t.Errorf("Repo should be preserved, got %q", loaded.Repo)
	}
}

func TestUpdateLastSyncMissingManifest(t *testing.T) {
	m := NewManifestManager(filepath.Join(t.TempDir(), "missing.yml"))
	if err := m.UpdateLastSync(); err == nil {
		t.Fatal("UpdateLastSync should fail for missing manifest")
	}
}

func TestGetManifestPathAndCAFDir(t *testing.T) {
	root := filepath.Join("some", "root")
	if got := GetManifestPath(root); got != filepath.Join(root, ".caf", "manifest.yml") {
		t.Errorf("GetManifestPath = %q", got)
	}
	if got := GetCAFDir(root); got != filepath.Join(root, ".caf") {
		t.Errorf("GetCAFDir = %q", got)
	}
}

func TestInitializeCAF(t *testing.T) {
	tempDir := t.TempDir()
	if err := InitializeCAF(tempDir); err != nil {
		t.Fatalf("InitializeCAF: %v", err)
	}

	cafDir := GetCAFDir(tempDir)
	for _, sub := range []string{"sessions", "decisions", "mental-model", "cache", "metrics", "reviews"} {
		info, err := os.Stat(filepath.Join(cafDir, sub))
		if err != nil || !info.IsDir() {
			t.Errorf("expected subdir %s to exist", sub)
		}
	}

	m := NewManifestManager(GetManifestPath(tempDir))
	manifest, err := m.LoadManifest()
	if err != nil {
		t.Fatalf("LoadManifest after InitializeCAF: %v", err)
	}
	if manifest.Repo != "unknown" {
		t.Errorf("Repo = %q, want unknown for non-git dir", manifest.Repo)
	}
}

func TestGetGitRepoName(t *testing.T) {
	tempDir := t.TempDir()

	if got := getGitRepoName(tempDir); got != "" {
		t.Errorf("no .git: got %q, want empty", got)
	}

	gitDir := filepath.Join(tempDir, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}
	if got := getGitRepoName(tempDir); got != "" {
		t.Errorf("no config: got %q, want empty", got)
	}

	config := "[remote \"origin\"]\n\turl = https://github.com/Reider85/OryaCodeAuthorityFirst.git\n"
	if err := os.WriteFile(filepath.Join(gitDir, "config"), []byte(config), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	if got := getGitRepoName(tempDir); got != "Reider85/OryaCodeAuthorityFirst" {
		t.Errorf("github url: got %q, want Reider85/OryaCodeAuthorityFirst", got)
	}

	noRemote := "[core]\n\trepositoryformatversion = 0\n"
	if err := os.WriteFile(filepath.Join(gitDir, "config"), []byte(noRemote), 0644); err != nil {
		t.Fatalf("rewrite config: %v", err)
	}
	if got := getGitRepoName(tempDir); got != "" {
		t.Errorf("no remote url: got %q, want empty", got)
	}
}
