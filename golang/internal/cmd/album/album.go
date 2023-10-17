package album

import (
	"github.com/spf13/cobra"
	"simonnordberg.com/ppcli/internal/config"
)

func NewCommand(state *config.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "album",
		Short:                 "Manage albums",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(&cobra.Command{
		Use:                   "import",
		Short:                 "Import albums",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  state.Wrap(runImport),
	})
	return cmd
}

func runImport(state *config.State, command *cobra.Command, strings []string) error {
	return nil
}
