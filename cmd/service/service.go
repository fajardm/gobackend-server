package service

import (
	"github.com/fajardm/gobackend-server/config"
	"github.com/fajardm/gobackend-server/internal/adapter/rest"
	"github.com/spf13/cobra"
)

var (
	configFile string
)

var command = &cobra.Command{
	Use:     "service",
	Aliases: []string{"svc"},
	Short:   "Run rest service",
	Run: func(c *cobra.Command, args []string) {
		conf, err := config.Load(configFile)
		if err != nil {
			panic(err)
		}

		rest.Serve(conf)
	},
}

func init() {
	command.Flags().StringVar(&configFile, "config", "./config.json", "Set config file path")
}

func GetCommand() *cobra.Command {
	return command
}
