package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cqroot/domic/pkg/stdpath"
)

var initEnableSsh bool

func init() {
	initCmd.PersistentFlags().BoolVarP(&initEnableSsh, "ssh", "s", false, "Use ssh instead of https when guessing repo url")

	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the local dotfiles repository",
	Long: `Initialize the local dotfiles repository

# Clone the dotfiles repository from https://github.com/YOURNAME/dotfiles
domic init YOURNAME
domic init YOURNAME/mydotfiles
domic init https://github.com/YOURNAME/dotfiles`,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run:  runInitCmd,
}

func runInitCmd(cmd *cobra.Command, args []string) {
	if len(args) == 1 {
		GitClone(args[0])
		return
	}
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

	repoDir := stdpath.BaseDir()

	_, err := os.Stat(repoDir)
	if err != nil {
		if !os.IsNotExist(err) {
			cobra.CheckErr(err)
		}
	} else {
		cobra.CheckErr("Dotfiles " + repoDir + " already exists")
	}

	gitArgs := []string{"clone", repo, repoDir}

	err = ExecCmd("git", gitArgs...)
	cobra.CheckErr(err)
}

func ExecCmd(name string, arg ...string) error {
	fmt.Println(name, strings.Join(arg, " "))

	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}
