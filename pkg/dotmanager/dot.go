package dotmanager

import (
	"github.com/cqroot/dotm/internal/config"
)

type Dot interface {
	Source() string
	Target() string
	Type() string

	CheckExec() error
	Check() error
	Apply() error
	Revoke() error
}

func GetDot(dotItem config.DotItem) Dot {
	cd := &commonDot{
		source:  dotItem.Source,
		target:  dotItem.Target,
		exec:    dotItem.Exec,
		dotType: dotItem.Type,
	}

	switch dotItem.Type {
	case "symlink_one":
		return SymlinkOneDot{
			commonDot: cd,
		}
	case "symlink_each":
		return SymlinkEachDot{
			commonDot: cd,
		}
	}

	return nil
}
