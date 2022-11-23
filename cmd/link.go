package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(linkCmd)
}

var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Create all symlinks to config directory",
	Long:  "Create all symlinks to config directory",
	Run:   RunLinkCmd,
}

func RunLinkCmd(cmd *cobra.Command, args []string) {
	cfg, err := readConfig()
	cobra.CheckErr(err)

	checkLink(cfg, OpLink)
}
