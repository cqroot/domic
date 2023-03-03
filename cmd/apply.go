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
	dm := dotmanager.Default()

	err := dm.Range(func(name string, dot dotmanager.Dot) {
		ok, err := dm.Check(name)
		if ok || err != nil {
			return
		}

		err = dm.Apply(name)
		if err != nil {
			fmt.Printf("%s: %s\n", name, err.Error())
		}
	})
	cobra.CheckErr(err)
}
