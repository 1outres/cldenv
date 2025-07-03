package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/1outres/cldenv/internal/context"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove <context>",
	Short: "Remove a context",
	Long: `Remove a context directory and all its configuration files.
You cannot remove the default context or the currently active context.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		contextName := args[0]
		
		manager, err := context.NewManager()
		if err != nil {
			return fmt.Errorf("failed to create context manager: %w", err)
		}

		if err := manager.LoadContexts(); err != nil {
			return fmt.Errorf("failed to load contexts: %w", err)
		}

		// Check if context exists
		if !manager.ContextExists(contextName) {
			fmt.Printf("Context '%s' not found.\n\n", contextName)
			
			// Show available contexts
			contexts := manager.GetContexts()
			if len(contexts) > 0 {
				fmt.Println("Available contexts:")
				for _, ctx := range contexts {
					fmt.Printf("  %s\n", ctx.Name)
				}
			} else {
				fmt.Println("No contexts available.")
			}
			
			return nil
		}

		// Remove the context
		if err := manager.RemoveContext(contextName); err != nil {
			return fmt.Errorf("failed to remove context '%s': %w", contextName, err)
		}

		fmt.Printf("✓ Removed context '%s'\n", contextName)
		return nil
	},
}