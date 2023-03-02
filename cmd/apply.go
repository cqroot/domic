package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/cqroot/gmdots/internal/configs"
	"github.com/cqroot/gmdots/internal/dot"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "",
	Long:  "",
	Run:   runApplyCmd,
}

func runApplyCmd(cmd *cobra.Command, args []string) {
	err := configs.RangeDotConfigs(func(dotName string, dotConfig dot.DotConfig) {
		err := configs.Apply(
			filepath.Join(dot.DotsDir(), dotConfig.Src),
			dotConfig.Dest,
		)
		if err != nil {
			fmt.Print(dotName, ": ")
			cobra.CheckErr(err)
		}
	})
	cobra.CheckErr(err)
}
