package spi

import (
	"encoding/xml"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const identityKey = "uid"

var AuthMiddleware = &jwt.GinJWTMiddleware{
	Realm:       "slotopol",
	Key:         []byte("secret key"),
	Timeout:     time.Hour * 8,
	MaxRefresh:  time.Hour * 72,
	IdentityKey: identityKey,
	PayloadFunc: func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*User); ok {
			return jwt.MapClaims{
				identityKey: v.UID,
			}
		}
		return jwt.MapClaims{}
	},
	IdentityHandler: func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		var uid = uint64(claims[identityKey].(float64))
		var user, _ = Users.Get(uid)
		return user
	},
	Authenticator: func(c *gin.Context) (interface{}, error) {
		var arg struct {
			XMLName  xml.Name `json:"-" yaml:"-" xml:"arg"`
			Username string   `json:"username" form:"username" binding:"required"`
			Password string   `json:"password" form:"password" binding:"required"`
		}
		if err := c.ShouldBind(&arg); err != nil {
			return "", jwt.ErrMissingLoginValues
		}

		var reguser *User
		Users.Range(func(uid uint64, user *User) bool {
			if (user.Email == arg.Username || user.Name == arg.Username) && user.Secret == arg.Password {
				reguser = user
				return false
			}
			return true
		})
		if reguser == nil {
			return nil, jwt.ErrFailedAuthentication
		}
		return reguser, nil
	},
	LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
		var ret struct {
			XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
			Token   string   `json:"token" yaml:"token" xml:"token"`
			Expire  string   `json:"expire" yaml:"expire" xml:"expire"`
		}
		ret.Token = token
		ret.Expire = expire.Format(time.RFC3339)
		RetOk(c, ret)
	},
	RefreshResponse: func(c *gin.Context, code int, token string, expire time.Time) {
		var ret struct {
			XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
			Token   string   `json:"token" yaml:"token" xml:"token"`
			Expire  string   `json:"expire" yaml:"expire" xml:"expire"`
		}
		ret.Token = token
		ret.Expire = expire.Format(time.RFC3339)
		RetOk(c, ret)
	},
	Unauthorized: func(c *gin.Context, code int, message string) {
		Negotiate(c, http.StatusUnauthorized, ajaxerr{
			What: message,
			Code: SEC_unauthorized,
		})
	},
	TokenLookup:   "header: Authorization, query: token, cookie: jwt",
	TokenHeadName: "Bearer",
	TimeFunc:      time.Now,
}

func Handle404(c *gin.Context) {
	Ret404(c, SEC_nourl, Err404)
}

func SpiAuthHello(c *gin.Context) {
	var claims = jwt.ExtractClaims(c)
	var user, _ = c.Get(identityKey)
	c.JSON(200, gin.H{
		"userID":   claims[identityKey],
		"userName": user.(*User).Name,
		"text":     "Hello World.",
	})
}
