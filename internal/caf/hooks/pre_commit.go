package hooks

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/config"
	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/storage"
)

// PreCommitHook executes the pre-commit hook
// Validates commit has CAF-Session trailer + decision link (or CAF-Trivial: true)
// This is ALPHA scope in MVP, but we implement the hook infrastructure
func PreCommitHook(ctx HookContext, commitMessage string) error {
	log.Printf("🔧 Pre-commit hook for session %s", ctx.SessionID)

	// Load config
	configManager := config.NewManager(storage.GetConfigPath(ctx.RepoRoot))
	cfg, err := configManager.LoadConfig()
	if err != nil {
		log.Printf("⚠ Failed to load config: %v", err)
		// Continue anyway
	}

	// Check if commit has CAF-Trivial: true
	if strings.Contains(commitMessage, "CAF-Trivial: true") {
		log.Printf("✓ Trivial commit detected, skipping validation")
		return nil
	}

	// Check if commit has CAF-Session trailer
	sessionID := extractCAFSessionTrailer(commitMessage)
	if sessionID == "" {
		return fmt.Errorf("commit missing CAF-Session trailer or CAF-Trivial: true")
	}

	// Load session to validate
	sessionManager := storage.NewManager(filepath.Join(ctx.RepoRoot, ".caf", "sessions"))
	session, err := sessionManager.LoadSession(sessionID)
	if err != nil {
		return fmt.Errorf("failed to load session %s: %w", sessionID, err)
	}

	// Check if session has unrecorded decisions
	if len(session.DecisionsMade) == 0 && len(session.Iterations) > 0 {
		log.Printf("⚠ Session %s has iterations but no decisions recorded", sessionID)
		// Don't block, just warn
	}

	log.Printf("✓ Pre-commit hook validated for session %s", sessionID)
	return nil
}

// PreCommitCommand handles the caf hook pre-commit command
func PreCommitCommand(repoRoot, sessionID string, commitMessage string) error {
	ctx := HookContext{
		RepoRoot: repoRoot,
		SessionID: sessionID,
		EventType: "pre-commit",
	}

	return ExecuteHook("pre-commit", ctx)
}

// extractCAFSessionTrailer extracts the CAF-Session trailer from commit message
func extractCAFSessionTrailer(commitMessage string) string {
	lines := strings.Split(commitMessage, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "CAF-Session: ") {
			return strings.TrimPrefix(line, "CAF-Session: ")
		}
	}
	return ""
}

// ValidateCommitMessage validates a commit message for CAF compliance
func ValidateCommitMessage(commitMessage string) error {
	// Check for CAF-Trivial: true
	if strings.Contains(commitMessage, "CAF-Trivial: true") {
		return nil
	}

	// Check for CAF-Session trailer
	sessionID := extractCAFSessionTrailer(commitMessage)
	if sessionID == "" {
		return fmt.Errorf("commit must have CAF-Session trailer or CAF-Trivial: true")
	}

	return nil
}

// AddCAFSessionTrailer adds a CAF-Session trailer to a commit message
func AddCAFSessionTrailer(commitMessage, sessionID string) string {
	// Check if already has CAF trailer
	if extractCAFSessionTrailer(commitMessage) != "" {
		return commitMessage
	}

	// Add CAF-Session trailer
	return strings.TrimSpace(commitMessage) + "\n\nCAF-Session: " + sessionID
}

// AddCAFTrivialTrailer adds a CAF-Trivial: true trailer to a commit message
func AddCAFTrivialTrailer(commitMessage string) string {
	// Check if already has CAF trailer
	if strings.Contains(commitMessage, "CAF-") {
		return commitMessage
	}

	// Add CAF-Trivial: true trailer
	return strings.TrimSpace(commitMessage) + "\n\nCAF-Trivial: true"
}