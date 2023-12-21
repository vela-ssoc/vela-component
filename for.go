package component

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
)

func mapL(co *lua.LState, m *lua.Map, cp lua.P) {
	keys := m.Keys()
	n := len(keys)
	if n == 0 {
		return
	}
	for i := 0; i < n; i++ {
		k := keys[i]
		v, ok := m.Get(k)
		if !ok {
			continue
		}

		if e := co.CallByParam(cp, lua.S2L(k), v); e != nil {
			xEnv.Errorf("for lua map call function fail %v", e)
		}

	}

}

func forL(L *lua.LState) int {
	lv := L.Get(1)
	fn := L.IsFunc(1)
	if fn == nil {
		return 0
	}

	cp := xEnv.P(fn)
	co := xEnv.Clone(L)
	defer xEnv.Free(co)

	switch lv.Type() {
	case lua.LTSlice:
		sliceL(co, lv.(lua.Slice), cp)

	case lua.LTMap:
		mapL(co, lv.(*lua.Map), cp)
	default:

	}

	return 0
}

/*

	local arr = vela.slice{}
	vela.for(arr , _(i , v)
		print(i , v)
	end)

	local map = vela.map{}
	vela.for(map , _(i , v)
		print(i , v)
	end)


*/

func newLuaForIndex(env vela.Environment) {
	env.Set("for", lua.NewFunction(forL))
}
