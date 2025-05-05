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
	"runtime"

	"github.com/cqroot/domic/internal/utils"
	"gopkg.in/yaml.v3"
)

type dotfileEntry struct {
	Exec         string            `yaml:"exec"`
	Files        map[string]string `yaml:"files"`
	FilesWindows map[string]string `yaml:"files_windows"`
	Dirs         map[string]string `yaml:"dirs"`
	DirsWindows  map[string]string `yaml:"dirs_windows"`
}

type Dotfile struct {
	Files map[string]string
	Dirs  map[string]string
}

func AdaptForRuntimeOs(dotfileItem dotfileEntry, dotfile *Dotfile) {
	switch runtime.GOOS {
	case "windows":
		if len(dotfileItem.FilesWindows) != 0 {
			dotfile.Files = dotfileItem.FilesWindows
		}
		if len(dotfileItem.DirsWindows) != 0 {
			dotfile.Dirs = dotfileItem.DirsWindows
		}
	}

	for name := range dotfile.Files {
		dotfile.Files[name] = utils.ExpandPath(dotfile.Files[name])
	}

	for name := range dotfile.Dirs {
		dotfile.Dirs[name] = utils.ExpandPath(dotfile.Dirs[name])
	}
}

func LoadConfig(configFile string) (map[string]Dotfile, error) {
	content, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	dotfileItems := map[string]dotfileEntry{}
	err = yaml.Unmarshal(content, &dotfileItems)
	if err != nil {
		return nil, err
	}

	dotfiles := map[string]Dotfile{}
	for name, dotfileItem := range dotfileItems {
		if dotfileItem.Exec != "" && !utils.CommandExists(dotfileItem.Exec) {
			continue
		}

		dotfile := Dotfile{
			Files: dotfileItem.Files,
			Dirs:  dotfileItem.Dirs,
		}
		AdaptForRuntimeOs(dotfileItem, &dotfile)

		if len(dotfile.Files) == 0 && len(dotfile.Dirs) == 0 {
			continue
		}

		dotfiles[name] = dotfile
	}

	return dotfiles, nil
}
