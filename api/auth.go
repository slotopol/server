package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/slotopol/server/core"
	"github.com/slotopol/server/util"
	"xorm.io/xorm"
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

const (
	sqlnewprops = `INSERT INTO props (cid,uid) SELECT cid,? FROM club`
)

var (
	ErrNoJwtID  = errors.New("jwt-token does not have user id")
	ErrBadUID   = errors.New("jwt-token id does not refer to registered user")
	ErrBadEmail = errors.New("jwt-token id does not refer to registered user email")
	ErrJwtOld   = errors.New("jwt-token is revocated by user password change or account update")
	ErrNoAuth   = errors.New("authorization is required")
	ErrNoScheme = errors.New("authorization does not have expected scheme")
	ErrNoSecret = errors.New("expected password or SHA256 hash on it and current time as a nonce")
	ErrSmallKey = errors.New("password too small")
	ErrNoCred   = errors.New("user with given credentials does not registered")
	ErrActivate = errors.New("activation required for this account")
	ErrOldCode  = errors.New("verification code expired")
	ErrBadCode  = errors.New("verification code does not pass")
	ErrNotPass  = errors.New("password is incorrect")
	ErrSigTime  = errors.New("signing time can not been recognized (time in RFC3339 expected)")
	ErrSigOut   = errors.New("nonce is expired")
	ErrBadHash  = errors.New("hash cannot be decoded in hexadecimal")
	ErrNoSess   = errors.New("session to refresh does not found")
	ErrTokAlien = errors.New("session to refresh is belongs to another user")
	ErrTokRevok = errors.New("token is revoked")
	ErrTokTheft = errors.New("session to refresh required from not congruent device")
	ErrTokTwice = errors.New("this token is already used to refresh")
)

// Claims of JWT-tokens. `Subject` field contains
// profile identifier or user email.
type Claims struct {
	jwt.RegisteredClaims
	SID uint64 `json:"sid,omitempty"`
}

func (c *Claims) GetUser() (*core.User, error) {
	if c.Subject == "" {
		return nil, ErrNoJwtID
	}
	var uid, err = strconv.ParseUint(c.Subject, 10, 64)
	if err == nil {
		if u, ok := core.Users.Get(uid); ok {
			return u, nil
		}
		return nil, ErrBadUID
	}
	for _, u := range core.Users.Items() {
		if c.Subject == u.Email {
			return u, nil
		}
	}
	return nil, ErrBadEmail
}

func (c *Claims) Validate() error {
	var user, err = c.GetUser()
	if err != nil {
		return err
	}
	// Its possible to revoke all tokens of user by changing
	// his password or by change vtime only.
	// Before change the password it can get new refresh token,
	// but after change the old tokens will be invalidated.
	if c.NotBefore.Add(Cfg.GracePeriod).Before(user.VTime) {
		return ErrJwtOld
	}
	return nil
}

var hints = strings.Join([]string{
	"Sec-CH-UA-Form-Factors",     // "Desktop", "Automotive", "Mobile", "Tablet", "XR", "EInk", "Watch"
	"Sec-CH-UA-Platform",         // "Android", "Chrome OS", "Chromium OS", "iOS", "Linux", "macOS", "Windows", or "Unknown"
	"Sec-CH-UA-Platform-Version", // "11.0.0"
	"Sec-CH-UA-Model",            // "Pixel 3"
	"Sec-CH-UA-Mobile",           // "?1" or "?0"
}, ", ")

func AuthHints() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Accept-CH", hints)
		c.Header("Critical-CH", hints)
		c.Next()
	}
}

// Auth is authorization middleware, sets User object associated
// with authorization to gin context. `required` parameter tells
// to continue if authorization is absent.
func Auth(required bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var code int
		var user *core.User
		if tokenstr, _ := GetTokenstr(c); tokenstr != "" {
			if user, code, err = GetBearerAuth(tokenstr); err != nil {
				Ret401(c, code, err)
				return
			}
		}
		if user == nil && Cfg.UseBasic {
			if credentials, _ := GetCredentials(c); credentials != "" {
				if user, code, err = GetBasicAuth(credentials); err != nil {
					Ret401(c, code, err)
					return
				}
			}
		}

		if user != nil {
			c.Set(userKey, user)
		} else if required {
			Ret401(c, AEC_auth_absent, ErrNoAuth)
			return
		}

		c.Next()
	}
}

