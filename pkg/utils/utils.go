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

package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ExpandPath expands a relative path into an absolute path.
func ExpandPath(path string, workDir string) (string, error) {
	newPath := os.ExpandEnv(path)

	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path, err
		}
		newPath = filepath.Join(home, path[2:])
	}

	if !filepath.IsAbs(newPath) {
		newPath = filepath.Join(workDir, newPath)
	}
	return filepath.Abs(newPath)
}

// CommandExists checks if a specified command is available in the system's executable path.
// It verifies the existence of the command by searching through directories listed in the PATH environment variable.
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
