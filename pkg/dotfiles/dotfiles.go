package dotfiles

import (
	"path/filepath"
	"runtime"

	"github.com/cqroot/doter/pkg/dotfile"
	"github.com/cqroot/doter/pkg/stdpath"
)

var Dotfiles = map[string]dotfile.Dotfile{
	// ************************************************************
	// *  Custom                                                  *
	// ************************************************************
	"bash": {
		Exec:    "bash",
		RelSrc:  "bash",
		DstFunc: func(string) string { return stdpath.DotConfigPath("bash") },
		Doc:     "You can put some bash files in this directory and source them in `.bashrc`.",
	},

	"bin": {
		RelSrc:  "bin",
		DstFunc: func(string) string { return stdpath.HomeDotPath("bin") },
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
				return stdpath.DotConfigPath("alacritty")
			case "windows":
				return stdpath.WindowsAppDataPath("alacritty")
			}
			return ""
		},
		Doc: "https://github.com/alacritty/alacritty#configuration",
	},

	"awesome": {
		Exec:   "awesome",
		RelSrc: "awesome",
		DstFunc: func(goos string) string {
			switch goos {
			case "linux":
				return stdpath.DotConfigPath("awesome")
			}
			return ""
		},
		Doc: "https://awesomewm.org/doc/api/documentation/07-my-first-awesome.md.html",
	},

	"bspwm": {
		Exec:   "bspwm",
		RelSrc: "bspwm",
		DstFunc: func(goos string) string {
			switch goos {
			case "linux":
				return stdpath.DotConfigPath("bspwm")
			}
			return ""
		},
		Doc: "https://github.com/baskerville/bspwm#configuration",
	},

	"dunst": {
		Exec:   "dunst",
		RelSrc: "dunst",
		DstFunc: func(goos string) string {
			switch goos {
			case "linux":
				return stdpath.DotConfigPath("dunst")
			}
			return ""
		},
		Doc: "https://github.com/dunst-project/dunst/wiki#getting-started",
	},

	"fontconfig": {
		RelSrc: "fontconfig",
		DstFunc: func(goos string) string {
			switch goos {
			case "linux":
				return stdpath.DotConfigPath("fontconfig")
			}
			return ""
		},
		Doc: "https://www.freedesktop.org/software/fontconfig/fontconfig-user.html",
	},

	"git": {
		Exec:    "git",
		RelSrc:  "git",
		DstFunc: func(string) string { return stdpath.DotConfigPath("git") },
	},

	"gitui": {
		Exec:   "gitui",
		RelSrc: "gitui",
		DstFunc: func(goos string) string {
			switch goos {
			case "linux", "darwin":
				return stdpath.DotConfigPath("gitui")
			case "windows":
				return stdpath.WindowsAppDataPath("gitui")
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
				return stdpath.DotConfigPath("go")
			case "darwin":
				return stdpath.ApplicationSupportPath("go")
			case "windows":
				return stdpath.WindowsAppDataPath("go")
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
				return stdpath.DotConfigPath("i3")
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
			return stdpath.DotConfigPath("joplin-desktop")
		},
	},

	"lf": {
		Exec:   "lf",
		RelSrc: "lf",
		DstFunc: func(goos string) string {
			switch goos {
			case "linux", "darwin":
				return stdpath.DotConfigPath("lf")
			case "windows":
				return stdpath.WindowsLocalAppDataPath("lf")
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
				return stdpath.DotConfigPath("mpv")
			case "windows":
				return stdpath.WindowsAppDataPath("mpv")
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
				return stdpath.WindowsAppDataPath("nushell")
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
				return stdpath.DotConfigPath("nvim")
			case "windows":
				return stdpath.WindowsLocalAppDataPath("nvim")
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
				return stdpath.DotConfigPath("picom")
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
				return stdpath.DotConfigPath("pip/pip.conf")
			case "darwin":
				return stdpath.DotConfigPath("pip/pip.conf") // `$HOME/Library/Application Support/pip/pip.conf` is not supported.
			case "windows":
				return stdpath.WindowsAppDataPath("pip/pip.ini")
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
				return stdpath.DotConfigPath("polybar")
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
				return stdpath.DotConfigPath("powershell")
			case "windows":
				return filepath.Join(stdpath.HomeDir(), "Documents/PowerShell")
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
				return stdpath.DotConfigPath("rofi")
			}
			return ""
		},
		Doc: "https://github.com/davatorium/rofi#configuration",
	},

	"sqlite": {
		Exec:    "sqlite3",
		RelSrc:  "sqlite/sqliterc",
		DstFunc: func(string) string { return filepath.Join(stdpath.HomeDir(), ".sqliterc") },
	},

	"starship": {
		Exec:    "starship",
		RelSrc:  "starship/starship.toml",
		DstFunc: func(string) string { return stdpath.DotConfigPath("starship.toml") },
		Doc:     "https://starship.rs/config/#configuration",
	},

	"tmux": {
		Exec:   "tmux",
		RelSrc: "tmux",
		DstFunc: func(goos string) string {
			switch goos {
			case "linux", "darwin":
				return stdpath.DotConfigPath("tmux")
			}
			return ""
		},
		Doc: "https://github.com/tmux/tmux/blob/master/tmux.1#L146",
	},

	"wezterm": {
		Exec:    "wezterm",
		RelSrc:  "wezterm",
		DstFunc: func(string) string { return stdpath.DotConfigPath("wezterm") },
		Doc:     "https://wezfurlong.org/wezterm/config/files.html",
	},
}
