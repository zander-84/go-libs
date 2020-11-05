package jwt

import (
	jwt2 "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
	"time"
)

type Jwt struct {
	conf           Conf
	AuthMiddleware *jwt2.GinJWTMiddleware
	mutex          sync.Mutex
}

func New(conf Conf, user User) *Jwt {
	this := new(Jwt)
	this.conf = conf

	return this.Init(user)
}

func (j *Jwt) Init(user User) *Jwt {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	var err error

	j.AuthMiddleware, err = jwt2.New(&jwt2.GinJWTMiddleware{
		Realm:            j.conf.Realm,
		SigningAlgorithm: j.conf.SigningAlgorithm,
		Key:              []byte(j.conf.Key),
		Timeout:          j.conf.Timeout,
		MaxRefresh:       j.conf.MaxRefresh,
		Authenticator:    user.Authenticator,
		Authorizator:     user.Authorizator,
		PayloadFunc:      user.PayloadFunc,
		Unauthorized:     user.Unauthorized,
		LoginResponse:    user.LoginResponse,
		LogoutResponse:   user.LogoutResponse,
		RefreshResponse:  user.RefreshResponse,
		IdentityHandler:  user.IdentityHandler,
		IdentityKey:      j.conf.IdentityKey, // c.set(IdentityKey, )
		TokenLookup:      j.conf.TokenLookup,
		TokenHeadName:    j.conf.TokenHeadName,
		TimeFunc:         time.Now,

		HTTPStatusMessageFunc: nil,
		PrivKeyFile:           "",
		PubKeyFile:            "",
		SendCookie:            false,
		SecureCookie:          false,
		CookieHTTPOnly:        false,
		CookieDomain:          "",
		SendAuthorization:     false,
		DisabledAbort:         false,
		CookieName:            "",
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return j
}

func (j *Jwt) Middleware() gin.HandlerFunc {
	return j.AuthMiddleware.MiddlewareFunc()
}
