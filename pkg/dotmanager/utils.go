package dotmanager

import (
	"os"

	"github.com/cqroot/doter/pkg/path"
)

func DotNames() ([]string, error) {
	files, err := os.ReadDir(path.DotsDir())
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
		dotfile, ok := m.DotMap[name]
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
