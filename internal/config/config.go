package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
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

// Prefs holds optional top-level user preferences loaded from
// domic.config.toml in the config directory.
type Prefs struct {
	Diff []string `toml:"diff"`
}

type appConfig struct {
	Path     string    `toml:"path"`
	Target   AppTarget `toml:"target"`
	Template bool      `toml:"template"`
	Bin      string    `toml:"bin"`
}

type App struct {
	Name     string
	Path     string
	Target   AppTarget
	Template bool
	Bin      string
}

var (
	loaded     bool
	apps       []App
	configPath string
	sourceBase string
	prefs      Prefs
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
			Name:     name,
			Path:     a.Path,
			Target:   a.Target,
			Template: a.Template,
			Bin:      a.Bin,
		})
	}
	sort.Slice(apps, func(i, j int) bool {
		return apps[i].Name < apps[j].Name
	})

	if err := loadPrefs(); err != nil {
		return err
	}

	loaded = true
	return nil
}

func loadPrefs() error {
	path := filepath.Join(xdg.ConfigHome, "domic", "domic.config.toml")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read %s: %w", path, err)
	}
	if _, err := toml.Decode(string(data), &prefs); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}

// DiffCommand returns the configured external diff command and its arguments.
// The first element is the program name, the rest are arguments. Defaults to
// `["diff", "-u"]`.
func DiffCommand() []string {
	if len(prefs.Diff) > 0 {
		return prefs.Diff
	}
	return []string{"diff", "-u"}
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

func ConfigDir() string {
	return filepath.Join(xdg.ConfigHome, "domic")
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

func (a App) BinAvailable() bool {
	if a.Bin == "" {
		return true
	}
	_, err := exec.LookPath(a.Bin)
	return err == nil
}
