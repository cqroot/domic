package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/pelletier/go-toml/v2"
)

type Dot struct {
	Name   string
	Source string
	Path   string
	Target string
	Type   string
}

type Config struct {
	BasePath   string
	ConfigPath string
	Dots       []Dot
}

func New(basePath string, configPath string) (*Config, error) {
	config := Config{
		BasePath:   basePath,
		ConfigPath: configPath,
	}
	err := config.Read()
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Config) Read() error {
	data, err := ioutil.ReadFile(c.ConfigPath)
	if err != nil {
		return err
	}

	var cfg map[string]Dot
	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}

	for name, dot := range cfg {
		if dot.Path == "" {
			dot.Path = name
		}
		dot.Name = path.Base(dot.Path)
		dot.Source = path.Join(c.BasePath, dot.Path)

		// If user does not specify a target, the target will be generated
		// according to the type.
		if dot.Target == "" {
			if dot.Type == "" {
				dot.Type = "config"
			}
			switch dot.Type {
			case "config":
				dir, err := os.UserConfigDir()
				if err != nil {

				}
				dot.Target = filepath.Join(dir, dot.Name)
			case "home":
				dir, err := os.UserHomeDir()
				if err != nil {

				}
				dot.Target = filepath.Join(dir, fmt.Sprintf(".%s", dot.Name))
			}
		}

		c.Dots = append(c.Dots, dot)
	}

	// Sort slice to ensure that the order of each output result is the same
	sort.Slice(c.Dots, func(i, j int) bool {
		return c.Dots[i].Name < c.Dots[j].Name
	})

	return nil
}
