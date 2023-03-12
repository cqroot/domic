package dotfile

import (
	"github.com/cqroot/domic/pkg/symlink"
)

type State int

const (
	StateApplied State = iota
	StateUnapplied
	StateIgnored
	StateUnsupported
	StateTargetAlreadyExists
)

func (dot Dotfile) State() State {
	if dot.IsIgnored() {
		return StateIgnored
	}

	ok, err := symlink.Check(dot.Src, dot.Dst)
	if err != nil {
		return StateTargetAlreadyExists
	}
	if ok {
		return StateApplied
	}
	return StateUnapplied
}

func (dot Dotfile) Apply() error {
	return symlink.Apply(dot.Src, dot.Dst)
}

func (dot Dotfile) Revoke() error {
	return symlink.Revoke(dot.Src, dot.Dst)
}
