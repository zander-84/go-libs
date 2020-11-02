package validate

type Conf struct {
	Locale      string //选择语言
	ValidateTag string //validate
	CommentTag  string //comment    比如 user      user:用户
}

func (c *Conf) SetDefault() Conf {
	c.SetDefaultBasic()
	return *c
}

func (c *Conf) SetDefaultBasic() {
	//if c.JsonTag == "" {
	//	c.JsonTag = "json"
	//}
	if c.CommentTag == "" {
		c.CommentTag = "comment"
	}
	if c.ValidateTag == "" {
		c.ValidateTag = "validate"
	}
	if c.Locale == "" {
		c.Locale = "zh"
	}
}
