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
package fileutil

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func MD5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func SameContents(a, b string) (bool, error) {
	ha, err := MD5(a)
	if err != nil {
		return false, fmt.Errorf("hash %s: %w", a, err)
	}
	hb, err := MD5(b)
	if err != nil {
		return false, fmt.Errorf("hash %s: %w", b, err)
	}
	return ha == hb, nil
}
