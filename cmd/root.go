/*
Copyright (C) 2025 Keith Chu <cqroot@outlook.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package cmd

import (
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
			cobra.CheckErr(m.Check())
		},
	}

	cmd.AddCommand(NewApplyCmd())

	return &cmd
}

func Execute() {
	cobra.CheckErr(NewRootCmd().Execute())
}
