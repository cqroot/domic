package dotfile

import (
	"os"

	"github.com/cqroot/dotm/pkg/common"
)

func CheckState(dot *Dot) (State, string) {
	if dot.Exec != "" && !common.CommandExists(dot.Exec) {
		return StateIgnored, "Ignored"
	}

	destination, err := os.Readlink(dot.Target)
	if err != nil {
		return StateLinkEmpty, "Not linked"
	}

	if destination == dot.Source {
		return StateLinkNormal, "Linked"
	} else {
		return StateExisted, "Target file already exists"
	}
}
