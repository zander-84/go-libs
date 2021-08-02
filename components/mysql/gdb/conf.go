package gdb

import "github.com/zander-84/go-libs/think"

type Conf struct {
	Port                string
	User                string
	Pwd                 string
	Host                string
	Database            string
	Charset             string
	MaxIdleconns        int //MaxIdleconns>=MaxOpenconns
	MaxOpenconns        int
	ConnMaxLifetime     int
	Debug               bool
	TimeZone            string
	RemoveSomeCallbacks bool
}

func (c *Conf) SetDefault() Conf {
	c.SetDefaultBasic()
	return *c
}

func (c *Conf) SetDefaultBasic() {
	if c.Port == "" {
		c.Port = "3306"
	}

	if c.User == "" {
		c.User = "root"
	}

	if c.Pwd == "" {
		c.Pwd = "123456"
	}

	if c.Host == "" {
		c.Host = "127.0.0.1"
	}

	if c.Database == "" {
		c.Database = "test"
	}

	if c.Charset == "" {
		c.Charset = "utf8mb4"
	}

	if c.MaxIdleconns == 0 {
		c.MaxIdleconns = 2000
	}

	if c.MaxOpenconns == 0 {
		c.MaxOpenconns = 2000
	}

	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = 300
	}

	if c.TimeZone == "" {
		c.TimeZone = think.DefaultTimeZone
	}
}
