package opencode

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf"
)

// withTempHome подменяет домашний каталог и XDG на временный, чтобы тесты
// никогда не трогали реальный ~/.config/opencode.
func withTempHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("XDG_CONFIG_HOME", "")
	return home
}

// stubBinary кладёт исполняемый стаб бинарника агента в dir — так Detect()
// находит его через $PATH (на Windows LookPath берёт .bat из PATHEXT).
func stubBinary(t *testing.T, dir, name string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		path := filepath.Join(dir, name+".bat")
		if err := os.WriteFile(path, []byte("@echo off\r\n"), 0644); err != nil {
			t.Fatalf("failed to write stub: %v", err)
		}
		return
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatalf("failed to write stub: %v", err)
	}
}

// mustRead читает файл или валит тест.
func mustRead(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	return string(data)
}

func TestIDAndCapabilities(t *testing.T) {
	a := &OpencodeAdapter{}
	if a.ID() != "opencode" {
		t.Errorf("ID() = %q, want opencode", a.ID())
	}
	// ARCHITECT.md §1: opencode — Tier 1, все события native, MCP да (FR-ADP-002).
	want := caf.CapabilityMatrix{
		PreSession: caf.TierNative,
		PostEdit:   caf.TierNative,
		PreCommit:  caf.TierNative,
		MCP:        true,
		Tier:       1,
	}
	if got := a.Capabilities(); got != want {
		t.Errorf("Capabilities() = %+v, want %+v", got, want)
	}
	if err := caf.ValidateCapabilityMatrix(a.Capabilities()); err != nil {
		t.Errorf("ValidateCapabilityMatrix: %v", err)
	}
}

func TestConfigPathUsesHome(t *testing.T) {
	home := withTempHome(t)
	got, err := ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath: %v", err)
	}
	want := filepath.Join(home, ".config", "opencode", "hooks.toml")
	if got != want {
		t.Errorf("ConfigPath() = %q, want %q", got, want)
	}
}

func TestConfigPathHonorsXDG(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("XDG_CONFIG_HOME is a unix convention")
	}
	withTempHome(t)
	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)
	got, err := ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath: %v", err)
	}
	want := filepath.Join(xdg, "opencode", "hooks.toml")
	if got != want {
		t.Errorf("ConfigPath() = %q, want %q", got, want)
	}
}

func TestDetectFindsBinaryInPath(t *testing.T) {
	dir := t.TempDir()
	stubBinary(t, dir, "opencode")
	t.Setenv("PATH", dir)
	if !(&OpencodeAdapter{}).Detect() {
		t.Error("Detect() = false, want true when opencode is in $PATH")
	}
}

func TestDetectFalseWithoutBinary(t *testing.T) {
	t.Setenv("PATH", "")
	if (&OpencodeAdapter{}).Detect() {
		t.Error("Detect() = true, want false when $PATH is empty")
	}
}

func TestInstallCreatesManagedBlock(t *testing.T) {
	withTempHome(t)
	a := &OpencodeAdapter{}
	if err := a.Install(); err != nil {
		t.Fatalf("Install: %v", err)
	}
	path, err := ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath: %v", err)
	}
	body := mustRead(t, path)
	for _, want := range []string{
		cafBlockStart, cafBlockEnd,
		"caf hook pre-session", "caf hook post-edit", "caf hook pre-commit",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("hooks.toml missing %q:\n%s", want, body)
		}
	}
	if strings.Count(body, cafBlockStart) != 1 {
		t.Errorf("expected exactly one managed block, got %d", strings.Count(body, cafBlockStart))
	}
	if !a.Configured() {
		t.Error("Configured() = false after Install, want true")
	}
}

