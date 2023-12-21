package component

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
	"github.com/vela-ssoc/vela-kit/vela"
	"strings"
)

/*

	local w = vela.to(es).dict("abc" , "abc").sub("abc" , "abc")



*/

type To struct {
	co      *lua.LState
	chains  *pipe.Chains
	extract []func(v lua.IndexEx, data *lua.Map)
}

func (t *To) String() string                         { return fmt.Sprintf("vela.to.%p", t) }
func (t *To) Type() lua.LValueType                   { return lua.LTObject }
func (t *To) AssertFloat64() (float64, bool)         { return 0, false }
func (t *To) AssertString() (string, bool)           { return "", false }
func (t *To) AssertFunction() (*lua.LFunction, bool) { return lua.NewFunction(t.LCall), true }
func (t *To) Peek() lua.LValue                       { return t }

func (t *To) LCall(L *lua.LState) int {
	return 0
}

func (t *To) PCall(args ...interface{}) error {
	if len(args) == 0 {
		return nil
	}

	val := args[0]
	if val == nil {
		return nil
	}

	n := len(t.extract)
	if n == 0 {
		t.chains.Do(val, t.co, func(err error) {
			//todo
		})
		return nil
	}

	vex, ok := args[0].(lua.IndexEx)
	if !ok {
		return fmt.Errorf("to must be index object")
	}

	data := lua.NewMap(32, false)
	for i := 0; i < n; i++ {
		fn := t.extract[i]
		fn(vex, data)
	}

	data.Set("minion_id", lua.LString(xEnv.ID()))
	data.Set("minion_inet", lua.LString(xEnv.Inet()))

	t.chains.Do(data, t.co, func(err error) {
		//todo

	})
	return nil
}

func (t *To) keys(L *lua.LState) []string {
	var keys []string

	L.Callback(func(v lua.LValue) (stop bool) {
		item := lua.IsString(v)
		if item == "" {
			L.RaiseError("invalid key got empty")
			return true
		}

		keys = append(keys, item)
		return false
	})
	return keys
}

func (t *To) Setter(data *lua.Map, k string, v lua.LValue) {
	k = strings.ToLower(k)

	if v == nil || v.Type() == lua.LTNil {
		data.Set(k, nil)
		return
	}

	data.Set(k, v)
}

func (t *To) DictL(L *lua.LState) int {
	keys := t.keys(L)
	n := len(keys)
	if n == 0 {
		L.Push(t)
		return 1
	}

	extract := func(v lua.IndexEx, data *lua.Map) {
		for i := 0; i < n; i++ {
			item := keys[i]
			lv := v.Index(t.co, keys[i])

			if lv != nil && lv.Type() != lua.LTNil {
				data.Set(item, lv)
				continue
			}
			data.Set(item, lua.S2L(""))
		}
	}

	t.extract = append(t.extract, extract)
	L.Push(t)
	return 1
}

func (t *To) SubL(L *lua.LState) int {
	keys := t.keys(L)
	n := len(keys)
	if n < 2 {
		L.Push(t)
		return 1
	}

	extract := func(v lua.IndexEx, data *lua.Map) {
		name := keys[0]
		lv := v.Index(t.co, name)

		getter := func(string) lua.LValue {
			return nil
		}

		ex, ok := lv.(lua.IndexEx)
		if ok {
			getter = func(s string) lua.LValue {
				return ex.Index(t.co, s)
			}
		}

		for i := 1; i < n; i++ {
			k := keys[i]
			t.Setter(data, name+"_"+k, getter(k))
		}
	}
	t.extract = append(t.extract, extract)
	L.Push(t)
	return 1
}

func (t *To) formatL(L *lua.LState) int {
	return 0
}

func (t *To) pipeL(L *lua.LState) int {
	t.chains.CheckMany(L)
	return 0
}

func (t *To) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "dict":
		return lua.NewFunction(t.DictL)
	case "sub":
		return lua.NewFunction(t.SubL)
	case "format":
		return lua.NewFunction(t.formatL)
	case "pipe":
		return lua.NewFunction(t.pipeL)
	}
	return lua.LNil
}

func newToL(L *lua.LState) int {
	L.Push(&To{
		co:     xEnv.Clone(L),
		chains: pipe.NewByLua(L),
	})
	return 1
}

func newLuaToIndex(env vela.Environment) {
	env.Set("to", lua.NewExport("lua.to.export", lua.WithFunc(newToL)))
}
