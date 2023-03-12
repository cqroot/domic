package dotfiles

import (
	"github.com/cqroot/domic/pkg/dotfile"
)

func ForEach(handleFunc func(name string, dot dotfile.Dotfile)) error {
	dots, err := Dotfiles()
	if err != nil {
		return err
	}

	for name, dot := range dots {
		handleFunc(name, dot)
	}
	return nil
}
