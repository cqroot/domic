package cmd

import (
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/cqroot/gmdots/pkg/dotmanager"
)

var rootCmd = &cobra.Command{
	Use:   "gmdots",
	Short: "Dotfiles Manager for Gopher",
	Long:  "Dotfiles Manager for Gopher",
	Run:   runRootCmd,
}

func printStatus() {
	t := table.NewWriter()
	t.SetStyle(table.StyleRounded)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Dot", "Src", "Dest", "Status"})

	dm := dotmanager.Default()

	err := dm.Range(func(name string, dot dotmanager.Dot) {
		ok, err := dm.Check(name)

		if ok {
			t.AppendRow(table.Row{name, dot.Src, dot.Dest, text.FgGreen.Sprint("OK")})
		} else if err != nil {
			if strings.HasPrefix(err.Error(), "Skip") {
				t.AppendRow(table.Row{name, dot.Src, dot.Dest, text.FgYellow.Sprint(err.Error())})
			} else {
				t.AppendRow(table.Row{name, dot.Src, dot.Dest, text.FgRed.Sprint(err.Error())})
			}
		} else {
			t.AppendRow(table.Row{name, dot.Src, dot.Dest, "Not applied"})
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
