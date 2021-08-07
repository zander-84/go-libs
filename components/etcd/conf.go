package etcd

type Conf struct {
	Endpoints []string
	TlsPem    string
	TlsKey    string
	TlsCa     string
}

func (c *Conf) SetDefault() Conf {
	c.SetDefaultBasic()
	return *c
}

func (c *Conf) SetDefaultBasic() {

}
