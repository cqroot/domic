package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/cqroot/domic/pkg/dotfile"
	"github.com/cqroot/domic/pkg/dotfiles"
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

# Print initialization guide based on current user.
domic init

# Clone the dotfiles repository from https://github.com/YOURNAME/dotfiles
domic init YOURNAME
domic init YOURNAME/mydotfiles
domic init https://github.com/YOURNAME/dotfiles`,
	Args: cobra.MatchAll(cobra.RangeArgs(0, 1), cobra.OnlyValidArgs),
	Run:  runInitCmd,
}

func runInitCmd(cmd *cobra.Command, args []string) {
	if len(args) == 1 {
		GitClone(args[0])
		return
	}

	PrintInitGuide()
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

func PrintInitGuide() {
	hasDotfiles := false

	printPrefix := func() {
		fmt.Print(`
  # Execute the following command to initialize your dotfiles repository
  # Before executing each command, you should make sure that the command is ok.

`)
		fmt.Print(
			text.FgYellow.Sprint("  mkdir -p "),
			text.FgGreen.Sprint(strconv.Quote(strings.ReplaceAll(stdpath.DotsDir(), "\\", "/"))),
			"\n\n",
		)
	}

	for _, dot := range dotfiles.Dotfiles {
		if dot.State() == dotfile.StateApplied {
			continue
		}

		_, err := os.Stat(dot.Dst())
		if err != nil {
			continue
		}

		dotPath := dot.Src()
		if strings.HasPrefix(filepath.Base(dotPath), ".") {
			dotPath = filepath.Join(filepath.Dir(dotPath), filepath.Base(dotPath)[1:])
		}

		if !hasDotfiles {
			hasDotfiles = true
			printPrefix()
		}

		if strings.Contains(dot.RelSrc, "/") {
			fmt.Print(
				text.FgYellow.Sprint("  mkdir -p "),
				text.FgGreen.Sprint(strconv.Quote(strings.ReplaceAll(filepath.Dir(dotPath), "\\", "/"))),
				text.FgYellow.Sprint(" && \\\n"),
			)
		}
		fmt.Println(
			text.FgYellow.Sprint("  mv"),
			text.FgGreen.Sprintf("%-50s", strconv.Quote(dot.Dst())),
			text.FgGreen.Sprint(strconv.Quote(strings.ReplaceAll(dotPath, "\\", "/"))),
		)

	}

	if hasDotfiles {
		fmt.Print(
			text.FgYellow.Sprint("\n  cd "),
			text.FgGreen.Sprint(strconv.Quote(strings.ReplaceAll(stdpath.BaseDir(), "\\", "/"))),
			"\n\n",
		)
	} else {
		fmt.Println(text.FgRed.Sprint(`No manageable dotfiles detected.
You can view supported applications at https://github.com/cqroot/domic/blob/main/docs/supported_applications.md`))
	}
}
