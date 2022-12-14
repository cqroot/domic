package dotmanager

import (
	"errors"
	"os"
)

type SymlinkOneDot struct {
	*commonDot
}

func (d SymlinkOneDot) Check() error {
	return checkDotSymlinkOne(d.exec, d.Source(), d.Target())
}

func checkDotSymlinkOne(exec, source, target string) error {
	if err := checkExec(exec); err != nil {
		return err
	}

	destination, err := os.Readlink(target)
	if err != nil {
		return errors.New("not linked")
	}

	if destination == source {
		return nil
	} else {
		return errors.New("target file already exists")
	}
}

func (d SymlinkOneDot) Apply() error {
	return applyDotSymlinkOne(d.Exec(), d.Source(), d.Target())
}

func applyDotSymlinkOne(exec, source, target string) error {
	err := checkDotSymlinkOne(exec, source, target)

	if err == nil || errors.Is(err, DotIgnoreError) {
		return DotIgnoreError
	}

	return os.Symlink(source, target)
}

func (d SymlinkOneDot) Revoke() error {
	return revokeDotSymlinkOne(d.Exec(), d.Source(), d.Target())
}

func revokeDotSymlinkOne(exec, source, target string) error {
	err := checkDotSymlinkOne(exec, source, target)

	if errors.Is(err, DotIgnoreError) {
		return DotIgnoreError
	} else if err != nil {
		return err
	}

	return os.Remove(target)
}
