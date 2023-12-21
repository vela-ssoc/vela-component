package component

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	strutil "github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	"io"
	"reflect"
	"time"
)

var (
	copyTypeOf = reflect.TypeOf((*copyGo)(nil)).String()
)

type copyGo struct {
	lua.SuperVelaData

	name string

	size int

	ctx  context.Context
	stop context.CancelFunc

	src lua.Reader
	dst lua.Writer

	vm string
}

//rock.copy(src , dst , tag)

func key(s1, s2, tag string) string {
	if e := strutil.Name(tag); e == nil {
		return tag
	}

	h := md5.New()
	io.WriteString(h, s1)
	io.WriteString(h, "_")
	io.WriteString(h, s2)
	io.WriteString(h, "-")
	return "copy_" + hex.EncodeToString(h.Sum(nil))
}

func newCopy(src lua.Reader, dst lua.Writer, tag string) *copyGo {
	obj := &copyGo{
		src: src,
		dst: dst,
	}

	obj.name = key(src.Name(), dst.Name(), tag)

	obj.V(lua.VTInit, copyTypeOf)

	return obj
}

func checkSrc(L *lua.LState, idx int) lua.Reader {
	src := lua.CheckReader(L.CheckVelaData(idx))
	if src == nil {
		L.RaiseError("invalid reader")
		return nil
	}

	return src
}

func checkDst(L *lua.LState, idx int) lua.Writer {
	dst := lua.CheckWriter(L.CheckVelaData(idx))
	if dst == nil {
		L.RaiseError("invalid writer")
		return nil
	}
	return dst
}

func checkTag(L *lua.LState, idx int) string {
	if tv := L.Get(idx); tv.Type() == lua.LTString {
		return strutil.CheckProcName(tv, L)
	}

	return ""
}

func newLuaCopy(L *lua.LState) int {
	src := checkSrc(L, 1)
	dst := checkDst(L, 2)
	tag := checkTag(L, 3)

	obj := newCopy(src, dst, tag)

	proc := L.NewVelaData(obj.Name(), copyTypeOf)
	proc.Set(obj)

	L.Push(proc)
	return 1
}

func (c *copyGo) Name() string {
	return c.name
}

func (c *copyGo) Type() string {
	return c.TypeOf
}

func (c *copyGo) handle() {

	_, err := strutil.Copy(c.ctx, c.dst, c.src)
	if err != nil {
		xEnv.Errorf("%s %s flow close error %v", c.vm, c.Name(), err)
	}

	select {
	case <-c.ctx.Done():
		xEnv.Errorf("%s %s exit", c.vm, c.Name())
		return
	default:
		xEnv.Errorf("%s %s copyGo restart", c.vm, c.Name())
	}
}

func (c *copyGo) Start() error {
	ctx, stop := context.WithCancel(context.Background())
	c.ctx = ctx
	c.stop = stop

	xEnv.Spawn(0, c.handle)
	c.V(lua.VTRun, time.Now())
	return nil
}

func (c *copyGo) Close() error {
	c.stop()
	return nil
}

func (c *copyGo) spawn(L *lua.LState) int {
	c.ctx, c.stop = context.WithCancel(context.Background())
	xEnv.Spawn(0, c.handle)
	c.V(lua.VTRun, time.Now())
	return 0
}

func (c *copyGo) run(L *lua.LState) int {
	c.ctx, c.stop = context.WithCancel(context.Background())
	c.V(lua.VTRun, time.Now())
	c.handle()
	return 0
}

func (c *copyGo) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "spawn":
		return L.NewFunction(c.spawn)
	case "run":
		return L.NewFunction(c.run)
	}
	return lua.LNil
}

func newLuaCopyIndex() lua.LValue {
	return lua.NewFunction(newLuaCopy)
}
