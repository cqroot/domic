package dotmanager

import (
	"fmt"
	"os"
	"path"
)

type SymlinkEachDot struct {
	*commonDot
}

func (d SymlinkEachDot) Check() error {
	if err := d.CheckExec(); err != nil {
		return err
	}

	files, err := os.ReadDir(d.Source())
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		err := checkDotSymlinkOne(
			d.Exec(),
			path.Join(d.Source(), file.Name()),
			path.Join(d.Target(), file.Name()),
		)
		if err != nil {
			return fmt.Errorf("%s: %w", file.Name(), err)
		}
	}

	return nil
}

func (d SymlinkEachDot) Apply() error {
	files, err := os.ReadDir(d.Source())
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		err := applyDotSymlinkOne(
			d.Exec(),
			path.Join(d.Source(), file.Name()),
			path.Join(d.Target(), file.Name()),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d SymlinkEachDot) Revoke() error {
	files, err := os.ReadDir(d.Source())
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		err := revokeDotSymlinkOne(
			d.Exec(),
			path.Join(d.Source(), file.Name()),
			path.Join(d.Target(), file.Name()),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
