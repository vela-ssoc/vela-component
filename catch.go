package component

import "github.com/vela-ssoc/vela-kit/lua"

func raise(L *lua.LState) int {
	narg := L.GetTop()
	if narg == 0 {
		return 0
	}

	for i := 1; i <= narg; i++ {
		val := L.Get(1)
		if val.Type() != lua.LTNil {
			L.RaiseError("%v", val)
			return 0
		}
	}
	return 0
}

func newLuaCatchIndex() lua.LValue {
	return lua.NewFunction(raise)
}
