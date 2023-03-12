package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/cqroot/domic/pkg/dotfiles"
)

func init() {
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		dots, err := dotfiles.Dotfiles()
		cobra.CheckErr(err)

		names := make([]string, 0, len(dots))
		for name := range dots {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			dot := dots[name]

			fmt.Println(name + ": {")

			fmt.Println("    Src  :", dot.Src)
			fmt.Println("    Dst  :", dot.Dst)
			fmt.Println("    Exec :", dot.Exec)

			fmt.Println("}")
			fmt.Println()
		}
	},
}
