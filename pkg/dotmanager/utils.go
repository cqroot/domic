package dotmanager

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

func BaseDir() string {
	return filepath.Join(xdg.DataHome, "gmdots")
}

func DotsDir() string {
	return filepath.Join(BaseDir(), "dots")
}

func DotNames() ([]string, error) {
	files, err := os.ReadDir(DotsDir())
	if err != nil {
		return nil, err
	}

	dots := make([]string, 0, len(files))

	for _, file := range files {
		if file.IsDir() {
			dots = append(dots, file.Name())
		}
	}

	return dots, nil
}

func (m DotManager) Range(handleFunc func(name string, dot Dot)) error {
	names, err := DotNames()
	if err != nil {
		return err
	}

	for _, name := range names {
		dotfile, ok := m.dotMap[name]
		if !ok {
			continue
		}

		if dotfile.Dest == "-" {
			continue
		}

		handleFunc(name, dotfile)
	}

	return nil
}
