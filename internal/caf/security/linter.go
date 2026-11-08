package security

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// SecurityLintRule defines a security linting rule
type SecurityLintRule struct {
	ID          string
	Name        string
	Description string
	Severity    string // low, medium, high, critical
	Category    string // input_validation, crypto, auth, error_handling, etc.
	Pattern     *regexp.Regexp
	CheckFunc   func(*ast.File, *token.FileSet, []SecurityIssue) []SecurityIssue
	Enabled     bool
}

// SecurityIssue represents a security issue found by linting
type SecurityIssue struct {
	RuleID      string
	Message     string
	Severity    string
	Category    string
	File        string
	Line        int
	Column      int
	Context     string
	Suggestion  string
	Confidence  float64
}

// SecurityLinter performs security-focused static analysis
type SecurityLinter struct {
	fset        *token.FileSet
	rules       []SecurityLintRule
	enabledOnly bool
	filePath    string
}

// NewSecurityLinter creates a new security linter
func NewSecurityLinter() *SecurityLinter {
	linter := &SecurityLinter{
		fset:  token.NewFileSet(),
		rules: make([]SecurityLintRule, 0),
	}
	
	// Initialize default rules
	linter.initializeDefaultRules()
	
	return linter
}

// initializeDefaultRules sets up the default security linting rules
func (l *SecurityLinter) initializeDefaultRules() {
	// Input validation rules
	l.AddRule(SecurityLintRule{
		ID:          "input_validation_missing",
		Name:        "Missing Input Validation",
		Description: "Functions that accept external input should validate it",
		Severity:    "high",
		Category:    "input_validation",
		CheckFunc:   l.checkInputValidation,
		Enabled:     true,
	})

	// Cryptographic rules
	l.AddRule(SecurityLintRule{
		ID:          "insecure_crypto",
		Name:        "Insecure Cryptographic Usage",
		Description: "Use of weak or deprecated cryptographic algorithms",
		Severity:    "critical",
		Category:    "crypto",
		Pattern:     regexp.MustCompile(`(?i)(md5|sha1|des|rc4)`),
		CheckFunc:   l.checkInsecureCrypto,
		Enabled:     true,
	})

	// Error handling rules
	l.AddRule(SecurityLintRule{
		ID:          "sensitive_info_exposure",
		Name:        "Sensitive Information Exposure",
		Description: "Potential exposure of sensitive information in error messages",
		Severity:    "medium",
		Category:    "error_handling",
		Pattern:     regexp.MustCompile(`(?i)(password|secret|key|token).*(error|fail)`),
		CheckFunc:   l.checkSensitiveInfoExposure,
		Enabled:     true,
	})

	// Authentication rules
	l.AddRule(SecurityLintRule{
		ID:          "hardcoded_credentials",
		Name:        "Hardcoded Credentials",
		Description: "Hardcoded passwords, API keys, or secrets in source code",
		Severity:    "critical",
		Category:    "auth",
		Pattern:     regexp.MustCompile(`(?i)(password|api_key|secret|token).*=.*["'][^"']{8,}["']`),
		CheckFunc:   l.checkHardcodedCredentials,
		Enabled:     true,
	})

	// SQL injection rules
	l.AddRule(SecurityLintRule{
		ID:          "sql_injection_risk",
		Name:        "SQL Injection Risk",
		Description: "Potential SQL injection through string concatenation",
		Severity:    "high",
		Category:    "injection",
		Pattern:     regexp.MustCompile(`(?i)(fmt\.Sprintf|string\.Concat).*(SELECT|INSERT|UPDATE|DELETE)`),
		CheckFunc:   l.checkSQLInjection,
		Enabled:     true,
	})

	// Path traversal rules
	l.AddRule(SecurityLintRule{
		ID:          "path_traversal_risk",
		Name:        "Path Traversal Risk",
		Description: "Potential path traversal through user input",
		Severity:    "high",
		Category:    "injection",
		Pattern:     regexp.MustCompile(`(?i)(os\.Open|filepath\.Join).*(user|input|request)`),
		CheckFunc:   l.checkPathTraversal,
		Enabled:     true,
	})

	// Insecure random rules
	l.AddRule(SecurityLintRule{
		ID:          "insecure_random",
		Name:        "Insecure Random Number Generation",
		Description: "Use of math/rand for security-sensitive purposes",
		Severity:    "medium",
		Category:    "crypto",
		Pattern:     regexp.MustCompile(`math\.rand`),
		CheckFunc:   l.checkInsecureRandom,
		Enabled:     true,
	})

	// Unchecked error rules
	l.AddRule(SecurityLintRule{
		ID:          "unchecked_error",
		Name:        "Unchecked Error",
		Description: "Function calls that return errors are not checked",
		Severity:    "low",
		Category:    "error_handling",
		Pattern:     regexp.MustCompile(`.*\(\)\s*$`),
		CheckFunc:   l.checkUncheckedError,
		Enabled:     true,
	})
}

