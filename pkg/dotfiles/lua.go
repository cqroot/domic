package dotfiles

import (
	"errors"
	"path/filepath"

	"github.com/yuin/gopher-lua"

	"github.com/cqroot/domic/internal/domiclua"
	"github.com/cqroot/domic/pkg/dotfile"
	"github.com/cqroot/domic/pkg/stdpath"
)

func Dotfiles() (map[string]dotfile.Dotfile, error) {
	return DotfilesFromLua(filepath.Join(stdpath.BaseDir(), "domic.lua"))
}

func DotfilesFromLua(luapath string) (map[string]dotfile.Dotfile, error) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("domic", domiclua.Loader)

	if err := L.DoFile(luapath); err != nil {
		panic(err)
	}

	ret := L.Get(-1)
	table, ok := ret.(*lua.LTable)
	if !ok {
		return nil, errors.New("parse error")
	}
	dots := dotfiles(filepath.Dir(luapath), table)
	return dots, nil
}

func dotfiles(baseDir string, table *lua.LTable) map[string]dotfile.Dotfile {
	dots := make(map[string]dotfile.Dotfile)

	table.ForEach(func(_ lua.LValue, item lua.LValue) {
		if dotTable, ok := item.(*lua.LTable); ok {
			nameLV := dotTable.RawGet(lua.LString("name"))

			name, ok := mustString(nameLV)
			if !ok {
				return
			}

			dot := dotfile.Dotfile{}

			dotTable.ForEach(func(keyLV lua.LValue, valueLV lua.LValue) {
				key, ok := mustString(keyLV)
				if !ok {
					return
				}

				value, ok := mustString(valueLV)
				if !ok {
					return
				}

				switch key {
				case "src":
					dot.Src = filepath.Join(baseDir, value)
				case "dst":
					dot.Dst = value
				case "exec":
					dot.Exec = value
				}
			})

			if dot.Src == "" {
				dot.Src = filepath.Join(baseDir, name)
			}
			if dot.Dst == "" {
				return
			}

			dots[name] = dot
		}
	})

	return dots
}

func mustString(value lua.LValue) (string, bool) {
	valueLS, ok := value.(lua.LString)
	if !ok {
		return "", false
	}
	return valueLS.String(), true
}
