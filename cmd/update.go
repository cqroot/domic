package cmd

import (
	"github.com/spf13/cobra"

	"github.com/cqroot/doter/pkg/stdpath"
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
	repoDir := stdpath.BaseDir()

	gitArgs := []string{"-C", repoDir, "pull"}

	err := ExecCmd("git", gitArgs...)
	cobra.CheckErr(err)
}
