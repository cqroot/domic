package dotmanager

import (
	"github.com/cqroot/dotm/internal/config"
)

type DotManager struct {
	Dots []Dot
}

func New(basePath string, configPath string, tag string) (*DotManager, error) {
	cfg, err := config.New(basePath, configPath)
	if err != nil {
		return nil, err
	}

    // var dots []Dot
    dots := make([]Dot, 0)

	for _, item := range cfg.DotItems {
		if tag != "" && !containsTag(tag, item.Tags) {
			continue
		}

        dot := GetDot(item)
        if dot == nil {
            continue
        }

		dots = append(dots, dot)
	}

	return &DotManager{
		Dots: dots,
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
