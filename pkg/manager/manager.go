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

	"github.com/cqroot/domic/pkg/config"
	"github.com/fatih/color"
)

type CheckResult struct {
	IsOk bool
	Err  error
}

type Manager struct {
	config  config.Config
	workDir string
	verbose bool
}

type Operation int

const (
	OperationCheck Operation = iota
	OperationApply
)

func New(opts ...Option) (*Manager, error) {
	mgr := Manager{}
	for _, opt := range opts {
		opt(&mgr)
	}

	if mgr.workDir == "" {
		mgr.workDir = "."
	}
	os.Setenv("DOMIC_WORK_DIR", mgr.workDir)

	configFile, err := filepath.Abs(filepath.Join(mgr.workDir, "domic.toml"))
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return nil, err
	}
	mgr.config = cfg

	return &mgr, nil
}

func (mgr Manager) Execute(op Operation) error {
	for _, name := range mgr.config.Names {
		pkg := mgr.config.Dotfiles[name]

		formattedName := fmt.Sprintf("%s%s", name, strings.Repeat(" ", mgr.config.MaxNameLen-len(name)+1))
		formattedNameHighlight := fmt.Sprintf("%s*%s", name, strings.Repeat(" ", mgr.config.MaxNameLen-len(name)))
		formattedSource := fmt.Sprintf("%s%s", pkg.Source, strings.Repeat(" ", mgr.config.MaxSourceLen-len(pkg.Source)))

		res := CheckPackage(pkg.Source, pkg.Target)
		if res.Err != nil {
			fmt.Printf("%s ERROR: %s.\n", color.RedString(formattedName), res.Err)
			continue
		}
		if res.IsOk {
			fmt.Printf("%s %s => %s\n", color.GreenString(formattedName), formattedSource, pkg.Target)
			continue
		}

		switch op {
		case OperationCheck:
			fmt.Printf("%s %s => %s\n", color.YellowString(formattedName), formattedSource, pkg.Target)
		case OperationApply:
			err := os.Symlink(pkg.Source, pkg.Target)
			if err != nil {
				fmt.Printf("%s ERROR: %s.\n", color.RedString(formattedName), err)
			} else {
				fmt.Printf("%s %s => %s\n", color.GreenString(formattedNameHighlight), formattedSource, pkg.Target)
			}
		}
	}
	return nil
}

func (mgr Manager) Check() error {
	return mgr.Execute(OperationCheck)
}

func (mgr Manager) Apply() error {
	return mgr.Execute(OperationApply)
}

func CheckPackage(source, target string) CheckResult {
	// Check if source exists
	_, err := os.Lstat(source)
	if err != nil {
		if os.IsNotExist(err) {
			// Target doesn't exist
			return CheckResult{IsOk: false, Err: fmt.Errorf("source does not exist (%s)", source)}
		}
		return CheckResult{IsOk: false, Err: err}
	}

	// Check if target exists
	targetInfo, err := os.Lstat(target)
	if err != nil {
		if os.IsNotExist(err) {
			// Target doesn't exist
			return CheckResult{IsOk: false, Err: nil}
		}
		return CheckResult{IsOk: false, Err: err}
	}

	// Check if target is a symlink
	if targetInfo.Mode()&os.ModeSymlink != 0 {
		// Read the symlink destination
		linkDest, err := os.Readlink(target)
		if err != nil {
			return CheckResult{IsOk: false, Err: err}
		}

		// Check if symlink points to the source
		if linkDest == source {
			// Correct symlink already exists - no action needed
			return CheckResult{IsOk: true, Err: nil}
		}
		return CheckResult{IsOk: false, Err: fmt.Errorf("target symlink points to different location (%s)", linkDest)}
	}

	// Target exists but is not a symlink
	return CheckResult{IsOk: false, Err: errors.New("target already exists and is not a symlink")}
}
