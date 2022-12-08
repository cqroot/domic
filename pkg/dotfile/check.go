package dotfile

import (
	"os"
	"path"

	"github.com/cqroot/dotm/pkg/common"
)

func CheckState(dot *Dot) (State, string) {
	if dot.Exec != "" && !common.CommandExists(dot.Exec) {
		return StateIgnored, "Ignored"
	}

	switch dot.Type {
	case "symlink_each":
		return checkStateSymlinkEach(dot.Source, dot.Target)

	default:
		return checkStateSymlinkOne(dot.Source, dot.Target)
	}
}

func checkStateSymlinkOne(source, target string) (State, string) {
	destination, err := os.Readlink(target)
	if err != nil {
		return StateLinkEmpty, "Not linked"
	}

	if destination == source {
		return StateLinkNormal, "Linked"
	} else {
		return StateExisted, "Target file already exists"
	}
}

func checkStateSymlinkEach(source, target string) (State, string) {
    files, err := os.ReadDir(source)
	if err != nil {
		return StateError, err.Error()
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		state, descr := checkStateSymlinkOne(
			path.Join(source, file.Name()), path.Join(target, file.Name()),
		)
		if state != StateLinkNormal {
			return state, descr
		}
	}
	return StateLinkNormal, "Linked"
}
