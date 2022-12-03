package cmd

import (
	"os"
	"path"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/cqroot/dotm/pkg/dotfile"
	"github.com/cqroot/dotm/pkg/dotmanager"
)

var (
	Tag     string
	rootCmd = &cobra.Command{
		Use:   "dm",
		Short: "DotM - Manage dotfiles more easily.",
		Long:  "DotM - Manage dotfiles more easily.",
		Run:   runRootCmd,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&Tag, "tag", "t", "", "use dotfiles with specified tags")
}

func Execute() {
	err := rootCmd.Execute()
	cobra.CheckErr(err)
}

func runRootCmd(cmd *cobra.Command, args []string) {
	t := newTable()

	t.AppendHeader(table.Row{"#", "Source Path", "Target Path", "Status"})
	for idx, dot := range dots() {
		state, descr := dotfile.CheckState(&dot)
		switch state {
		case dotfile.StateIgnored:
			t.AppendRow([]interface{}{idx, dot.Source, dot.Target, descr})
		case dotfile.StateExisted:
			t.AppendRow([]interface{}{idx, dot.Source, dot.Target, text.FgRed.Sprint(descr)})
		case dotfile.StateLinkNormal:
			t.AppendRow([]interface{}{idx, dot.Source, dot.Target, text.FgGreen.Sprint(descr)})
		case dotfile.StateLinkEmpty:
			t.AppendRow([]interface{}{idx, dot.Source, dot.Target, text.FgRed.Sprint(descr)})
		}
	}

	t.Render()
}

func newTable() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)
	t.Style().Format.Header = text.FormatDefault

	return t
}

func getBaseDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return path.Join(configDir, "dotfiles"), nil
}

func dots() []dotmanager.Dot {
	baseDir, err := getBaseDir()
	cobra.CheckErr(err)

	dm, err := dotmanager.New(baseDir, path.Join(baseDir, "dotm.toml"))
	cobra.CheckErr(err)

	return dm.DotsWithTag(Tag)
}
