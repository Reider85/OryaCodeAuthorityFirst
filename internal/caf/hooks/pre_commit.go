package hooks

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/storage"
)

// PreCommitHook — universal pre-commit event (FR-ADP-003, ADR-003, ADR-010).
//
// Я выбираю advisory/warn-only поведение для MVP: хук НИКОГДА не блокирует
// коммит (INV-ARCH-004, FR-REL-001, SC-8 = 100% никогда не блокирует).
// Блокирующая валидация CAF-Session/decision-trailer — это ALPHA
// (FR-DEC-001, ARCHITECT.md §4.1), поэтому в MVP все находки уходят
// только в лог как предупреждения; сигнатура error оставлена под ALPHA.
//
// Известное ограничение (честно логируется, а не замалчивается): git
// запускает pre-commit ДО записи .git/COMMIT_EDITMSG, поэтому при вызове
// через `git commit` сообщение коммита обычно недоступно — проверка
// трейлеров пропускается, коммит всё равно идёт.
func PreCommitHook(ctx HookContext, commitMessage string) error {
	log.Printf("🔧 pre-commit: repo=%s session=%q", ctx.RepoRoot, ctx.SessionID)

	// 1) Сообщение недоступно на этой стадии (см. комментарий выше).
	if strings.TrimSpace(commitMessage) == "" {
		log.Printf("⚠ pre-commit: commit message unavailable at this stage " +
			"(git writes .git/COMMIT_EDITMSG after pre-commit); trailer check " +
			"skipped — advisory only, never blocks (INV-ARCH-004)")
		return nil
	}

	// 2) Явный bypass: CAF-Trivial: true — трейлеры не проверяем.
	if strings.Contains(commitMessage, "CAF-Trivial: true") {
		log.Printf("✓ pre-commit: CAF-Trivial: true found — advisory trailer check skipped")
		return nil
	}

	// 3) CAF-Session trailer → сверяем, что сессия реально записана в .caf/.
	sessionID := extractCAFSessionTrailer(commitMessage)
	if sessionID == "" {
		log.Printf("⚠ pre-commit: no CAF-Session trailer and no CAF-Trivial: true " +
			"(advisory in MVP; enforcement is ALPHA, FR-DEC-001) — commit NOT blocked")
		return nil
	}

	manager := storage.NewManager(filepath.Join(ctx.RepoRoot, ".caf", "sessions"))
	if _, err := manager.LoadSession(sessionID); err != nil {
		log.Printf("⚠ pre-commit: CAF-Session %q not found in .caf/sessions (%v) — "+
			"advisory only, never blocks", sessionID, err)
		return nil
	}

	log.Printf("✓ pre-commit: CAF-Session %s verified (advisory; MVP never blocks a commit)", sessionID)

	// NEW: Traceability validation (INV-PROMPT-MVP-005)
	if err := PreCommitHookWithTraceability(ctx, commitMessage); err != nil {
		// Advisory only - never blocks commits in MVP
		log.Printf("⚠ pre-commit: traceability validation completed with issues (advisory only)")
	}

	return nil
}

// PreCommitCommand — точка входа `caf hook pre-commit`.
// Я вызываю PreCommitHook напрямую (не через ExecuteHook): ExecuteHook
// перезапускает сам caf, из-за чего PreCommitCommand вызывал себя снова
// и снова — рекурсия, баг которой был отложен на шаг 6 и исправлен здесь.
func PreCommitCommand(repoRoot, sessionID string, commitMessage string) error {
	ctx := HookContext{
		RepoRoot:  repoRoot,
		SessionID: sessionID,
		EventType: "pre-commit",
	}
	return PreCommitHook(ctx, commitMessage)
}

// extractCAFSessionTrailer извлекает значение трейлера CAF-Session.
func extractCAFSessionTrailer(commitMessage string) string {
	for _, line := range strings.Split(commitMessage, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "CAF-Session:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "CAF-Session:"))
		}
	}
	return ""
}

// ValidateCommitMessage — строгая проверка трейлеров коммита.
// В MVP этот результат НЕ используется на путях исполнения (блокировка —
// ALPHA, FR-DEC-001); функция оставлена как подготовка к ALPHA и покрыта
// тестами, чтобы ALPHA просто подключила её к решению о блокировке.
func ValidateCommitMessage(commitMessage string) error {
	if strings.Contains(commitMessage, "CAF-Trivial: true") {
		return nil
	}
	if extractCAFSessionTrailer(commitMessage) == "" {
		return fmt.Errorf("commit must have CAF-Session trailer or CAF-Trivial: true")
	}
	return nil
}

// AddCAFSessionTrailer добавляет CAF-Session trailer (идемпотентно).
func AddCAFSessionTrailer(commitMessage, sessionID string) string {
	if extractCAFSessionTrailer(commitMessage) != "" {
		return commitMessage
	}
	return strings.TrimSpace(commitMessage) + "\n\nCAF-Session: " + sessionID
}

// PreCommitHookWithTraceability extends PreCommitHook with traceability validation.
// Validates that all staged files have proper traceability in their commit history.
// Advisory only in MVP - never blocks commits (INV-PROMPT-MVP-005).
func PreCommitHookWithTraceability(ctx HookContext, commitMessage string) error {
	log.Printf("🔍 Starting traceability validation for staged files")

	// Initialize traceability validator
	validator, err := NewTraceabilityValidator(ctx.RepoRoot)
	if err != nil {
		log.Printf("⚠ Failed to initialize traceability validator: %v (advisory only)", err)
		return nil // Advisory only - never blocks commits
	}

	// Skip validation if disabled in config
	if !validator.config.Enabled {
		log.Printf("ℹ Traceability validation disabled in configuration")
		return nil
	}

	// Validate all staged files
	results := validator.ValidateStagedFiles(ctx.RepoRoot)

	// Generate and report results
	report := validator.GenerateReport(results)
	log.Printf("📋 Traceability validation report:\n%s", report)

	// Advisory logging for issues
	var hasIssues bool
	for _, result := range results {
		if !result.HasTraceability && len(result.TraceabilityIssues) > 0 {
			hasIssues = true
			break
		}
	}

	if hasIssues {
		log.Printf("⚠ Traceability validation completed with issues - advisory only, commit not blocked")
	} else {
		log.Printf("✅ Traceability validation passed")
	}

	return nil // Advisory only - never blocks commits
}

// AddCAFTrivialTrailer добавляет CAF-Trivial: true trailer (идемпотентно).
func AddCAFTrivialTrailer(commitMessage string) string {
	if strings.Contains(commitMessage, "CAF-") {
		return commitMessage
	}
	return strings.TrimSpace(commitMessage) + "\n\nCAF-Trivial: true"
}
