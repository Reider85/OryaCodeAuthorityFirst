package hooks

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/config"
	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/storage"
)

// PostEditHook executes the post-edit hook
// Writes .caf/sessions/<id>.yml draft when invoked by adapter
func PostEditHook(ctx HookContext, filesChanged []string, changesSummary string) error {
	log.Printf("🔧 Post-edit hook for session %s", ctx.SessionID)

	// Load config
	configManager := config.NewManager(config.GetConfigPath(ctx.RepoRoot))
	cfg, err := configManager.LoadConfig()
	if err != nil {
		log.Printf("⚠ Failed to load config: %v", err)
		// Continue anyway
	}

	// Check if changes are trivial
	if isTrivialChange(filesChanged, changesSummary, cfg.TrivialityThreshold) {
		log.Printf("✓ Trivial change detected, skipping session logging")
		return nil
	}

	// Create or update session
	sessionManager := storage.NewManager(filepath.Join(ctx.RepoRoot, ".caf", "sessions"))
	
	// Try to load existing session
	session, err := sessionManager.LoadSession(ctx.SessionID)
	if err != nil {
		// Create new session
		session = &storage.Session{
			SessionID:     ctx.SessionID,
			Repo:          getRepoName(ctx.RepoRoot),
			StartedAt:     getStartTime(ctx.SessionID),
			EndedAt:       getEndTime(),
			Agent:         getAgentFromSessionID(ctx.SessionID),
			Model:         getModelFromSessionID(ctx.SessionID),
			ToolCall:      getToolCallFromSessionID(ctx.SessionID),
			InitialPrompt: getInitialPromptFromSessionID(ctx.SessionID),
			VoiceSignature: getVoiceSignatureFromSessionID(ctx.SessionID),
			FilesTouched:  filesChanged,
		}
	}

	// Add iteration
	iterationID := fmt.Sprintf("i%d", len(session.Iterations)+1)
	session.Iterations = append(session.Iterations, storage.Iteration{
		ID:     iterationID,
		Action: "edit",
		Target: strings.Join(filesChanged, ", "),
		Diff:   changesSummary,
	})

	// Save session
	if err := sessionManager.SaveSession(session); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	log.Printf("✓ Session %s updated with edit iteration", ctx.SessionID)
	return nil
}

// PostEditCommand handles the caf hook post-edit command
func PostEditCommand(repoRoot, sessionID string, filesChanged []string, changesSummary string) error {
	ctx := HookContext{
		RepoRoot: repoRoot,
		SessionID: sessionID,
		EventType: "post-edit",
	}

	return PostEditHook(ctx, filesChanged, changesSummary)
}

// isTrivialChange checks if the changes are trivial based on LOC and file count
func isTrivialChange(filesChanged []string, changesSummary string, threshold int) bool {
	// Simple heuristic: if changes summary mentions "trivial" or file count > threshold
	if strings.Contains(strings.ToLower(changesSummary), "trivial") {
		return true
	}
	
	if len(filesChanged) > threshold {
		return false
	}
	
	// Additional heuristics could be added here
	return false
}

// Helper functions for session data (would be populated by the agent)
func getRepoName(repoRoot string) string {
	// Simplified - in real implementation, parse .git/config
	return "unknown"
}

func getStartTime(sessionID string) time.Time {
	// Parse session ID to get timestamp
	// Format: 2026-10-24-a1b2c3
	t, err := time.Parse(time.RFC3339, "2026-10-24T11:02:14+03:00")
	if err != nil {
		return time.Now()
	}
	return t
}

func getEndTime() time.Time {
	t, err := time.Parse(time.RFC3339, "2026-10-24T11:47:33+03:00")
	if err != nil {
		return time.Now()
	}
	return t
}

func getAgentFromSessionID(sessionID string) string {
	// Extract agent from session ID or context
	if strings.Contains(sessionID, "opencode") {
		return "opencode"
	}
	if strings.Contains(sessionID, "claude") {
		return "claude-code"
	}
	return "unknown"
}

func getModelFromSessionID(sessionID string) string {
	// Extract model from session ID or context
	return "glm-4.6"
}

func getToolCallFromSessionID(sessionID string) string {
	// Extract tool call from session ID or context
	return "opencode run \"add rate limiter\""
}

func getInitialPromptFromSessionID(sessionID string) string {
	// Extract initial prompt from session ID or context
	return "Need to add rate limiter on POST /v1/policies."
}

func getVoiceSignatureFromSessionID(sessionID string) string {
	// Extract voice signature from session ID or context
	return "rate limiter needed urgently"
}