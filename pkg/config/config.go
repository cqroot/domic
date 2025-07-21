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

	"github.com/cqroot/domic/pkg/utils"
	"github.com/pelletier/go-toml/v2"
)

type Package struct {
	// Package source directory; if not provided, defaults to a directory named `name`
	// under the directory where domic.toml is located.
	Source      string   `toml:"source"`
	Target      string   `toml:"target"`
	Exec        string   `toml:"exec"`
	SupportedOs []string `toml:"os"`
}

type Config struct {
	Dotfiles     map[string]*Package `toml:"dotfiles"`
	WorkDir      string
	Names        []string
	MaxNameLen   int // maxNameLen stores the maximum length of the 'name' field for alignment in formatted output
	MaxSourceLen int // maxSourceLen stores the maximum length of the 'source' field for alignment in formatted output
}

func FillConfig(baseDir string, cfg *Config) error {
	var err error
	for name := range cfg.Dotfiles {
		pkg := cfg.Dotfiles[name]
		if pkg.Source == "" {
			pkg.Source = filepath.Join(cfg.WorkDir, name)
		}

		pkg.Source, err = utils.ExpandPath(pkg.Source, cfg.WorkDir)
		if err != nil {
			return err
		}
		pkg.Target, err = utils.ExpandPath(pkg.Target, cfg.WorkDir)
		if err != nil {
			return err
		}

		cfg.Names = append(cfg.Names, name)
		if len(name) > cfg.MaxNameLen {
			cfg.MaxNameLen = len(name)
		}
		if len(pkg.Source) > cfg.MaxSourceLen {
			cfg.MaxSourceLen = len(pkg.Source)
		}
	}

	sort.Strings(cfg.Names)
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
	err = FillConfig(filepath.Dir(configFile), &cfg)
	return cfg, err
}
