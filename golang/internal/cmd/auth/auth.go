package auth

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"simonnordberg.com/ppcli/internal/client"
	"simonnordberg.com/ppcli/internal/config"
)

var hostname string
var username string
var password string
var context string

func init() {
}

func NewCommand(state *config.State) *cobra.Command {
	loginCommand := &cobra.Command{
		Use:                   "login",
		Short:                 "Login",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  state.Wrap(runLogin),
	}

	loginCommand.Flags().StringVar(&context, "context", "CONTEXT", "Context to use")
	loginCommand.Flags().StringVar(&hostname, "hostname", "HOSTNAME", "Hostname, e.g. http://localhost:2342")
	loginCommand.Flags().StringVar(&username, "username", "USERNAME", "Username")
	loginCommand.Flags().StringVar(&password, "password", "PASSWORD", "Password")
	loginCommand.MarkFlagsRequiredTogether("hostname", "username", "password", "context")

	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(loginCommand)
	return cmd
}

func runLogin(state *config.State, command *cobra.Command, args []string) error {
	ctx := state.Config.ContextByName(context)
	if ctx == nil {
		return errors.New("invalid context")
	}

	c := client.NewClient(hostname)
	token, err := c.Authenticate(username, password)

	if err != nil {
		fmt.Println(err)
	}

	ctx.Url = hostname
	ctx.Token = token

	return state.WriteConfig()
}