func GetTokenstr(c *gin.Context) (tokenstr string, err error) {
	if tokenstr = c.GetHeader("Authorization"); tokenstr != "" {
		if strings.HasPrefix(tokenstr, "Bearer ") {
			tokenstr = tokenstr[7:]
			return
		}
	}
	if tokenstr = c.Query("token"); tokenstr != "" {
		return
	}
	if tokenstr = c.Query("jwt"); tokenstr != "" {
		return
	}
	if tokenstr, _ = c.Cookie("token"); tokenstr != "" {
		return
	}
	if tokenstr, _ = c.Cookie("jwt"); tokenstr != "" {
		return
	}
	if tokenstr = c.PostForm("token"); tokenstr != "" {
		return
	}
	if tokenstr = c.PostForm("jwt"); tokenstr != "" {
		return
	}
	return "", io.EOF
}

func GetCredentials(c *gin.Context) (credentials string, err error) {
	if credentials = c.GetHeader("Authorization"); credentials != "" {
		if strings.HasPrefix(credentials, "Basic ") {
			credentials = credentials[6:]
			return
		}
	}
	if credentials = c.Query("cred"); credentials != "" {
		return
	}
	if credentials, _ = c.Cookie("cred"); credentials != "" {
		return
	}
	if credentials = c.PostForm("cred"); credentials != "" {
		return
	}
	return "", io.EOF
}

func GetBasicAuth(credentials string) (user *core.User, code int, err error) {
	var decoded []byte
	if decoded, err = base64.RawURLEncoding.DecodeString(credentials); err != nil {
		return nil, AEC_basic_decode, err
	}
	var parts = strings.Split(util.B2S(decoded), ":")

	var email = util.ToLower(parts[0])
	for _, u := range core.Users.Items() {
		if u.Email == email {
			user = u
			break
		}
	}
	if user == nil {
		err, code = ErrNoCred, AEC_basic_nouser
		return
	}
	if user.Secret != parts[1] {
		err, code = ErrNotPass, AEC_basic_deny
		return
	}
	return
}

func GetBearerClaims(tokenstr string, refresh bool) (claims Claims, user *core.User, code int, err error) {
	_, err = jwt.ParseWithClaims(tokenstr, &claims, func(*jwt.Token) (any, error) {
		if refresh {
			return util.S2B(Cfg.RefreshKey), nil
		} else {
			return util.S2B(Cfg.AccessKey), nil
		}
	},
		jwt.WithNotBeforeRequired(),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(jwtIssuer),
		jwt.WithLeeway(5*time.Second))

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			code = AEC_token_malform
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			code = AEC_token_notsign
		case errors.Is(err, jwt.ErrTokenInvalidClaims):
			code = AEC_token_badclaims
		case errors.Is(err, jwt.ErrTokenExpired):
			code = AEC_token_expired
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			code = AEC_token_notyet
		case errors.Is(err, jwt.ErrTokenInvalidIssuer):
			code = AEC_token_issuer
		default:
			code = AEC_token_error
		}
		return
	}
	if user, err = claims.GetUser(); err != nil {
		code = AEC_token_nouser
	}
	return
}

func GetBearerAuth(tokenstr string) (user *core.User, code int, err error) {
	_, user, code, err = GetBearerClaims(tokenstr, false)
	return
}

type AuthResp struct {
	XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
	UID     uint64   `json:"uid" yaml:"uid" xml:"uid"`          // user identifier
	SID     uint64   `json:"sid" yaml:"sid" xml:"sid"`          // session identifier
	Email   string   `json:"email" yaml:"email" xml:"email"`    // user email
	Access  string   `json:"access" yaml:"access" xml:"access"` // access token
	Refrsh  string   `json:"refrsh" yaml:"refrsh" xml:"refrsh"` // refresh token
	Expire  string   `json:"expire" yaml:"expire" xml:"expire"` // access token expiration time
	Period  string   `json:"period" yaml:"period" xml:"period"` // period until new tokens pair can be get
}

