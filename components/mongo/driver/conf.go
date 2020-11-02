package driver

type Conf struct {
	Host            string
	Port            string
	MaxPoolSize     uint64
	MinPoolSize     uint64
	MaxConnIdleTime int
	Database        string

	User         string
	Pwd          string
	Charset      string
	MaxIdleconns int //MaxIdleconns>=MaxOpenconns
	MaxOpenconns int64
	Debug        bool
	TimeZone     string
}

func (c *Conf) SetDefault() Conf {
	c.SetDefaultBasic()
	return *c
}

func (c *Conf) SetDefaultBasic() {

}
