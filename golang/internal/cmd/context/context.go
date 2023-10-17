package context

import (
	"fmt"
	"github.com/spf13/cobra"
	"simonnordberg.com/ppcli/internal/config"
)

func NewCommand(state *config.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Work with profiles",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "use",
		Short: "Use context",
		Args:  cobra.ExactArgs(1),
		RunE:  state.Wrap(runUse),
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List profiles",
		Args:  cobra.NoArgs,
		RunE:  state.Wrap(runList),
	})
	return cmd
}

func runUse(state *config.State, cmd *cobra.Command, args []string) error {
	name := args[0]
	context := state.Config.ContextByName(name)
	if context == nil {
		state.Config.Contexts = append(state.Config.Contexts, &config.ConfigContext{
			Name: name,
		})
	}

	return state.WriteConfig()
}

func runList(state *config.State, cmd *cobra.Command, args []string) error {
	for _, context := range state.Config.Contexts {
		fmt.Println(context.Name)
	}
	return nil
}
