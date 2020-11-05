package rate_limit

import "github.com/zander-84/go-libs/components/response"

type Conf struct {
	Max           float64
	Burst         int //实际并发
	RemoveHeader  bool
	IPLookups     []string
	ErrorResponse *response.Data
}

func (c Conf) setDefault() Conf {
	if c.Max < 1 {
		c.Max = 10
	}
	if c.ErrorResponse == nil {
		c.ErrorResponse = response.UserTooManyRequestsError(response.NewData())
	}
	if c.IPLookups == nil {
		c.IPLookups = []string{"RemoteAddr", "X-Real-IP", "X-Forwarded-For"}
	}
	return c
}
