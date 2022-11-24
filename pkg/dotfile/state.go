package dotfile

type State int

const (
	StateIgnored State = iota
	StateExisted

	StateLinkNormal
	StateLinkEmpty
)
