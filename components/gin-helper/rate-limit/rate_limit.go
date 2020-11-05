package rate_limit

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"sync"
)

type RateLimiter struct {
	conf   Conf
	Engine *limiter.Limiter
	mutex  sync.Mutex
}

func NewRateLimiter(conf Conf) *RateLimiter {
	this := new(RateLimiter)
	this.conf = conf.setDefault()
	this.Engine = tollbooth.NewLimiter(this.conf.Max, nil)
	this.Engine.SetIPLookups(this.conf.IPLookups)
	this.Engine.SetBurst(this.conf.Burst)
	return this
}

func (r *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(r.Engine, c.Writer, c.Request)

		if r.conf.RemoveHeader {
			c.Writer.Header().Del("X-Rate-Limit-Request-Forwarded-For")
			c.Writer.Header().Del("X-Rate-Limit-Request-Remote-Addr")
			c.Writer.Header().Del("X-Rate-Limit-Duration")
			c.Writer.Header().Del("X-Rate-Limit-Limit")
		}

		if httpError != nil {
			r := r.conf.ErrorResponse
			c.JSON(r.HttpCode, r)
			c.Abort()
			return
		} else {
			c.Next()
		}

	}
}
