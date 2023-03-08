package dotfiles

import "github.com/cqroot/doter/pkg/dotfile"

func ForEach(names []string, handleFunc func(name string, dot dotfile.Dotfile)) {
	for _, name := range names {
		dot, ok := Dotfiles[name]
		if !ok {
			continue
		}
		if dot.Dst() == "" {
			continue
		}

		handleFunc(name, dot)
	}
}
