package cmd

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/cqroot/domic/pkg/dotfile"
	"github.com/cqroot/domic/pkg/dotfiles"
)

func init() {
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Create symlinks to all dotfiles in basedir",
	Long:  "Create symlinks to all dotfiles in basedir",
	Run:   runApplyCmd,
}

func runApplyCmd(cmd *cobra.Command, args []string) {
	err := dotfiles.ForEach(func(name string, dot dotfile.Dotfile) {
		if dot.IsIgnored() {
			return
		}

		err := dot.Apply()
		if err != nil {
			fmt.Print(text.FgRed.Sprintf("%s: %s\n", name, err.Error()))
		}
	})
	cobra.CheckErr(err)
}
