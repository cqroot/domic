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
	"runtime"
	"slices"

	"github.com/cqroot/domic/pkg/config"
	"github.com/cqroot/domic/pkg/utils"
)

type CheckResult struct {
	IsOk bool
	Err  error
}

type PackageChecker func(config.Config, config.Package) CheckResult

func CheckPackageOs(cfg config.Config, pkg config.Package) CheckResult {
	if len(pkg.SupportedOs) == 0 {
		return CheckResult{IsOk: true, Err: nil}
	}
	if !slices.Contains(pkg.SupportedOs, runtime.GOOS) {
		return CheckResult{IsOk: false, Err: fmt.Errorf("operating system (%s) does not support", runtime.GOOS)}
	}
	return CheckResult{IsOk: true, Err: nil}
}

func CheckPackageExec(cfg config.Config, pkg config.Package) CheckResult {
	if pkg.Exec == "" {
		return CheckResult{IsOk: true, Err: nil}
	}
	if !utils.CommandExists(pkg.Exec) {
		return CheckResult{IsOk: false, Err: fmt.Errorf("command (%s) does not exist", pkg.Exec)}
	}
	return CheckResult{IsOk: true, Err: nil}
}

func CheckPackageSymlink(cfg config.Config, pkg config.Package) CheckResult {
	// Check if source exists
	_, err := os.Lstat(pkg.Source)
	if err != nil {
		if os.IsNotExist(err) {
			// Target doesn't exist
			return CheckResult{IsOk: false, Err: fmt.Errorf("source does not exist (%s)", pkg.Source)}
		}
		return CheckResult{IsOk: false, Err: err}
	}

	// Check if target exists
	targetInfo, err := os.Lstat(pkg.Target)
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
		linkDest, err := os.Readlink(pkg.Target)
		if err != nil {
			return CheckResult{IsOk: false, Err: err}
		}

		// Check if symlink points to the source
		if linkDest == pkg.Source {
			// Correct symlink already exists - no action needed
			return CheckResult{IsOk: true, Err: nil}
		}
		return CheckResult{IsOk: false, Err: fmt.Errorf("target symlink points to different location (%s)", linkDest)}
	}

	// Target exists but is not a symlink
	return CheckResult{IsOk: false, Err: errors.New("target already exists and is not a symlink")}
}

func CheckPackage(cfg config.Config, pkg config.Package) CheckResult {
	checkers := []PackageChecker{
		CheckPackageOs,
		CheckPackageExec,
		CheckPackageSymlink,
	}

	for _, checker := range checkers {
		res := checker(cfg, pkg)
		if res.Err != nil || !res.IsOk {
			return res
		}
	}
	return CheckResult{IsOk: true, Err: nil}
}
