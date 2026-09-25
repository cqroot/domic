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
package apply

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cqroot/domic/internal/config"
	"github.com/cqroot/domic/internal/distribute"
	"github.com/cqroot/domic/internal/fileutil"
)

type Options struct {
	Verbose bool
}

func Run(opts Options) error {
	apps := config.Apps()
	if len(apps) == 0 {
		fmt.Println("no apps configured")
		return nil
	}

	base := config.SourceBase()
	cache := loadCache()
	var errs []string

	for _, app := range apps {
		if !app.BinAvailable() {
			if opts.Verbose {
				fmt.Printf("%-*s %s (binary %s not found)\n", nameColWidth, app.Name, paint(ansiGray, "skipped"), app.Bin)
			}
			continue
		}
		rawTarget, ok := config.TargetForOS(app)
		if !ok {
			if opts.Verbose {
				fmt.Printf("%-*s %s (no target for current OS)\n", nameColWidth, app.Name, paint(ansiGray, "skipped"))
			}
			continue
		}
		target, err := config.ExpandHome(rawTarget)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: expand target: %v", app.Name, err))
			continue
		}
		source := filepath.Join(base, app.Path)
		if err := applyEntry(app, source, target, cache); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", app.Name, err))
		}
	}

	if err := cache.save(); err != nil {
		errs = append(errs, fmt.Sprintf("save cache: %v", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("apply completed with errors:\n  %s", strings.Join(errs, "\n  "))
	}
	return nil
}

func applyEntry(app config.App, source, target string, cache cache) error {
	info, err := os.Stat(source)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("source path %s does not exist", source)
		}
		return fmt.Errorf("stat source %s: %w", source, err)
	}
	if info.IsDir() {
		return walkAndApply(app, source, target, cache)
	}
	return processFile(app, source, target, cache)
}

func walkAndApply(app config.App, sourceDir, targetDir string, cache cache) error {
	return filepath.WalkDir(sourceDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		target := distribute.Target(app, path, filepath.Join(targetDir, rel))
		return processFile(app, path, target, cache)
	})
}

func processFile(app config.App, source, target string, cache cache) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("create parent dir: %w", err)
	}

	distributed, err := distribute.Build(app, source)
	if err != nil {
		return err
	}
	sourceMD5 := distribute.MD5(distributed)

	targetInfo, statErr := os.Stat(target)
	if statErr != nil {
		if !errors.Is(statErr, os.ErrNotExist) {
			return fmt.Errorf("stat target: %w", statErr)
		}
		if err := os.WriteFile(target, distributed, 0o644); err != nil {
			return fmt.Errorf("write: %w", err)
		}
		cache.set(target, sourceMD5)
		fmt.Printf("%-*s %s %s\n", nameColWidth, app.Name, paint(ansiGreen, "applied"), target)
		return nil
	}

	if targetInfo.IsDir() {
		return fmt.Errorf("target %s is a directory, expected file", target)
	}

	targetMD5, err := fileutil.MD5(target)
	if err != nil {
		return fmt.Errorf("hash target: %w", err)
	}

	if targetMD5 == sourceMD5 {
		cache.set(target, sourceMD5)
		fmt.Printf("%-*s %s %s\n", nameColWidth, app.Name, paint(ansiGreen, "up to date"), target)
		return nil
	}

	if recorded, ok := cache.get(target); ok && recorded == targetMD5 {
		if err := os.WriteFile(target, distributed, 0o644); err != nil {
			return fmt.Errorf("write: %w", err)
		}
		cache.set(target, sourceMD5)
		fmt.Printf("%-*s %s %s\n", nameColWidth, app.Name, paint(ansiGreen, "applied"), target)
		return nil
	}

	fmt.Printf("%-*s %s\n  source: %s\n  target: %s\n", nameColWidth, app.Name, paint(ansiYellow, "target already exists and differs from source"), source, target)
	return nil
}

const (
	ansiReset  = "\033[0m"
	ansiRed    = "\033[31m"
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
	ansiGray   = "\033[90m"

	nameColWidth = 15
)

var colorEnabled = os.Getenv("NO_COLOR") == ""

func paint(color, s string) string {
	if !colorEnabled {
		return s
	}
	return color + s + ansiReset
}
