package dotmanager

import (
	"path/filepath"
	"runtime"

	"github.com/cqroot/gmdots/pkg/path"
)

func (dm DotManager) defaultDotfileMap() map[string]Dot {
	return map[string]Dot{
		// ************************************************************
		// *  Custom                                                  *
		// ************************************************************
		"bash": {
			Exec: "bash",
			Src:  "bash",
			Dest: path.DotConfigPath("bash"),
		},

		// ************************************************************
		// *  Standard                                                *
		// ************************************************************

		// https://github.com/alacritty/alacritty#configuration
		"alacritty": {
			Exec: "alacritty",
			Src:  "alacritty",
			Dest: func() string {
				if runtime.GOOS == "linux" {
					return path.DotConfigPath("alacritty")
				} else if runtime.GOOS == "windows" {
					return path.WindowsAppDataPath("alacritty")
				}
				return ""
			}(),
		},

		"git": {
			Exec: "git",
			Src:  "git",
			Dest: path.DotConfigPath("git"),
		},

		// https://github.com/extrawurst/gitui/blob/master/KEY_CONFIG.md#key-config
		"gitui": {
			Exec: "gitui",
			Src:  "gitui",
			Dest: func() string {
				switch runtime.GOOS {
				case "linux", "darwin":
					return path.DotConfigPath("gitui")
				case "windows":
					return path.WindowsAppDataPath("gitui")
				}
				return ""
			}(),
		},

		// go env GOENV
		"go": {
			Exec: "go",
			Src:  "go",
			Dest: func() string {
				if runtime.GOOS == "linux" {
					return path.DotConfigPath("go")
				} else if runtime.GOOS == "windows" {
					return path.WindowsAppDataPath("go")
				}
				return ""
			}(),
		},

		// https://github.com/gokcehan/lf/blob/master/docstring.go#L240
		"lf": {
			Exec: "lf",
			Src:  "lf",
			Dest: func() string {
				switch runtime.GOOS {
				case "linux", "darwin":
					return path.DotConfigPath("lf")
				case "windows":
					return path.WindowsLocalAppDataPath("lf")
				}
				return ""
			}(),
		},

		// https://neovim.io/doc/user/starting.html#standard-path
		"nvim": {
			Exec: "nvim",
			Src:  "nvim",
			Dest: func() string {
				if runtime.GOOS == "linux" {
					return path.DotConfigPath("nvim")
				} else if runtime.GOOS == "windows" {
					return path.WindowsLocalAppDataPath("nvim")
				}
				return ""
			}(),
		},

		// https://pip.pypa.io/en/stable/topics/configuration/#location
		"pip": {
			Exec: "pip",
			Src:  "pip/pip.conf",
			Dest: func() string {
				if runtime.GOOS == "linux" {
					return path.DotConfigPath("pip/pip.conf")
				} else if runtime.GOOS == "windows" {
					return path.WindowsAppDataPath("pip/pip.ini")
				}
				return ""
			}(),
		},

		// https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_profiles?view=powershell-7.3#the-profile-files
		"powershell": {
			Exec: "pwsh",
			Src:  "powershell",
			Dest: func() string {
				switch runtime.GOOS {
				case "linux", "darwin":
					return path.DotConfigPath("powershell")
				case "windows":
					return filepath.Join(path.HomeDir(), "Documents/PowerShell")
				}
				return ""
			}(),
		},

		"sqlite": {
			Exec: "sqlite3",
			Src:  "sqlite/sqliterc",
			Dest: func() string {
				return filepath.Join(path.HomeDir(), ".sqliterc")
			}(),
		},
	}
}
