package dotmanager

import (
	"github.com/cqroot/dotm/internal/config"
)

type Dot = config.Dot

type DotManager struct {
	dots []Dot
}

func New(basePath string, configPath string, tag string) (*DotManager, error) {
	cfg, err := config.New(basePath, configPath)
	if err != nil {
		return nil, err
	}

	var dots []Dot

	if tag == "" {
		dots = cfg.Dots
	} else {
		dots = make([]Dot, 0)

		for _, dot := range cfg.Dots {
			if containsTag(tag, dot.Tags) {
				dots = append(dots, dot)
				continue
			}
		}
	}

	return &DotManager{
		dots: dots,
	}, nil
}

func containsTag(tag string, tags []string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (d DotManager) Dots() []Dot {
	return d.dots
}
