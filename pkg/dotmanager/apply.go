package dotmanager

import (
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/cqroot/gmdots/pkg/file"
	"github.com/cqroot/gmdots/pkg/symlink"
)

func (m DotManager) AbsSrcPath(name string, dot Dot) string {
	src := dot.Src
	if src == "" {
		src = name
	}
	return filepath.Join(DotsDir(), src)
}

func (m DotManager) CheckSkip(name string) bool {
	dot, ok := m.dotMap[name]
	if !ok {
		return true
	}

	if dot.Dest == "" {
		return true
	}

	if dot.Exec != "" {
		_, err := exec.LookPath(dot.Exec)
		if err != nil {
			return true
		}
	}

	return false
}

// returns `false, nil` if skipped.
// returns `true,  nil` if applied normally.
// returns `false, err` if an error occurs.
func (m DotManager) Check(name string) (bool, error) {
	if m.CheckSkip(name) {
		return false, nil
	}

	dot := m.dotMap[name]
	src := m.AbsSrcPath(name, dot)
	dest := dot.Dest

	if runtime.GOOS == "windows" {
		return file.Check(src, dest)
	}
	return symlink.Check(src, dest)
}

func (m DotManager) Apply(name string) error {
	if m.CheckSkip(name) {
		return nil
	}

	dot := m.dotMap[name]
	src := m.AbsSrcPath(name, dot)
	dest := dot.Dest

	if runtime.GOOS == "windows" {
		return file.Apply(src, dest)
	}
	return symlink.Apply(src, dest)
}
