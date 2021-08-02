package robfig

import (
	"github.com/zander-84/go-libs/think"
)

type Conf struct {
	TimeZone string
}

func (c *Conf) SetDefault() Conf {
	c.SetDefaultBasic()
	return *c
}

func (c *Conf) SetDefaultBasic() {
	if c.TimeZone == "" {
		c.TimeZone = think.DefaultTimeZone
	}
}