// AddRule adds a security linting rule
func (l *SecurityLinter) AddRule(rule SecurityLintRule) {
	l.rules = append(l.rules, rule)
}

// RemoveRule removes a rule by ID
func (l *SecurityLinter) RemoveRule(id string) {
	var newRules []SecurityLintRule
	for _, rule := range l.rules {
		if rule.ID != id {
			newRules = append(newRules, rule)
		}
	}
	l.rules = newRules
}

// SetEnabledOnly sets whether to only run enabled rules
func (l *SecurityLinter) SetEnabledOnly(enabled bool) {
	l.enabledOnly = enabled
}

// SetFilePath sets the file path for linting
func (l *SecurityLinter) SetFilePath(path string) {
	l.filePath = path
}

// LintFile lints a single Go file
func (l *SecurityLinter) LintFile(filePath string) ([]SecurityIssue, error) {
	l.filePath = filePath
	
	// Parse the Go file
	node, err := parser.ParseFile(l.fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	// Run all enabled rules
	var issues []SecurityIssue
	for _, rule := range l.rules {
		if l.enabledOnly && !rule.Enabled {
			continue
		}
		
		if rule.CheckFunc != nil {
			ruleIssues := rule.CheckFunc(node, l.fset, issues)
			issues = append(issues, ruleIssues...)
		}
		
		if rule.Pattern != nil {
			patternIssues := l.checkPattern(node, rule)
			issues = append(issues, patternIssues...)
		}
	}

	return issues, nil
}

// LintDirectory lints all Go files in a directory
func (l *SecurityLinter) LintDirectory(dirPath string) ([]SecurityIssue, error) {
	var allIssues []SecurityIssue
	
	// Walk through directory
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip non-Go files
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			issues, err := l.LintFile(path)
			if err != nil {
				return err
			}
			allIssues = append(allIssues, issues...)
		}
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}
	
	return allIssues, nil
}

// checkPattern checks for regex patterns in the AST
func (l *SecurityLinter) checkPattern(node *ast.File, rule SecurityLintRule) []SecurityIssue {
	var issues []SecurityIssue
	
	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		
		// Check for string literals
		if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			if rule.Pattern.MatchString(lit.Value) {
				pos := l.fset.Position(lit.Pos())
				issue := SecurityIssue{
					RuleID:     rule.ID,
					Message:    fmt.Sprintf("%s: %s", rule.Name, rule.Description),
					Severity:   rule.Severity,
					Category:   rule.Category,
					File:       l.filePath,
					Line:       pos.Line,
					Column:     pos.Column,
					Context:    lit.Value,
					Confidence: 0.8,
				}
				issues = append(issues, issue)
			}
		}
		
		return true
	})
	
	return issues
}

