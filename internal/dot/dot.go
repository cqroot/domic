package dot

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

type DotConfig struct {
	Src  string
	Dest string
}

func BaseDir() string {
	return filepath.Join(xdg.DataHome, "gmdots")
}

func DotsDir() string {
	return filepath.Join(BaseDir(), "dots")
}

func Dots() ([]string, error) {
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
