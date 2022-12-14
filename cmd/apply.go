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
	if Verbose {
		t.AppendHeader(table.Row{
			"#", "Name", "Type", "Source Path", "Target Path", "Status",
		})
	} else {
		t.AppendHeader(table.Row{
			"#", "Name", "Type", "Status",
		})
	}

	doNothing := true

	for idx, dot := range dm.Dots {
		if apply {
			err = dot.Apply()
		} else {
			err = dot.Revoke()
		}

		msg := ""

		switch err {
		case nil:
			msg = text.FgGreen.Sprint("OK")
		case dotmanager.DotIgnoreError:
			continue
		default:
			msg = text.FgRed.Sprint(err.Error())
		}

		doNothing = false
		if Verbose {
			t.AppendRow([]interface{}{
				idx, dot.Name(), dot.Type(), dot.Source(), dot.Target(), msg,
			})
		} else {
			t.AppendRow([]interface{}{
				idx, dot.Name(), dot.Type(), msg,
			})
		}
	}

	if doNothing {
		fmt.Println("There are no dotfiles to process here.")
	} else {
		t.Render()
	}
}