func (r *AuthResp) Setup(user *core.User, sid uint64, now time.Time) {
	var err error
	var token *jwt.Token
	var exp = now.Add(Cfg.AccessTTL)
	var age = now.Add(Cfg.RefreshTTL)
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			Issuer:    jwtIssuer,
			Subject:   strconv.FormatUint(user.UID, 10),
		},
		SID: sid,
	})
	if r.Access, err = token.SignedString(util.S2B(Cfg.AccessKey)); err != nil {
		panic(err)
	}
	r.Expire = exp.Format(time.RFC3339)
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(age),
			Issuer:    jwtIssuer,
			Subject:   strconv.FormatUint(user.UID, 10),
		},
		SID: sid,
	})
	if r.Refrsh, err = token.SignedString(util.S2B(Cfg.RefreshKey)); err != nil {
		panic(err)
	}
	r.Period = age.Format(time.RFC3339)
	r.UID = user.UID
	r.SID = sid
	r.Email = user.Email
}

func sendcode(name, email string, code uint32) (err error) {
	type person struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email"`
	}
	type content struct {
		Sender  person   `json:"sender"`
		ReplyTo person   `json:"replyTo"`
		To      []person `json:"to"`
		Subject string   `json:"subject,omitempty"`
		Html    string   `json:"htmlContent,omitempty"`
		Text    string   `json:"textContent,omitempty"`
		Tags    []string `json:"tags,omitempty"`
	}

	const ct = "application/json"
	var m = content{
		Sender: person{
			Name:  Cfg.SenderName,
			Email: Cfg.SenderEmail,
		},
		ReplyTo: person{
			Email: Cfg.ReplytoEmail,
		},
		To: []person{
			{
				Name:  name,
				Email: email,
			},
		},
		Subject: Cfg.EmailSubject,
		Html:    fmt.Sprintf(Cfg.EmailHtmlContent, code),
	}

	var body []byte
	if body, err = json.Marshal(m); err != nil {
		return err
	}

	var req *http.Request
	if req, err = http.NewRequest("POST", Cfg.BrevoEmailEndpoint, bytes.NewReader(body)); err != nil {
		return err
	}
	req.Header.Set("api-key", Cfg.BrevoApiKey)
	req.Header.Set("Content-Type", ct)
	req.Header.Set("Accept", ct)

	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
		return err
	}
	defer resp.Body.Close()

	if body, err = io.ReadAll(resp.Body); err != nil {
		return err
	}
	var rm map[string]string
	if err = json.Unmarshal(body, &rm); err != nil {
		return err
	}
	if msg, ok := rm["message"]; ok {
		return errors.New(msg)
	}
	return
}

