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

func LinuxXdgConfigHome(path string) string {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		xdgConfigHome = filepath.Join(HomeDir(), ".config")
	}
	return filepath.Join(xdgConfigHome, path)
}

func WindowsAppData(path string) string {
	return filepath.Join(os.Getenv("APPDATA"), path)
}

func WindowsLocalAppData(path string) string {
	return filepath.Join(os.Getenv("LOCALAPPDATA"), path)
}
