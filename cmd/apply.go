package cmd

import (
	"fmt"

	"github.com/cqroot/gmdots/pkg/dotmanager"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "",
	Long:  "",
	Run:   runApplyCmd,
}

func runApplyCmd(cmd *cobra.Command, args []string) {
	dm, err := dotmanager.Default()
	cobra.CheckErr(err)

	err = dm.Range(func(name string, dot dotmanager.Dot) {
		err := dm.Apply(name)
		if err != nil {
			fmt.Print(name, ": ")
			cobra.CheckErr(err)
		}
	})
	cobra.CheckErr(err)
}
