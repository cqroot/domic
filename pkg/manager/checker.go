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
	ErrCheckResultOk               CheckResultError = nil
	ErrCheckResultOsNotSupport     CheckResultError = errors.New("operating system does not support")
	ErrCheckResultCommandNotExist  CheckResultError = errors.New("command does not exist")
	ErrCheckResultSourceNotExist   CheckResultError = errors.New("source does not exist")
	ErrCheckResultDifferentSymlink CheckResultError = errors.New("target symlink points to different location")
	ErrCheckResultTargetExist      CheckResultError = errors.New("target already exists and is not a symlink")
	ErrCheckResultTargetNotExist   CheckResultError = errors.New("target does not exist")
)

type PackageChecker func(config.Config, config.Package) error

func CheckPackageOs(cfg config.Config, pkg config.Package) error {
	if len(pkg.SupportedOs) == 0 {
		return ErrCheckResultOk
	}
	if !slices.Contains(pkg.SupportedOs, runtime.GOOS) {
		return fmt.Errorf("%w: %s", ErrCheckResultOsNotSupport, runtime.GOOS)
	}
	return ErrCheckResultOk
}

func CheckPackageExec(cfg config.Config, pkg config.Package) error {
	if pkg.Exec == "" {
		return ErrCheckResultOk
	}
	if !utils.CommandExists(pkg.Exec) {
		return fmt.Errorf("%w: %s", ErrCheckResultCommandNotExist, pkg.Exec)
	}
	return ErrCheckResultOk
}

func CheckPackageSymlink(cfg config.Config, pkg config.Package) error {
	// Check if source exists
	_, err := os.Lstat(pkg.Source)
	if err != nil {
		if os.IsNotExist(err) {
			// Target doesn't exist
			return fmt.Errorf("%w: %s", ErrCheckResultSourceNotExist, pkg.Source)
		}
		return err
	}

	// Check if target exists
	targetInfo, err := os.Lstat(pkg.Target)
	if err != nil {
		if os.IsNotExist(err) {
			// Target doesn't exist
			return ErrCheckResultTargetNotExist
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
			return ErrCheckResultOk
		}
		return fmt.Errorf("%w: %s", ErrCheckResultDifferentSymlink, linkDest)
	}

	// Target exists but is not a symlink
	return ErrCheckResultTargetExist
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
	return ErrCheckResultOk
}
