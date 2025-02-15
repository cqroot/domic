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

	"github.com/fatih/color"

	"github.com/cqroot/domic/internal/utils"
	"gopkg.in/yaml.v3"
)

type DotfileConfig struct {
	Exec  string            `yaml:"exec"`
	Files map[string]string `json:"files"`
}

type Manager struct {
	workingDir string
	configFile string
	dotfiles   map[string]DotfileConfig
}

func New(configFile string) *Manager {
	return &Manager{
		workingDir: filepath.Dir(configFile),
		configFile: configFile,
		dotfiles:   make(map[string]DotfileConfig),
	}
}

func (m *Manager) LoadConfig() error {
	content, err := os.ReadFile(m.configFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &m.dotfiles)
	if err != nil {
		return err
	}

	return nil
}

type CheckResult error

var (
	CheckResultOk             CheckResult = errors.New("OK")
	CheckResultBinaryNotFound CheckResult = errors.New("binary not found")
	CheckResultSrcNotFound    CheckResult = errors.New("source file not found")
	CheckResultDstNotFound    CheckResult = errors.New("destination file not found")
	CheckResultFilesDifferent CheckResult = errors.New("files are different")
	CheckResultGetFileHashErr CheckResult = errors.New("get file hash error")
)

func (m *Manager) checkDotfile(name string, config DotfileConfig) CheckResult {
	for src, dst := range config.Files {
		src = filepath.Join(m.workingDir, name, src)
		dst = utils.ExpandPath(dst)

		if len(config.Exec) != 0 && !utils.CommandExists(config.Exec) {
			return fmt.Errorf("%w: %s", CheckResultBinaryNotFound, config.Exec)
		}

		if _, err := os.Stat(src); err != nil {
			return fmt.Errorf("%w: %s", CheckResultSrcNotFound, src)
		}

		if _, err := os.Stat(dst); err != nil {
			return fmt.Errorf("%w: %s", CheckResultDstNotFound, dst)
		}

		srcHash, err := utils.GetFileHash(src)
		if err != nil {
			return fmt.Errorf("%w: %s", CheckResultGetFileHashErr, err)
		}

		dstHash, err := utils.GetFileHash(dst)
		if err != nil {
			return fmt.Errorf("%w: %s", CheckResultGetFileHashErr, err)
		}

		if srcHash != dstHash {
			return fmt.Errorf("%w: %s -> %s", CheckResultFilesDifferent, src, dst)
		}
	}

	return CheckResultOk
}

func (m *Manager) Check() error {
	if err := m.LoadConfig(); err != nil {
		return err
	}

	result := make(map[string]CheckResult)
	maxKeyLen := 0
	for name, config := range m.dotfiles {
		if len(name) > maxKeyLen {
			maxKeyLen = len(name)
		}
		result[name] = m.checkDotfile(name, config)
	}

	for name, result := range result {
		output := ""
		if errors.Is(result, CheckResultOk) {
			output = color.GreenString(result.Error())
		} else if errors.Is(result, CheckResultBinaryNotFound) {
			output = color.HiBlackString(result.Error())
		} else {
			output = color.RedString(result.Error())
		}
		fmt.Printf("%s %s: %s\n", color.CyanString(name), strings.Repeat(" ", maxKeyLen-len(name)), output)
	}

	return nil
}

func (m *Manager) Apply() error {
	if err := m.LoadConfig(); err != nil {
		return err
	}

	for name, config := range m.dotfiles {
		err := m.checkDotfile(name, config)
		if errors.Is(err, CheckResultOk) {
			continue
		}

		if errors.Is(err, CheckResultBinaryNotFound) {
			fmt.Printf("%s %s\n", color.CyanString(name), color.HiBlackString("Ignored"))
			continue
		}

		for src, dst := range config.Files {
			src := filepath.Join(m.workingDir, name, src)
			dst := utils.ExpandPath(dst)

			if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
				return fmt.Errorf("error creating directory for %s: %v", name, err)
			}

			fmt.Printf("%s %s %s %s\n", color.CyanString(name), src, color.YellowString("->"), dst)
			if err := utils.CopyFile(src, dst); err != nil {
				return fmt.Errorf("error applying %s: %v", name, err)
			}
		}
	}
	return nil
}