// Check whether email address is available for registration.
// Returns status 200, or 404.
func ApiSignis(c *gin.Context) {
	var err error
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		Email   string   `json:"email" yaml:"email" xml:"email" form:"email" binding:"omitempty,email"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid"`
		Email   string   `json:"email" yaml:"email" xml:"email"`
		Name    string   `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_signis_nobind, err)
		return
	}

	var email = util.ToLower(arg.Email)
	for _, user := range core.Users.Items() {
		if user.Email == email {
			ret.UID = user.UID
			ret.Email = email
			ret.Name = user.Name
			break
		}
	}
	if ret.UID == 0 {
		Ret404(c, AEC_signis_noemail, ErrNoCred)
		return
	}

	RetOk(c, ret)
}

func ApiSendCode(c *gin.Context) {
	var err error
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid" form:"uid" binding:"required_without=Email"`
		Email   string   `json:"email" yaml:"email" xml:"email" form:"email" binding:"omitempty,email"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_sendcode_nobind, err)
		return
	}

	var email = util.ToLower(arg.Email)

	var user *core.User
	if arg.UID != 0 {
		user, _ = core.Users.Get(arg.UID)
	} else {
		for _, u := range core.Users.Items() {
			if u.Email == email {
				user = u
				break
			}
		}
	}
	if user == nil {
		Ret403(c, AEC_sendcode_nouser, ErrNoCred)
		return
	}

	var code = rand.N[uint32](1000000) // 6 digits

	var u = map[string]interface{}{
		"code":  code,
		"utime": time.Now(),
	}
	if _, err = core.XormStorage.Table("user").Where("uid=?", user.UID).Update(u); err != nil {
		Ret500(c, AEC_sendcode_update, err)
		return
	}

	if err = sendcode(user.Name, user.Email, code); err != nil {
		Ret500(c, AEC_sendcode_code, err)
		return
	}

	user.Code = code

	Ret204(c)
}

func ApiActivate(c *gin.Context) {
	var err error
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid" form:"uid" binding:"required_without=Email"`
		Email   string   `json:"email" yaml:"email" xml:"email" form:"email" binding:"omitempty,email"`
		Code    uint32   `json:"code,omitempty" yaml:"code,omitempty" xml:"code,omitempty" form:"code"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_activate_nobind, err)
		return
	}

	var email = util.ToLower(arg.Email)

	var user *core.User
	if arg.UID != 0 {
		user, _ = core.Users.Get(arg.UID)
	} else {
		for _, u := range core.Users.Items() {
			if u.Email == email {
				user = u
				break
			}
		}
	}
	if user == nil {
		Ret403(c, AEC_activate_nouser, ErrNoCred)
		return
	}

	if _, al := core.GetAdmin(c, 0); al&core.ALadmin == 0 {
		if time.Since(user.UTime) > Cfg.CodeTimeout {
			Ret403(c, AEC_activate_oldcode, ErrOldCode)
			return
		}
		if arg.Code != user.Code {
			Ret403(c, AEC_activate_badcode, ErrBadCode)
			return
		}
	}

	var u = map[string]interface{}{
		"status": user.Status | core.UFactivated,
		"utime":  time.Now(),
	}
	if _, err = core.XormStorage.Table("user").Where("uid=?", user.UID).Update(u); err != nil {
		Ret500(c, AEC_activate_update, err)
		return
	}

	user.Status |= core.UFactivated

	Ret204(c)
}

func ApiSignup(c *gin.Context) {
	var err error
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		Email   string   `json:"email" yaml:"email" xml:"email" form:"email" binding:"required,email"`
		Secret  string   `json:"secret" yaml:"secret" xml:"secret" form:"secret" binding:"required"`
		Name    string   `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid"`
		Email   string   `json:"email" yaml:"email" xml:"email"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_signup_nobind, err)
		return
	}
	if len(arg.Secret) < 6 {
		Ret400(c, AEC_signup_smallsec, ErrSmallKey)
		return
	}

	var email = util.ToLower(arg.Email)

	var code uint32
	var status core.UF
	if _, al := core.GetAdmin(c, 0); al&core.ALadmin != 0 || !Cfg.UseActivation {
		status = core.UFactivated
	} else {
		code = rand.N[uint32](1000000) // 6 digits
		if err = sendcode(arg.Name, email, code); err != nil {
			Ret500(c, AEC_signup_code, err)
			return
		}
	}

	var user = &core.User{
		VTime:  time.Now(),
		Email:  email,
		Secret: arg.Secret,
		Name:   arg.Name,
		Code:   code,
		Status: status,
	}
	if err = core.SafeTransaction(core.XormStorage, func(session *xorm.Session) (err error) {
		if _, err = session.Insert(user); err != nil {
			return
		}

		if _, err = session.Exec(sqlnewprops, user.UID); err != nil {
			return
		}

		var props = make([]core.Props, len(core.PropMaster))
		copy(props, core.PropMaster)
		for i := range props {
			props[i].UID = user.UID
			if _, err = session.Where("cid=? AND uid=?", props[i].CID, user.UID).Update(&props[i]); err != nil {
				return
			}
		}

		user.Init()
		for cid := range core.Clubs.Items() {
			user.InsertProps(&core.Props{
				CID: cid,
				UID: user.UID,
			})
		}
		for i := range props {
			user.InsertProps(&props[i])
		}
		core.Users.Set(user.UID, user)
		return
	}); err != nil {
		Ret500(c, AEC_signup_sql, err)
		return
	}

	ret.UID = user.UID
	ret.Email = email
	RetOk(c, ret)
}

func ApiSignin(c *gin.Context) {
	var err error
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid" form:"uid" binding:"required_without=Email"`
		Email   string   `json:"email" yaml:"email" xml:"email" form:"email" binding:"omitempty,email"`
		Secret  string   `json:"secret" yaml:"secret,omitempty" xml:"secret,omitempty" form:"secret"`
		HS256   string   `json:"hs256,omitempty" yaml:"hs256,omitempty" xml:"hs256,omitempty" form:"hs256"`
		SigTime string   `json:"sigtime,omitempty" yaml:"sigtime,omitempty" xml:"sigtime,omitempty" form:"sigtime"`
		Code    uint32   `json:"code,omitempty" yaml:"code,omitempty" xml:"code,omitempty" form:"code"`
	}
	var ret AuthResp

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_signin_nobind, err)
		return
	}
	if len(arg.SigTime) == 0 && len(arg.Secret) == 0 {
		Ret400(c, AEC_signin_nosecret, ErrNoSecret)
		return
	}
	if len(arg.Secret) > 0 && len(arg.Secret) < 6 {
		Ret400(c, AEC_signin_smallsec, ErrSmallKey)
		return
	}

	var email = util.ToLower(arg.Email)

	var user *core.User
	if arg.UID != 0 {
		user, _ = core.Users.Get(arg.UID)
	} else {
		for _, u := range core.Users.Items() {
			if u.Email == email {
				user = u
				break
			}
		}
	}
	if user == nil {
		Ret403(c, AEC_signin_nouser, ErrNoCred)
		return
	}

	if user.Status&core.UFactivated == 0 {
		Ret403(c, AEC_signin_activate, ErrActivate)
		return
	}

	if user.Status&core.UFsigncode != 0 {
		if time.Since(user.UTime) > Cfg.CodeTimeout {
			Ret403(c, AEC_signin_oldcode, ErrOldCode)
			return
		}
		if arg.Code != user.Code {
			Ret403(c, AEC_signin_badcode, ErrBadCode)
			return
		}
	}

	if len(arg.Secret) > 0 {
		if arg.Secret != user.Secret {
			Ret403(c, AEC_signin_denypass, ErrNotPass)
			return
		}
	} else {
		var sigtime time.Time
		if sigtime, err = time.Parse(time.RFC3339, arg.SigTime); err != nil {
			Ret400(c, AEC_signin_sigtime, ErrSigTime)
			return
		}
		if time.Since(sigtime) > Cfg.NonceTimeout {
			Ret403(c, AEC_signin_timeout, ErrSigOut)
			return
		}

		var hs256 []byte
		if hs256, err = hex.DecodeString(arg.HS256); err != nil {
			Ret400(c, AEC_signin_hs256, ErrBadHash)
			return
		}

		var h = hmac.New(sha256.New, util.S2B(arg.SigTime))
		h.Write(util.S2B(user.Secret))
		var master = h.Sum(nil)
		if !hmac.Equal(master, hs256) {
			Ret403(c, AEC_signin_denyhash, ErrNotPass)
			return
		}
	}

	var sess core.Session
	sess.UID = user.UID
	sess.Status = 1
	sess.Device.Setup(c)
	if _, err = core.XormStorage.Table("session").Insert(&sess); err != nil {
		Ret500(c, AEC_signin_session, err)
		return
	}

	ret.Setup(user, sess.SID, sess.CTime)
	RetOk(c, ret)
}

