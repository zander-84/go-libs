package memory

type Conf struct {
	Expiration      int //分钟 Create a cache with a default expiration time of 5 minutes, and which
	CleanupInterval int //分钟 purges expired items every 10 minutes
}

func (c *Conf) SetDefault() Conf {
	c.SetDefaultBasic()
	return *c
}

func (c *Conf) SetDefaultBasic() {
	if c.Expiration == 0 {
		c.Expiration = -1
	}
	if c.CleanupInterval == 0 {
		c.Expiration = 10
	}
}
