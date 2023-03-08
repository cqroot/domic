package dotfile

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cqroot/domic/pkg/stdpath"
)

type Dotfile struct {
	RelSrc  string
	DstFunc func(goos string) string
	Exec    string // Applied only if exec is present. Don't check if empty.
	Doc     string // Documentation for configuration file paths. Usually a link to somewhere on the official website.
}

func (dot Dotfile) Src() string {
	src := filepath.Join(stdpath.DotsDir(), dot.RelSrc)
	return strings.ReplaceAll(src, "\\", "/")
}

func (dot Dotfile) Dst() string {
	dst := dot.DstFunc(runtime.GOOS)
	return strings.ReplaceAll(dst, "\\", "/")
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
