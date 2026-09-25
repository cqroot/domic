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
package initcmd

import "testing"

func TestResolveRepo(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"username only", "alice", "https://github.com/alice/dotfiles"},
		{"owner/repo", "alice/dots", "https://github.com/alice/dots"},
		{"https URL", "https://github.com/alice/dots", "https://github.com/alice/dots"},
		{"ssh URL", "git@github.com:alice/dots.git", "git@github.com:alice/dots.git"},
		{"ssh scheme", "ssh://git@github.com/alice/dots", "ssh://git@github.com/alice/dots"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := resolveRepo(c.in)
			if got != c.want {
				t.Errorf("resolveRepo(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}
