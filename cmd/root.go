package cmd

import (
	"fmt"
	"os"

	"github.com/cqroot/domic/internal/config"
	"github.com/cqroot/domic/internal/status"
	"github.com/spf13/cobra"
)

func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "domic",
		Short:         "A config-based cross-platform dotfiles manager",
		Long:          "A config-based cross-platform dotfiles manager",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return status.Run()
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			file, _ := cmd.Flags().GetString("config")
			return config.Load(file)
		},
	}
	cmd.PersistentFlags().String("config", "", "path to config file (default is $XDG_CONFIG_HOME/domic/config.toml)")
	cmd.AddCommand(newInstallCmd())
	cmd.AddCommand(newStatusCmd())
	return cmd
}
