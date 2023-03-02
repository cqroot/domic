package symlink

import (
	"errors"
	"os"
	"path/filepath"
)

func Check(source, target string) (bool, error) {
	destination, err := os.Readlink(target)
	if err != nil {
		return false, nil
	}

	if destination == source {
		return true, nil
	} else {
		return false, errors.New("target file already exists")
	}
}

func Apply(source, target string) error {
	ok, err := Check(source, target)

	if ok {
		return nil
	}

	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(target), os.ModePerm)
	if err != nil {
		return err
	}

	return os.Symlink(source, target)
}
