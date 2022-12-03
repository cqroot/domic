package cmd

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/cqroot/dotm/pkg/dotfile"
)

func init() {
	rootCmd.AddCommand(linkCmd)
}

var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Create all symlinks to config directory",
	Long:  "Create all symlinks to config directory",
	Run:   RunLinkCmd,
}

func RunLinkCmd(cmd *cobra.Command, args []string) {
	LinkOrUnlink(true)
}

func LinkOrUnlink(link bool) {
	t := newTable()

	t.AppendHeader(table.Row{"#", "Source Path", "Target Path", "Status"})
	doNothing := true
	for idx, dot := range dots() {
		var hasOp bool
		var err error

		if link {
			hasOp, err = dotfile.Link(&dot)
		} else {
			hasOp, err = dotfile.Unlink(&dot)
		}
		if err != nil {
			doNothing = false
			t.AppendRow([]interface{}{idx, dot.Source, dot.Target, text.FgRed.Sprint(err.Error())})
			continue
		}

		if hasOp {
			doNothing = false
			if link {
				t.AppendRow([]interface{}{idx, dot.Source, dot.Target, text.FgGreen.Sprint("Linked!")})
			} else {
				t.AppendRow([]interface{}{idx, dot.Source, dot.Target, text.FgGreen.Sprint("Unlinked!")})
			}
			continue
		}
	}

	if doNothing {
		fmt.Println("There are no dotfiles to process here.")
	} else {
		t.Render()
	}
}
