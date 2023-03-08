package dotfiles

import (
	"os"

	"github.com/cqroot/domic/pkg/stdpath"
)

func LocalDotNames() ([]string, error) {
	files, err := os.ReadDir(stdpath.DotsDir())
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
