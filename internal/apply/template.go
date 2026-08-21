package apply

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

const tmplExt = ".tmpl"

func isTemplateSource(appName string, templateEnabled bool, source string) bool {
	return templateEnabled && strings.HasSuffix(source, tmplExt)
}

func templateTarget(target string) string {
	return strings.TrimSuffix(target, tmplExt)
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
	if err := t.Execute(&buf, templateVars()); err != nil {
		return nil, fmt.Errorf("execute template %s: %w", source, err)
	}
	return buf.Bytes(), nil
}

func templateVars() map[string]any {
	return map[string]any{
		"domic": map[string]any{
			"os": runtime.GOOS,
		},
	}
}

func md5Bytes(b []byte) string {
	h := md5.Sum(b)
	return fmt.Sprintf("%x", h)
}