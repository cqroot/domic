package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(revokeCmd)
}

var revokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke all dotfiles",
	Long:  "Revoke all dotfiles",
	Run:   RunRevokeCmd,
}

func RunRevokeCmd(cmd *cobra.Command, args []string) {
	ApplyOrRevoke(false)
}
