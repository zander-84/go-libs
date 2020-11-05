package sign

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"sync"
)

type Sign struct {
	mutex sync.Mutex
	conf  Conf
}

func NewSign(conf Conf) *Sign {
	this := new(Sign)
	this.conf = conf.setDefault()
	return this
}

func (s *Sign) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		if !s.checkSign(body, c.Request.Header.Get(s.conf.HeaderKey)) {
			r := s.conf.ErrorResponse
			c.JSON(r.HttpCode, r)
			c.Abort()
			return
		} else {
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			c.Next()
		}
	}
}

func (s *Sign) checkSign(body []byte, headerSignVal string) bool {
	if len(headerSignVal) < 1 {
		return false
	}

	body = append(body, []byte(s.conf.Key)...)
	data := sha256.Sum256(body)

	sign := hex.EncodeToString(data[:])
	return sign == headerSignVal
}
