package cmd

import (
	"github.com/cqroot/domic/internal/apply"
	"github.com/spf13/cobra"
)

func newApplyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "apply",
		Short: "Apply configured apps to their target paths",
		RunE: func(cmd *cobra.Command, args []string) error {
			return apply.Run()
		},
	}
}
