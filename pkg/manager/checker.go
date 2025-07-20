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

type CheckResultError error

var (
	CheckResultOk               CheckResultError = nil
	CheckResultOsNotSupport     CheckResultError = errors.New("operating system does not support")
	CheckResultCommandNotExist  CheckResultError = errors.New("command does not exist")
	CheckResultSourceNotExist   CheckResultError = errors.New("source does not exist")
	CheckResultDifferentSymlink CheckResultError = errors.New("target symlink points to different location")
	CheckResultTargetExist      CheckResultError = errors.New("target already exists and is not a symlink")
	CheckResultTargetNotExist   CheckResultError = errors.New("target does not exist")
)

type PackageChecker func(config.Config, config.Package) error

func CheckPackageOs(cfg config.Config, pkg config.Package) error {
	if len(pkg.SupportedOs) == 0 {
		return CheckResultOk
	}
	if !slices.Contains(pkg.SupportedOs, runtime.GOOS) {
		return fmt.Errorf("%w: %s", CheckResultOsNotSupport, runtime.GOOS)
	}
	return CheckResultOk
}

func CheckPackageExec(cfg config.Config, pkg config.Package) error {
	if pkg.Exec == "" {
		return CheckResultOk
	}
	if !utils.CommandExists(pkg.Exec) {
		return fmt.Errorf("%w: %s", CheckResultCommandNotExist, pkg.Exec)
	}
	return CheckResultOk
}

func CheckPackageSymlink(cfg config.Config, pkg config.Package) error {
	// Check if source exists
	_, err := os.Lstat(pkg.Source)
	if err != nil {
		if os.IsNotExist(err) {
			// Target doesn't exist
			return fmt.Errorf("%w: %s", CheckResultSourceNotExist, pkg.Source)
		}
		return err
	}

	// Check if target exists
	targetInfo, err := os.Lstat(pkg.Target)
	if err != nil {
		if os.IsNotExist(err) {
			// Target doesn't exist
			return CheckResultTargetNotExist
		}
		return err
	}

	// Check if target is a symlink
	if targetInfo.Mode()&os.ModeSymlink != 0 {
		// Read the symlink destination
		linkDest, err := os.Readlink(pkg.Target)
		if err != nil {
			return err
		}

		// Check if symlink points to the source
		if linkDest == pkg.Source {
			// Correct symlink already exists - no action needed
			return CheckResultOk
		}
		return fmt.Errorf("%w: %s", CheckResultDifferentSymlink, linkDest)
	}

	// Target exists but is not a symlink
	return CheckResultTargetExist
}

func CheckPackage(cfg config.Config, pkg config.Package) error {
	checkers := []PackageChecker{
		CheckPackageOs,
		CheckPackageExec,
		CheckPackageSymlink,
	}

	for _, checker := range checkers {
		err := checker(cfg, pkg)
		if err != nil {
			return err
		}
	}
	return CheckResultOk
}
