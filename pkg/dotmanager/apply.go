package dotmanager

import (
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

func (m DotManager) Check(name string) (bool, error) {
	dot, ok := m.dotMap[name]
	if !ok {
		return false, nil
	}

	src := m.AbsSrcPath(name, dot)
	dest := dot.Dest

	if runtime.GOOS == "windows" {
		return file.Check(src, dest)
	}
	return symlink.Check(src, dest)
}

func (m DotManager) Apply(name string) error {
	dot, ok := m.dotMap[name]
	if !ok {
		return nil
	}

	src := m.AbsSrcPath(name, dot)
	dest := dot.Dest

	if runtime.GOOS == "windows" {
		return file.Apply(src, dest)
	}
	return symlink.Apply(src, dest)
}
