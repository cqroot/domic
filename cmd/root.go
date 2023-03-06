package cmd

import (
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/cqroot/doter/pkg/dotmanager"
)

var DotManager = dotmanager.Default()

var rootCmd = &cobra.Command{
	Use:   "doter",
	Short: "Dotfiles Manager for Gopher",
	Long:  "Dotfiles Manager for Gopher",
	Run:   runRootCmd,
}

func printStatus() {
	t := table.NewWriter()
	t.SetStyle(table.StyleRounded)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Dot", "Src", "Dest", "Status"})

	err := DotManager.Range(func(name string, dot dotmanager.Dot) {
		ok, err := DotManager.Check(name)
		dest := strings.ReplaceAll(dot.Dest, "\\", "/")

		if ok {
			t.AppendRow(table.Row{name, dot.Src, dest, text.FgGreen.Sprint("✔")})
		} else {
			if err != nil {
				if strings.HasPrefix(err.Error(), "Skip") {
					t.AppendRow(table.Row{name, dot.Src, dest, text.FgYellow.Sprint(err.Error())})
				} else {
					t.AppendRow(table.Row{name, dot.Src, dest, text.FgRed.Sprint(err.Error())})
				}
			} else {
				t.AppendRow(table.Row{name, dot.Src, dest, "✖"})
			}
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
