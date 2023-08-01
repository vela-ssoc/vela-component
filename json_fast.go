package component

import (
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-kit/lua"
)

func newLuaFastJson(L *lua.LState) int {
	val := L.CheckString(1)
	fast := &kind.Fast{}
	e := fast.ParseBytes(lua.S2B(val))
	if e != nil {
		L.RaiseError("%v", e)
		return 0
	}
	L.Push(fast)
	return 1
}
