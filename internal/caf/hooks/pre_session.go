package hooks

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/config"
	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/storage"
)

// PreSessionHook executes the pre-session hook
// Reads .caf/context.md and injects into agent's system prompt
func PreSessionHook(ctx HookContext) error {
	log.Printf("🔧 Pre-session hook for session %s", ctx.SessionID)

	// Load config
	configManager := config.NewManager(config.GetConfigPath(ctx.RepoRoot))
	_, err := configManager.LoadConfig()
	if err != nil {
		log.Printf("⚠ Failed to load config: %v", err)
		// Continue anyway
	}

	// Generate context if needed
	contextFile := filepath.Join(ctx.RepoRoot, ".caf", "context.md")
	if _, err := os.Stat(contextFile); os.IsNotExist(err) {
		log.Printf("📝 Generating context for session %s", ctx.SessionID)
		contextGen := storage.NewContextGenerator(
			filepath.Join(ctx.RepoRoot, ".caf", "sessions"),
			filepath.Join(ctx.RepoRoot, ".caf", "decisions"),
			10, // max sessions
			5,  // max decisions
		)
		if err := contextGen.GenerateContext(ctx.RepoRoot); err != nil {
			log.Printf("⚠ Failed to generate context: %v", err)
			// Continue anyway
		}
	}

	// Execute the hook (this would inject context into agent's prompt)
	// For MVP, this is a placeholder - actual injection depends on agent
	log.Printf("✓ Pre-session hook completed for session %s", ctx.SessionID)
	return nil
}

// PreSessionCommand handles the caf hook pre-session command
func PreSessionCommand(repoRoot, sessionID string) error {
	ctx := HookContext{
		RepoRoot: repoRoot,
		SessionID: sessionID,
		EventType: "pre-session",
	}

	return ExecuteHook("pre-session", ctx)
}