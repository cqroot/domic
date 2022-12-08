package dotmanager

import (
	"github.com/cqroot/dotm/internal/config"
)

type Dot = config.Dot

type DotManager struct {
	dots []Dot
	tag  string
}

func New(basePath string, configPath string, tag string) (*DotManager, error) {
	cfg, err := config.New(basePath, configPath)
	if err != nil {
		return nil, err
	}

	return &DotManager{
		dots: cfg.Dots,
		tag:  tag,
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

func (d *DotManager) Dots() []Dot {
	if d.tag == "" {
		return d.dots
	}

	result := make([]Dot, 0)

	for _, dot := range d.dots {
		if containsTag(d.tag, dot.Tags) {
			result = append(result, dot)
			continue
		}
	}

	return result
}
