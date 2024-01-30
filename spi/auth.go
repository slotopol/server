package spi

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/util"
)

const (
	// "iss" field for this tokens.
	jwtIssuer = "slotopol"

	// Pointer to User object stored at gin context
	// after successful authorization.
	userKey = "user"

	realmBasic  = `Basic realm="slotopol", charset="UTF-8"`
	realmBearer = `JWT realm="slotopol", charset="UTF-8"`
)

var (
	ErrNoJwtID  = errors.New("jwt-token does not have user id")
	ErrBadJwtID = errors.New("jwt-token id does not refer to registered user")
	ErrNoAuth   = errors.New("authorization is required")
	ErrNoScheme = errors.New("authorization does not have expected scheme")
	ErrNoCred   = errors.New("user with given credentials does not registered")
	ErrBadPass  = errors.New("password is incorrect")
)

var (
	Cfg = cfg.Cfg
)

// Claims of JWT-tokens. Contains additional profile identifier.
type Claims struct {
	jwt.RegisteredClaims
	UID uint64 `json:"uid,omitempty"`
}

func (c *Claims) Validate() error {
	if c.UID == 0 {
		return ErrNoJwtID
	}
	return nil
}

type AuthGetter func(c *gin.Context) (*User, int, error)

// AuthGetters is the list of functions to extract the authorization
// data from the parts of request. List and order in it can be changed.
var AuthGetters = []AuthGetter{
	UserFromHeader, UserFromQuery, UserFromCookie,
}

// Auth is authorization middleware, sets User object associated
// with authorization to gin context. `required` parameter tells
// to continue if authorization is absent.
func Auth(required bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var code int
		var user *User
		for _, getter := range AuthGetters {
			if user, code, err = getter(c); err != nil {
				Ret401(c, code, err)
				return
			}
			if user != nil {
				break
			}
		}

		if user != nil {
			c.Set(userKey, user)
		} else if required {
			Ret401(c, SEC_auth_absent, ErrNoAuth)
			return
		}

		c.Next()
	}
}

func UserFromHeader(c *gin.Context) (*User, int, error) {
	if hdr := c.Request.Header.Get("Authorization"); hdr != "" {
		if strings.HasPrefix(hdr, "Basic ") {
			return GetBasicAuth(hdr[6:])
		} else if strings.HasPrefix(hdr, "Bearer ") {
			return GetBearerAuth(hdr[7:])
		} else {
			return nil, SEC_auth_scheme, ErrNoScheme
		}
	}
	return nil, 0, nil
}

func UserFromQuery(c *gin.Context) (*User, int, error) {
	if credentials := c.Query("cred"); credentials != "" {
		return GetBasicAuth(credentials)
	} else if tokenstr := c.Query("token"); tokenstr != "" {
		return GetBearerAuth(tokenstr)
	} else if tokenstr := c.Query("jwt"); tokenstr != "" {
		return GetBearerAuth(tokenstr)
	}
	return nil, 0, nil
}

func UserFromCookie(c *gin.Context) (*User, int, error) {
	if credentials, _ := c.Cookie("cred"); credentials != "" {
		return GetBasicAuth(credentials)
	} else if tokenstr, _ := c.Cookie("token"); tokenstr != "" {
		return GetBearerAuth(tokenstr)
	} else if tokenstr, _ := c.Cookie("jwt"); tokenstr != "" {
		return GetBearerAuth(tokenstr)
	}
	return nil, 0, nil
}

func UserFromForm(c *gin.Context) (*User, int, error) {
	if credentials := c.PostForm("cred"); credentials != "" {
		return GetBasicAuth(credentials)
	} else if tokenstr := c.PostForm("token"); tokenstr != "" {
		return GetBearerAuth(tokenstr)
	} else if tokenstr := c.PostForm("jwt"); tokenstr != "" {
		return GetBearerAuth(tokenstr)
	}
	return nil, 0, nil
}

func GetBasicAuth(credentials string) (user *User, code int, err error) {
	var decoded []byte
	if decoded, err = base64.RawURLEncoding.DecodeString(credentials); err != nil {
		return nil, SEC_basic_decode, err
	}
	var parts = strings.Split(util.B2S(decoded), ":")

	var email = util.ToLower(parts[0])
	Users.Range(func(uid uint64, u *User) bool {
		if u.Email != email {
			return true
		}
		user = u
		return false
	})
	if user == nil {
		err, code = ErrNoCred, SEC_basic_nouser
		return
	}
	if user.Secret != parts[1] {
		err, code = ErrBadPass, SEC_basic_deny
		return
	}
	return
}

