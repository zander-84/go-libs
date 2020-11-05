package sign

import "github.com/zander-84/go-libs/components/response"

type Conf struct {
	Key           string
	HeaderKey     string
	ErrorResponse *response.Data
}

func (c Conf) setDefault() Conf {
	if c.Key == "" {
		c.Key = "hello gin"
	}

	if c.HeaderKey == "" {
		c.HeaderKey = "sign"
	}

	if c.ErrorResponse == nil {
		c.ErrorResponse = response.UserForbiddenError(response.NewData())
	}

	return c
}
