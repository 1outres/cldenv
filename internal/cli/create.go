package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/1outres/cldenv/internal/context"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <context>",
	Short: "Create a new context",
	Long: `Create a new context directory with empty configuration files.
After creating a context, you can switch to it using 'cldenv use <context>'.`,
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

		// Check if context already exists
		if manager.ContextExists(contextName) {
			fmt.Printf("Context '%s' already exists.\n\n", contextName)
			
			// Show available contexts
			contexts := manager.GetContexts()
			if len(contexts) > 0 {
				fmt.Println("Existing contexts:")
				for _, ctx := range contexts {
					marker := "  "
					if ctx.Name == manager.GetActiveContext() {
						marker = "* "
					}
					fmt.Printf("%s%s\n", marker, ctx.Name)
				}
			}
			
			return nil
		}

		// Create the context
		if err := manager.CreateContext(contextName); err != nil {
			return fmt.Errorf("failed to create context '%s': %w", contextName, err)
		}

		fmt.Printf("✓ Created context '%s'\n", contextName)
		fmt.Printf("Use 'cldenv use %s' to switch to this context\n", contextName)
		return nil
	},
}