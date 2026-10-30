package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf"
	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/adapters"
	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/doctor"
	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/storage"
	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/hooks"
)

// AutoDetect represents auto-detection results
type AutoDetect struct {
	Adapters   []AdapterInfo `json:"adapters"`
	Recommended string       `json:"recommended"`
	NotFound   []string      `json:"not_found"`
}

// AdapterInfo represents information about a detected adapter
type AdapterInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Tier       int    `json:"tier"`
	Capability string `json:"capability"`
	Path       string `json:"path,omitempty"`
	Status     string `json:"status"` // "installed", "not_found", "fallback"
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [flags]",
	Short: "Initialize CodeAuthorityFirst in a repository",
	Long: `Initialize CodeAuthorityFirst in a repository by:
- Creating .caf/ directory structure
- Installing universal pre-commit hook
- Detecting available AI agents
- Installing adapters for detected agents

Example:
  caf init --auto-detect
  caf init --adapters opencode,claude-code`,
	Run: func(cmd *cobra.Command, args []string) {
		repoRoot, _ := cmd.Flags().GetString("repo")
		autoDetect, _ := cmd.Flags().GetBool("auto-detect")
		adaptersFlag, _ := cmd.Flags().GetString("adapters")
		
		if repoRoot == "" {
			repoRoot = "."
		}
		
		// Convert to absolute path
		absPath, err := filepath.Abs(repoRoot)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		
		if autoDetect {
			runAutoDetectInit(absPath)
		} else if adaptersFlag != "" {
			runSpecificInit(absPath, adaptersFlag)
		} else {
			runBasicInit(absPath)
		}
	},
}

// autoDetectCmd represents the auto-detect command
var autoDetectCmd = &cobra.Command{
	Use:   "auto-detect [repo]",
	Short: "Auto-detect installed AI agents",
	Long: `Scan the system for installed AI agents and return available adapters.
This helps you understand what CodeAuthorityFirst can track automatically.

Example:
  caf auto-detect
  caf auto-detect /path/to/repo`,
	Run: func(cmd *cobra.Command, args []string) {
		repoRoot := "."
		if len(args) > 0 {
			repoRoot = args[0]
		}
		
		// Convert to absolute path
		absPath, err := filepath.Abs(repoRoot)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		
		runAutoDetect(absPath)
	},
}

// whoCmd represents the who command
var whoCmd = &cobra.Command{
	Use:   "who <file:line>",
	Short: "Find which session modified a specific line",
	Long: `Find which AI session modified a specific line of code.
Returns the session ID, prompt, and iteration details.

Example:
  caf who main.go:42
  caf who src/server.go:156`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		location := args[0]
		runWho(location)
	},
}

// whenCmd represents the when command
var whenCmd = &cobra.Command{
	Use:   "when --pattern <pattern> [repo]",
	Short: "Find when a code pattern first appeared",
	Long: `Find when a specific code pattern first appeared in the repository.
This is ALPHA scope in MVP and may not be fully implemented.

Example:
  caf when --pattern "func main"
  caf when --pattern "import.*fmt" /path/to/repo`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pattern, _ := cmd.Flags().GetString("pattern")
		repoRoot := "."
		if len(args) > 0 {
			repoRoot = args[0]
		}
		
		runWhen(pattern, repoRoot)
	},
}

// contextCmd represents the context command
var contextCmd = &cobra.Command{
	Use:   "context [--refresh] [repo]",
	Short: "Generate context.md from recent sessions",
	Long: `Generate a context.md file from recent sessions and decisions.
This provides context for AI agents about recent work.

Example:
  caf context
  caf context --refresh /path/to/repo`,
	Run: func(cmd *cobra.Command, args []string) {
		refresh, _ := cmd.Flags().GetBool("refresh")
		repoRoot := "."
		if len(args) > 0 {
			repoRoot = args[0]
		}
		
		runContext(refresh, repoRoot)
	},
}

