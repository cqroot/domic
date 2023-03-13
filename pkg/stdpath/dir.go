package stdpath

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
	return filepath.Join(HomeDir(), ".dotfiles")
}

func DotsDir() string {
	return BaseDir()
}

// $HOME/.config/{path}
func DotConfigPath(path string) string {
	return filepath.Join(HomeDir(), ".config", path)
}

// $HOME/.{path}
func HomeDotPath(path string) string {
	return filepath.Join(HomeDir(), "."+path)
}

// $HOME/Library/Application Support/{path}
func ApplicationSupportPath(path string) string {
	return filepath.Join(HomeDir(), "Library/Application Support", path)
}

// %APPDATA%/{path}
func WindowsAppDataPath(path string) string {
	appDataDir := os.Getenv("APPDATA")
	return filepath.Join(appDataDir, path)
}

// %LOCALAPPDATA%/{path}
func WindowsLocalAppDataPath(path string) string {
	localAppDataDir := os.Getenv("LOCALAPPDATA")
	return filepath.Join(localAppDataDir, path)
}
