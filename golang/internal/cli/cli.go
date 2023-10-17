package cli

import (
	"github.com/spf13/cobra"
	"simonnordberg.com/ppcli/internal/cmd/album"
	"simonnordberg.com/ppcli/internal/cmd/auth"
	"simonnordberg.com/ppcli/internal/cmd/context"
	"simonnordberg.com/ppcli/internal/config"
)

func NewRootCommand(state *config.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ppcli",
		Short:                 "Client to interact with a PhotoPrism instance",
		TraverseChildren:      true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		album.NewCommand(state),
		auth.NewCommand(state),
		context.NewCommand(state),
	)
	return cmd
}
