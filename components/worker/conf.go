package worker

type Conf struct {
	MaxWorkers int
	MaxQueues  int
}

func (c *Conf) SetDefault() Conf {
	c.SetDefaultBasic()
	return *c
}

func (c *Conf) SetDefaultBasic() {
	if c.MaxWorkers == 0 {
		c.MaxWorkers = 20
	}

	if c.MaxQueues == 0 {
		c.MaxQueues = 100000
	}
}
