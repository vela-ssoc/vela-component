package component

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/audit"
	"github.com/vela-ssoc/vela-kit/lua"
	"reflect"
	"sync/atomic"
	"time"
)

var (
	//name subscript
	subscript int32 = 0

	chanTypeof = reflect.TypeOf((*channel)(nil)).String()
)

type channel struct {
	lua.SuperVelaData
	name  string
	queue chan lua.LValue
}

func newChannelByBuffer(n int) *channel {
	sub := atomic.AddInt32(&subscript, 1)
	chn := &channel{
		name:  fmt.Sprintf("channel.%d", sub),
		queue: make(chan lua.LValue, n),
	}
	chn.V(lua.VTRun, time.Now())
	return chn
}

func (c *channel) Name() string {
	return c.name
}

func (c *channel) Type() string {
	return chanTypeof
}

func (c *channel) Start() error {
	return nil
}

func (c *channel) Close() error {
	if c.IsClose() || c.queue == nil {
		return nil
	}

	defer c.recover("close channel fail")
	close(c.queue)
	c.V(lua.VTClose)
	return nil
}

func (c *channel) recover(sub string) {
	audit.Recover(audit.NewEvent("channel").Subject(sub))
}
