package symlink

import (
	"os"
)

const (
	StatusOK        int = 0
	StatusError         = 1
	StatusLinked        = 2
	StatusNotLinked     = 3
	StatusExisted       = 4
)

func CheckStatus(source string, target string) (int, string) {
	destination, err := os.Readlink(target)
	if err != nil {
		return StatusNotLinked, "Not linked"
	}

	if destination == source {
		return StatusLinked, "Linked"
	} else {
		return StatusExisted, "Target file already exists"
	}
}

func LinkAndCheckStatus(source string, target string) (int, string) {
	code, descr := CheckStatus(source, target)
	if code == StatusNotLinked {
		err := os.Symlink(source, target)
		if err != nil {
			return StatusError, err.Error()
		} else {
			return StatusOK, "OK"
		}
	}
	return code, descr
}

func UnlinkAndCheckStatus(source string, target string) (int, string) {
	code, descr := CheckStatus(source, target)
	if code == StatusLinked {
		err := os.Remove(target)
		if err != nil {
			return StatusError, err.Error()
		} else {
			return StatusOK, "Removed"
		}
	}
	return code, descr
}
