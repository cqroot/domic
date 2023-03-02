package configs

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/cqroot/gmdots/internal/dot"
)

func RangeDotConfigs(handleFunc func(dotName string, dotConfig dot.DotConfig)) error {
	dots, err := dot.Dots()
	if err != nil {
		return err
	}

	defDotConfigs, err := DefaultDotConfigs()
	if err != nil {
		return err
	}

	for _, dotName := range dots {
		dotConfig, ok := defDotConfigs[dotName]
		if !ok {
			continue
		}

		if dotConfig.Dest == "-" {
			continue
		}

		handleFunc(dotName, dotConfig)
	}

	return nil
}

func DefaultDotConfigs() (map[string]dot.DotConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// homePath := func(path string) string {
	// 	return filepath.Join(homeDir, path)
	// }

	xdgConfigPath := func(path string) string {
		return filepath.Join(homeDir, ".config", path)
	}

	appDataPath := func(path string) string {
		return filepath.Join(os.Getenv("APPDATA"), path)
	}

	localAppDataPath := func(path string) string {
		return filepath.Join(os.Getenv("LOCALAPPDATA"), path)
	}

	return map[string]dot.DotConfig{
		// https://neovim.io/doc/user/starting.html#standard-path
		"nvim": {
			Src: "nvim",
			Dest: func() string {
				if runtime.GOOS == "linux" {
					return xdgConfigPath("nvim")
				} else if runtime.GOOS == "windows" {
					return localAppDataPath("nvim")
				}
				return "-"
			}(),
		},

		// https://pip.pypa.io/en/stable/topics/configuration/#location
		"pip": {
			Src: "pip/pip.conf",
			Dest: func() string {
				if runtime.GOOS == "linux" {
					return xdgConfigPath("pip/pip.conf")
				} else if runtime.GOOS == "windows" {
					return appDataPath("pip/pip.ini")
				}
				return "-"
			}(),
		},
	}, nil
}
