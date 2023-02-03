package component

import (
	"github.com/vela-ssoc/vela-kit/vela"
	"github.com/vela-ssoc/vela-kit/grep"
	"github.com/vela-ssoc/vela-kit/lua"
	"regexp"
	"strings"
)

type cond struct {
	value bool
}

func (c *cond) String() string                         { return "true" }
func (c *cond) Type() lua.LValueType                   { return lua.LTObject }
func (c *cond) AssertFloat64() (float64, bool)         { return 0, false }
func (c *cond) AssertString() (string, bool)           { return "", false }
func (c *cond) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (c *cond) Peek() lua.LValue                       { return c }

func (c *cond) do(L *lua.LState) {
	n := L.GetTop()
	if n < 1 {
		return
	}
	L.PCall(n-1, 0, nil)
}

func (c *cond) Y(L *lua.LState) int {
	if c.value {
		c.do(L)
	}
	L.Push(c)
	return 1
}

func (c *cond) N(L *lua.LState) int {
	if !c.value {
		c.do(L)
	}
	L.Push(c)
	return 1
}

func (c *cond) ok(L *lua.LState) int {
	L.Push(lua.LBool(c.value))
	return 1
}

func (c *cond) Index(L *lua.LState, key string) lua.LValue {

	switch key {
	case "Y":
		return L.NewFunction(c.Y)
	case "N":
		return L.NewFunction(c.N)
	case "OK":
		return lua.LBool(c.value)
	}

	return lua.LNil
}

func regexL(L *lua.LState, src lua.LValue, re lua.LValue) bool {
	r, err := regexp.Compile(re.String())
	if err != nil {
		return false
	}
	return r.MatchString(src.String())
}

func equalL(L *lua.LState, src lua.LValue, val lua.LValue) bool {
	return L.Equal(src, val)
}

func containL(L *lua.LState, src lua.LValue, val lua.LValue) bool {
	return strings.Contains(src.String(), val.String())
}

func suffixL(L *lua.LState, src lua.LValue, val lua.LValue) bool {

	if src.Type() != val.Type() {
		return false
	}

	if src.Type() != lua.LTString {
		return false
	}

	return strings.HasSuffix(src.String(), val.String())
}

func prefixL(_ *lua.LState, src lua.LValue, v lua.LValue) bool {
	if src.Type() != v.Type() {
		return false
	}

	if src.Type() != lua.LTString {
		return false
	}

	return strings.HasPrefix(src.String(), v.String())
}

func grepL(L *lua.LState, src lua.LValue, val lua.LValue) bool {
	if val.Type() != lua.LTString {
		return false
	}

	g, err := grep.Compile(val.String(), nil)
	if err != nil {
		return false
	}

	return g.Match(src.String())

}

func doFunc(L *lua.LState, fn func(*lua.LState, lua.LValue, lua.LValue) bool) bool {
	n := L.GetTop()

	switch n {
	case 0:
		return false

	case 1:
		return L.Get(1).Type() != lua.LTNil

	case 2:
		return fn(L, L.Get(1), L.Get(2))

	default:
		lv := L.Get(1)
		for i := 2; i <= n; i++ {
			if fn(L, lv, L.Get(i)) {
				return true
			}
		}
		return false
	}
}

type cFunc func(*lua.LState, lua.LValue, lua.LValue) bool

func newLuaFn(fn cFunc) *lua.LFunction {
	return lua.NewFunction(func(L *lua.LState) int {
		L.Push(&cond{value: doFunc(L, fn)})
		return 1
	})
}

func newLuaCondIndex(env vela.Environment) {
	env.Set("regex", newLuaFn(regexL))
	env.Set("equal", newLuaFn(equalL))
	env.Set("contain", newLuaFn(containL))
	//env.Set("suffix", newLuaFn(suffixL))
	//env.Set("prefix", newLuaFn(prefixL))
	env.Set("grep", newLuaFn(grepL))
}
