package path

import (
	"os"
	"path/filepath"
)

var (
	replacedHomeDir         string = ""
	replacedAppDataDir      string = ""
	replacedLocalAppDataDir string = ""
)

func ReplaceHomeDir(homeDir string) {
	replacedHomeDir = homeDir
}

func ReplaceAppDataDir(AppDataDir string) {
	replacedAppDataDir = AppDataDir
}

func ReplaceLocalAppDataDir(LocalAppDataDir string) {
	replacedLocalAppDataDir = LocalAppDataDir
}

func HomeDir() string {
	if replacedHomeDir != "" {
		return replacedHomeDir
	}

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

// $HOME/.config/{path}
func DotConfigPath(path string) string {
	return filepath.Join(HomeDir(), ".config", path)
}

// $HOME/.{path}
func HomeDotPath(path string) string {
	return filepath.Join(HomeDir(), "."+path)
}

// %APPDATA%/{path}
func WindowsAppDataPath(path string) string {
	appDataDir := os.Getenv("APPDATA")
	if replacedAppDataDir != "" {
		appDataDir = replacedAppDataDir
	}
	return filepath.Join(appDataDir, path)
}

// %LOCALAPPDATA%/{path}
func WindowsLocalAppDataPath(path string) string {
	localAppDataDir := os.Getenv("LOCALAPPDATA")
	if replacedLocalAppDataDir != "" {
		localAppDataDir = replacedLocalAppDataDir
	}
	return filepath.Join(localAppDataDir, path)
}
