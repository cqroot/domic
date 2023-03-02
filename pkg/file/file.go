package file

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"
)

func Check(source, target string) (bool, error) {
	_, err := os.Stat(target)
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, errors.New("target file already exists")
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

	return copy.Copy(source, target)
}
