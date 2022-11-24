package cmd

import (
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
	cfg, err := readConfig()
	cobra.CheckErr(err)

	t := newTable()

	t.AppendHeader(table.Row{"#", "Source Path", "Target Path", "Status"})
	for idx, dot := range cfg.Dots {
		hasOp, err := dotfile.LinkAll(&dot)
		if err != nil {
			t.AppendRow([]interface{}{idx, dot.Source, dot.Target, text.FgRed.Sprint(err.Error())})
			continue
		}

		if hasOp {
			t.AppendRow([]interface{}{idx, dot.Source, dot.Target, text.FgGreen.Sprint("Linked!")})
		}
	}

	t.Render()
}
