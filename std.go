package component

import (
	auxlib2 "github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	"io"
	"os"
	"time"
)

type stdio struct {
	lua.SuperVelaData
	name string
	w    io.Writer
}

func newStdio(name string, w io.Writer) *stdio {
	s := &stdio{name: name, w: w}
	s.V(lua.VTRun, time.Now(), "std."+name)
	return s
}

func (st *stdio) Name() string {
	return "std" + st.name
}

func (st *stdio) Close() error {
	return nil
}
func (st *stdio) Start() error {
	return nil
}

func (st *stdio) Type() string {
	return "std." + st.name
}

func (st *stdio) toLValue() *lua.VelaData {
	return lua.NewVelaData(st)
}

func (st *stdio) Write(p []byte) (int, error) {
	return os.Stderr.Write(p)
}

func (st *stdio) Push(v interface{}) error {
	_, err := auxlib2.Push(st, v)
	return err
}

func (st *stdio) printL(L *lua.LState) int {
	st.Push(auxlib2.Format(L, 0))
	return 0
}

func (st *stdio) printlnL(L *lua.LState) int {
	n := L.GetTop()
	if n == 0 {
		return 0
	}

	for i := 1; i <= n; i++ {
		st.Push(L.Get(1).String())
	}
	st.Push("\n")
	return 0
}

func (st *stdio) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "print":
		return L.NewFunction(st.printL)
	case "println":
		return L.NewFunction(st.printlnL)

	default:
		return lua.LNil
	}
}

func newLuaStdIndex() lua.LValue {
	std := lua.NewUserKV()
	std.Set("out", newStdio("out", os.Stdout).toLValue())
	std.Set("err", newStdio("err", os.Stderr).toLValue())
	return std
}
