package hooks

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/config"
)

// TraceabilityValidator handles traceability validation for files.
type TraceabilityValidator struct {
	config *config.TraceabilityConfig
}

// TraceabilityResult contains the results of traceability validation.
type TraceabilityResult struct {
	File         string   `yaml:"file"`
	Commits      []Commit `yaml:"commits"`
	HasTraceability bool   `yaml:"has_traceability"`
	TraceabilityIssues []string `yaml:"traceability_issues"`
}

// Commit represents a git commit with traceability information.
type Commit struct {
	Hash        string   `yaml:"hash"`
	Message     string   `yaml:"message"`
	Author      string   `yaml:"author"`
	Date        string   `yaml:"date"`
	HasTraceability bool `yaml:"has_traceability"`
	Requirements []string `yaml:"requirements"`
}

// TraceabilityValidator handles traceability validation for files.
type TraceabilityValidator struct {
	config TraceabilityConfig
}

// NewTraceabilityValidator creates a new traceability validator.
func NewTraceabilityValidator(repoRoot string) (*TraceabilityValidator, error) {
	configPath := config.GetTraceabilityConfigPath(repoRoot)
	traceabilityConfig, err := config.LoadTraceabilityConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load traceability config: %w", err)
	}
	
	if err := traceabilityConfig.Validate(); err != nil {
		return nil, fmt.Errorf("invalid traceability config: %w", err)
	}
	
	return &TraceabilityValidator{
		config: traceabilityConfig,
	}, nil
}

// NewTraceabilityValidatorWithConfig creates a new traceability validator with custom config.
func NewTraceabilityValidatorWithConfig(config *config.TraceabilityConfig) *TraceabilityValidator {
	return &TraceabilityValidator{
		config: config,
	}
}

// ValidateFileTraceability validates traceability for a single file.
// Returns advisory warnings for missing traceability, never blocks.
func (v *TraceabilityValidator) ValidateFileTraceability(repoRoot, filePath string) *TraceabilityResult {
	log.Printf("🔍 Validating traceability for: %s", filePath)

	// Check if file should be validated based on patterns
	if !v.shouldValidateFile(filePath) {
		log.Printf("⏭️  File %s excluded from traceability validation", filePath)
		return &TraceabilityResult{
			File:         filePath,
			HasTraceability: true, // Not applicable for excluded files
		}
	}

	// Get commit history for the file
	commits := v.getFileCommitHistory(repoRoot, filePath)
	if len(commits) == 0 {
		log.Printf("⚠️  No commit history found for: %s", filePath)
		return &TraceabilityResult{
			File:         filePath,
			Commits:      commits,
			HasTraceability: false,
			TraceabilityIssues: []string{"No commit history found"},
		}
	}

	// Check traceability for each commit
	var issues []string
	var hasTraceability bool

	for _, commit := range commits {
		if commit.HasTraceability {
			hasTraceability = true
		} else {
			issues = append(issues, fmt.Sprintf("Commit %s: Missing 'Traces to:' footer", commit.Hash[:7]))
		}
	}

	// Check if file is new and within grace period
	if !hasTraceability && v.isNewFile(repoRoot, filePath, commits) {
		log.Printf("ℹ️  File %s is new, within grace period", filePath)
		return &TraceabilityResult{
			File:         filePath,
			Commits:      commits,
			HasTraceability: false,
			TraceabilityIssues: []string{"New file within grace period"},
		}
	}

	result := &TraceabilityResult{
		File:         filePath,
		Commits:      commits,
		HasTraceability: hasTraceability,
		TraceabilityIssues: issues,
	}

	if !hasTraceability {
		log.Printf("⚠️  Traceability issues found for: %s", filePath)
		for _, issue := range issues {
			log.Printf("   - %s", issue)
		}
	} else {
		log.Printf("✅ Traceability validated for: %s", filePath)
	}

	return result
}

// shouldValidateFile checks if a file should be included in traceability validation.
func (v *TraceabilityValidator) shouldValidateFile(filePath string) bool {
	// Check exclude patterns first
	for _, pattern := range v.config.ExcludePatterns {
		if matchPattern(filePath, pattern) {
			return false
		}
	}

	// Check include patterns
	for _, pattern := range v.config.IncludePatterns {
		if matchPattern(filePath, pattern) {
			return true
		}
	}

	return false
}

// matchPattern checks if a file path matches a glob pattern.
func matchPattern(filePath, pattern string) bool {
	// Simple pattern matching - can be enhanced with full glob support
	// For now, handle basic wildcards and ** for recursive matching
	
	if pattern == "**" {
		return true
	}
	
	if strings.Contains(pattern, "**") {
		// Handle ** as recursive wildcard
		patternParts := strings.Split(pattern, "**")
		if len(patternParts) == 2 {
			prefix := patternParts[0]
			suffix := patternParts[1]
			
			if strings.HasPrefix(filePath, prefix) && strings.HasSuffix(filePath, suffix) {
				return true
			}
		}
	}
	
	// Simple wildcard matching
	if strings.Contains(pattern, "*") {
		patternRegex := strings.ReplaceAll(pattern, ".", "\\.")
		patternRegex = strings.ReplaceAll(patternRegex, "*", ".*")
		matched, _ := regexp.MatchString("^"+patternRegex+"$", filePath)
		return matched
	}
	
	// Exact match
	return filePath == pattern
}

