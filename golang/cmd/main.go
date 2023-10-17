/*
Copyright © 2023 Simon Nordberg <simon@simonnordberg.com>
*/
package main

import (
	"log"
	"os"
	"simonnordberg.com/ppcli/internal/cli"
	"simonnordberg.com/ppcli/internal/config"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("ppcli: ")
	log.SetOutput(os.Stderr)

	//	cobra.OnInitialize(config.InitConfig)
	//	cobra.OnFinalize(config.WriteConfig)
}

func main() {
	state := config.New()

	err := state.ReadConfig()
	if err != nil {
		return
	}

	rootCmd := cli.NewRootCommand(state)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
