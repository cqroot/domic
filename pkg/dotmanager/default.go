package dotmanager

import (
	"path/filepath"

	"github.com/cqroot/gmdots/pkg/path"
)

func DefaultDotMap(goos string) map[string]Dot {
	return map[string]Dot{
		// ************************************************************
		// *  Custom                                                  *
		// ************************************************************
		"bash": {
			Exec: "bash",
			Src:  "bash",
			Dest: path.DotConfigPath("bash"),
			Doc:  "You can put some bash files in this directory and source them in `.bashrc`.",
		},

		"bin": {
			Src:  "bin",
			Dest: path.HomeDotPath("bin"),
			Doc:  "You can put some executable scripts in this directory.",
		},

		// ************************************************************
		// *  Standard                                                *
		// ************************************************************
		"alacritty": {
			Exec: "alacritty",
			Src:  "alacritty",
			Dest: func() string {
				if goos == "linux" {
					return path.DotConfigPath("alacritty")
				} else if goos == "windows" {
					return path.WindowsAppDataPath("alacritty")
				}
				return ""
			}(),
			Doc: "https://github.com/alacritty/alacritty#configuration",
		},

		"git": {
			Exec: "git",
			Src:  "git",
			Dest: path.DotConfigPath("git"),
		},

		"gitui": {
			Exec: "gitui",
			Src:  "gitui",
			Dest: func() string {
				switch goos {
				case "linux", "darwin":
					return path.DotConfigPath("gitui")
				case "windows":
					return path.WindowsAppDataPath("gitui")
				}
				return ""
			}(),
			Doc: "https://github.com/extrawurst/gitui/blob/master/KEY_CONFIG.md#key-config",
		},

		"go": {
			Exec: "go",
			Src:  "go",
			Dest: func() string {
				if goos == "linux" {
					return path.DotConfigPath("go")
				} else if goos == "windows" {
					return path.WindowsAppDataPath("go")
				}
				return ""
			}(),
			Doc: "go env GOENV",
		},

		"lf": {
			Exec: "lf",
			Src:  "lf",
			Dest: func() string {
				switch goos {
				case "linux", "darwin":
					return path.DotConfigPath("lf")
				case "windows":
					return path.WindowsLocalAppDataPath("lf")
				}
				return ""
			}(),
			Doc: "https://github.com/gokcehan/lf/blob/master/docstring.go#L240",
		},

		"nvim": {
			Exec: "nvim",
			Src:  "nvim",
			Dest: func() string {
				if goos == "linux" {
					return path.DotConfigPath("nvim")
				} else if goos == "windows" {
					return path.WindowsLocalAppDataPath("nvim")
				}
				return ""
			}(),
			Doc: "https://neovim.io/doc/user/starting.html#standard-path",
		},

		"pip": {
			Exec: "pip",
			Src:  "pip/pip.conf",
			Dest: func() string {
				if goos == "linux" {
					return path.DotConfigPath("pip/pip.conf")
				} else if goos == "windows" {
					return path.WindowsAppDataPath("pip/pip.ini")
				}
				return ""
			}(),
			Doc: "https://pip.pypa.io/en/stable/topics/configuration/#location",
		},

		"powershell": {
			Exec: "pwsh",
			Src:  "powershell",
			Dest: func() string {
				switch goos {
				case "linux", "darwin":
					return path.DotConfigPath("powershell")
				case "windows":
					return filepath.Join(path.HomeDir(), "Documents/PowerShell")
				}
				return ""
			}(),
			Doc: "https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_profiles?view=powershell-7.3#the-profile-files",
		},

		"sqlite": {
			Exec: "sqlite3",
			Src:  "sqlite/sqliterc",
			Dest: func() string {
				return filepath.Join(path.HomeDir(), ".sqliterc")
			}(),
		},

		"starship": {
			Exec: "starship",
			Src:  "starship/starship.toml",
			Dest: func() string {
				return path.DotConfigPath("starship.toml")
			}(),
			Doc: "https://starship.rs/config/#configuration",
		},

		"wezterm": {
			Exec: "wezterm",
			Src:  "wezterm",
			Dest: path.DotConfigPath("wezterm"),
			Doc:  "https://wezfurlong.org/wezterm/config/files.html",
		},
	}
}
