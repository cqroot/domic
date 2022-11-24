package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(unlinkCmd)
}

var unlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "Remove all symlinks from config directory",
	Long:  "Remove all symlinks from config directory",
	Run:   RunUnlinkCmd,
}

func RunUnlinkCmd(cmd *cobra.Command, args []string) {
	LinkOrUnlink(false)
}
