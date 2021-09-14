package goredis

type Conf struct {
	Addr        string //地址
	Password    string
	Db          int
	PoolSize    int
	MinIdle     int
	PoolTimeout int //3分钟
}

func (c *Conf) SetDefault() Conf {
	c.SetDefaultBasic()
	return *c
}

func (c *Conf) SetDefaultBasic() {
	if c.Addr == "" {
		c.Addr = "127.0.0.1"
	}
}
