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
package initcmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
)

func Run(arg string) error {
	target := filepath.Join(xdg.ConfigHome, "domic")
	repo := resolveRepo(arg)
	return clone(repo, target)
}

func resolveRepo(arg string) string {
	if isFullURL(arg) {
		return arg
	}
	if !strings.Contains(arg, "/") {
		return "https://github.com/" + arg + "/dotfiles"
	}
	return "https://github.com/" + arg
}

func isFullURL(s string) bool {
	return strings.Contains(s, "://") || strings.HasPrefix(s, "git@")
}

func clone(repo, target string) error {
	cmd := exec.Command("git", "clone", repo, target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}
	return nil
}
