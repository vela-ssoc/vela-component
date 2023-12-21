package component

import (
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
)

func newLuaJsonEncode(L *lua.LState) int {
	lv := L.CheckAny(1)

	data, err := kind.MarshalJson(lv)
	if err != nil {
		L.Push(lua.LNil)
		L.Pushf("%v", err)
		return 2
	} else {
		L.Push(lua.B2L(data))
		return 1
	}
}

func newLuaJsonDecode(L *lua.LState) int {
	str := L.CheckString(1)

	value, err := kind.Decode(L, lua.S2B(str))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(value)
	return 1
}

func newLuaJsonKind(L *lua.LState) int {
	enc := &JsonKind{inline: kind.NewJsonEncoder()}
	L.Push(enc)
	return 1
}

/*
	local v = vela.json("{name:123 , pass:123}")
	print(v.name)
	print(v.name)

	local kind = vela.json.kind()
	kind.tab("")
	kind.kv("name" , "vela")
	kind.kv("pass" , "123456")
	kind.kv("pass" , "123456")
	kind.kv("pass" , "123456")
	kind.kv("pass" , "123456")
	kind.kv("pass" , "123456")
	kind.end("}")
*/

func newLuaJsonIndex(env vela.Environment) {
	tuple := lua.NewUserKV()
	tuple.Set("encode", lua.NewFunction(newLuaJsonEncode))
	tuple.Set("decode", lua.NewFunction(newLuaJsonDecode))
	tuple.Set("kind", lua.NewFunction(newLuaJsonKind))
	env.Set("json", lua.NewExport("vela.json.export", lua.WithTable(tuple), lua.WithFunc(newLuaFastJson)))
}
