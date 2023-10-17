package config

import "github.com/spf13/cobra"

func (c *State) Wrap(f func(*State, *cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return f(c, cmd, args)
	}
}
