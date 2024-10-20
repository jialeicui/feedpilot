package main

import (
	"github.com/spf13/cobra"

	"github.com/jialeicui/feedpilot/pkg/config"
	"github.com/jialeicui/feedpilot/pkg/server"
)

func main() {
	cmd := rootCmd()
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feedpilot",
		Short: "feedpilot is a social media aggregator",
	}

	cmd.AddCommand(serveCmd())

	return cmd
}

func serveCmd() *cobra.Command {
	var (
		flagConfig string
	)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run the server",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load(flagConfig)
			if err != nil {
				panic(err)
			}
			srv := server.NewServer(cfg)
			if err := srv.Start(); err != nil {
				panic(err)
			}
		},
	}

	cmd.Flags().StringVarP(&flagConfig, "config", "c", "config.yaml", "Path to the config file")
	return cmd
}
