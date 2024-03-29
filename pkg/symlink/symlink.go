package symlink

import (
	"errors"
	"os"
	"path/filepath"
)

func Check(source, target string) (bool, error) {
	_, err := os.Lstat(target)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	destination, err := os.Readlink(target)
	if err != nil {
		return false, errors.New("destination dotfile already exists")
	}

	if filepath.FromSlash(destination) == filepath.FromSlash(source) {
		return true, nil
	} else {
		return false, errors.New("destination dotfile already exists")
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

func Revoke(source, target string) error {
	ok, err := Check(source, target)
	if err != nil {
		return nil
	}
	if !ok {
		return nil
	}
	return os.Remove(target)
}
