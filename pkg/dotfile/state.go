package dotfile

import "github.com/cqroot/dotm/internal/config"

type Dot = config.Dot

type State int

const (
	StateIgnored State = iota
	StateExisted
    StateError

	StateLinkNormal
	StateLinkEmpty
)
