package jwt

import "time"

type Conf struct {
	Realm            string //zander zone
	SigningAlgorithm string //默认 HS256
	Key              string // 秘钥
	Timeout          time.Duration
	MaxRefresh       time.Duration
	IdentityKey      string //identity

	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	// - "param:<name>"
	TokenLookup string

	// TokenHeadName is a string in the header. Default value is "Bearer"
	TokenHeadName string
}
