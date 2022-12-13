package dotmanager

import (
	"errors"
	"os"
)

type SymlinkOneDot struct {
	source string
	target string

	*dotCommon
}

func (d SymlinkOneDot) Source() string {
	return d.source
}

func (d SymlinkOneDot) Target() string {
	return d.target
}

func (d SymlinkOneDot) Check() error {
	if err := d.CheckExec(); err != nil {
		return err
	}

	destination, err := os.Readlink(d.target)
	if err != nil {
		return errors.New("not linked")
	}

	if destination == d.source {
		return nil
	} else {
		return errors.New("target file already exists")
	}
}

func (d SymlinkOneDot) Apply() error {
	err := d.Check()

	if err == nil || errors.Is(err, DotIgnoreError) {
		return DotIgnoreError
	}

	return os.Symlink(d.source, d.target)
}

func (d SymlinkOneDot) Revoke() error {
	err := d.Check()

	if errors.Is(err, DotIgnoreError) {
		return DotIgnoreError
	}

	return os.Remove(d.target)
}
