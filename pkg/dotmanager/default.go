package dotmanager

import "runtime"

func (dm DotManager) defaultDotfileMap() map[string]Dot {
	return map[string]Dot{
		// https://github.com/alacritty/alacritty#configuration
		"alacritty": {
			Exec: "alacritty",
			Src:  "alacritty",
			Dest: func() string {
				if runtime.GOOS == "linux" {
					return dm.LinuxXdgConfigHome("alacritty")
				} else if runtime.GOOS == "windows" {
					return dm.WindowsAppData("alacritty")
				}
				return ""
			}(),
		},

		// go env GOENV
		"go": {
			Exec: "go",
			Dest: func() string {
				if runtime.GOOS == "linux" {
					return dm.LinuxXdgConfigHome("go")
				} else if runtime.GOOS == "windows" {
					return dm.WindowsAppData("go")
				}
				return ""
			}(),
		},

		// https://neovim.io/doc/user/starting.html#standard-path
		"nvim": {
			Exec: "nvim",
			Dest: func() string {
				if runtime.GOOS == "linux" {
					return dm.LinuxXdgConfigHome("nvim")
				} else if runtime.GOOS == "windows" {
					return dm.WindowsLocalAppData("nvim")
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
					return dm.LinuxXdgConfigHome("pip/pip.conf")
				} else if runtime.GOOS == "windows" {
					return dm.WindowsAppData("pip/pip.ini")
				}
				return ""
			}(),
		},
	}
}
