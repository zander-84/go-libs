package validator

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
)

type ValidateCustom interface {
	Validate(c *gin.Context) map[string]interface{}
}

func ShouldBind(c *gin.Context, form ValidateCustom, opts ...func(form ValidateCustom) ValidateCustom) error {
	if err := c.ShouldBind(form); err != nil {
		return err
	} else {
		if len(opts) > 0 {
			for _, v := range opts {
				form = v(form)
			}
		}
		errMap := form.Validate(c)
		if len(errMap) > 0 {
			mjson, _ := json.Marshal(errMap)
			return errors.New(string(mjson))
		} else {
			return nil
		}
	}
}

func Map2Err(errMap map[string]interface{}) error {
	mjson, _ := json.Marshal(errMap)
	return errors.New(string(mjson))
}
