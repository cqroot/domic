package cmd

import (
	"github.com/cqroot/domic/internal/initcmd"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init <source>",
		Short: "Clone a dotfiles repo into $XDG_CONFIG_HOME/domic",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return initcmd.Run(args[0])
		},
	}
}
