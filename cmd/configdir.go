package cmd

import (
	"fmt"

	"github.com/cqroot/domic/internal/config"
	"github.com/spf13/cobra"
)

func newConfigDirCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "configdir",
		Short: "Print the OS-specific domic config directory ($XDG_CONFIG_HOME/domic)",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(config.ConfigDir())
			return nil
		},
	}
}