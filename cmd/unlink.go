package cmd

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/cqroot/dotm/pkg/dotfile"
)

func init() {
	rootCmd.AddCommand(unlinkCmd)
}

var unlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "Remove all symlinks from config directory",
	Long:  "Remove all symlinks from config directory",
	Run:   RunUnlinkCmd,
}

func RunUnlinkCmd(cmd *cobra.Command, args []string) {
	cfg, err := readConfig()
	cobra.CheckErr(err)

	t := newTable()

	t.AppendHeader(table.Row{"#", "Source Path", "Target Path", "Status"})
	for idx, dot := range cfg.Dots {
		hasOp, err := dotfile.UnlinkAll(&dot)
		if err != nil {
			t.AppendRow([]interface{}{idx, dot.Source, dot.Target, text.FgRed.Sprint(err.Error())})
			continue
		}

		if hasOp {
			t.AppendRow([]interface{}{idx, dot.Source, dot.Target, text.FgGreen.Sprint("Unlinked!")})
		}
	}

	t.Render()
}
