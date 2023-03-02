package cmd

import (
	"os"
	"path/filepath"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/cqroot/gmdots/configs"
	"github.com/cqroot/gmdots/internal/dot"
)

var rootCmd = &cobra.Command{
	Use:   "dg",
	Short: "Dotfiles Manager for Gopher",
	Long:  "Dotfiles Manager for Gopher",
	Run:   runRootCmd,
}

func printStatus() {
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Dot", "Src", "Dest", "Status"})

	err := configs.RangeDotConfigs(func(dotName string, dotConfig dot.DotConfig) {
		ok, err := configs.Check(
			filepath.Join(dot.DotsDir(), dotConfig.Src),
			dotConfig.Dest,
		)
		if ok {
			t.AppendRow(table.Row{dotName, dotConfig.Src, dotConfig.Dest, "OK"})
		} else if err != nil {
			t.AppendRow(table.Row{dotName, dotConfig.Src, dotConfig.Dest, err.Error()})
		} else {
			t.AppendRow(table.Row{dotName, dotConfig.Src, dotConfig.Dest, ""})
		}
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
