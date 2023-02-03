package component

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"sync"
)

func newLuaMutex(L *lua.LState) int {
	mu := &sync.Mutex{}
	ud := L.NewAnyData(mu)
	ud.Meta("lock", L.NewFunction(func(_ *lua.LState) int {
		mu.Lock()
		return 0
	}))
	ud.Meta("unlock", L.NewFunction(func(_ *lua.LState) int {
		mu.Unlock()
		return 0
	}))

	L.Push(ud)
	return 1
}

func newLuaRwMutex(L *lua.LState) int {
	rw := &sync.RWMutex{}
	ud := L.NewAnyData(rw)

	ud.Meta("lock", L.NewFunction(func(_ *lua.LState) int {
		rw.Lock()
		return 0
	}))

	ud.Meta("unlock", L.NewFunction(func(_ *lua.LState) int {
		rw.Unlock()
		return 0
	}))

	ud.Meta("read", L.NewFunction(func(_ *lua.LState) int {
		rw.RLock()
		return 0
	}))

	ud.Meta("unread", L.NewFunction(func(_ *lua.LState) int {
		rw.RUnlock()
		return 0
	}))

	L.Push(ud)
	return 1
}

func newLuaWaitGroup(L *lua.LState) int {
	wg := &sync.WaitGroup{}
	lv := L.Get(1)
	if lv.Type() != lua.LTNil {
		if n, ok := lv.AssertFloat64(); ok {
			wg.Add(int(n))
		}
	}

	ud := L.NewAnyData(wg)
	ud.Meta("wait", L.NewFunction(func(_ *lua.LState) int {
		wg.Wait()
		return 0
	}))

	ud.Meta("done", L.NewFunction(func(_ *lua.LState) int {
		wg.Done()
		return 0
	}))

	ud.Meta("add", L.NewFunction(func(co *lua.LState) int {
		n := co.CheckInt(1)
		wg.Add(n)
		return 0
	}))

	L.Push(ud)
	return 1
}

func newLuaSyncIndex() lua.LValue {
	uv := lua.NewUserKV()
	uv.Set("mutex", lua.NewFunction(newLuaMutex))
	uv.Set("rw_mutex", lua.NewFunction(newLuaRwMutex))
	uv.Set("wait_group", lua.NewFunction(newLuaWaitGroup))
	return uv
}