func GetBearerAuth(tokenstr string) (user *User, code int, err error) {
	var claims Claims
	_, err = jwt.ParseWithClaims(tokenstr, &claims, func(*jwt.Token) (any, error) {
		var keys = jwt.VerificationKeySet{
			Keys: []jwt.VerificationKey{
				util.S2B(Cfg.AccessKey),
				util.S2B(Cfg.RefreshKey),
			},
		}
		return keys, nil
	}, jwt.WithExpirationRequired(), jwt.WithIssuer(jwtIssuer), jwt.WithLeeway(5*time.Second))

	if err == nil {
		var ok bool
		if user, ok = Users.Get(claims.UID); !ok {
			err, code = ErrBadJwtID, SEC_token_baduid
		}
		return
	}
	switch {
	case errors.Is(err, jwt.ErrTokenMalformed):
		code = SEC_token_malform
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		code = SEC_token_notsign
	case errors.Is(err, jwt.ErrTokenInvalidClaims):
		code = SEC_token_badclaims
	case errors.Is(err, jwt.ErrTokenExpired):
		code = SEC_token_expired
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		code = SEC_token_notyet
	case errors.Is(err, jwt.ErrTokenInvalidIssuer):
		code = SEC_token_issuer
	default:
		code = SEC_token_error
	}
	return
}

type AuthResp struct {
	XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
	UID     uint64   `json:"uid" yaml:"uid" xml:"uid"`
	Access  string   `json:"access" yaml:"access" xml:"access"`
	Refrsh  string   `json:"refrsh" yaml:"refrsh" xml:"refrsh"`
	Expire  string   `json:"expire" yaml:"expire" xml:"expire"`
	Living  string   `json:"living" yaml:"living" xml:"living"`
}

func (r *AuthResp) Setup(user *User) {
	var err error
	var token *jwt.Token
	var now = jwt.NewNumericDate(time.Now())
	var exp = jwt.NewNumericDate(time.Now().Add(Cfg.AccessTTL))
	var age = jwt.NewNumericDate(time.Now().Add(Cfg.RefreshTTL))
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: now,
			ExpiresAt: exp,
			Issuer:    jwtIssuer,
		},
		UID: user.UID,
	})
	if r.Access, err = token.SignedString([]byte(Cfg.AccessKey)); err != nil {
		panic(err)
	}
	r.Expire = exp.Format(time.RFC3339)
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: now,
			ExpiresAt: age,
			Issuer:    jwtIssuer,
		},
		UID: user.UID,
	})
	if r.Refrsh, err = token.SignedString([]byte(Cfg.AccessKey)); err != nil {
		panic(err)
	}
	r.Living = age.Format(time.RFC3339)
	r.UID = user.UID
}

func Handle404(c *gin.Context) {
	Ret404(c, SEC_nourl, Err404)
}

func SpiSignup(c *gin.Context) {
	var err error
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		Email   string   `json:"email" yaml:"email" xml:"email" form:"email" binding:"required"`
		Secret  string   `json:"secret" yaml:"secret" xml:"secret" form:"secret" binding:"required"`
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

func SpiSignin(c *gin.Context) {
	var err error
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		Email   string   `json:"email" yaml:"email" xml:"email" form:"email" binding:"required"`
		Secret  string   `json:"secret" yaml:"secret" xml:"secret" form:"secret" binding:"required"`
	}
	var ret AuthResp

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_signin_nobind, err)
		return
	}
	if len(arg.Secret) < 6 {
		Ret400(c, SEC_signin_smallsec, ErrSmallKey)
		return
	}

	var email = util.ToLower(arg.Email)

	var user *User
	Users.Range(func(uid uint64, u *User) bool {
		if u.Email != email {
			return true
		}
		user = u
		return false
	})
	if user == nil {
		Ret403(c, SEC_signin_nouser, ErrNoCred)
		return
	}
	if arg.Secret != user.Secret {
		Ret403(c, SEC_signin_deny, ErrBadPass)
		return
	}

	ret.Setup(user)
	RetOk(c, ret)
}

func SpiRefresh(c *gin.Context) {
	var ret AuthResp

	var user = c.MustGet(userKey).(*User)
	ret.Setup(user)
	RetOk(c, ret)
}
