package domiclua

import (
	"os"
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
	"join_path":                   JoinPath,
	"goos":                        GOOS,
	"home_dir":                    HomeDir,
	"home_path":                   HomePath,
	"dot_config_dir":              DotConfigDir,
	"dot_config_path":             DotConfigPath,
	"windows_app_data_dir":        WindowsAppDataDir,
	"windows_app_data_path":       WindowsAppDataPath,
	"windows_local_app_data_dir":  WindowsLocalAppDataDir,
	"windows_local_app_data_path": WindowsLocalAppDataPath,
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

func HomePath(L *lua.LState) int {
	lv := L.ToString(1)

	L.Push(lua.LString(filepath.Join(stdpath.HomeDir(), lv)))
	return 1 // number of results
}

func DotConfigDir(L *lua.LState) int {
	L.Push(lua.LString(filepath.Join(stdpath.HomeDir(), ".config")))
	return 1 // number of results
}

func DotConfigPath(L *lua.LState) int {
	lv := L.ToString(1)

	L.Push(lua.LString(filepath.Join(stdpath.HomeDir(), ".config", lv)))
	return 1 // number of results
}

func WindowsAppDataDir(L *lua.LState) int {
	L.Push(lua.LString(os.Getenv("APPDATA")))
	return 1 // number of results
}

func WindowsAppDataPath(L *lua.LState) int {
	lv := L.ToString(1)

	L.Push(lua.LString(filepath.Join(os.Getenv("APPDATA"), lv)))
	return 1 // number of results
}

func WindowsLocalAppDataDir(L *lua.LState) int {
	L.Push(lua.LString(os.Getenv("LOCALAPPDATA")))
	return 1 // number of results
}

func WindowsLocalAppDataPath(L *lua.LState) int {
	lv := L.ToString(1)

	L.Push(lua.LString(filepath.Join(os.Getenv("LOCALAPPDATA"), lv)))
	return 1 // number of results
}