func ApiToken(c *gin.Context) {
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		Refrsh  string   `json:"refrsh" yaml:"refrsh" xml:"refrsh"`
	}
	var ret AuthResp

	if err := c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_token_nobind, err)
		return
	}

	var claims, user, code, err = GetBearerClaims(arg.Refrsh, true)
	if err != nil {
		Ret401(c, code, err)
		return
	}

	var sess core.Session
	var ok bool
	sess.SID = claims.SID
	if ok, err = core.XormStorage.Get(&sess); err != nil {
		Ret500(c, AEC_token_session, err)
		return
	}
	if !ok {
		Ret401(c, AEC_token_nosess, ErrNoSess)
		return
	}
	if sess.UID != user.UID {
		Ret401(c, AEC_token_alien, ErrTokAlien)
		return
	}
	if sess.Status == 0 {
		Ret401(c, AEC_token_revoked, ErrTokRevok)
		return
	}

	var d2 core.Device
	d2.Setup(c)
	if !sess.Device.IsCongruent(&d2) {
		Ret401(c, AEC_token_device, ErrTokTheft)
		return
	}

	if sess.UTime.Sub(claims.NotBefore.Time) > Cfg.GracePeriod {
		Ret401(c, AEC_token_already, ErrTokTwice)
		return
	}

	var now = time.Now()
	if now.Sub(sess.UTime) > Cfg.GracePeriod { // prevent refresh race between browser tabs
		sess.UTime = now
		if _, err = core.XormStorage.ID(sess.SID).Cols("utime").Update(&sess); err != nil {
			Ret500(c, AEC_token_update, err)
			return
		}
	}

	ret.Setup(user, sess.SID, sess.UTime)
	RetOk(c, ret)
}

