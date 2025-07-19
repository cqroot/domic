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

package config

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type Package struct {
	// Package source directory; if not provided, defaults to a directory named `name`
	// under the directory where domic.toml is located.
	Source string `toml:"source"`
	Target string `toml:"target"`
}

type Config struct {
	Dotfiles     map[string]*Package `toml:"dotfiles"`
	WorkDir      string
	Names        []string
	MaxNameLen   int // maxNameLen stores the maximum length of the 'name' field for alignment in formatted output
	MaxSourceLen int // maxSourceLen stores the maximum length of the 'source' field for alignment in formatted output
}

// ExpandPath expands a relative path into an absolute path.
func ExpandPath(path string) (string, error) {
	newPath := os.ExpandEnv(path)

	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path, err
		}
		newPath = filepath.Join(home, path[2:])
	}

	return filepath.Abs(newPath)
}

func FillConfig(baseDir string, config *Config) error {
	var err error
	for name := range config.Dotfiles {
		pkg := config.Dotfiles[name]
		if pkg.Source == "" {
			pkg.Source = filepath.Join(config.WorkDir, name)
		}

		pkg.Source, err = ExpandPath(pkg.Source)
		if err != nil {
			return err
		}
		pkg.Target, err = ExpandPath(pkg.Target)
		if err != nil {
			return err
		}

		config.Names = append(config.Names, name)
		if len(name) > config.MaxNameLen {
			config.MaxNameLen = len(name)
		}
		if len(pkg.Source) > config.MaxSourceLen {
			config.MaxSourceLen = len(pkg.Source)
		}
	}

	sort.Strings(config.Names)
	return nil
}

func LoadConfig(configFile string) (Config, error) {
	content, err := os.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = toml.Unmarshal([]byte(content), &cfg)
	if err != nil {
		return Config{}, err
	}

	cfg.WorkDir = filepath.Dir(configFile)
	FillConfig(filepath.Dir(configFile), &cfg)
	return cfg, nil
}
