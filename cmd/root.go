package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/cqroot/dotmanager/internal/configs"
	"github.com/cqroot/dotmanager/internal/dot"
)

var rootCmd = &cobra.Command{
	Use:   "dm",
	Short: "Dotfiles Manager for Gopher",
	Long:  "Dotfiles Manager for Gopher",
	Run:   runRootCmd,
}

func printStatus() {
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Dot", "Src", "Dest"})

	err := configs.RangeDotConfigs(func(dotName string, dotConfig dot.DotConfig) {
		t.AppendRow(table.Row{dotName, dotConfig.Src, dotConfig.Dest})
	})
	cobra.CheckErr(err)

	t.Render()
}

func runRootCmd(cmd *cobra.Command, args []string) {
	printStatus()
}

func Execute() {
	err := rootCmd.Execute()
	cobra.CheckErr(err)
}
