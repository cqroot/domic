package cmd

import (
	"fmt"
	"path"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/cqroot/dotm/pkg/dotmanager"
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
	baseDir, err := getBaseDir()
	cobra.CheckErr(err)

	dm, err := dotmanager.New(baseDir, path.Join(baseDir, "dotm.toml"), Tag)
	cobra.CheckErr(err)

	t := newTable()
	t.AppendHeader(table.Row{
		"#", "Type", "Source Path", "Target Path", "Status",
	})

	doNothing := true

	var results []dotmanager.ExecuteResult

	if link {
		results = dm.Apply()
	} else {
		results = dm.Revoke()
	}

	for idx, result := range results {
		if result.Err != nil {
			doNothing = false
			t.AppendRow([]interface{}{
				idx, result.Dot.Type,
				result.Dot.Source, result.Dot.Target,
				text.FgRed.Sprint(result.Err.Error()),
			})
			continue
		}

		if result.HasOp {
			doNothing = false
			if link {
				t.AppendRow([]interface{}{
					idx, result.Dot.Type,
					result.Dot.Source, result.Dot.Target,
					text.FgGreen.Sprint("Linked!"),
				})
			} else {
				t.AppendRow([]interface{}{
					idx, result.Dot.Type,
					result.Dot.Source, result.Dot.Target,
					text.FgGreen.Sprint("Unlinked!"),
				})
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
