package cmd

import (
	"os"
	"path"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

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
	rootCmd.PersistentFlags().StringVarP(
		&Tag, "tag", "t", "", "use dotfiles with specified tags")
}

func Execute() {
	err := rootCmd.Execute()
	cobra.CheckErr(err)
}

func runRootCmd(cmd *cobra.Command, args []string) {
	baseDir, err := getBaseDir()
	cobra.CheckErr(err)

	dm, err := dotmanager.New(baseDir, path.Join(baseDir, "dotm.toml"), Tag)
	cobra.CheckErr(err)

	t := newTable()
	t.AppendHeader(table.Row{
		"#", "Type", "Source Path", "Target Path", "Status",
	})

	for idx, result := range dm.Check() {
		switch result.Level {
		case dotmanager.Ignored:
			t.AppendRow([]interface{}{
				idx, result.Dot.Type,
				result.Dot.Source, result.Dot.Target,
				result.Description,
			})
		case dotmanager.Info:
			t.AppendRow([]interface{}{
				idx, result.Dot.Type,
				result.Dot.Source, result.Dot.Target,
				text.FgGreen.Sprint(result.Description),
			})
		default:
			t.AppendRow([]interface{}{
				idx, result.Dot.Type,
				result.Dot.Source, result.Dot.Target,
				text.FgRed.Sprint(result.Description),
			})
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

	dm, err := dotmanager.New(baseDir, path.Join(baseDir, "dotm.toml"), Tag)
	cobra.CheckErr(err)

	return dm.Dots()
}
