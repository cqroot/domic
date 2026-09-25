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
package distribute

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/cqroot/domic/internal/config"
)

const tmplExt = ".tmpl"

// IsTemplate reports whether the source should be rendered as a template
// for the given app.
func IsTemplate(app config.App, source string) bool {
	return app.Template && strings.HasSuffix(source, tmplExt)
}

// Target returns the effective target path for a source file, stripping the
// .tmpl suffix when the source is rendered as a template.
func Target(app config.App, source, target string) string {
	if IsTemplate(app, source) {
		return strings.TrimSuffix(target, tmplExt)
	}
	return target
}

// Build returns the bytes that would be written to the target for the given
// source (rendered template for .tmpl sources, file contents otherwise).
func Build(app config.App, source string) ([]byte, error) {
	if IsTemplate(app, source) {
		return renderTemplate(source)
	}
	data, err := os.ReadFile(source)
	if err != nil {
		return nil, fmt.Errorf("read source %s: %w", source, err)
	}
	return data, nil
}

// MD5 returns the hex MD5 digest of b.
func MD5(b []byte) string {
	h := md5.Sum(b)
	return fmt.Sprintf("%x", h)
}

// Vars returns the variables exposed to templates.
func Vars() map[string]any {
	return map[string]any{
		"domic": map[string]any{
			"os": runtime.GOOS,
		},
	}
}

func renderTemplate(source string) ([]byte, error) {
	data, err := os.ReadFile(source)
	if err != nil {
		return nil, fmt.Errorf("read template %s: %w", source, err)
	}
	t, err := template.New(filepath.Base(source)).Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("parse template %s: %w", source, err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, Vars()); err != nil {
		return nil, fmt.Errorf("execute template %s: %w", source, err)
	}
	return buf.Bytes(), nil
}
