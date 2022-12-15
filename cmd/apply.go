package cmd

import (
	"fmt"
	"path"
	"sync"

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

	var wg sync.WaitGroup
	rowChan := make(chan table.Row)

	for idx, dot := range dm.Dots {
		wg.Add(1)

		go func(idx int, dot dotmanager.Dot) {
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
				rowChan <- nil
				return
			default:
				msg = text.FgRed.Sprint(err.Error())
			}

			if Verbose {
				rowChan <- []interface{}{
					idx, dot.Name(), dot.Type(), dot.Source(), dot.Target(), msg,
				}
			} else {
				rowChan <- []interface{}{
					idx, dot.Name(), dot.Type(), msg,
				}
			}
		}(idx, dot)
	}

	rows := make([]table.Row, 0)

	go func() {
		for row := range rowChan {
			if row != nil {
				rows = append(rows, row)
			}
			wg.Done()
		}
	}()

	wg.Wait()
	close(rowChan)

	if len(rows) == 0 {
		fmt.Println("There are no dotfiles to process here.")
	} else {
		t.AppendRows(rows)
		t.Render()
	}
}
