package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/cqroot/doter/pkg/dotfiles"
	"github.com/cqroot/doter/pkg/path"
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
	Args:  cobra.MatchAll(cobra.RangeArgs(0, 1), cobra.OnlyValidArgs),
	Run:   runInitCmd,
}

func runInitCmd(cmd *cobra.Command, args []string) {
	if len(args) == 1 {
		GitClone(args[0])
		return
	}

	fmt.Println(`
  # Execute the following command to initialize your dotfiles repository
  # Before executing each command, you should make sure that the command is ok.`)
	fmt.Println()
	fmt.Println("  mkdir -p", strconv.Quote(strings.ReplaceAll(path.DotsDir(), "\\", "/")))

	t := table.NewWriter()
	t.SetStyle(table.Style{
		Box: table.BoxStyle{
			BottomLeft:       " ",
			BottomRight:      " ",
			BottomSeparator:  " ",
			Left:             " ",
			LeftSeparator:    " ",
			MiddleHorizontal: " ",
			MiddleSeparator:  " ",
			MiddleVertical:   " ",
			PaddingLeft:      " ",
			PaddingRight:     "",
			PageSeparator:    "\n",
			Right:            " ",
			RightSeparator:   " ",
			TopLeft:          " ",
			TopRight:         " ",
			TopSeparator:     " ",
			UnfinishedRow:    " ",
		},
		Options: table.Options{
			DrawBorder:      true,
			SeparateColumns: true,
			SeparateFooter:  true,
			SeparateHeader:  true,
			SeparateRows:    false,
		},
	})
	t.SetOutputMirror(os.Stdout)

	for _, dot := range dotfiles.Dotfiles {
		_, err := os.Stat(dot.Dst())
		if err != nil {
			continue
		}

		dotPath := dot.Src()
		if strings.HasPrefix(filepath.Base(dotPath), ".") {
			dotPath = filepath.Join(filepath.Dir(dotPath), filepath.Base(dotPath)[1:])
		}

		if strings.Contains(dot.RelSrc, "/") {
			fmt.Println("  mkdir -p", strconv.Quote(strings.ReplaceAll(filepath.Dir(dotPath), "\\", "/")))
		}

		t.AppendRow(table.Row{
			"mv", strconv.Quote(dot.Dst()),
			strconv.Quote(strings.ReplaceAll(dotPath, "\\", "/")),
		})
	}

	t.Render()

	fmt.Println("  cd", strconv.Quote(strings.ReplaceAll(path.BaseDir(), "\\", "/")))
	fmt.Println()
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
