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
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply all dotfiles",
	Long:  "Apply all dotfiles",
	Run:   RunApplyCmd,
}

func RunApplyCmd(cmd *cobra.Command, args []string) {
	ApplyOrRevoke(true)
}

func ApplyOrRevoke(apply bool) {
	baseDir, err := getBaseDir()
	cobra.CheckErr(err)

	dm, err := dotmanager.New(baseDir, path.Join(baseDir, "dotm.toml"), Tag)
	cobra.CheckErr(err)

	t := newTable()
	t.AppendHeader(table.Row{
		"#", "Type", "Source Path", "Target Path", "Status",
	})

	doNothing := true

	for idx, dot := range dm.Dots {
		if apply {
			err = dot.Apply()
		} else {
			err = dot.Revoke()
		}

		switch err {
		case nil:
			doNothing = false
			t.AppendRow([]interface{}{
				idx, dot.Type(), dot.Source(), dot.Target(),
				text.FgGreen.Sprint("OK"),
			})
		case dotmanager.DotIgnoreError:
		default:
			doNothing = false
			t.AppendRow([]interface{}{
				idx, dot.Type(), dot.Source(), dot.Target(),
				text.FgRed.Sprint(err.Error()),
			})
		}
	}

	if doNothing {
		fmt.Println("There are no dotfiles to process here.")
	} else {
		t.Render()
	}
}
