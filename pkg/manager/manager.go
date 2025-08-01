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
	"sync"

	"github.com/cqroot/domic/pkg/config"
	"github.com/cqroot/domic/pkg/utils"
	"github.com/fatih/color"
)

type Manager struct {
	config  config.Config
	workDir string
	verbose bool
}

type Operation int

const (
	OperationCheck Operation = iota
	OperationApply
	OperationRemove
)

func New(opts ...Option) (*Manager, error) {
	mgr := Manager{}
	for _, opt := range opts {
		opt(&mgr)
	}

	if mgr.workDir == "" {
		mgr.workDir = "~/.dotfiles"
	}

	var err error
	mgr.workDir, err = utils.ExpandPath(mgr.workDir, ".")
	if err != nil {
		return nil, err
	}
	_ = os.Setenv("DOMIC_WORK_DIR", mgr.workDir)

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

func (mgr Manager) ExecutePackage(name string, op Operation) string {
	pkg := mgr.config.Dotfiles[name]

	formattedName := fmt.Sprintf("%s%s", name, strings.Repeat(" ", mgr.config.MaxNameLen-len(name)+1))
	formattedNameHighlight := fmt.Sprintf("%s*%s", name, strings.Repeat(" ", mgr.config.MaxNameLen-len(name)))
	formattedSource := fmt.Sprintf("%s%s", pkg.Source, strings.Repeat(" ", mgr.config.MaxSourceLen-len(pkg.Source)))

	separator := color.BlueString("=>")

	errString := func(err error) string {
		return fmt.Sprintf("%s %s\n",
			color.RedString(formattedName), color.RedString(err.Error()))
	}
	hlString := func() string {
		return fmt.Sprintf("%s %s %s %s\n",
			color.GreenString(formattedNameHighlight), formattedSource, separator, pkg.Target)
	}

	err := CheckPackage(mgr.config, *pkg)
	if err == nil {
		if op == OperationRemove {
			err := os.Remove(pkg.Target)
			if err != nil {
				return errString(err)
			} else {
				return hlString()
			}
		}
		return fmt.Sprintf("%s %s %s %s\n",
			color.GreenString(formattedName), formattedSource, separator, pkg.Target)
	}

	if !errors.Is(err, ErrCheckResultTargetNotExist) {
		if !mgr.verbose && errors.Is(err, ErrCheckResultCommandNotExist) {
			return ""
		}
		return errString(err)
	}

	// Only if the target does not exist does the change need to be applied
	switch op {
	case OperationCheck:
		return fmt.Sprintf("%s %s %s %s\n", color.YellowString(formattedName), formattedSource, separator, pkg.Target)
	case OperationApply:
		err := os.Symlink(pkg.Source, pkg.Target)
		if err != nil {
			return errString(err)
		} else {
			return hlString()
		}
	}
	return ""
}

type execResult struct {
	name   string
	output string
}

func (mgr Manager) Execute(op Operation) error {
	results := make(map[string]string)
	resultCh := make(chan execResult)

	wgWorker := sync.WaitGroup{}
	for _, name := range mgr.config.Names {
		wgWorker.Add(1)
		go func() {
			defer wgWorker.Done()
			resultCh <- execResult{name, mgr.ExecutePackage(name, op)}
		}()
	}

	wgResult := sync.WaitGroup{}
	wgResult.Add(1)
	go func() {
		idx := 0
		for result := range resultCh {
			results[result.name] = result.output
			idx++
		}
		wgResult.Done()
	}()

	wgWorker.Wait()
	close(resultCh)
	wgResult.Wait()

	for _, name := range mgr.config.Names {
		fmt.Print(results[name])
	}
	return nil
}

func (mgr Manager) Check() error {
	return mgr.Execute(OperationCheck)
}

func (mgr Manager) Apply() error {
	return mgr.Execute(OperationApply)
}

func (mgr Manager) Remove() error {
	return mgr.Execute(OperationRemove)
}
