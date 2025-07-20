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
	"github.com/cqroot/domic/pkg/manager"
	"github.com/spf13/cobra"
)

var (
	workDir string
	verbose bool
)

// NewManager creates and returns a new Manager instance configured with
// command-line derived settings.
func NewManager() *manager.Manager {
	mgr, err := manager.New(
		manager.WithWorkDir(workDir),
		manager.WithVerbose(verbose),
	)
	cobra.CheckErr(err)

	return mgr
}

func NewRootCmd() *cobra.Command {
	rootCmd := cobra.Command{
		Use:   "domic",
		Short: "Manage your dotfiles more easily.",
		Long:  "Manage your dotfiles more easily.",
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(NewManager().Check())
		},
	}

	rootCmd.PersistentFlags().StringVarP(&workDir, "directory", "d", "~/.dotfiles", "working directory")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.AddCommand(NewApplyCmd())

	return &rootCmd
}

func Execute() {
	cobra.CheckErr(NewRootCmd().Execute())
}
