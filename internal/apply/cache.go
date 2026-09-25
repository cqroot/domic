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
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/cqroot/domic/internal/config"
)

type cache map[string]string

func cachePath() string {
	return filepath.Join(config.ConfigDir(), ".domic-cache.json")
}

func loadCache() cache {
	c := cache{}
	data, err := os.ReadFile(cachePath())
	if err != nil {
		return c
	}
	_ = json.Unmarshal(data, &c)
	return c
}

func (c cache) save() error {
	if err := os.MkdirAll(config.ConfigDir(), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cachePath(), data, 0o644)
}

func (c cache) get(target string) (string, bool) {
	v, ok := c[target]
	return v, ok
}

func (c cache) set(target, md5sum string) {
	c[target] = md5sum
}
