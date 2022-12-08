package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/cqroot/dotm/pkg/common"
	"github.com/pelletier/go-toml/v2"
)

type Dot struct {
	Name   string // The name of the file or folder itself
	Source string // Absolute path to source file
	Target string // Absolute path to target file
	Type   string // symlink_one (default) / symlink_each

	// Relative path to source file
	RelativePath string `toml:"relative_path"`

	// config: ~/.config/xxx or home: ~/.xxx
	TargetType         string `toml:"target_type"`
	TargetRelativePath string `toml:"target_relative_path"`

	Exec string
	Tags []string
}

type Config struct {
	BasePath   string
	ConfigPath string
	Dots       []Dot
}

func New(basePath string, configPath string) (*Config, error) {
	basePath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}
	configPath, err = filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	config := Config{
		BasePath:   basePath,
		ConfigPath: configPath,
	}
	err = config.Read()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) Read() error {
	data, err := os.ReadFile(c.ConfigPath)
	if err != nil {
		return err
	}

	var cfg map[string]Dot
	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}

	for name, dot := range cfg {
		if dot.RelativePath == "" {
			dot.RelativePath = name
		}
		dot.Name = path.Base(dot.RelativePath)
		dot.Source = path.Join(c.BasePath, dot.RelativePath)

		if dot.Type != "symlink_each" {
			dot.Type = "symlink_one"
		}

		// If user does not specify a target, the target will be generated
		// according to the type.
		if dot.Target == "" {
			if dot.TargetType == "" {
				dot.TargetType = "config"
			}
			dir, err := common.DotDir(dot.TargetType)
			if err != nil {
				fmt.Println(err)
				continue
			}
			switch dot.TargetType {
			case "config":
				if dot.TargetRelativePath != "" {
					dot.Target = filepath.Join(dir, dot.TargetRelativePath)
				} else {
					dot.Target = filepath.Join(dir, dot.Name)
				}
			case "home":
				if dot.TargetRelativePath != "" {
					dot.Target = filepath.Join(dir, dot.TargetRelativePath)
				} else {
					dot.Target = filepath.Join(dir, fmt.Sprintf(".%s", dot.Name))
				}
			}
		}

		c.Dots = append(c.Dots, dot)
	}

	// Sort slice to ensure that the order of each output result is the same
	sort.Slice(c.Dots, func(i, j int) bool {
		if c.Dots[i].TargetType != c.Dots[j].TargetType {
			return c.Dots[i].TargetType < c.Dots[j].TargetType
		}
		return c.Dots[i].Name < c.Dots[j].Name
	})

	return nil
}
