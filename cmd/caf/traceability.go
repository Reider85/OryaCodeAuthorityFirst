package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/hooks"
	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/config"
)

// traceabilityCmd represents the traceability command
var traceabilityCmd = &cobra.Command{
	Use:   "traceability",
	Short: "Validate traceability for files in the repository",
	Long: `Validate traceability for files in the repository.

This command checks that all files have proper 'Traces to:' links in their
commit history according to the configured requirements.

Examples:
  caf traceability check              # Check traceability for staged files
  caf traceability validate            # Validate all tracked files
  caf traceability report --format json # Generate JSON report`,
	Run: func(cmd *cobra.Command, args []string) {
		repoRoot, _ := cmd.Flags().GetString("repo")
		if repoRoot == "" {
			// Use current working directory
			repoRoot = "."
		}

		// Get subcommand
		subcommand := "check"
		if len(args) > 0 {
			subcommand = args[0]
		}

		switch subcommand {
		case "check":
			checkTraceability(repoRoot)
		case "validate":
			validateTraceability(repoRoot)
		case "report":
			generateTraceabilityReport(repoRoot, cmd.Flags())
		default:
			fmt.Printf("Unknown subcommand: %s\n", subcommand)
			fmt.Println("Available subcommands: check, validate, report")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(traceabilityCmd)
	
	// Add flags
	traceabilityCmd.Flags().StringP("repo", "r", "", "Repository path (default: current directory)")
	traceabilityCmd.Flags().String("format", "cli", "Output format: cli, json, markdown")
	traceabilityCmd.Flags().Bool("verbose", false, "Enable verbose logging")
}

func checkTraceability(repoRoot string) {
	fmt.Println("🔍 Checking traceability for staged files...")

	validator, err := hooks.NewTraceabilityValidator(repoRoot)
	if err != nil {
		fmt.Printf("❌ Failed to initialize traceability validator: %v\n", err)
		os.Exit(1)
	}

	results := validator.ValidateStagedFiles(repoRoot)
	report := validator.GenerateReport(results)
	fmt.Println(report)
}

func validateTraceability(repoRoot string) {
	fmt.Println("🔍 Validating traceability for all tracked files...")

	validator, err := hooks.NewTraceabilityValidator(repoRoot)
	if err != nil {
		fmt.Printf("❌ Failed to initialize traceability validator: %v\n", err)
		os.Exit(1)
	}

	// Get all tracked files
	cmd := "git ls-files"
	output, err := hooks.BashCommand(repoRoot, cmd)
	if err != nil {
		fmt.Printf("❌ Failed to get tracked files: %v\n", err)
		os.Exit(1)
	}

	var files []string
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			files = append(files, line)
		}
	}

	var results []hooks.TraceabilityResult
	for _, file := range files {
		result := validator.ValidateFileTraceability(repoRoot, file)
		results = append(results, *result)
	}

	report := validator.GenerateReport(results)
	fmt.Println(report)
}

func generateTraceabilityReport(repoRoot string, cmd *cobra.Command) {
	format, _ := cmd.Flags().GetString("format")
	verbose, _ := cmd.Flags().GetBool("verbose")

	fmt.Printf("📋 Generating traceability report (%s format)...\n", format)

	validator, err := hooks.NewTraceabilityValidator(repoRoot)
	if err != nil {
		fmt.Printf("❌ Failed to initialize traceability validator: %v\n", err)
		os.Exit(1)
	}

	// Get all tracked files
	cmdStr := "git ls-files"
	output, err := hooks.BashCommand(repoRoot, cmdStr)
	if err != nil {
		fmt.Printf("❌ Failed to get tracked files: %v\n", err)
		os.Exit(1)
	}

	var files []string
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			files = append(files, line)
		}
	}

	var results []hooks.TraceabilityResult
	for _, file := range files {
		result := validator.ValidateFileTraceability(repoRoot, file)
		results = append(results, *result)
	}

	switch format {
	case "json":
		generateJSONReport(results)
	case "markdown":
		generateMarkdownReport(results)
	case "cli":
		report := validator.GenerateReport(results)
		fmt.Println(report)
	default:
		fmt.Printf("❌ Unknown format: %s\n", format)
		os.Exit(1)
	}
}

func generateJSONReport(results []hooks.TraceabilityResult) {
	type Report struct {
		Summary struct {
			TotalFiles int `json:"total_files"`
			Passed     int `json:"passed"`
			Failed     int `json:"failed"`
			TotalIssues int `json:"total_issues"`
		} `json:"summary"`
		Results []hooks.TraceabilityResult `json:"results"`
	}

	report := Report{}
	
	for _, result := range results {
		report.Summary.TotalFiles++
		if result.HasTraceability {
			report.Summary.Passed++
		} else {
			report.Summary.Failed++
		}
		report.Summary.TotalIssues += len(result.TraceabilityIssues)
		report.Results = append(report.Results, result)
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Printf("❌ Failed to generate JSON report: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
}

func generateMarkdownReport(results []hooks.TraceabilityResult) {
	fmt.Println("# Traceability Validation Report")
	fmt.Println()
	fmt.Println("## Summary")
	
	totalFiles := len(results)
	passed := 0
	failed := 0
	totalIssues := 0
	
	for _, result := range results {
		if result.HasTraceability {
			passed++
		} else {
			failed++
		}
		totalIssues += len(result.TraceabilityIssues)
	}
	
	fmt.Printf("- **Total Files**: %d\n", totalFiles)
	fmt.Printf("- **Passed**: %d\n", passed)
	fmt.Printf("- **Failed**: %d\n", failed)
	fmt.Printf("- **Total Issues**: %d\n", totalIssues)
	fmt.Println()
	
	if failed > 0 {
		fmt.Println("## Files with Traceability Issues")
		fmt.Println()
		
		for _, result := range results {
			if !result.HasTraceability && len(result.TraceabilityIssues) > 0 {
				fmt.Printf("### %s\n", result.File)
				fmt.Println()
				
				for _, issue := range result.TraceabilityIssues {
					fmt.Printf("- %s\n", issue)
				}
				fmt.Println()
			}
		}
	}
}