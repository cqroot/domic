package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cqroot/doter/pkg/stdpath"
)

func init() {
	rootCmd.AddCommand(basedirCmd)
}

var basedirCmd = &cobra.Command{
	Use:   "basedir",
	Short: "",
	Long:  "",
	Run:   runBasedirCmd,
}

func runBasedirCmd(cmd *cobra.Command, args []string) {
	fmt.Println(strings.ReplaceAll(stdpath.BaseDir(), "\\", "/"))
}
