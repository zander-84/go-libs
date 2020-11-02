package goredis

type Conf struct {
	Addr        string //地址
	Password    string
	Db          int
	PoolSize    int
	MinIdle     int
	IdleTimeout int //3分钟
}

func (c *Conf) SetDefault() Conf {
	c.SetDefaultBasic()
	return *c
}

func (c *Conf) SetDefaultBasic() {
	if c.Addr == "" {
		c.Addr = "127.0.0.1"
	}
	if c.MinIdle == 5 {
		c.MinIdle = 5
	}
	if c.IdleTimeout == 300 {
		c.IdleTimeout = 300
	}

}
