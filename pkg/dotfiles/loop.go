package dotfiles

import (
	"sort"

	"github.com/cqroot/domic/pkg/dotfile"
)

func ForEach(handleFunc func(name string, dot dotfile.Dotfile)) error {
	dots, err := Dotfiles()
	if err != nil {
		return err
	}

	names := make([]string, 0, len(dots))
	for name := range dots {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		dot := dots[name]

		handleFunc(name, dot)
	}
	return nil
}
