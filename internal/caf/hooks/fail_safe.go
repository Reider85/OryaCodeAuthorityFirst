package hooks

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// HookContext contains context for hook execution
type HookContext struct {
	RepoRoot string
	SessionID string
	EventType string
}

// ExecuteHook executes a hook with fail-safe guarantees
// INV-ARCH-004: Hook execution is fail-safe: caf binary missing → warn-and-continue, never block session/commit
func ExecuteHook(event string, ctx HookContext) error {
	// Look for caf binary in PATH
	cafPath, err := exec.LookPath("caf")
	if err != nil {
		// caf binary missing — warn and continue, never block
		log.Printf("⚠ caf binary not in $PATH; hook %s skipped", event)
		return nil
	}

	// Build command
	args := []string{fmt.Sprintf("hook %s", event)}
	if ctx.RepoRoot != "" {
		args = append(args, "--repo", ctx.RepoRoot)
	}
	if ctx.SessionID != "" {
		args = append(args, "--session", ctx.SessionID)
	}

	cmd := exec.Command(cafPath, args...)

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("hook %s failed: %w, output: %s", event, err, string(output))
	}

	log.Printf("✓ Hook %s executed successfully: %s", event, string(output))
	return nil
}

// IsHookInstalled checks if the universal pre-commit hook is installed
func IsHookInstalled(repoRoot string) (bool, error) {
	hookPath := fmt.Sprintf("%s/.git/hooks/pre-commit", repoRoot)
	
	// Check if hook exists and is executable
	info, err := os.Stat(hookPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	// Check if it's executable (Unix permissions)
	return info.Mode().Perm()&0111 != 0, nil
}

// InstallUniversalHook installs the universal pre-commit hook
func InstallUniversalHook(repoRoot string) error {
	hookPath := fmt.Sprintf("%s/.git/hooks/pre-commit", repoRoot)
	
	// Create hook content
	hookContent := `#!/bin/sh
# .git/hooks/pre-commit — installed by caf init
# Universal: fires for all AI agents writing to this repo.
# Fail-safe: if caf binary is missing, warn and continue (never block commit).
CAF_PATH=$(command -v caf 2>/dev/null)
if [ -z "$CAF_PATH" ]; then
  echo "⚠ caf binary not in $PATH; pre-commit hook skipped (sessions continue)."
  exit 0
fi
exec "$CAF_PATH" hook pre-commit --repo "$(pwd)" "$@"
`

	// Write hook file
	if err := os.WriteFile(hookPath, []byte(hookContent), 0755); err != nil {
		return fmt.Errorf("failed to write pre-commit hook: %w", err)
	}

	return nil
}

// RemoveUniversalHook removes the universal pre-commit hook
func RemoveUniversalHook(repoRoot string) error {
	hookPath := fmt.Sprintf("%s/.git/hooks/pre-commit", repoRoot)
	return os.Remove(hookPath)
}