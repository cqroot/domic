package cmd

import (
	"github.com/cqroot/gmdots/pkg/path"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "",
	Long:  "",
	Run:   runUpdateCmd,
}

func runUpdateCmd(cmd *cobra.Command, args []string) {
	GitPull()
}

func GitPull() {
	repoDir := path.BaseDir()

	gitArgs := []string{"-C", repoDir, "pull"}

	err := ExecCmd("git", gitArgs...)
	cobra.CheckErr(err)
}
