package gotoggle

import (
	"context"

	"github.com/ahmetb/go-linq/v3"
	"go.chromium.org/luci/common/clock"
)

type Flag struct {
	Flag string `mapstructure:"flag" json:"flag,omitempty"`
	On   int64  `mapstructure:"on" json:"on,omitempty"`
	Off  int64  `mapstructure:"off" json:"off,omitempty"`
}

func (t Flag) IsActive(ctx context.Context, flag string) (returnData bool) {
	nowUnix := clock.Get(ctx).Now().Unix()
	returnData = t.Flag == flag

	if !returnData {
		return
	}

	var on, off int64

	if t.On <= 0 {
		on = 0
	} else {
		on = t.On
	}

	if t.Off <= 0 {
		off = 0
	} else {
		off = t.Off
	}

	if on == 0 && off == 0 {
		returnData = false
	} else if on > 0 && off == 0 {
		returnData = on > nowUnix
	} else if on == 0 && off > 0 {
		returnData = off <= nowUnix
	} else if on > 0 && off > 0 {
		if on > off {
			if off < nowUnix && on < nowUnix {
				returnData = false
			} else if off <= nowUnix && on > nowUnix {
				returnData = true
			} else if off < nowUnix && on == nowUnix {
				returnData = false
			} else if off > nowUnix {
				returnData = true
			}
		} else if on == off {
			returnData = false
		} else if on < off {
			if off < nowUnix && on < nowUnix {
				returnData = true
			} else if off > nowUnix && on > nowUnix {
				returnData = true
			} else if off > nowUnix && on <= nowUnix {
				returnData = false
			} else if off <= nowUnix {
				returnData = true
			}
		}
	}

	return
}

type Toggles struct {
	Flags []Flag `mapstructure:"flags" json:"flags,omitempty"`
}

func NewToggles(flags ...Flag) (returnData Toggles) {
	return Toggles{
		Flags: flags,
	}
}

func (ts Toggles) IsActive(ctx context.Context, flag string) (returnData bool) {
	countResonsToInactive := linq.From(ts.Flags).
		Where(func(c interface{}) bool {
			return c.(Flag).IsActive(ctx, flag)
		}).
		Count()

	returnData = countResonsToInactive <= 0

	return
}
