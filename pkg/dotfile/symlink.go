package dotfile

import (
	"os"
	"path"
)

func Link(dot *Dot) (bool, error) {
	state, _ := CheckState(dot)

	switch state {
	case StateLinkEmpty:
		var err error

		switch dot.Type {
		case "symlink_each":
			err = LinkEach(dot.Source, dot.Target)
		default:
			err = LinkOne(dot.Source, dot.Target)
		}

		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}

	return false, nil
}

func LinkOne(source, target string) error {
	return os.Symlink(source, target)
}

func LinkEach(source, target string) error {
	files, err := os.ReadDir(source)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		err := LinkOne(
			path.Join(source, file.Name()), path.Join(target, file.Name()),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func Unlink(dot *Dot) (bool, error) {
	state, _ := CheckState(dot)

	switch state {
	case StateLinkNormal:
		var err error

		switch dot.Type {
		case "symlink_each":
			err = UnlinkEach(dot.Source, dot.Target)
		default:
			err = UnlinkOne(dot.Source, dot.Target)
		}

		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
	return false, nil
}

func UnlinkOne(source, target string) error {
	// TODO
	return os.Remove(target)
}

func UnlinkEach(source, target string) error {
	files, err := os.ReadDir(source)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		err := UnlinkOne(
			path.Join(source, file.Name()), path.Join(target, file.Name()),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
