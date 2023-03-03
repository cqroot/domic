package dotmanager

import (
	"path/filepath"
	"runtime"

	"github.com/cqroot/gmdots/pkg/path"
)

func (dm DotManager) defaultDotfileMap() map[string]Dot {
	return map[string]Dot{
		// https://github.com/alacritty/alacritty#configuration
		"alacritty": {
			Exec: "alacritty",
			Src:  "alacritty",
			Dest: func() string {
				if runtime.GOOS == "linux" {
					return path.LinuxXdgConfigHome("alacritty")
				} else if runtime.GOOS == "windows" {
					return path.WindowsAppData("alacritty")
				}
				return ""
			}(),
		},

		"git": {
			Exec: "git",
			Src:  "git",
			Dest: func() string {
				return path.LinuxXdgConfigHome("git")
			}(),
		},

		// go env GOENV
		"go": {
			Exec: "go",
			Src:  "go",
			Dest: func() string {
				if runtime.GOOS == "linux" {
					return path.LinuxXdgConfigHome("go")
				} else if runtime.GOOS == "windows" {
					return path.WindowsAppData("go")
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
					return path.LinuxXdgConfigHome("nvim")
				} else if runtime.GOOS == "windows" {
					return path.WindowsLocalAppData("nvim")
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
					return path.LinuxXdgConfigHome("pip/pip.conf")
				} else if runtime.GOOS == "windows" {
					return path.WindowsAppData("pip/pip.ini")
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
