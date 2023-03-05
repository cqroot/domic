package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cqroot/gmdots/pkg/path"
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
	fmt.Println(path.BaseDir())
}
