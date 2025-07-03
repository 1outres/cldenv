package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/1outres/cldenv/internal/context"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use <context>",
	Short: "Switch to a different context",
	Long: `Switch to a different context by creating symbolic links from ~/.claude/ 
to the specified context directory.`,
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
				fmt.Println("Use 'cldenv create <context>' to create a new context")
			}
			
			return nil
		}

		// Check if already active
		if manager.GetActiveContext() == contextName {
			fmt.Printf("Already using context '%s'\n", contextName)
			return nil
		}

		// Switch to the context
		if err := manager.SwitchContext(contextName); err != nil {
			return fmt.Errorf("failed to switch to context '%s': %w", contextName, err)
		}

		fmt.Printf("✓ Switched to context '%s'\n", contextName)
		return nil
	},
}