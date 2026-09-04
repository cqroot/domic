package cmd

import (
	"github.com/cqroot/domic/internal/diffcmd"
	"github.com/spf13/cobra"
)

func newDiffCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "diff [app...]",
		Short: "Show diff between source and target for each configured app",
		RunE: func(cmd *cobra.Command, args []string) error {
			return diffcmd.Run(args)
		},
	}
}