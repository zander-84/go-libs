package jwt

import (
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"time"
)

type User interface {
	PayloadFunc(data interface{}) jwt.MapClaims
	//func(data interface{}) jwt.MapClaims {
	//	if v, ok := data.(*User); ok {
	//		return jwt.MapClaims{
	//			identityKey: v.UserName,
	//		}
	//	}
	//	return jwt.MapClaims{}
	//}

	IdentityHandler(c *gin.Context) interface{}
	//	func(c *gin.Context) interface{} {
	//		claims := jwt.ExtractClaims(c)
	//		return &User{
	//			UserName: claims[identityKey].(string),
	//		}
	//	},

	Authenticator(c *gin.Context) (interface{}, error)
	//func(c *gin.Context) (interface{}, error) {
	//	var loginVals login
	//	if err := c.ShouldBind(&loginVals); err != nil {
	//		return "", jwt.ErrMissingLoginValues
	//	}
	//	userID := loginVals.Username
	//	password := loginVals.Password
	//
	//	if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
	//		return &User{
	//			UserName:  userID,
	//			LastName:  "Bo-Yi",
	//			FirstName: "Wu",
	//		}, nil
	//	}
	//
	//	return nil, jwt.ErrFailedAuthentication
	//},

	Authorizator(data interface{}, c *gin.Context) bool
	//func(data interface{}, c *gin.Context) bool {
	//	if v, ok := data.(*User); ok && v.UserName == "admin" {
	//		return true
	//	}
	//
	//	return false
	//},
	//LoginResponse(context *gin.Context, i int, s string, i2 time.Time)

	Unauthorized(c *gin.Context, code int, message string)

	LoginResponse(c *gin.Context, i int, token string, expire time.Time)

	RefreshResponse(*gin.Context, int, string, time.Time)

	LogoutResponse(*gin.Context, int)
}
