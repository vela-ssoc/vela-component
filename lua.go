package component

import (
	"github.com/vela-ssoc/vela-kit/vela"
)

var xEnv vela.Environment

func WithEnv(env vela.Environment) {
	xEnv = env

	//cond
	newLuaCondIndex(env)

	//fmt
	newLuaFmtIndex(env)

	//for
	newLuaForIndex(env)

	//json
	newLuaJsonIndex(env)

	//args
	env.Set("args", newLuaArgsIndex())

	//copy
	env.Set("copy", newLuaCopyIndex())

	//channel
	env.Set("channel", newLuaChannelIndex())

	//catch
	env.Set("catch", newLuaCatchIndex())

	//sync
	env.Global("sync", newLuaSyncIndex())

	//atomic
	env.Global("atomic", newLuaAtomicIndex())

	//std
	env.Global("std", newLuaStdIndex())
}
