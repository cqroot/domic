package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/cqroot/domic/internal/manager"
	"github.com/spf13/cobra"
)

var configFile string = "./domic.yaml"

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "./"
	}
	dotfilesDir := filepath.Join(homeDir, ".dotfiles")

	switch runtime.GOOS {
	case "windows":
		configFile = filepath.Join(dotfilesDir, "./domic_windows.yaml")
	case "darwin":
		configFile = filepath.Join(dotfilesDir, "/domic_darwin.yaml")
	default:
		configFile = filepath.Join(dotfilesDir, "domic.yaml")
	}
}

func NewApplyCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "apply",
		Short: "Update all dotfiles to the target path",
		Run: func(cmd *cobra.Command, args []string) {
			m := manager.New(configFile)
			cobra.CheckErr(m.Apply())
		},
	}
	return &cmd
}

func NewRootCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "domic",
		Short: "Manage your dotfiles more easily.",
		Long:  "Manage your dotfiles more easily.",
		Run: func(cmd *cobra.Command, args []string) {
			m := manager.New(configFile)
			checkResult, err := m.Check()
			cobra.CheckErr(err)

			for name, result := range checkResult {
				fmt.Printf("%s: %s\n", name, result)
			}
		},
	}

	cmd.AddCommand(NewApplyCmd())

	return &cmd
}

func Execute() {
	cobra.CheckErr(NewRootCmd().Execute())
}