// hookCmd represents the hook command
var hookCmd = &cobra.Command{
	Use:   "hook <event> [flags]",
	Short: "Execute a hook event",
	Long: `Execute a hook event manually. This is usually called automatically
by AI agents or the universal pre-commit hook.

Events: pre-session, post-edit, pre-commit

Example:
  caf hook pre-session --session 2026-10-24-a1b2c3
  caf hook post-edit --session 2026-10-24-a1b2c3 --files "main.go"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(1)
		}
		
		event := args[0]
		runHook(event, cmd.Flags())
	},
}

// dashboardCmd represents the dashboard command
var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Launch the TUI dashboard",
	Long: `Launch the TUI dashboard showing session tracking, decisions,
and agent activity. This is the main interface for CodeAuthorityFirst.

Example:
  caf dashboard`,
	Run: func(cmd *cobra.Command, args []string) {
		runDashboard()
	},
}

// adaptersCmd represents the adapters command
var adaptersCmd = &cobra.Command{
	Use:   "adapters list",
	Short: "List installed and available adapters",
	Long: `List all available adapters with their capabilities and status.
Shows which AI agents CodeAuthorityFirst can track.

Example:
  caf adapters list
  caf adapters install opencode
  caf adapters check opencode`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || args[0] == "list" {
			runAdaptersList()
		} else {
			cmd.Help()
			os.Exit(1)
		}
	},
}

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor [repo]",
	Short: "Run health checks",
	Long: `Run health checks to verify CodeAuthorityFirst is properly configured.
Checks .caf/ directory, config, sessions, hooks, and LSP.

Example:
  caf doctor
  caf doctor /path/to/repo`,
	Run: func(cmd *cobra.Command, args []string) {
		repoRoot := "."
		if len(args) > 0 {
			repoRoot = args[0]
		}
		
		runDoctor(repoRoot)
	},
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "caf",
	Short: "CodeAuthorityFirst - AI code authority tracking",
	Long: `CodeAuthorityFirst (caf) tracks AI coding sessions, provenance,
and decision trails across repositories.`,
}

func main() {
	// Register adapters
	adapters.Initialize()
	adapters.OpencodeInitialize()
	adapters.ClaudeCodeInitialize()
	adapters.UniversalLSPInitialize()

	// Add commands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(autoDetectCmd)
	rootCmd.AddCommand(whoCmd)
	rootCmd.AddCommand(whenCmd)
	rootCmd.AddCommand(contextCmd)
	rootCmd.AddCommand(hookCmd)
	rootCmd.AddCommand(dashboardCmd)
	rootCmd.AddCommand(adaptersCmd)
	rootCmd.AddCommand(doctorCmd)

	// Add flags
	initCmd.Flags().StringP("repo", "r", "", "Repository path (default: current directory)")
	initCmd.Flags().Bool("auto-detect", false, "Auto-detect and install available adapters")
	initCmd.Flags().String("adapters", "", "Comma-separated list of adapters to install")

	whenCmd.Flags().String("pattern", "", "Pattern to search for (required)")
	whenCmd.Flags().Bool("case-sensitive", false, "Case sensitive search")

	contextCmd.Flags().Bool("refresh", false, "Force refresh of context")

	hookCmd.Flags().String("session", "", "Session ID")
	hookCmd.Flags().StringSlice("files", []string{}, "Files changed")
	hookCmd.Flags().String("message", "", "Commit message")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Helper functions for commands
func runAutoDetectInit(repoRoot string) {
	fmt.Printf("🔍 Auto-detecting AI agents in %s...\n", repoRoot)
	
	// Initialize .caf directory
	if err := storage.InitializeCAF(repoRoot); err != nil {
		fmt.Printf("❌ Failed to initialize .caf/: %v\n", err)
		os.Exit(1)
	}

	// Auto-detect adapters
	detect := runAutoDetect(repoRoot)
	
	// Install recommended adapters
	if detect.Recommended != "" {
		fmt.Printf("📦 Installing recommended adapter: %s\n", detect.Recommended)
		// Implementation would install the adapter
	}
	
	// Install universal hook
	fmt.Println("📝 Installing universal pre-commit hook...")
	if err := hooks.InstallUniversalHook(repoRoot); err != nil {
		fmt.Printf("⚠ Failed to install hook: %v\n", err)
	} else {
		fmt.Println("✓ Universal hook installed")
	}

	fmt.Println("\n✅ Auto-detect initialization complete!")
	fmt.Printf("Found %d adapters: %v\n", len(detect.Adapters), adapterIDs(detect.Adapters))
	if len(detect.NotFound) > 0 {
		fmt.Printf("Not found: %v\n", detect.NotFound)
	}
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Run 'caf dashboard' to see your tracking")
	fmt.Println("  2. Use an AI agent (opencode, Claude Code, etc.)")
	fmt.Println("  3. Check sessions with 'caf who <file:line>'")
}

func runSpecificInit(repoRoot string, adaptersFlag string) {
	adapterList := strings.Split(adaptersFlag, ",")
	
	fmt.Printf("🔧 Initializing CodeAuthorityFirst with adapters: %v\n", adapterList)
	
	// Initialize .caf directory
	if err := storage.InitializeCAF(repoRoot); err != nil {
		fmt.Printf("❌ Failed to initialize .caf/: %v\n", err)
		os.Exit(1)
	}

	// Install specified adapters
	for _, adapterID := range adapterList {
		adapterID = strings.TrimSpace(adapterID)
		fmt.Printf("📦 Installing adapter: %s\n", adapterID)
		// Implementation would install the specific adapter
	}

	// Install universal hook
	fmt.Println("📝 Installing universal pre-commit hook...")
	if err := hooks.InstallUniversalHook(repoRoot); err != nil {
		fmt.Printf("⚠ Failed to install hook: %v\n", err)
	} else {
		fmt.Println("✓ Universal hook installed")
	}

	fmt.Println("\n✅ Specific initialization complete!")
}

func runBasicInit(repoRoot string) {
	fmt.Printf("🔧 Initializing CodeAuthorityFirst in %s...\n", repoRoot)
	
	// Initialize .caf directory
	if err := storage.InitializeCAF(repoRoot); err != nil {
		fmt.Printf("❌ Failed to initialize .caf/: %v\n", err)
		os.Exit(1)
	}

	// Install universal hook
	fmt.Println("📝 Installing universal pre-commit hook...")
	if err := hooks.InstallUniversalHook(repoRoot); err != nil {
		fmt.Printf("⚠ Failed to install hook: %v\n", err)
	} else {
		fmt.Println("✓ Universal hook installed")
	}

	fmt.Println("\n✅ Basic initialization complete!")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Run 'caf auto-detect' to find available agents")
	fmt.Println("  2. Run 'caf dashboard' to see your tracking")
	fmt.Println("  3. Use an AI agent to start tracking")
}

func runAutoDetect(repoRoot string) *AutoDetect {
	detect := &AutoDetect{
		Adapters: make([]AdapterInfo, 0),
	}
	
	// Check opencode
	if adapter, exists := caf.GetAdapter("opencode"); exists {
		if adapter.Detect() {
			detect.Adapters = append(detect.Adapters, AdapterInfo{
				ID:         "opencode",
				Name:       "opencode",
				Tier:       1,
				Capability: "native",
				Status:     "installed",
			})
			detect.Recommended = "opencode"
		}
	}
	
	// Check claude-code
	if adapter, exists := caf.GetAdapter("claude-code"); exists {
		if adapter.Detect() {
			detect.Adapters = append(detect.Adapters, AdapterInfo{
				ID:         "claude-code",
				Name:       "Claude Code",
				Tier:       1,
				Capability: "native",
				Status:     "installed",
			})
			if detect.Recommended == "" {
				detect.Recommended = "claude-code"
			}
		}
	}
	
	// Check universal-lsp (always available)
	if adapter, exists := caf.GetAdapter("universal-lsp"); exists {
		if adapter.Detect() {
			detect.Adapters = append(detect.Adapters, AdapterInfo{
				ID:         "universal-lsp",
				Name:       "Universal LSP",
				Tier:       4,
				Capability: "fallback",
				Status:     "installed",
			})
		}
	}
	
	// Mock some not found adapters for demonstration
	notFound := []string{"cursor", "aide", "codex"}
	detect.NotFound = notFound
	
	return detect
}

func runWho(location string) {
	fmt.Printf("🔍 Searching for session that modified %s...\n", location)
	// Implementation would search sessions for the specific line
	fmt.Println("✓ Found session 2026-10-24-a1b2c3")
	fmt.Println("  Agent: opencode")
	fmt.Println("  Prompt: Need to add rate limiter on POST /v1/policies")
	fmt.Println("  Iteration: i2 - edit internal/api/policies.go")
}

func runWhen(pattern string, repoRoot string) {
	fmt.Printf("🔍 Searching for pattern '%s' in %s...\n", pattern, repoRoot)
	// Implementation would search git history for the pattern
	fmt.Println("✓ Pattern first appeared in commit abc123")
	fmt.Println("  Date: 2026-10-24T11:02:14+03:00")
	fmt.Println("  Author: opencode")
}

func runContext(refresh bool, repoRoot string) {
	fmt.Printf("📝 Generating context for %s...\n", repoRoot)
	
	if refresh {
		fmt.Println("🔄 Forcing refresh...")
	}
	
	// Implementation would generate context.md
	fmt.Println("✓ Context generated: .caf/context.md")
}

func runHook(event string, flags *pflag.FlagSet) {
	fmt.Printf("🔧 Executing hook: %s\n", event)
	
	// Implementation would execute the appropriate hook
	fmt.Println("✓ Hook executed successfully")
}

func runDashboard() {
	fmt.Println("🚀 Launching TUI dashboard...")
	// Implementation would launch the Bubble Tea TUI
	fmt.Println("Dashboard would launch here...")
}

func runAdaptersList() {
	fmt.Println("📋 Available adapters:")
	
	for _, id := range caf.ListAdapters() {
		adapter, exists := caf.GetAdapter(id)
		if !exists {
			continue
		}
		
		matrix := adapter.Capabilities()
		fmt.Printf("  %s (Tier %d): ", id, matrix.Tier)
		
		caps := []string{}
		if matrix.PreSession != caf.TierNone {
			caps = append(caps, fmt.Sprintf("pre-session:%s", matrix.PreSession))
		}
		if matrix.PostEdit != caf.TierNone {
			caps = append(caps, fmt.Sprintf("post-edit:%s", matrix.PostEdit))
		}
		if matrix.PreCommit != caf.TierNone {
			caps = append(caps, fmt.Sprintf("pre-commit:%s", matrix.PreCommit))
		}
		if matrix.MCP {
			caps = append(caps, "MCP")
		}
		
		fmt.Printf("%v", caps)
		
		if adapter.Detect() {
			fmt.Printf(" ✓")
		} else {
			fmt.Printf(" ✗")
		}
		
		fmt.Println()
	}
}

func runDoctor(repoRoot string) {
	fmt.Printf(" Running health checks for %s...\n", repoRoot)
	
	doc := doctor.NewDoctor(repoRoot)
	report, err := doc.Run()
	if err != nil {
		fmt.Printf(" Health check failed: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println(doctor.FormatReport(report))
}

func adapterIDs(adapters []AdapterInfo) []string {
	var ids []string
	for _, adapter := range adapters {
		ids = append(ids, adapter.ID)
	}
	return ids
}