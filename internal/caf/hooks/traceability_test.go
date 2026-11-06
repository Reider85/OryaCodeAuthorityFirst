package hooks

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// captureLog перенаправляет стандартный logger в буфер на время fn и
// возвращает записанный вывод — логи и есть контракт advisory-хука.
func captureLog(t *testing.T, fn func()) string {
	t.Helper()
	var buf bytes.Buffer
	prev := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(prev)
	fn()
	return buf.String()
}

func TestTraceabilityValidatorCreation(t *testing.T) {
	validator := NewTraceabilityValidator()
	
	if validator == nil {
		t.Fatal("NewTraceabilityValidator() returned nil")
	}
	
	if len(validator.config.IncludePatterns) == 0 {
		t.Error("Default include patterns should not be empty")
	}
	
	if len(validator.config.RequirementPrefixes) == 0 {
		t.Error("Default requirement prefixes should not be empty")
	}
	
	if validator.config.LookbackCommits <= 0 {
		t.Error("Lookback commits should be positive")
	}
}

func TestShouldValidateFile(t *testing.T) {
	validator := NewTraceabilityValidator()
	
	tests := []struct {
		name     string
		filePath string
		expected bool
	}{
		{"Go file", "main.go", true},
		{"Markdown file", "README.md", true},
		{"YAML file", "config.yaml", true},
		{"Test file (excluded)", "test_test.go", false},
		{"Vendor file (excluded)", "vendor/main.go", false},
		{"Node modules (excluded)", "node_modules/test.js", false},
		{"Unknown extension", "test.txt", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.shouldValidateFile(tt.filePath)
			if result != tt.expected {
				t.Errorf("shouldValidateFile(%q) = %v, want %v", tt.filePath, result, tt.expected)
			}
		})
	}
}

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		filePath  string
		pattern   string
		expected  bool
	}{
		{"main.go", "*.go", true},
		{"main.go", "main.go", true},
		{"src/main.go", "**/*.go", true},
		{"src/main.go", "src/*.go", true},
		{"src/main.go", "*.go", false},
		{"test.txt", "*.md", false},
		{"README.md", "**/*.md", true},
		{"docs/README.md", "**/*.md", true},
		{"vendor/main.go", "**/*.go", false}, // Should be excluded by config, not pattern
	}
	
	for _, tt := range tests {
		t.Run(tt.pattern+" -> "+tt.filePath, func(t *testing.T) {
			result := matchPattern(tt.filePath, tt.pattern)
			if result != tt.expected {
				t.Errorf("matchPattern(%q, %q) = %v, want %v", tt.filePath, tt.pattern, result, tt.expected)
			}
		})
	}
}

func TestExtractTraceability(t *testing.T) {
	validator := NewTraceabilityValidator()
	
	tests := []struct {
		name        string
		message     string
		expectTrace bool
		requirements []string
	}{
		{"No traceability", "feat: add feature", false, []string{}},
		{"With traceability", "feat: add feature\n\nTraces to: FR-001", true, []string{"FR-001"}},
		{"Multiple requirements", "feat: add feature\n\nTraces to: FR-001, ADR-002, INV-003", true, []string{"FR-001", "ADR-002", "INV-003"}},
		{"Invalid requirement", "feat: add feature\n\nTraces to: INVALID-001", false, []string{}},
		{"Mixed valid/invalid", "feat: add feature\n\nTraces to: FR-001, INVALID-002", true, []string{"FR-001"}},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasTrace, requirements := validator.extractTraceability(tt.message)
			
			if hasTrace != tt.expectTrace {
				t.Errorf("extractTraceability(%q) hasTrace = %v, want %v", tt.message, hasTrace, tt.expectTrace)
			}
			
			if len(requirements) != len(tt.requirements) {
				t.Errorf("extractTraceability(%q) requirements = %v, want %v", tt.message, requirements, tt.requirements)
			} else {
				for i, req := range requirements {
					if req != tt.requirements[i] {
						t.Errorf("extractTraceability(%q) requirements[%d] = %v, want %v", tt.message, i, req, tt.requirements[i])
					}
				}
			}
		})
	}
}

