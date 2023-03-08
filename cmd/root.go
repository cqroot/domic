package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/cqroot/doter/pkg/dotfile"
	"github.com/cqroot/doter/pkg/dotfiles"
)

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
	t.AppendHeader(table.Row{"Dot", "Dst", "Status"})

	names, err := dotfiles.LocalDotNames()
	cobra.CheckErr(err)

	dotfiles.ForEach(names, func(name string, dot dotfile.Dotfile) {
		switch dot.State() {
		case dotfile.StateApplied:
			t.AppendRow(table.Row{name, dot.Dst(), text.FgGreen.Sprint("✔")})
		case dotfile.StateUnapplied:
			t.AppendRow(table.Row{name, dot.Dst(), "✖"})
		case dotfile.StateIgnored:
			t.AppendRow(table.Row{name, dot.Dst(), text.FgYellow.Sprint("Ignored")})
		case dotfile.StateTargetAlreadyExists:
			t.AppendRow(table.Row{name, dot.Dst(), text.FgRed.Sprint("Destination dotfile already exists")})
		}
	})

	t.Render()
}

func runRootCmd(cmd *cobra.Command, args []string) {
	printStatus()
}

func Execute() {
	err := rootCmd.Execute()
	cobra.CheckErr(err)
}
