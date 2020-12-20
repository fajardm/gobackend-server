package cmd

import (
	"github.com/fajardm/gobackend-server/cmd/service"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "server",
}

func init() {
	rootCmd.AddCommand(service.GetCommand())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
