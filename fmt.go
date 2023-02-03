package component

import (
	"github.com/tidwall/pretty"
	"github.com/vela-ssoc/vela-kit/vela"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
)

func newLuaFmtSprintf(L *lua.LState) int {
	L.Push(lua.S2L(auxlib.Format(L, 0)))
	return 1
}

func prettyJsonL(L *lua.LState) int {
	chunk := lua.S2B(L.Get(1).String())
	L.Push(lua.B2L(pretty.PrettyOptions(chunk, nil)))
	return 1
}

func newLuaFmtIndex(env vela.Environment) {
	env.Set("format", lua.NewFunction(newLuaFmtSprintf))
	env.Set("pretty", lua.NewFunction(prettyJsonL))
}
