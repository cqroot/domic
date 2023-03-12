package dotfile

import (
	"os/exec"
)

type Dotfile struct {
	// RelSrc string
	// DstFunc func(goos string) string
	Src  string
	Dst  string
	Exec string // Applied only if exec is present. Don't check if empty.
	Doc  string // Documentation for configuration file paths. Usually a link to somewhere on the official website.
}

func (dot Dotfile) IsIgnored() bool {
	if dot.Exec != "" {
		_, err := exec.LookPath(dot.Exec)
		if err != nil {
			return true
		}
	}

	return false
}
