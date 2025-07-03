package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/1outres/cldenv/internal/context"
)

var (
	version = "dev"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cldenv",
	Short: "Manage Claude Code environments",
	Long: `cldenv is a CLI tool for managing multiple Claude Code environments by switching 
between different ~/.claude/CLAUDE.md and ~/.claude/settings.json configurations 
using symbolic links.`,
	Version: version,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listContexts()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	
	// Add subcommands
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(removeCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Handle first run migration
	if context.IsFirstRun() {
		if err := context.MigrateToDefault(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to migrate existing files: %v\n", err)
		}
	}
	
	// Ensure default context exists
	if err := context.EnsureDefaultContext(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to ensure default context: %v\n", err)
	}
}

// listContexts lists all available contexts
func listContexts() error {
	manager, err := context.NewManager()
	if err != nil {
		return fmt.Errorf("failed to create context manager: %w", err)
	}

	if err := manager.LoadContexts(); err != nil {
		return fmt.Errorf("failed to load contexts: %w", err)
	}

	contexts := manager.GetContexts()
	activeContext := manager.GetActiveContext()

	if len(contexts) == 0 {
		fmt.Println("No contexts available.")
		fmt.Println("Use 'cldenv create <context>' to create a new context")
		return nil
	}

	fmt.Println("Available contexts:")
	for _, ctx := range contexts {
		marker := "  "
		if ctx.Name == activeContext {
			marker = "* "
		}
		
		status := ""
		if ctx.Name == activeContext {
			status = " (active)"
		}
		
		fmt.Printf("%s%s%s\n", marker, ctx.Name, status)
	}

	fmt.Println()
	fmt.Println("Use 'cldenv use <context>' to switch context")
	fmt.Println("Use 'cldenv create <context>' to create new context")
	fmt.Println("Use 'cldenv remove <context>' to remove context")

	return nil
}