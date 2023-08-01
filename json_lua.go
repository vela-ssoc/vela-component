package component

import (
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-kit/lua"
)

type JsonKind struct {
	inline *kind.JsonEncoder
}

func (enc *JsonKind) String() string                         { return lua.B2S(enc.inline.Bytes()) }
func (enc *JsonKind) Type() lua.LValueType                   { return lua.LTObject }
func (enc *JsonKind) AssertFloat64() (float64, bool)         { return 0, false }
func (enc *JsonKind) AssertString() (string, bool)           { return "", false }
func (enc *JsonKind) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (enc *JsonKind) Peek() lua.LValue                       { return enc }

func (enc *JsonKind) convert(key string, val lua.LValue) {

	switch val.Type() {
	case lua.LTBool:
		enc.inline.KV(key, bool(val.(lua.LBool)))
	case lua.LTString:
		enc.inline.KV(key, val.String())
	case lua.LTNumber:
		enc.inline.KV(key, float64(val.(lua.LNumber)))
	case lua.LTInt:
		enc.inline.KV(key, int(val.(lua.LInt)))
	case lua.LTAnyData:
		enc.inline.KV(key, val.(*lua.AnyData).Data)
	case lua.LTUserData:
		enc.inline.KV(key, val.(*lua.LUserData).Value)
	case lua.LTVelaData:
		enc.inline.KV(key, val.(*lua.VelaData).Data)
	default:
		enc.inline.KV(key, val.String())
	}

}

func (enc *JsonKind) metaKV(L *lua.LState) int {
	key := L.CheckString(1)
	val := L.CheckAny(2)
	enc.convert(key, val)
	return 0
}

func (enc *JsonKind) metaInt64(L *lua.LState) int {
	key := L.CheckString(1)
	val := L.CheckNumber(2)
	enc.inline.KL(key, int64(val))
	return 0
}

func (enc *JsonKind) metaUint64(L *lua.LState) int {
	key := L.CheckString(1)
	val := L.CheckNumber(2)
	enc.inline.KUL(key, uint64(val))
	return 0
}

func (enc *JsonKind) metaEof(L *lua.LState) int {
	val := L.CheckString(1)
	enc.inline.End(val)
	return 0
}

func (enc *JsonKind) metaTab(L *lua.LState) int {
	val := L.CheckString(1)
	enc.inline.Tab(val)
	return 0
}

func (enc *JsonKind) metaArr(L *lua.LState) int {
	val := L.CheckString(1)
	enc.inline.Arr(val)
	return 0
}

func (enc *JsonKind) metaStr(L *lua.LState) int {
	L.Push(lua.B2L(enc.inline.Bytes()))
	return 1
}

func (enc *JsonKind) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "kv":
		return L.NewFunction(enc.metaKV)
	case "int64":
		return L.NewFunction(enc.metaInt64)
	case "uint64":
		return L.NewFunction(enc.metaUint64)
	case "tab":
		return L.NewFunction(enc.metaTab)
	case "arr":
		return L.NewFunction(enc.metaArr)
	case "eof":
		return L.NewFunction(enc.metaEof)
	case "byte":
		return L.NewFunction(enc.metaStr)
	}

	return L.NewFunction(func(_ *lua.LState) int {
		xEnv.Debugf("not found %s json kind function", key)
		return 0
	})
}
