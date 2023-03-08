package cmd

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/cqroot/doter/pkg/dotfile"
	"github.com/cqroot/doter/pkg/dotfiles"
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
	names, err := dotfiles.LocalDotNames()
	cobra.CheckErr(err)

	dotfiles.ForEach(names, func(name string, dot dotfile.Dotfile) {
		if dot.IsIgnored() {
			return
		}

		err := dot.Apply()
		if err != nil {
			fmt.Print(text.FgRed.Sprintf("%s: %s\n", name, err.Error()))
		}
	})
}
