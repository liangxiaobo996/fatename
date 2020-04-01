package main

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
)

func serverCmd() *cobra.Command {
	serverCmd := &cobra.Command{
		Use: "runserver",
		RunE: func(cmd *cobra.Command, args []string) error {
			runServer()
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			initConfig(cfgFile)
			return nil
		},
	}

	serverCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")

	return serverCmd
}

func main() {
	var RootCmd = cobra.Command{
		Use: "fatename",
	}

	RootCmd.AddCommand(serverCmd())

	if err := RootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
