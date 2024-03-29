package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/cqroot/domic/pkg/dotfile"
	"github.com/cqroot/domic/pkg/dotfiles"
)

var (
	rootFlagAll     bool
	rootFlagIgnored bool
	rootCmd         = &cobra.Command{
		Use:   "domic",
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

	err := dotfiles.ForEach(func(name string, dot dotfile.Dotfile) {
		switch dot.State() {
		case dotfile.StateApplied:
			if rootFlagIgnored {
				return
			}
			t.AppendRow(table.Row{name, dot.Dst, text.FgGreen.Sprint("✔")})

		case dotfile.StateUnapplied:
			if rootFlagIgnored {
				return
			}
			t.AppendRow(table.Row{name, dot.Dst, "✖"})

		case dotfile.StateIgnored:
			if !rootFlagIgnored && !rootFlagAll {
				return
			}
			t.AppendRow(table.Row{name, dot.Dst, text.FgYellow.Sprint("Ignored")})

		case dotfile.StateTargetAlreadyExists:
			if rootFlagIgnored {
				return
			}
			t.AppendRow(table.Row{name, dot.Dst, text.FgRed.Sprint("Destination dotfile already exists")})
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
