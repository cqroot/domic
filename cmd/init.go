package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cqroot/gmdots/pkg/path"
	"github.com/spf13/cobra"
)

var initEnableSsh bool

func init() {
	initCmd.PersistentFlags().BoolVarP(&initEnableSsh, "ssh", "s", false, "Use ssh instead of https when guessing repo url")

	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "",
	Long:  "",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run:   runInitCmd,
}

func runInitCmd(cmd *cobra.Command, args []string) {
	GitClone(args[0])
}

func GitClone(repo string) {
	switch strings.Count(repo, "/") {
	case 0:
		if initEnableSsh {
			repo = "git@github.com:" + repo + "/dotfiles.git"
		} else {
			repo = "https://github.com/" + repo + "/dotfiles.git"
		}

	case 1:
		if initEnableSsh {
			repo = "git@github.com:" + repo
		} else {
			repo = "https://github.com/" + repo
		}
	}

	repoDir := path.BaseDir()

	_, err := os.Stat(repoDir)
	if err != nil {
		if !os.IsNotExist(err) {
			cobra.CheckErr(err)
		}
	} else {
		cobra.CheckErr("Dotfiles " + repoDir + " already exists")
	}

	gitArgs := []string{"clone", repo, repoDir}
	fmt.Println("git", strings.Join(gitArgs, " "))

	err = ExecCmd("git", gitArgs...)
	cobra.CheckErr(err)
}

func ExecCmd(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}
