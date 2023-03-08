package dotfiles

import (
	"path/filepath"
	"runtime"

	"github.com/cqroot/doter/pkg/dotfile"
	"github.com/cqroot/doter/pkg/path"
)

var Dotfiles = map[string]dotfile.Dotfile{
	// ************************************************************
	// *  Custom                                                  *
	// ************************************************************
	"bash": {
		Exec:    "bash",
		RelSrc:  "bash",
		DstFunc: func(string) string { return path.DotConfigPath("bash") },
		Doc:     "You can put some bash files in this directory and source them in `.bashrc`.",
	},

	"bin": {
		RelSrc:  "bin",
		DstFunc: func(string) string { return path.HomeDotPath("bin") },
		Doc:     "You can put some executable scripts in this directory.",
	},

	// ************************************************************
	// *  Standard                                                *
	// ************************************************************
	"alacritty": {
		Exec:   "alacritty",
		RelSrc: "alacritty",
		DstFunc: func(goos string) string {
			switch goos {
			case "linux", "darwin":
				return path.DotConfigPath("alacritty")
			case "windows":
				return path.WindowsAppDataPath("alacritty")
			}
			return ""
		},
		Doc: "https://github.com/alacritty/alacritty#configuration",
	},

	"dunst": {
		Exec:   "dunst",
		RelSrc: "dunst",
		DstFunc: func(goos string) string {
			switch goos {
			case "linux":
				return path.DotConfigPath("dunst")
			}
			return ""
		},
		Doc: "https://github.com/dunst-project/dunst/wiki#getting-started",
	},

	"git": {
		Exec:    "git",
		RelSrc:  "git",
		DstFunc: func(string) string { return path.DotConfigPath("git") },
	},

	"gitui": {
		Exec:   "gitui",
		RelSrc: "gitui",
		DstFunc: func(goos string) string {
			switch goos {
			case "linux", "darwin":
				return path.DotConfigPath("gitui")
			case "windows":
				return path.WindowsAppDataPath("gitui")
			}
			return ""
		},
		Doc: "https://github.com/extrawurst/gitui/blob/master/KEY_CONFIG.md#key-config",
	},

	"go": {
		Exec:   "go",
		RelSrc: "go",
		DstFunc: func(goos string) string {
			switch goos {
			case "linux":
				return path.DotConfigPath("go")
			case "darwin":
				return path.ApplicationSupportPath("go")
			case "windows":
				return path.WindowsAppDataPath("go")
			}
			return ""
		},
		Doc: "go env GOENV",
	},

	"i3": {
		Exec:   "i3",
		RelSrc: "i3",
		DstFunc: func(goos string) string {
			switch goos {
			case "linux":
				return path.DotConfigPath("i3")
			}
			return ""
		},
		Doc: "https://i3wm.org/docs/userguide.html#configuring",
	},

	"joplin-desktop": {
		Exec: func() string {
			switch runtime.GOOS {
			case "linux":
				return "joplin-desktop"
			case "windows":
				return "joplin"
			}
			return ""
		}(),
		RelSrc: "joplin-desktop",
		DstFunc: func(goos string) string {
			return path.DotConfigPath("joplin-desktop")
		},
	},

	"lf": {
		Exec:   "lf",
		RelSrc: "lf",
		DstFunc: func(goos string) string {
			switch goos {
			case "linux", "darwin":
				return path.DotConfigPath("lf")
			case "windows":
				return path.WindowsLocalAppDataPath("lf")
			}
			return ""
		},
		Doc: "https://github.com/gokcehan/lf/blob/master/docstring.go#L240",
	},

	"mpv": {
		Exec:   "mpv",
		RelSrc: "mpv",
		DstFunc: func(goos string) string {
			switch goos {
			case "linux", "darwin":
				return path.DotConfigPath("mpv")
			case "windows":
				return path.WindowsAppDataPath("mpv")
			}
			return ""
		},
		Doc: "https://mpv.io/manual/stable/#configuration-files",
	},

	"nushell": {
		Exec:   "nushell",
		RelSrc: "nushell",
		DstFunc: func(goos string) string {
			switch goos {
			case "windows":
				return path.WindowsAppDataPath("nushell")
			}
			return ""
		},
		Doc: "https://www.nushell.sh/book/configuration.html",
	},

	"nvim": {
		Exec:   "nvim",
		RelSrc: "nvim",
		DstFunc: func(goos string) string {
			switch goos {
			case "linux", "darwin":
				return path.DotConfigPath("nvim")
			case "windows":
				return path.WindowsLocalAppDataPath("nvim")
			}
			return ""
		},
		Doc: "https://neovim.io/doc/user/starting.html#standard-path",
	},

	"picom": {
		Exec:   "picom",
		RelSrc: "picom",
		DstFunc: func(goos string) string {
			switch goos {
			case "linux":
				return path.DotConfigPath("picom")
			}
			return ""
		},
	},

	"pip": {
		Exec:   "pip",
		RelSrc: "pip/pip.conf",
		DstFunc: func(goos string) string {
			switch goos {
			case "linux":
				return path.DotConfigPath("pip/pip.conf")
			case "darwin":
				return path.DotConfigPath("pip/pip.conf") // `$HOME/Library/Application Support/pip/pip.conf` is not supported.
			case "windows":
				return path.WindowsAppDataPath("pip/pip.ini")
			}
			return ""
		},
		Doc: "https://pip.pypa.io/en/stable/topics/configuration/#location",
	},

	"polybar": {
		Exec:   "polybar",
		RelSrc: "polybar",
		DstFunc: func(goos string) string {
			switch goos {
			case "linux":
				return path.DotConfigPath("polybar")
			}
			return ""
		},
		Doc: "https://github.com/polybar/polybar/wiki#where-to-start",
	},

	"powershell": {
		Exec:   "pwsh",
		RelSrc: "powershell",
		DstFunc: func(goos string) string {
			switch goos {
			case "linux", "darwin":
				return path.DotConfigPath("powershell")
			case "windows":
				return filepath.Join(path.HomeDir(), "Documents/PowerShell")
			}
			return ""
		},
		Doc: "https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_profiles?view=powershell-7.3#the-profile-files",
	},

	"rofi": {
		Exec:   "rofi",
		RelSrc: "rofi",
		DstFunc: func(goos string) string {
			switch goos {
			case "linux":
				return path.DotConfigPath("rofi")
			}
			return ""
		},
		Doc: "https://github.com/davatorium/rofi#configuration",
	},

	"sqlite": {
		Exec:    "sqlite3",
		RelSrc:  "sqlite/sqliterc",
		DstFunc: func(string) string { return filepath.Join(path.HomeDir(), ".sqliterc") },
	},

	"starship": {
		Exec:    "starship",
		RelSrc:  "starship/starship.toml",
		DstFunc: func(string) string { return path.DotConfigPath("starship.toml") },
		Doc:     "https://starship.rs/config/#configuration",
	},

	"tmux": {
		Exec:   "tmux",
		RelSrc: "tmux",
		DstFunc: func(goos string) string {
			switch goos {
			case "linux", "darwin":
				return path.DotConfigPath("tmux")
			}
			return ""
		},
		Doc: "https://github.com/tmux/tmux/blob/master/tmux.1#L146",
	},

	"wezterm": {
		Exec:    "wezterm",
		RelSrc:  "wezterm",
		DstFunc: func(string) string { return path.DotConfigPath("wezterm") },
		Doc:     "https://wezfurlong.org/wezterm/config/files.html",
	},
}