// getFileCommitHistory gets the commit history for a file.
func (v *TraceabilityValidator) getFileCommitHistory(repoRoot, filePath string) []Commit {
	// Get the git log for the file
	cmd := fmt.Sprintf("git log --format=\"%%H|%%s|%%an|%%ad\" --date=short -n %d %s", 
		v.config.LookbackCommits, filePath)
	
	output, err := BashCommand(repoRoot, cmd)
	if err != nil {
		log.Printf("⚠️  Failed to get commit history for %s: %v", filePath, err)
		return []Commit{}
	}

	var commits []Commit
	lines := strings.Split(strings.TrimSpace(output), "\n")
	
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		
		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			continue
		}
		
		commit := Commit{
			Hash:   parts[0],
			Message: parts[1],
			Author: parts[2],
			Date:   parts[3],
		}
		
		// Check for traceability in commit message
		commit.HasTraceability, commit.Requirements = v.extractTraceability(parts[1])
		
		commits = append(commits, commit)
	}

	return commits
}

// extractTraceability extracts requirement information from commit message.
func (v *TraceabilityValidator) extractTraceability(message string) (bool, []string) {
	lines := strings.Split(message, "\n")
	var requirements []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Traces to:") {
			requirements = v.parseRequirements(strings.TrimPrefix(line, "Traces to:"))
			return true, requirements
		}
	}
	
	return false, []string{}
}

// parseRequirements parses requirement IDs from a traceability string.
func (v *TraceabilityValidator) parseRequirements(traceString string) []string {
	var requirements []string
	requirementRegex := regexp.MustCompile(`\b([A-Z]+-\d+)\b`)
	
	requirementIDs := requirementRegex.FindAllString(traceString, -1)
	for _, id := range requirementIDs {
		// Check if requirement has valid prefix
		for _, prefix := range v.config.RequirementPrefixes {
			if strings.HasPrefix(id, prefix) {
				requirements = append(requirements, id)
				break
			}
		}
	}
	
	return requirements
}

// isNewFile checks if a file is new and within the grace period.
func (v *TraceabilityValidator) isNewFile(repoRoot, filePath string, commits []Commit) bool {
	if len(commits) == 0 {
		return true
	}
	
	// Check if file was created recently
	if len(commits) <= v.config.NewFileGraceCommits {
		// Check if this is the first commit for this file
		cmd := fmt.Sprintf("git log --format=%%H --follow -n 1 %s", filePath)
		firstCommit, err := BashCommand(repoRoot, cmd)
		if err != nil {
			return false
		}
		
		return commits[0].Hash == strings.TrimSpace(firstCommit)
	}
	
	return false
}

// ValidateStagedFiles validates traceability for all staged files.
func (v *TraceabilityValidator) ValidateStagedFiles(repoRoot string) []TraceabilityResult {
	// Get staged files
	cmd := "git diff --cached --name-only"
	output, err := BashCommand(repoRoot, cmd)
	if err != nil {
		log.Printf("⚠️  Failed to get staged files: %v", err)
		return []TraceabilityResult{}
	}

	var files []string
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			files = append(files, line)
		}
	}

	var results []TraceabilityResult
	for _, file := range files {
		result := v.ValidateFileTraceability(repoRoot, file)
		results = append(results, *result)
	}

	return results
}

// GenerateReport generates a traceability validation report.
func (v *TraceabilityValidator) GenerateReport(results []TraceabilityResult) string {
	var report strings.Builder
	
	report.WriteString("📋 Traceability Validation Report\n")
	report.WriteString("==============================\n\n")
	
	var totalFiles, passedFiles, failedFiles int
	var totalIssues int
	
	for _, result := range results {
		totalFiles++
		
		if result.HasTraceability {
			passedFiles++
			report.WriteString("✅ ")
		} else {
			failedFiles++
			report.WriteString("❌ ")
		}
		
		report.WriteString(result.File)
		if len(result.TraceabilityIssues) > 0 {
			report.WriteString(fmt.Sprintf(" (%d issues)", len(result.TraceabilityIssues)))
			totalIssues += len(result.TraceabilityIssues)
		}
		report.WriteString("\n")
	}
	
	report.WriteString("\n📊 Summary\n")
	report.WriteString("---------\n")
	report.WriteString(fmt.Sprintf("Total files checked: %d\n", totalFiles))
	report.WriteString(fmt.Sprintf("Passed: %d\n", passedFiles))
	report.WriteString(fmt.Sprintf("Failed: %d\n", failedFiles))
	report.WriteString(fmt.Sprintf("Total issues: %d\n", totalIssues))
	
	if failedFiles > 0 {
		report.WriteString("\n⚠️  Files with traceability issues:\n")
		for _, result := range results {
			if !result.HasTraceability && len(result.TraceabilityIssues) > 0 {
				report.WriteString(fmt.Sprintf("\n%s:\n", result.File))
				for _, issue := range result.TraceabilityIssues {
					report.WriteString(fmt.Sprintf("  - %s\n", issue))
				}
			}
		}
	}
	
	return report.String()
}

// BashCommand executes a bash command in the specified directory.
func BashCommand(workDir, command string) (string, error) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("cd %s && %s", workDir, command))
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}