package cmd

import (
	"github.com/cqroot/domic/internal/status"
	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "List all configured apps and their sync status",
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, _ := cmd.Flags().GetBool("verbose")
			return status.Run(status.Options{ShowFiles: verbose})
		},
	}
}
