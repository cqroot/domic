package dotfiles

import (
	"os"

	"github.com/cqroot/doter/pkg/path"
)

func LocalDotNames() ([]string, error) {
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
