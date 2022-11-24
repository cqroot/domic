package cmd

import (
	"os"
	"path"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/cqroot/dotm/pkg/config"
	"github.com/cqroot/dotm/pkg/dotfile"
)

var (
	verbose bool
	rootCmd = &cobra.Command{
		Use:   "dm",
		Short: "DotM - Manage dotfiles more easily.",
		Long:  "DotM - Manage dotfiles more easily.",
		Run:   runRootCmd,
	}
)

func Execute() {
	err := rootCmd.Execute()
	cobra.CheckErr(err)
}

func runRootCmd(cmd *cobra.Command, args []string) {
	cfg, err := readConfig()
	cobra.CheckErr(err)

	t := newTable()

	t.AppendHeader(table.Row{"#", "Source Path", "Target Path", "Status"})
	for idx, dot := range cfg.Dots {
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
	t.SetStyle(table.StyleLight)
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

func readConfig() (*config.Config, error) {
	baseDir, err := getBaseDir()
	cobra.CheckErr(err)
	return config.New(baseDir, path.Join(baseDir, "dotm.toml"))
}
