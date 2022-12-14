package dotmanager

import "github.com/cqroot/dotm/pkg/common"

type commonDot struct {
	source  string
	target  string
	exec    string
	dotType string
}

func (dc commonDot) Source() string {
	return dc.source
}

func (dc commonDot) Target() string {
	return dc.target
}

func (dc commonDot) Exec() string {
    return dc.exec
}

func (dc commonDot) Type() string {
	return dc.dotType
}

func (dc commonDot) CheckExec() error {
	if dc.exec != "" && !common.CommandExists(dc.exec) {
		return DotIgnoreError
	}

	return nil
}

func checkExec(exec string) error {
	if exec != "" && !common.CommandExists(exec) {
		return DotIgnoreError
	}

	return nil
}
