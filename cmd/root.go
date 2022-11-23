package cmd

import (
	"os"
	"path"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/cqroot/dotm/pkg/common"
	"github.com/cqroot/dotm/pkg/config"
	"github.com/cqroot/dotm/pkg/symlink"
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

	checkLink(cfg, OpNothing)
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

const (
	OpNothing int = 0
	OpLink        = 1
	OpUnlink      = 2
)

func checkLink(cfg *config.Config, op int) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.Style().Format.Header = text.FormatDefault

	t.AppendHeader(table.Row{"#", "Source Path", "Target Path", "Status"})
	for idx, dot := range cfg.Dots {
		var code int
		var descr string

		if dot.Exec != "" && !common.CommandExists(dot.Exec) {
			descr = text.FgYellow.Sprint("Ignored")
		} else {
			switch op {
			case OpNothing:
				code, descr = symlink.CheckStatus(dot.Source, dot.Target)
			case OpLink:
				code, descr = symlink.LinkAndCheckStatus(dot.Source, dot.Target)
			case OpUnlink:
				code, descr = symlink.UnlinkAndCheckStatus(dot.Source, dot.Target)
			}

			if code == symlink.StatusOK {
				descr = text.FgGreen.Sprint(descr)
			} else if code == symlink.StatusLinked {
				descr = text.FgYellow.Sprint(descr)
			} else if code == symlink.StatusError {
				descr = text.FgRed.Sprint(descr)
			}
		}

		t.AppendRow([]interface{}{idx, dot.Source, dot.Target, descr})
	}
	t.Render()
}
