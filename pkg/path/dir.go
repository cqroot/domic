package path

import (
	"os"
	"path/filepath"
)

func HomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/"
	}
	return homeDir
}

func BaseDir() string {
	return filepath.Join(HomeDir(), ".local/share/gmdots")
}

func DotsDir() string {
	return filepath.Join(BaseDir(), "dots")
}
