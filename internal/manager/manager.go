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

package manager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cqroot/domic/internal/checker"
	"github.com/cqroot/domic/internal/config"
	"github.com/cqroot/domic/internal/utils"
	"github.com/fatih/color"
)

type Manager struct {
	workingDir string
	configFile string
	dotfiles   map[string]config.Dotfile
}

func New(workDir string) (*Manager, error) {
	mgr := Manager{
		workingDir: workDir,
		configFile: filepath.Join(workDir, "domic.yaml"),
	}

	dotfiles, err := config.LoadConfig(mgr.configFile)
	if err != nil {
		return nil, err
	}

	mgr.dotfiles = dotfiles
	return &mgr, nil
}

func MaxNameLen(dotfilesResult []checker.DotfileResult) int {
	maxNameLen := 0
	for _, dotfileResult := range dotfilesResult {
		if len(dotfileResult.Name) > maxNameLen {
			maxNameLen = len(dotfileResult.Name)
		}
	}
	return maxNameLen
}

func (m *Manager) Check() error {
	dotfilesResult := checker.CheckDotfiles(m.workingDir, m.dotfiles)
	maxNameLen := MaxNameLen(dotfilesResult)

	for _, dotfileResult := range dotfilesResult {
		emoji := color.GreenString("🟢")
		if dotfileResult.OkCnt != dotfileResult.TotalCnt {
			emoji = color.RedString("🔴")
		}
		fmt.Printf("%s  %s   %s %s\n",
			emoji,
			color.CyanString(dotfileResult.Name),
			strings.Repeat(" ", maxNameLen-len(dotfileResult.Name)),
			color.HiCyanString("%d/%d", dotfileResult.OkCnt, dotfileResult.TotalCnt))

		for _, fileResult := range dotfileResult.FileResults {
			if errors.Is(fileResult.Error, checker.FileCheckSrcNotFound) {
				fmt.Printf("    ❗  %s\n", color.RedString("%s", fileResult.Error))
			} else if errors.Is(fileResult.Error, checker.FileCheckGetFileHashErr) {
				fmt.Printf("    ❌  %s\n", color.RedString("%s", fileResult.Error))
			} else {
				fmt.Printf("    ⭕  %s\n", color.RedString("%s", fileResult.Error))
			}
		}
	}

	return nil
}

func (m *Manager) Apply() error {
	dotfilesResult := checker.CheckDotfiles(m.workingDir, m.dotfiles)
	maxNameLen := MaxNameLen(dotfilesResult)

	for _, dotfileResult := range dotfilesResult {
		if len(dotfileResult.Name) > maxNameLen {
			maxNameLen = len(dotfileResult.Name)
		}
	}

	for _, dotfileResult := range dotfilesResult {
		for _, fileResult := range dotfileResult.FileResults {
			if !errors.Is(fileResult.Error, checker.FileCheckDstNotFound) &&
				!errors.Is(fileResult.Error, checker.FileCheckFilesDifferent) {
				continue
			}

			fmt.Printf("%s %s %s %s %s\n",
				color.CyanString(dotfileResult.Name), strings.Repeat(" ", maxNameLen-len(dotfileResult.Name)),
				fileResult.Src, color.YellowString("->"), fileResult.Dst)

			if err := os.MkdirAll(filepath.Dir(fileResult.Dst), 0o755); err != nil {
				return fmt.Errorf("error creating directory for %s: %v", dotfileResult.Name, err)
			}

			if err := utils.CopyFile(fileResult.Src, fileResult.Dst); err != nil {
				return fmt.Errorf("error applying %s: %v", dotfileResult.Name, err)
			}
		}
	}

	return nil
}
