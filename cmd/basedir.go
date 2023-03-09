package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cqroot/domic/pkg/stdpath"
)

func init() {
	rootCmd.AddCommand(basedirCmd)
}

var basedirCmd = &cobra.Command{
	Use:   "basedir",
	Short: "Print the dotfile base directory of Domic",
	Long:  "Print the dotfile base directory of Domic",
	Run:   runBasedirCmd,
}

func runBasedirCmd(cmd *cobra.Command, args []string) {
	fmt.Println(strings.ReplaceAll(stdpath.BaseDir(), "\\", "/"))
}