func TestParseRequirements(t *testing.T) {
	validator := NewTraceabilityValidator()
	
	tests := []struct {
		traceString string
		expected    []string
	}{
		{"FR-001", []string{"FR-001"}},
		{"FR-001, ADR-002, INV-003", []string{"FR-001", "ADR-002", "INV-003"}},
		{"FR-001 and some text ADR-002 INV-003", []string{"FR-001", "ADR-002", "INV-003"}},
		{"INVALID-001, FR-002", []string{"FR-002"}}, // Invalid prefix filtered out
		{"No requirements here", []string{}},
		{"", []string{}},
	}
	
	for _, tt := range tests {
		t.Run(tt.traceString, func(t *testing.T) {
			result := validator.parseRequirements(tt.traceString)
			if len(result) != len(tt.expected) {
				t.Errorf("parseRequirements(%q) = %v, want %v", tt.traceString, result, tt.expected)
			} else {
				for i, req := range result {
					if req != tt.expected[i] {
						t.Errorf("parseRequirements(%q)[%d] = %v, want %v", tt.traceString, i, req, tt.expected[i])
					}
				}
			}
		})
	}
}

func TestValidateFileTraceability(t *testing.T) {
	// Test with a temporary repository
	repo := t.TempDir()
	
	// Initialize git repo
	bashCommand(repo, "git init")
	bashCommand(repo, "git config user.email 'test@example.com'")
	bashCommand(repo, "git config user.name 'Test User'")
	
	// Create test file
	testFile := filepath.Join(repo, "main.go")
	err := write("package main\n\nfunc main() {\n}", testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	validator := NewTraceabilityValidator()
	
	t.Run("No commit history", func(t *testing.T) {
		result := validator.ValidateFileTraceability(repo, testFile)
		if result == nil {
			t.Fatal("ValidateFileTraceability returned nil")
		}
		
		if len(result.TraceabilityIssues) == 0 {
			t.Error("Expected traceability issues for file with no history")
		}
	})
	
	t.Run("With valid traceability", func(t *testing.T) {
		// Add commit with traceability
		bashCommand(repo, "git add .")
		bashCommand(repo, "git commit -m 'feat: implement main function\n\nTraces to: FR-001, ADR-002'")
		
		result := validator.ValidateFileTraceability(repo, testFile)
		if result == nil {
			t.Fatal("ValidateFileTraceability returned nil")
		}
		
		if !result.HasTraceability {
			t.Error("Expected traceability to be true for file with valid commit")
		}
		
		if len(result.TraceabilityIssues) > 0 {
			t.Errorf("Expected no traceability issues, got: %v", result.TraceabilityIssues)
		}
	})
	
	t.Run("With invalid traceability", func(t *testing.T) {
		// Modify file without traceability
		write("package main\n\nfunc main() {\n\tfmt.Println(\"test\")\n}", testFile)
		bashCommand(repo, "git add .")
		bashCommand(repo, "git commit -m 'feat: add print statement'")
		
		result := validator.ValidateFileTraceability(repo, testFile)
		if result == nil {
			t.Fatal("ValidateFileTraceability returned nil")
		}
		
		if result.HasTraceability {
			t.Error("Expected traceability to be false for commit without traceability")
		}
		
		if len(result.TraceabilityIssues) == 0 {
			t.Error("Expected traceability issues for commit without traceability")
		}
	})
}

func TestIsNewFile(t *testing.T) {
	repo := t.TempDir()
	
	// Initialize git repo
	bashCommand(repo, "git init")
	bashCommand(repo, "git config user.email 'test@example.com'")
	bashCommand(repo, "git config user.name 'Test User'")
	
	validator := NewTraceabilityValidator()
	
	t.Run("New file within grace period", func(t *testing.T) {
		// Create test file
		testFile := filepath.Join(repo, "new_file.go")
		write("package main", testFile)
		
		// Add commit (this should be the first commit for the file)
		bashCommand(repo, "git add .")
		bashCommand(repo, "git commit -m 'feat: add new file'")
		
		// Get commit history
		commits := validator.getFileCommitHistory(repo, testFile)
		
		isNew := validator.isNewFile(repo, testFile, commits)
		if !isNew {
			t.Error("Expected file to be considered new within grace period")
		}
	})
	
	t.Run("Old file beyond grace period", func(t *testing.T) {
		// Create test file
		testFile := filepath.Join(repo, "old_file.go")
		write("package main", testFile)
		
		// Add multiple commits to exceed grace period
		bashCommand(repo, "git add .")
		bashCommand(repo, "git commit -m 'feat: add old file'")
		bashCommand(repo, "git commit -m 'feat: modify old file'")
		bashCommand(repo, "git commit -m 'feat: modify old file again'")
		bashCommand(repo, "git commit -m 'feat: modify old file again'")
		bashCommand(repo, "git commit -m 'feat: modify old file again'")
		bashCommand(repo, "git commit -m 'feat: modify old file again'")
		bashCommand(repo, "git commit -m 'feat: modify old file again'")
		
		// Get commit history
		commits := validator.getFileCommitHistory(repo, testFile)
		
		isNew := validator.isNewFile(repo, testFile, commits)
		if isNew {
			t.Error("Expected file to be considered old beyond grace period")
		}
	})
}

func TestGenerateReport(t *testing.T) {
	validator := NewTraceabilityValidator()
	
	results := []TraceabilityResult{
		{
			File: "main.go",
			HasTraceability: true,
		},
		{
			File: "test.go",
			HasTraceability: false,
			TraceabilityIssues: []string{"Missing traceability in commit abc123"},
		},
	}
	
	report := validator.GenerateReport(results)
	
	if !strings.Contains(report, "Traceability Validation Report") {
		t.Error("Generated report should contain title")
	}
	
	if !strings.Contains(report, "Total files checked: 2") {
		t.Error("Generated report should contain total files count")
	}
	
	if !strings.Contains(report, "Passed: 1") {
		t.Error("Generated report should contain passed count")
	}
	
	if !strings.Contains(report, "Failed: 1") {
		t.Error("Generated report should contain failed count")
	}
	
	if !strings.Contains(report, "test.go") {
		t.Error("Generated report should contain failed file name")
	}
	
	if !strings.Contains(report, "Missing traceability in commit abc123") {
		t.Error("Generated report should contain traceability issue details")
	}
}

func TestPreCommitHookWithTraceability(t *testing.T) {
	repo := t.TempDir()
	
	// Initialize git repo
	bashCommand(repo, "git init")
	bashCommand(repo, "git config user.email 'test@example.com'")
	bashCommand(repo, "git config user.name 'Test User'")
	
	// Create test file
	testFile := filepath.Join(repo, "main.go")
	write("package main", testFile)
	
	// Stage the file
	bashCommand(repo, "git add .")
	
	ctx := HookContext{
		RepoRoot:  repo,
		SessionID: "test-session",
		EventType: "pre-commit",
	}
	
	t.Run("Traceability validation runs without error", func(t *testing.T) {
		err := PreCommitHookWithTraceability(ctx, "")
		if err != nil {
			t.Errorf("PreCommitHookWithTraceability should not return error, got: %v", err)
		}
	})
	
	t.Run("Traceability validation logs properly", func(t *testing.T) {
		out := captureLog(t, func() {
			_ = PreCommitHookWithTraceability(ctx, "")
		})
		
		if !strings.Contains(out, "Starting traceability validation") {
			t.Error("Expected traceability validation to start")
		}
		
		if !strings.Contains(out, "Traceability validation report") {
			t.Error("Expected traceability validation report")
		}
	})
}

// Test helper to write a file
func write(content, filePath string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}

// Test helper to execute bash command
func bashCommand(workDir, command string) (string, error) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("cd %s && %s", workDir, command))
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}