// checkInputValidation checks for missing input validation
func (l *SecurityLinter) checkInputValidation(node *ast.File, fset *token.FileSet, existingIssues []SecurityIssue) []SecurityIssue {
	var issues []SecurityIssue
	
	// Look for function parameters that might need validation
	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if fn.Type.Params != nil {
				for _, param := range fn.Type.Params.List {
					// Check if parameter name suggests external input
					if param.Names != nil {
						for _, name := range param.Names {
							if l.isInputParameter(name.Name) {
								// Check if there's validation code
								hasValidation := l.hasValidation(fn, name.Name)
								if !hasValidation {
									pos := fset.Position(fn.Pos())
									issue := SecurityIssue{
										RuleID:     "input_validation_missing",
										Message:    "Function accepts external input but lacks validation",
										Severity:   "high",
										Category:   "input_validation",
										File:       l.filePath,
										Line:       pos.Line,
										Column:     pos.Column,
										Context:    fmt.Sprintf("Parameter: %s", name.Name),
										Suggestion: "Add input validation for parameter: " + name.Name,
										Confidence: 0.7,
									}
									issues = append(issues, issue)
								}
							}
						}
					}
				}
			}
		}
		return true
	})
	
	return issues
}

// checkInsecureCrypto checks for insecure cryptographic usage
func (l *SecurityLinter) checkInsecureCrypto(node *ast.File, fset *token.FileSet, existingIssues []SecurityIssue) []SecurityIssue {
	var issues []SecurityIssue
	
	insecureAlgorithms := []string{"md5", "sha1", "des", "rc4"}
	
	ast.Inspect(node, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if ident, ok := call.Fun.(*ast.Ident); ok {
				for _, algo := range insecureAlgorithms {
					if strings.Contains(strings.ToLower(ident.Name), algo) {
						pos := fset.Position(call.Pos())
						issue := SecurityIssue{
							RuleID:     "insecure_crypto",
							Message:    fmt.Sprintf("Use of insecure cryptographic algorithm: %s", algo),
							Severity:   "critical",
							Category:   "crypto",
							File:       l.filePath,
							Line:       pos.Line,
							Column:     pos.Column,
							Context:    ident.Name,
							Suggestion: "Use a secure alternative like SHA-256 or bcrypt",
							Confidence: 0.9,
						}
						issues = append(issues, issue)
					}
				}
			}
		}
		return true
	})
	
	return issues
}

// checkSensitiveInfoExposure checks for sensitive information in error messages
func (l *SecurityLinter) checkSensitiveInfoExposure(node *ast.File, fset *token.FileSet, existingIssues []SecurityIssue) []SecurityIssue {
	var issues []SecurityIssue
	
	ast.Inspect(node, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if ident, ok := call.Fun.(*ast.Ident); ok {
				if strings.Contains(strings.ToLower(ident.Name), "error") || strings.Contains(strings.ToLower(ident.Name), "log") {
					// Check if error message contains sensitive information
					for _, arg := range call.Args {
						if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
							if l.containsSensitiveInfo(lit.Value) {
								pos := fset.Position(call.Pos())
								issue := SecurityIssue{
									RuleID:     "sensitive_info_exposure",
									Message:    "Potential sensitive information exposure in error message",
									Severity:   "medium",
									Category:   "error_handling",
									File:       l.filePath,
									Line:       pos.Line,
									Column:     pos.Column,
									Context:    lit.Value,
									Suggestion: "Remove sensitive information from error messages",
									Confidence: 0.6,
								}
								issues = append(issues, issue)
							}
						}
					}
				}
			}
		}
		return true
	})
	
	return issues
}

// checkHardcodedCredentials checks for hardcoded credentials
func (l *SecurityLinter) checkHardcodedCredentials(node *ast.File, fset *token.FileSet, existingIssues []SecurityIssue) []SecurityIssue {
	var issues []SecurityIssue
	
	ast.Inspect(node, func(n ast.Node) bool {
		if assign, ok := n.(*ast.AssignStmt); ok {
			for _, expr := range assign.Rhs {
				if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.STRING {
					if l.looksLikeCredential(lit.Value) {
						pos := fset.Position(assign.Pos())
						issue := SecurityIssue{
							RuleID:     "hardcoded_credentials",
							Message:    "Hardcoded credentials detected",
							Severity:   "critical",
							Category:   "auth",
							File:       l.filePath,
							Line:       pos.Line,
							Column:     pos.Column,
							Context:    lit.Value,
							Suggestion: "Move credentials to environment variables or secure configuration",
							Confidence: 0.9,
						}
						issues = append(issues, issue)
					}
				}
			}
		}
		return true
	})
	
	return issues
}

