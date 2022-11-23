package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(baseCmd)
}

var baseCmd = &cobra.Command{
	Use:   "base",
	Short: "Print base directory path",
	Long:  "Print base directory path",
	Run:   RunBaseCmd,
}

func RunBaseCmd(cmd *cobra.Command, args []string) {
	baseDir, err := getBaseDir()
	cobra.CheckErr(err)
	fmt.Println(baseDir)
}
