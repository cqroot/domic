package dotmanager

import (
	"github.com/cqroot/dotm/internal/config"
)

type Dot = config.Dot

type DotManager struct {
	dots []Dot
}

func New(basePath string, configPath string) (*DotManager, error) {
	cfg, err := config.New(basePath, configPath)
	if err != nil {
		return nil, err
	}

	return &DotManager{
		dots: cfg.Dots,
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

func (d *DotManager) DotsWithTag(tag string) []Dot {
	if tag == "" {
		return d.dots
	}

	result := make([]Dot, 0)

	for _, dot := range d.dots {
		if containsTag(tag, dot.Tags) {
			result = append(result, dot)
			continue
		}
	}

	return result
}
