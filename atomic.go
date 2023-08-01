package component

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"sync/atomic"
)

type av struct{ v int64 }

var atomicMetaTab *lua.LTable

func initAtomicMetaTab() {
	tab := lua.CreateTable(0, 4)
	tab.RawSetString("__add", lua.NewFunction(avMetaAdd))
	tab.RawSetString("__sub", lua.NewFunction(avMetaSub))
	tab.RawSetString("__index", lua.NewFunction(avMetaIndex))
	tab.RawSetString("__newindex", lua.NewFunction(avMetaNewIndex))
	atomicMetaTab = tab
}

func avMetaAdd(L *lua.LState) int {
	a := checkAv(L, 1)
	n := L.CheckInt64(2)
	atomic.AddInt64(&a.v, n)
	return 0
}

func avMetaSub(L *lua.LState) int {
	a := checkAv(L, 1)
	n := L.CheckInt64(2)
	atomic.AddInt64(&a.v, -1*n)
	return 0
}

func avMetaIndex(L *lua.LState) int {
	a := checkAv(L, 1)
	key := L.CheckString(2)
	switch key {
	case "value":
		L.Push(lua.LNumber(atomic.LoadInt64(&a.v)))
	default:
		L.Push(lua.LNil)
	}
	return 1
}

func avMetaNewIndex(L *lua.LState) int {
	a := checkAv(L, 1)
	key := L.CheckString(2)
	val := L.CheckNumber(3)
	switch key {
	case "value":
		atomic.StoreInt64(&a.v, int64(val))
	default:
	}
	return 0
}

func checkAv(L *lua.LState, idx int) *av {
	ud := L.CheckUserData(idx)

	a, ok := ud.Value.(*av)
	if ok {
		return a
	}

	L.RaiseError("invalid atomic userdata")
	return nil
}

func newLuaAtomicInt(L *lua.LState) int {
	ud := L.NewUserData()
	ud.Value = L.CheckInt64(1)
	ud.Metatable = atomicMetaTab
	L.Push(ud)
	return 1
}

func newLuaAtomicIndex() lua.LValue {
	initAtomicMetaTab()
	return lua.NewFunction(newLuaAtomicInt)
}
