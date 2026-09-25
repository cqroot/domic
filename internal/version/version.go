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
package version

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// Version info variables. Overridden at build time via -ldflags.
var (
	version   = "dev"
	commit    = "none"
	date      = "unknown"
	builtWith = fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
)

// Info contains version and build information.
type Info struct {
	Version   string
	Commit    string
	Date      string
	BuiltWith string
}

// Get returns the current version information.
func Get() Info {
	return Info{
		Version:   version,
		Commit:    commit,
		Date:      date,
		BuiltWith: builtWith,
	}
}

const (
	ansiReset = "\033[0m"
	ansiCyan  = "\033[36m"
)

var colorEnabled = os.Getenv("NO_COLOR") == ""

// paint wraps s in cyan when color is enabled.
func paint(s string) string {
	if !colorEnabled {
		return s
	}
	return ansiCyan + s + ansiReset
}

// String returns a formatted version info string.
func (i Info) String() string {
	var sb strings.Builder
	sb.WriteString("\n  ")
	sb.WriteString(paint("• Version:      "))
	sb.WriteString(i.Version)

	sb.WriteString("\n  ")
	sb.WriteString(paint("• Commit:       "))
	sb.WriteString(i.Commit)

	sb.WriteString("\n  ")
	sb.WriteString(paint("• Built at:     "))
	sb.WriteString(i.Date)

	sb.WriteString("\n  ")
	sb.WriteString(paint("• Built with:   "))
	sb.WriteString(i.BuiltWith)
	return sb.String()
}