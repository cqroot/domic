package domiclua

import (
	"path/filepath"
	"runtime"

	"github.com/cqroot/domic/pkg/stdpath"
	"github.com/yuin/gopher-lua"
)

func Loader(L *lua.LState) int {
	// register functions to the table
	mod := L.SetFuncs(L.NewTable(), exports)

	// returns the module
	L.Push(mod)
	return 1
}

var exports = map[string]lua.LGFunction{
	"joinpath": JoinPath,
	"goos":     GOOS,
	"homedir":  HomeDir,
}

func JoinPath(L *lua.LState) int {
	lv := L.ToTable(1) // get argument
	elems := make([]string, 0, L.ObjLen(lv))
	lv.ForEach(func(_ lua.LValue, elem lua.LValue) {
		if str, ok := elem.(lua.LString); ok {
			elems = append(elems, string(str))
		}
	})

	L.Push(lua.LString(filepath.Join(elems...)))
	return 1 // number of results
}

func GOOS(L *lua.LState) int {
	L.Push(lua.LString(runtime.GOOS))
	return 1 // number of results
}

func HomeDir(L *lua.LState) int {
	L.Push(lua.LString(stdpath.HomeDir()))
	return 1 // number of results
}
