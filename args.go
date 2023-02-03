package component

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"os"
)

func argsL(L *lua.LState) int {
	idx := L.CheckInt(1)
	if idx <= 0 || idx >= len(os.Args) {
		L.RaiseError("%d overflow os args", idx)
		return 0
	}

	L.Push(lua.S2L(os.Args[idx]))
	return 1
}

func newLuaArgsIndex() lua.LValue {
	return lua.NewFunction(argsL)
}
