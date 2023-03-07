package cmd

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

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
	dots := dotfiles.Dotfiles
	names, err := dotfiles.LocalDotNames()
	cobra.CheckErr(err)

	for _, name := range names {
		dot, ok := dots[name]
		if !ok {
			continue
		}

		if dot.IsIgnored() {
			continue
		}

		err := dot.Apply()
		if err != nil {
			fmt.Print(text.FgRed.Sprintf("%s: %s\n", name, err.Error()))
			continue
		}
	}
}
