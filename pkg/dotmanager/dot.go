package dotmanager

import (
	"github.com/cqroot/dotm/internal/config"
	"github.com/cqroot/dotm/pkg/common"
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

type dotCommon struct {
	exec    string
	dotType string
}

func (dc dotCommon) Type() string {
	return dc.dotType
}

func (dc dotCommon) CheckExec() error {
	if dc.exec != "" && !common.CommandExists(dc.exec) {
		return DotIgnoreError
	}

	return nil
}

func GetDot(dotItem config.DotItem) Dot {
	dc := &dotCommon{
		exec:    dotItem.Exec,
		dotType: dotItem.Type,
	}

	switch dotItem.Type {
	case "symlink_one":
		return SymlinkOneDot{
			source: dotItem.Source,
			target: dotItem.Target,

			dotCommon: dc,
		}
	case "symlink_each":
		return nil
	}

	return nil
}
