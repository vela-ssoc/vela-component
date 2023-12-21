package component

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
)

func sliceL(co *lua.LState, s lua.Slice, cp lua.P) {
	n := s.Len()
	if n == 0 {
		return
	}
	for i := 0; i < n; i++ {
		if e := co.CallByParam(cp, s[i], lua.LInt(i)); e != nil {
			xEnv.Errorf("for call function fail %v", e)
		}
	}
}

func NewSliceL(L *lua.LState) int {
	v := L.Get(1)

	switch v.Type() {
	case lua.LTTable:
		L.Push(lua.Slice(v.(*lua.LTable).Array()))
		return 1
	case lua.LTNumber:
		n := int(v.(lua.LNumber))
		L.Push(lua.NewSlice(n))
		return 1
	}

	L.Push(lua.Slice{})
	return 1
}

func newLuaSliceIndex(env vela.Environment) {
	export := lua.NewExport("lua.slice.export", lua.WithFunc(NewSliceL))
	env.Set("slice", export)
	env.Set("arr", export)
}
