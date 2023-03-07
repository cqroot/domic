package cmd

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/cqroot/doter/pkg/dotfiles"
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

		err := dot.Revoke()
		if err != nil {
			fmt.Print(text.FgRed.Sprintf("%s: %s\n", name, err.Error()))
			continue
		}
	}
}
