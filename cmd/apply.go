package cmd

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/cqroot/doter/pkg/dotmanager"
)

func init() {
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Create symlinks to all dotfiles in `basedir/dots`",
	Long:  "Create symlinks to all dotfiles in `basedir/dots`",
	Run:   runApplyCmd,
}

func runApplyCmd(cmd *cobra.Command, args []string) {
	err := DotManager.Range(func(name string, dot dotmanager.Dot) {
		ok, err := DotManager.Check(name)
		if ok || err != nil {
			return
		}

		err = DotManager.Apply(name)
		if err != nil {
			fmt.Printf("%s: %s\n", name, err.Error())
			return
		}

		fmt.Println(text.FgGreen.Sprint(name), "->", strings.ReplaceAll(dot.Dest, "\\", "/"))
	})
	cobra.CheckErr(err)
}
