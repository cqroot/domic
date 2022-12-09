package dotmanager

import "github.com/cqroot/dotm/pkg/dotfile"

type Level int

const (
	Info Level = iota
	Warn
	Error
	Ignored
)

type CheckResult struct {
	Dot         Dot
	Level       Level
	Description string
}

func (d DotManager) Check() []CheckResult {
	result := make([]CheckResult, 0)

	for _, dot := range d.dots {
		state, descr := dotfile.CheckState(&dot)
		var level Level
		switch state {
		case dotfile.StateIgnored:
			level = Ignored
		case dotfile.StateLinkNormal:
			level = Info
		default:
			level = Error
		}

		result = append(result, CheckResult{
			Dot:         dot,
			Level:       level,
			Description: descr,
		})
	}

	return result
}
