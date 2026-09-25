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

	"github.com/cqroot/domic/internal/config"
	"github.com/spf13/cobra"
)

func newConfigDirCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "configdir",
		Short: "Print the OS-specific domic config directory ($XDG_CONFIG_HOME/domic)",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(config.ConfigDir())
			return nil
		},
	}
}
