package component

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
)

func (c *channel) pop(L *lua.LState) int {
	thd := L.IsInt(1)
	pip := pipe.NewByLua(L, pipe.Env(xEnv), pipe.Seek(1))

	for i := 0; i <= thd; i++ {
		xEnv.Spawn(0, func() {
			co := xEnv.Clone(L)
			defer xEnv.Free(co)

			for lv := range c.queue {
				pip.Do(lv, co, func(err error) {
					xEnv.Errorf("channel pop handle error %v", err)
				})
			}
		})
	}
	return 0
}

func (c *channel) push(L *lua.LState) int {
	if c.IsClose() {
		return 0
	}

	n := L.GetTop()
	if n == 0 {
		return 0
	}

	defer c.recover("push channel fail")

	for i := 1; i <= n; i++ {
		lv := L.Get(i)
		if lv.Type() != lua.LTNil {
			c.queue <- lv
		}
	}

	return 0
}

func (c *channel) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "pop":
		return lua.NewFunction(c.pop)
	case "push":
		return lua.NewFunction(c.push)
	}

	return lua.LNil
}

func newLuaChannel(L *lua.LState) int {
	n := L.IsInt(1)
	ch := newChannelByBuffer(n)

	proc := L.NewVelaData(ch.name, chanTypeof)
	if proc.IsNil() {
		proc.Set(ch)
	} else {
		old := proc.Data.(*channel)
		old.Close()
		proc.Set(ch)
	}

	L.Push(proc)
	return 1
}

func newLuaChannelIndex() lua.LValue {
	return lua.NewFunction(newLuaChannel)
}
