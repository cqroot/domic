package dotmanager

import (
	"os"
)

type Dot struct {
	Src  string
	Dest string
	Exec string // Applied only if exec is present. Don't check if empty.
}

type DotManager struct {
	homeDir string
	dotMap  map[string]Dot
}

func New() (*DotManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return &DotManager{
		homeDir: homeDir,
		dotMap:  make(map[string]Dot),
	}, nil
}

func Default() (*DotManager, error) {
	dm, err := New()
	if err != nil {
		return nil, err
	}

	dm.dotMap = dm.defaultDotfileMap()
	return dm, nil
}
