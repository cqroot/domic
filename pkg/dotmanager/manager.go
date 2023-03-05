package dotmanager

import "runtime"

type Dot struct {
	Src  string
	Dest string
	Exec string // Applied only if exec is present. Don't check if empty.
	Doc  string // Documentation for configuration file paths. Usually a link to somewhere on the official website.
}

type DotManager struct {
	DotMap map[string]Dot
}

func New() *DotManager {
	return &DotManager{
		DotMap: make(map[string]Dot),
	}
}

func Default() *DotManager {
	dm := New()
	dm.DotMap = DefaultDotMap(runtime.GOOS)
	return dm
}
