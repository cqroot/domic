package dotfile

import (
	"os"
)

func Link(dot *Dot) (bool, error) {
	state, _ := CheckState(dot)

	switch state {
	case StateLinkEmpty:
		err := os.Symlink(dot.Source, dot.Target)
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}

	return false, nil
}

func Unlink(dot *Dot) (bool, error) {
	state, _ := CheckState(dot)

	switch state {
	case StateLinkNormal:
		err := os.Remove(dot.Target)
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
	return false, nil
}
