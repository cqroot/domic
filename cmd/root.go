package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/cqroot/gmdots/pkg/dotmanager"
)

var rootCmd = &cobra.Command{
	Use:   "dg",
	Short: "Dotfiles Manager for Gopher",
	Long:  "Dotfiles Manager for Gopher",
	Run:   runRootCmd,
}

func printStatus() {
	t := table.NewWriter()
	t.SetStyle(table.StyleRounded)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Dot", "Src", "Dest", "Status"})

	dm, err := dotmanager.Default()
	cobra.CheckErr(err)

	err = dm.Range(func(name string, dot dotmanager.Dot) {
		ok, err := dm.Check(name)

		if ok {
			t.AppendRow(table.Row{name, dot.Src, dot.Dest, "OK"})
		} else if err != nil {
			t.AppendRow(table.Row{name, dot.Src, dot.Dest, err.Error()})
		} else {
			t.AppendRow(table.Row{name, dot.Src, dot.Dest, "Skipped"})
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
