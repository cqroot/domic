package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/cqroot/doter/pkg/dotfile"
	"github.com/cqroot/doter/pkg/dotfiles"
)

var (
	rootFlagAll     bool
	rootFlagIgnored bool
	rootCmd         = &cobra.Command{
		Use:   "doter",
		Short: "Manage your dotfiles more easily.",
		Long:  "Manage your dotfiles more easily.",
		Run:   runRootCmd,
	}
)

func init() {
	rootCmd.Flags().BoolVarP(&rootFlagAll, "all", "a", false, "show all dotfiles")
	rootCmd.Flags().BoolVarP(&rootFlagIgnored, "ignored", "i", false, "show ignored dotfiles")
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
			if rootFlagIgnored {
				return
			}
			t.AppendRow(table.Row{name, dot.Dst(), text.FgGreen.Sprint("✔")})

		case dotfile.StateUnapplied:
			if rootFlagIgnored {
				return
			}
			t.AppendRow(table.Row{name, dot.Dst(), "✖"})

		case dotfile.StateIgnored:
			if !rootFlagIgnored && !rootFlagAll {
				return
			}
			t.AppendRow(table.Row{name, dot.Dst(), text.FgYellow.Sprint("Ignored")})

		case dotfile.StateTargetAlreadyExists:
			if rootFlagIgnored {
				return
			}
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
