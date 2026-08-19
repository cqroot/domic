package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
)

type AppTarget struct {
	Linux   string `toml:"linux"`
	Windows string `toml:"windows"`
	Darwin  string `toml:"darwin"`
}

type appConfig struct {
	Path   string    `toml:"path"`
	Target AppTarget `toml:"target"`
}

type App struct {
	Name   string
	Path   string
	Target AppTarget
}

var (
	loaded     bool
	apps       []App
	configPath string
	sourceBase string
)

func Load(file string) error {
	if file != "" {
		configPath = file
	} else {
		configPath = filepath.Join(xdg.ConfigHome, "domic", "domic.toml")
	}

	sourceBase = filepath.Dir(configPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("config file %s does not exist, please create it first", configPath)
		}
		return fmt.Errorf("read config %s: %w", configPath, err)
	}

	raw := map[string]appConfig{}
	if _, err := toml.Decode(string(data), &raw); err != nil {
		return fmt.Errorf("parse config %s: %w", configPath, err)
	}

	apps = make([]App, 0, len(raw))
	for name, a := range raw {
		apps = append(apps, App{
			Name:   name,
			Path:   a.Path,
			Target: a.Target,
		})
	}
	sort.Slice(apps, func(i, j int) bool {
		return apps[i].Name < apps[j].Name
	})

	loaded = true
	return nil
}

func Apps() []App {
	if !loaded {
		return nil
	}
	return apps
}

func ConfigPath() string {
	return configPath
}

func SourceBase() string {
	return sourceBase
}

func ExpandHome(p string) (string, error) {
	if p == "" || p[0] != '~' {
		return p, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home: %w", err)
	}
	if p == "~" {
		return home, nil
	}
	if len(p) > 1 && (p[1] == '/' || p[1] == filepath.Separator) {
		return filepath.Join(home, p[2:]), nil
	}
	return p, nil
}

func TargetForOS(app App) (string, bool) {
	var raw string
	switch runtime.GOOS {
	case "linux":
		raw = app.Target.Linux
	case "windows":
		raw = app.Target.Windows
	case "darwin":
		raw = app.Target.Darwin
	default:
		return "", false
	}
	if raw == "" {
		return "", false
	}
	return raw, true
}
