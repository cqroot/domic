package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cqroot/gmdots/pkg/path"
	"github.com/spf13/cobra"
)

func init() {
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
		repo = "https://github.com/" + repo + "/dotfiles.git"

	case 1:
		repo = "https://github.com/" + repo
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

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}

	err = cmd.Wait()
	return err
}
