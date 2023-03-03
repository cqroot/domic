package dotmanager

type Dot struct {
	Src  string
	Dest string
	Exec string // Applied only if exec is present. Don't check if empty.
}

type DotManager struct {
	dotMap map[string]Dot
}

func New() *DotManager {
	return &DotManager{
		dotMap: make(map[string]Dot),
	}
}

func Default() *DotManager {
	dm := New()
	dm.dotMap = dm.defaultDotfileMap()
	return dm
}