// checkSQLInjection checks for potential SQL injection
func (l *SecurityLinter) checkSQLInjection(node *ast.File, fset *token.FileSet, existingIssues []SecurityIssue) []SecurityIssue {
	var issues []SecurityIssue
	
	ast.Inspect(node, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if ident, ok := call.Fun.(*ast.Ident); ok {
				// Check for string formatting functions
				if strings.Contains(strings.ToLower(ident.Name), "sprintf") || strings.Contains(strings.ToLower(ident.Name), "concat") {
					// Check if it's building SQL queries
					for _, arg := range call.Args {
						if l.containsSQLKeywords(arg) {
							pos := fset.Position(call.Pos())
							issue := SecurityIssue{
								RuleID:     "sql_injection_risk",
								Message:    "Potential SQL injection through string concatenation",
								Severity:   "high",
								Category:   "injection",
								File:       l.filePath,
								Line:       pos.Line,
								Column:     pos.Column,
								Context:    ident.Name,
								Suggestion: "Use parameterized queries or prepared statements",
								Confidence: 0.7,
							}
							issues = append(issues, issue)
						}
					}
				}
			}
		}
		return true
	})
	
	return issues
}

// checkPathTraversal checks for potential path traversal
func (l *SecurityLinter) checkPathTraversal(node *ast.File, fset *token.FileSet, existingIssues []SecurityIssue) []SecurityIssue {
	var issues []SecurityIssue
	
	ast.Inspect(node, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := selector.X.(*ast.Ident); ok {
					// Check for file operations with user input
					if strings.Contains(strings.ToLower(selector.Sel.Name), "open") || 
					   strings.Contains(strings.ToLower(selector.Sel.Name), "join") {
						if l.containsUserInput(call.Args) {
							pos := fset.Position(call.Pos())
							issue := SecurityIssue{
								RuleID:     "path_traversal_risk",
								Message:    "Potential path traversal through user input",
								Severity:   "high",
								Category:   "injection",
								File:       l.filePath,
								Line:       pos.Line,
								Column:     pos.Column,
								Context:    selector.Sel.Name,
								Suggestion: "Sanitize file paths and validate input",
								Confidence: 0.6,
							}
							issues = append(issues, issue)
						}
					}
				}
			}
		}
		return true
	})
	
	return issues
}

// checkInsecureRandom checks for insecure random number generation
func (l *SecurityLinter) checkInsecureRandom(node *ast.File, fset *token.FileSet, existingIssues []SecurityIssue) []SecurityIssue {
	var issues []SecurityIssue
	
	ast.Inspect(node, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if ident, ok := call.Fun.(*ast.SelectorExpr); ok {
				if ident.Sel != nil && strings.Contains(strings.ToLower(ident.Sel.Name), "rand") {
					// Check if it's using math/rand
					if ident.X != nil {
						if xIdent, ok := ident.X.(*ast.Ident); ok && xIdent.Name == "math" {
							pos := fset.Position(call.Pos())
							issue := SecurityIssue{
								RuleID:     "insecure_random",
								Message:    "Use of math/rand for security-sensitive purposes",
								Severity:   "medium",
								Category:   "crypto",
								File:       l.filePath,
								Line:       pos.Line,
								Column:     pos.Column,
								Context:    ident.Sel.Name,
								Suggestion: "Use crypto/rand for security-sensitive random number generation",
								Confidence: 0.8,
							}
							issues = append(issues, issue)
						}
					}
				}
			}
		}
		return true
	})
	
	return issues
}

// checkUncheckedError checks for unchecked error returns
func (l *SecurityLinter) checkUncheckedError(node *ast.File, fset *token.FileSet, existingIssues []SecurityIssue) []SecurityIssue {
	var issues []SecurityIssue
	
	ast.Inspect(node, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			// Check if function returns an error but it's not handled
			if len(call.Args) > 0 {
				lastArg := call.Args[len(call.Args)-1]
				if ident, ok := lastArg.(*ast.Ident); ok && ident.Name == "err" {
					pos := fset.Position(call.Pos())
					issue := SecurityIssue{
						RuleID:     "unchecked_error",
						Message:    "Function call that returns error is not checked",
						Severity:   "low",
						Category:   "error_handling",
						File:       l.filePath,
						Line:       pos.Line,
						Column:     pos.Column,
						Context:    ident.Name,
						Suggestion: "Check the error return value",
						Confidence: 0.5,
					}
					issues = append(issues, issue)
				}
			}
		}
		return true
	})
	
	return issues
}

