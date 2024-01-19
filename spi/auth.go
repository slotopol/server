package spi

import (
	"encoding/xml"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/util"
)

const identityKey = "uid"

var AuthMiddleware = &jwt.GinJWTMiddleware{
	Realm:       "slotopol",
	Key:         []byte("secret key"), // loaded from config later
	Timeout:     time.Hour * 24,       // loaded from config later
	MaxRefresh:  time.Hour * 72,       // loaded from config later
	IdentityKey: identityKey,
	Authenticator: func(c *gin.Context) (interface{}, error) {
		var arg struct {
			XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
			Email   string   `json:"email" yaml:"email" xml:"email" form:"email" binding:"required"`
			Secret  string   `json:"secret" yaml:"secret" xml:"secret" form:"secret" binding:"required"`
		}
		if err := c.ShouldBind(&arg); err != nil {
			return "", jwt.ErrMissingLoginValues
		}

		var reguser *User
		Users.Range(func(uid uint64, user *User) bool {
			if user.Email == arg.Email && user.Secret == arg.Secret {
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
	Authorizator: func(data interface{}, c *gin.Context) bool {
		// on case if user was deleted by another workflow before request
		return data != nil
	},
	PayloadFunc: func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*User); ok {
			return jwt.MapClaims{
				identityKey: v.UID,
			}
		}
		return jwt.MapClaims{}
	},
	IdentityHandler: func(c *gin.Context) interface{} {
		var claims = jwt.ExtractClaims(c)
		var uid = uint64(claims[identityKey].(float64))
		if user, ok := Users.Get(uid); ok {
			return user
		}
		return nil
	},
	LoginResponse:   TokenResponse,
	RefreshResponse: TokenResponse,
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

func TokenResponse(c *gin.Context, code int, token string, expire time.Time) {
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		Token   string   `json:"token" yaml:"token" xml:"token"`
		Expire  string   `json:"expire" yaml:"expire" xml:"expire"`
	}
	ret.Token = token
	ret.Expire = expire.Format(time.RFC3339)
	RetOk(c, ret)
}

func Handle404(c *gin.Context) {
	Ret404(c, SEC_nourl, Err404)
}

func SpiSignup(c *gin.Context) {
	var err error
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		Email   string   `json:"email" yaml:"email" xml:"email" form:"email"`
		Secret  string   `json:"secret" yaml:"secret" xml:"secret" form:"secret"`
		Name    string   `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_signup_nobind, err)
		return
	}
	if len(arg.Secret) < 6 {
		Ret400(c, SEC_signup_smallsec, ErrSmallKey)
		return
	}
	if arg.Email == "" {
		Ret400(c, SEC_signup_noemail, ErrNoData)
		return
	}

	var email = util.ToLower(arg.Email)

	var user = &User{
		Email:  email,
		Secret: arg.Secret,
		Name:   arg.Name,
	}
	if _, err = cfg.XormStorage.Insert(user); err != nil {
		Ret500(c, SEC_signup_insert, err)
		return
	}

	Users.Set(user.UID, user)

	ret.UID = user.UID
	RetOk(c, ret)
}
