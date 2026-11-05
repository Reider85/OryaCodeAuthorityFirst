package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ContextGenerator generates context.md from recent sessions
type ContextGenerator struct {
	sessionsDir string
	maxSessions int
}

// NewContextGenerator creates a new context generator
func NewContextGenerator(sessionsDir string, maxSessions int) *ContextGenerator {
	return &ContextGenerator{
		sessionsDir: sessionsDir,
		maxSessions: maxSessions,
	}
}

// GenerateContext generates a context.md file from recent sessions
func (g *ContextGenerator) GenerateContext(repoRoot string) error {
	// Get recent sessions
	sessionIDs, err := g.listRecentSessions()
	if err != nil {
		return fmt.Errorf("failed to list recent sessions: %w", err)
	}

	// Build context
	var context strings.Builder
	context.WriteString("# CodeAuthorityFirst Context\n\n")
	context.WriteString(fmt.Sprintf("Generated at: %s\n\n", time.Now().Format(time.RFC3339)))
	context.WriteString("## Recent Sessions\n\n")

	// Add recent sessions
	for i, sessionID := range sessionIDs {
		if i >= g.maxSessions {
			break
		}

		session, err := LoadSession(filepath.Join(g.sessionsDir, sessionID))
		if err != nil {
			continue // Skip invalid sessions
		}

		context.WriteString(fmt.Sprintf("### Session %s (%s)\n", session.SessionID, session.Agent))
		context.WriteString(fmt.Sprintf("- **Started**: %s\n", session.StartedAt.Format(time.RFC3339)))
		context.WriteString(fmt.Sprintf("- **Model**: %s\n", session.Model))
		context.WriteString(fmt.Sprintf("- **Initial prompt**: %s\n", session.InitialPrompt))
		
		if len(session.Iterations) > 0 {
			context.WriteString("- **Iterations**:\n")
			for _, iteration := range session.Iterations {
				context.WriteString(fmt.Sprintf("  - %s: %s\n", iteration.Action, iteration.Target))
				if iteration.Diff != "" {
					context.WriteString(fmt.Sprintf("    *%s*\n", iteration.Diff))
				}
			}
		}

		if len(session.FilesTouched) > 0 {
			context.WriteString("- **Files touched**:\n")
			for _, file := range session.FilesTouched {
				context.WriteString(fmt.Sprintf("  - %s\n", file))
			}
		}



		context.WriteString("\n")
	}

	// Add recent decisions
	context.WriteString("## Recent Decisions\n\n")


	// Write context file
	contextFile := filepath.Join(repoRoot, ".caf", "context.md")
	if err := os.WriteFile(contextFile, []byte(context.String()), 0644); err != nil {
		return fmt.Errorf("failed to write context file: %w", err)
	}

	return nil
}

// listRecentSessions returns recent session IDs sorted by start time
func (g *ContextGenerator) listRecentSessions() ([]string, error) {
	// Параллельный bulk-load (§6 file ops): последовательный цикл LoadSession
	// на 10k сессиях давал ~3s (шаг 15, enforcement TestBudgetFileOps10k).
	sessions, err := NewManager(g.sessionsDir).LoadAllSessions()
	if err != nil {
		return nil, err
	}

	// sort.Slice вместо O(n²) вставок: 10k сессий = 50M сравнений (§22.2).
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].StartedAt.After(sessions[j].StartedAt)
	})

	sessionIDsSorted := make([]string, 0, len(sessions))
	for _, session := range sessions {
		sessionIDsSorted = append(sessionIDsSorted, session.SessionID)
	}

	return sessionIDsSorted, nil
}

