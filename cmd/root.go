// Copyright (C) 2026 Keith Chu <cqroot@outlook.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
package cmd

import (
	"fmt"
	"os"

	"github.com/cqroot/domic/internal/config"
	"github.com/cqroot/domic/internal/status"
	"github.com/cqroot/domic/internal/version"
	"github.com/spf13/cobra"
)

func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var verbose bool
	cmd := &cobra.Command{
		Use:           "domic",
		Short:         "A config-based cross-platform dotfiles manager",
		Long:          "A config-based cross-platform dotfiles manager",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return status.Run(status.Options{Verbose: resolveVerbose(cmd)})
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Name() == "init" {
				return nil
			}
			file, _ := cmd.Flags().GetString("config")
			return config.Load(file)
		},
	}
	cmd.PersistentFlags().String("config", "", "path to config file (default is $XDG_CONFIG_HOME/domic/domic.toml)")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "show per-file status rows")
	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newApplyCmd())
	cmd.AddCommand(newStatusCmd())
	cmd.AddCommand(newDiffCmd())
	cmd.AddCommand(newConfigDirCmd())
	cmd.Version = version.Get().String()
	return cmd
}

// resolveVerbose returns the verbose flag value with the user prefs from
// domic.config.toml as the fallback. CLI flag takes precedence.
func resolveVerbose(cmd *cobra.Command) bool {
	if cmd.Flags().Changed("verbose") {
		v, _ := cmd.Flags().GetBool("verbose")
		return v
	}
	return config.Verbose()
}
