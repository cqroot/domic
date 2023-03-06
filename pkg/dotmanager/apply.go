package dotmanager

import (
	"errors"
	"os/exec"
	"path/filepath"

	"github.com/cqroot/doter/pkg/path"
	"github.com/cqroot/doter/pkg/symlink"
)

func (m DotManager) AbsSrcPath(name string, dot Dot) string {
	return filepath.Join(path.DotsDir(), dot.Src)
}

func (m DotManager) CheckSkip(name string) error {
	dot, ok := m.DotMap[name]
	if !ok {
		return errors.New("Skip: dot config not found")
	}

	if dot.Dest == "" {
		return errors.New("Skip: empty dest")
	}

	if dot.Exec != "" {
		_, err := exec.LookPath(dot.Exec)
		if err != nil {
			return errors.New("Skip: exec not found")
		}
	}

	return nil
}

// returns `false, nil` if the target does not exist.
// returns `true,  nil` if applied normally.
// returns `false, err` if an error occurs.
func (m DotManager) Check(name string) (bool, error) {
	if err := m.CheckSkip(name); err != nil {
		return false, err
	}

	dot := m.DotMap[name]
	src := m.AbsSrcPath(name, dot)
	dest := dot.Dest

	return symlink.Check(src, dest)
}

func (m DotManager) Apply(name string) error {
	if err := m.CheckSkip(name); err != nil {
		return nil
	}

	dot := m.DotMap[name]
	src := m.AbsSrcPath(name, dot)
	dest := dot.Dest

	return symlink.Apply(src, dest)
}
