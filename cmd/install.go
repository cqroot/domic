package cmd

import (
	"github.com/cqroot/domic/internal/install"
	"github.com/spf13/cobra"
)

func newInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Install configured apps to their target paths",
		RunE: func(cmd *cobra.Command, args []string) error {
			return install.Run()
		},
	}
}