// Sign-out session. Revocation could be done by the user
// for any own session, or by admin for any session at all.
func ApiSignout(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		SID     uint64   `json:"sid,omitempty" yaml:"sid,omitempty" xml:"sid,omitempty,attr" form:"sid"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_signout_nobind, err)
		return
	}

	if arg.SID == 0 {
		var tokenstr string
		if tokenstr, err = GetTokenstr(c); err != nil {
			Ret401(c, AEC_signout_notoken, err)
			return
		}
		var claims, _, code, err = GetBearerClaims(tokenstr, false)
		if err != nil {
			Ret401(c, code, err)
			return
		}
		arg.SID = claims.SID
	}

	var sess core.Session
	sess.SID = arg.SID
	if ok, err = core.XormStorage.Get(&sess); err != nil {
		Ret500(c, AEC_signout_session, err)
		return
	}
	if !ok {
		Ret401(c, AEC_signout_nosess, ErrNoSess)
		return
	}
	if sess.Status == 0 {
		Ret401(c, AEC_signout_revoked, ErrTokRevok)
		return
	}

	var user *core.User
	if user, ok = core.Users.Get(sess.UID); !ok {
		Ret404(c, AEC_signout_nouser, ErrNoUser)
		return
	}

	var admin, al = core.MustAdmin(c, 0)
	if admin != user && al&core.ALadmin == 0 {
		Ret403(c, AEC_signout_noaccess, ErrNoAccess)
		return
	}

	var now = time.Now()
	sess.Status = 0
	sess.UTime = now
	if _, err = core.XormStorage.ID(sess.SID).Cols("status", "utime").Update(&sess); err != nil {
		Ret500(c, AEC_signout_update, err)
		return
	}

	Ret204(c)
}

func ApiSessions(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		UID     uint64   `json:"uid,omitempty" yaml:"uid,omitempty" xml:"uid,omitempty,attr" form:"uid"`
	}
	var ret struct {
		Sessions []core.Session `json:"sessions" yaml:"sessions" xml:"sessions"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_sessions_nobind, err)
		return
	}

	var user *core.User
	if arg.UID == 0 {
		var tokenstr string
		if tokenstr, err = GetTokenstr(c); err != nil {
			Ret401(c, AEC_sessions_notoken, err)
			return
		}
		var claims, _, code, err = GetBearerClaims(tokenstr, false)
		if err != nil {
			Ret401(c, code, err)
			return
		}
		if user, err = claims.GetUser(); err != nil {
			Ret401(c, AEC_sessions_claims, err)
			return
		}
	} else if user, ok = core.Users.Get(arg.UID); !ok {
		Ret404(c, AEC_sessions_nouser, ErrNoUser)
		return
	}

	var admin, al = core.MustAdmin(c, 0)
	if admin != user && al&core.ALadmin == 0 {
		Ret403(c, AEC_sessions_noaccess, ErrNoAccess)
		return
	}

	var from = time.Now().Add(-Cfg.RefreshTTL)
	if err = core.XormStorage.Where("uid=? AND utime>? AND status>0", arg.UID, from.Format(time.DateTime)).Find(&ret.Sessions); err != nil {
		Ret500(c, AEC_sessions_get, err)
		return
	}

	RetOk(c, ret)
}