// Helper functions for security checks

func (l *SecurityLinter) isInputParameter(name string) bool {
	inputKeywords := []string{"input", "request", "param", "arg", "data", "body", "query", "form"}
	lowerName := strings.ToLower(name)
	for _, keyword := range inputKeywords {
		if strings.Contains(lowerName, keyword) {
			return true
		}
	}
	return false
}

func (l *SecurityLinter) hasValidation(fn *ast.FuncDecl, paramName string) bool {
	// This is a simplified check - in a real implementation, you'd analyze the function body
	// for validation logic
	return false
}

func (l *SecurityLinter) containsSensitiveInfo(value string) bool {
	sensitivePatterns := []string{"password", "secret", "key", "token", "credential"}
	lowerValue := strings.ToLower(value)
	for _, pattern := range sensitivePatterns {
		if strings.Contains(lowerValue, pattern) {
			return true
		}
	}
	return false
}

func (l *SecurityLinter) looksLikeCredential(value string) bool {
	// Check if string looks like a credential (long, random-looking)
	if len(value) < 8 {
		return false
	}
	
	// Check for common credential patterns
	credentialPatterns := []string{
		`[a-zA-Z0-9+/]{16,}`,  // Base64-like
		`[a-zA-Z0-9_-]{20,}`,  // API key-like
		`[a-fA-F0-9]{32,}`,    // Hex-like
	}
	
	for _, pattern := range credentialPatterns {
		matched, _ := regexp.MatchString(pattern, value)
		if matched {
			return true
		}
	}
	
	return false
}

func (l *SecurityLinter) containsSQLKeywords(node ast.Node) bool {
	if lit, ok := node.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		sqlKeywords := []string{"SELECT", "INSERT", "UPDATE", "DELETE", "FROM", "WHERE"}
		upperValue := strings.ToUpper(lit.Value)
		for _, keyword := range sqlKeywords {
			if strings.Contains(upperValue, keyword) {
				return true
			}
		}
	}
	return false
}

func (l *SecurityLinter) containsUserInput(args []ast.Expr) bool {
	// Simplified check - in real implementation, analyze arguments
	return len(args) > 0
}

// FormatIssues formats security issues for display
func (l *SecurityLinter) FormatIssues(issues []SecurityIssue) string {
	if len(issues) == 0 {
		return "No security issues found"
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Found %d security issue(s):\n", len(issues)))
	
	for i, issue := range issues {
		builder.WriteString(fmt.Sprintf("\n%d. %s\n", i+1, issue.Message))
		builder.WriteString(fmt.Sprintf("   Severity: %s\n", issue.Severity))
		builder.WriteString(fmt.Sprintf("   Category: %s\n", issue.Category))
		builder.WriteString(fmt.Sprintf("   File: %s:%d:%d\n", issue.File, issue.Line, issue.Column))
		builder.WriteString(fmt.Sprintf("   Context: %s\n", issue.Context))
		if issue.Suggestion != "" {
			builder.WriteString(fmt.Sprintf("   Suggestion: %s\n", issue.Suggestion))
		}
		builder.WriteString(fmt.Sprintf("   Confidence: %.2f\n", issue.Confidence))
	}
	
	return builder.String()
}

// GetStats returns statistics about the linting results
func (l *SecurityLinter) GetStats(issues []SecurityIssue) map[string]interface{} {
	stats := map[string]interface{}{
		"total_issues": len(issues),
		"by_severity":  make(map[string]int),
		"by_category":  make(map[string]int),
	}
	
	for _, issue := range issues {
		stats["by_severity"].(map[string]int)[issue.Severity]++
		stats["by_category"].(map[string]int)[issue.Category]++
	}
	
	return stats
}