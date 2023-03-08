package cmd

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/cqroot/domic/pkg/dotfile"
	"github.com/cqroot/domic/pkg/dotfiles"
)

func init() {
	rootCmd.AddCommand(revokeCmd)
}

var revokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "",
	Long:  "",
	Run:   runRevokeCmd,
}

func runRevokeCmd(cmd *cobra.Command, args []string) {
	names, err := dotfiles.LocalDotNames()
	cobra.CheckErr(err)

	dotfiles.ForEach(names, func(name string, dot dotfile.Dotfile) {
		if dot.IsIgnored() {
			return
		}

		err := dot.Revoke()
		if err != nil {
			fmt.Print(text.FgRed.Sprintf("%s: %s\n", name, err.Error()))
		}
	})
}