func TestInstallIsIdempotent(t *testing.T) {
	withTempHome(t)
	a := &OpencodeAdapter{}
	if err := a.Install(); err != nil {
		t.Fatalf("first Install: %v", err)
	}
	path, _ := ConfigPath()
	first := mustRead(t, path)

	if err := a.Install(); err != nil {
		t.Fatalf("second Install must be idempotent: %v", err)
	}
	if second := mustRead(t, path); first != second {
		t.Errorf("Install is not idempotent:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
	if strings.Count(mustRead(t, path), cafBlockStart) != 1 {
		t.Error("second Install duplicated the managed block")
	}
}

func TestInstallPreservesForeignHooks(t *testing.T) {
	withTempHome(t)
	path, _ := ConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	foreign := "[[hooks]]\nevent = \"custom\"\ncommand = \"my-tool\"\n"
	if err := os.WriteFile(path, []byte(foreign), 0644); err != nil {
		t.Fatal(err)
	}

	a := &OpencodeAdapter{}
	if err := a.Install(); err != nil {
		t.Fatalf("Install: %v", err)
	}
	body := mustRead(t, path)
	if !strings.Contains(body, "my-tool") {
		t.Errorf("foreign hook was clobbered:\n%s", body)
	}
	if !strings.HasPrefix(body, foreign) {
		t.Errorf("foreign content must stay at the top of the file:\n%s", body)
	}
	// Повторная установка не плодит блоки.
	if err := a.Install(); err != nil {
		t.Fatalf("second Install: %v", err)
	}
	body = mustRead(t, path)
	if strings.Count(body, cafBlockStart) != 1 {
		t.Errorf("managed block duplicated after reinstall:\n%s", body)
	}
	if !strings.Contains(body, "my-tool") {
		t.Errorf("foreign hook lost after reinstall:\n%s", body)
	}
}

func TestInstallHealsTruncatedBlock(t *testing.T) {
	withTempHome(t)
	path, _ := ConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	// Оборванная запись: начало блока есть, конца и тела нет.
	if err := os.WriteFile(path, []byte(cafBlockStart+"\n"), 0644); err != nil {
		t.Fatal(err)
	}

	a := &OpencodeAdapter{}
	if err := a.Install(); err != nil {
		t.Fatalf("Install: %v", err)
	}
	body := mustRead(t, path)
	if strings.Count(body, cafBlockStart) != 1 || strings.Count(body, cafBlockEnd) != 1 {
		t.Errorf("truncated block was not healed:\n%s", body)
	}
	if !strings.Contains(body, "caf hook pre-session") {
		t.Errorf("healed block lost hook commands:\n%s", body)
	}
}

func TestInstallKeepsUnknownTailAfterDamagedMarker(t *testing.T) {
	withTempHome(t)
	path, _ := ConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	// Маркер наш, а содержимое после него — не managed-блок: файл
	// переписывается аккуратно, чужое содержимое сохраняется.
	if err := os.WriteFile(path, []byte(cafBlockStart+"\ncustom = 1\n"), 0644); err != nil {
		t.Fatal(err)
	}

	a := &OpencodeAdapter{}
	if err := a.Install(); err != nil {
		t.Fatalf("Install: %v", err)
	}
	body := mustRead(t, path)
	if !strings.Contains(body, "custom = 1") {
		t.Errorf("unknown content after marker was dropped:\n%s", body)
	}
	if !strings.Contains(body, cafBlockEnd) {
		t.Errorf("managed block is incomplete:\n%s", body)
	}
}

// TestInstallWithoutHomeFails — регрессия бага шага 7: раньше пустой
// Getenv("HOME") приводил к записи в относительный ".config/opencode"
// от текущей директории (на Windows USERPROFILE и HOME пусты).
func TestInstallWithoutHomeFails(t *testing.T) {
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", "")
	t.Setenv("XDG_CONFIG_HOME", "")

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := (&OpencodeAdapter{}).Install(); err == nil {
		t.Fatal("Install must fail when no home directory can be resolved")
	}
	if _, err := os.Stat(filepath.Join(wd, ".config")); !os.IsNotExist(err) {
		t.Errorf("Install must not create .config under the current directory (wd=%s)", wd)
	}
}

func TestUninstallRemovesOwnFileOnly(t *testing.T) {
	withTempHome(t)
	a := &OpencodeAdapter{}
	if err := a.Install(); err != nil {
		t.Fatalf("Install: %v", err)
	}
	if err := a.Uninstall(); err != nil {
		t.Fatalf("Uninstall: %v", err)
	}
	path, _ := ConfigPath()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("Uninstall must remove the caf-owned file at %s", path)
	}
	// Идемпотентность: повторный вызов — не ошибка.
	if err := a.Uninstall(); err != nil {
		t.Errorf("second Uninstall must be nil, got %v", err)
	}
	if a.Configured() {
		t.Error("Configured() = true after Uninstall, want false")
	}
}

func TestUninstallKeepsForeignContent(t *testing.T) {
	withTempHome(t)
	path, _ := ConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	foreign := "[[hooks]]\nevent = \"custom\"\ncommand = \"my-tool\"\n"
	if err := os.WriteFile(path, []byte(foreign), 0644); err != nil {
		t.Fatal(err)
	}

	a := &OpencodeAdapter{}
	if err := a.Install(); err != nil {
		t.Fatalf("Install: %v", err)
	}
	if err := a.Uninstall(); err != nil {
		t.Fatalf("Uninstall: %v", err)
	}

	body := mustRead(t, path)
	if !strings.Contains(body, "my-tool") {
		t.Errorf("foreign hook was removed by Uninstall:\n%s", body)
	}
	if strings.Contains(body, cafBlockStart) || strings.Contains(body, "caf hook") {
		t.Errorf("caf block survived Uninstall:\n%s", body)
	}
	if a.Configured() {
		t.Error("Configured() = true after Uninstall, want false")
	}
}

func TestUninstallRefusesForeignFile(t *testing.T) {
	withTempHome(t)
	path, _ := ConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	foreign := "[[hooks]]\nevent = \"custom\"\ncommand = \"my-tool\"\n"
	if err := os.WriteFile(path, []byte(foreign), 0644); err != nil {
		t.Fatal(err)
	}

	if err := (&OpencodeAdapter{}).Uninstall(); err == nil {
		t.Fatal("Uninstall must refuse a file caf did not write")
	}
	if got := mustRead(t, path); got != foreign {
		t.Errorf("refused file must stay untouched:\ngot:  %q\nwant: %q", got, foreign)
	}
}

func TestUninstallDamagedBlockErrors(t *testing.T) {
	withTempHome(t)
	path, _ := ConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	damaged := cafBlockStart + "\ncustom = 1\n"
	if err := os.WriteFile(path, []byte(damaged), 0644); err != nil {
		t.Fatal(err)
	}

	if err := (&OpencodeAdapter{}).Uninstall(); err == nil {
		t.Fatal("Uninstall must refuse to cut into content it cannot attribute")
	}
	if got := mustRead(t, path); got != damaged {
		t.Errorf("damaged file must stay untouched:\ngot:  %q\nwant: %q", got, damaged)
	}
}

func TestInitializeRegistersEventAdapterOnly(t *testing.T) {
	Initialize()
	if _, ok := caf.GetAdapter("opencode"); !ok {
		t.Fatal("Initialize() must register the opencode event adapter")
	}
	// Языковые адаптеры с выдуманными symbol graph удалены на шаге 7 —
	// в MVP граф строит gopls через universal LSP fallback (ADR-013).
	for _, id := range []string{"go", "python"} {
		if _, ok := caf.GetAdapter(id); ok {
			t.Errorf("language adapter %q must not be registered (removed in step 7)", id)
		}
	}
}

func TestRenderBlockRoundTrip(t *testing.T) {
	// Ядро идемпотентности: mergeBlock собственного блока возвращает его
	// байт-в-байт, иначе каждый `caf init --adapters` плодил бы блоки.
	block := renderBlock()
	if got := mergeBlock(block); got != block {
		t.Errorf("mergeBlock(renderBlock()) != renderBlock()\ngot:\n%s\nwant:\n%s", got, block)
	}
	// И пустой файл → ровно один блок, без лишних переносов.
	if got := mergeBlock(""); got != block {
		t.Errorf("mergeBlock(\\\"\\\") must be exactly the block\ngot:\n%s", got)
	}
	// Чужой файл без завершающего переноса получает разделитель.
	got := mergeBlock("a = 1")
	if !strings.HasPrefix(got, "a = 1\n\n") || !strings.HasSuffix(got, cafBlockEnd+"\n") {
		t.Errorf("mergeBlock must separate foreign content from the block:\n%s", got)
	}
	if strings.Count(got, cafBlockStart) != 1 {
		t.Errorf("expected exactly one block marker:\n%s", got)
	}
}
