package spi

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"xorm.io/xorm"
)

type Session = xorm.Session

var Offered = []string{
	binding.MIMEJSON,
	binding.MIMEXML,
	binding.MIMEYAML,
	binding.MIMETOML,
}

func Negotiate(c *gin.Context, code int, data any) {
	switch c.NegotiateFormat(Offered...) {
	case binding.MIMEJSON:
		c.JSON(code, data)
	case binding.MIMEXML:
		c.XML(code, data)
	case binding.MIMEYAML:
		c.YAML(code, data)
	case binding.MIMETOML:
		c.TOML(code, data)
	default:
		c.JSON(code, data)
	}
	c.Abort()
}

func RetOk(c *gin.Context, data any) {
	Negotiate(c, http.StatusOK, data)
}

type jerr struct {
	error
}

// Unwrap returns inherited error object.
func (err jerr) Unwrap() error {
	return err.error
}

// MarshalJSON is standard JSON interface implementation to stream errors on Ajax.
func (err jerr) MarshalJSON() ([]byte, error) {
	return json.Marshal(err.Error())
}

// MarshalYAML is YAML marshaler interface implementation to stream errors on Ajax.
func (err jerr) MarshalYAML() (any, error) {
	return err.Error(), nil
}

// MarshalXML is XML marshaler interface implementation to stream errors on Ajax.
func (err jerr) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(err.Error(), start)
}

type ajaxerr struct {
	XMLName xml.Name `json:"-" yaml:"-" xml:"error"`
	What    jerr     `json:"what" yaml:"what" xml:"what"`
	Code    int      `json:"code,omitempty" yaml:"code,omitempty" xml:"code,omitempty"`
	UID     uint64   `json:"uid,omitempty" yaml:"uid,omitempty" xml:"uid,omitempty,attr"`
}

func (err ajaxerr) Error() string {
	return fmt.Sprintf("what: %s, code: %d", err.What, err.Code)
}

func (err ajaxerr) Unwrap() error {
	return err.What.error
}

func RetErr(c *gin.Context, status, code int, err error) {
	var uid uint64
	if uv, ok := c.Get(userKey); ok {
		uid = uv.(*User).UID
	}
	Negotiate(c, status, ajaxerr{
		What: jerr{err},
		Code: code,
		UID:  uid,
	})
}

func Ret400(c *gin.Context, code int, err error) {
	RetErr(c, http.StatusBadRequest, code, err)
}

func Ret401(c *gin.Context, code int, err error) {
	c.Writer.Header().Add("WWW-Authenticate", realmBasic)
	c.Writer.Header().Add("WWW-Authenticate", realmBearer)
	RetErr(c, http.StatusUnauthorized, code, err)
}

func Ret403(c *gin.Context, code int, err error) {
	RetErr(c, http.StatusForbidden, code, err)
}

func Ret404(c *gin.Context, code int, err error) {
	RetErr(c, http.StatusNotFound, code, err)
}

func Ret500(c *gin.Context, code int, err error) {
	RetErr(c, http.StatusInternalServerError, code, err)
}

func Router(r *gin.Engine) {
	r.NoRoute(Auth(false), Handle404)
	r.GET("/ping", SpiPing)
	r.GET("/servinfo", SpiServInfo)
	r.GET("/memusage", SpiMemUsage)
	r.GET("/gamelist", SpiGameList)

	// authorization
	r.GET("/signis", SpiSignis)
	r.POST("/signup", SpiSignup)
	r.POST("/signin", SpiSignin)
	r.GET("/refresh", Auth(true), SpiRefresh)
	var ra = r.Group("/", Auth(true))

	//r.Use(gzip.Gzip(gzip.DefaultCompression))
	var rg = ra.Group("/game")
	rg.POST("/join", SpiGameJoin)
	rg.POST("/part", SpiGamePart)
	rg.POST("/state", SpiGameState)
	rg.POST("/bet/get", SpiGameBetGet)
	rg.POST("/bet/set", SpiGameBetSet)
	rg.POST("/sbl/get", SpiGameSblGet)
	rg.POST("/sbl/set", SpiGameSblSet)
	rg.POST("/reels/get", SpiGameReelsGet)
	rg.POST("/reels/set", SpiGameReelsSet)
	rg.POST("/spin", SpiGameSpin)
	rg.POST("/doubleup", SpiGameDoubleup)
	rg.POST("/collect", SpiGameCollect)
	var rp = ra.Group("/prop")
	rp.POST("/wallet/get", SpiPropsWalletGet)
	rp.POST("/wallet/add", SpiPropsWalletAdd)
	var ru = ra.Group("/user")
	ru.POST("/rename", SpiUserRename)
	ru.POST("/secret", SpiUserSecret)
	ru.POST("/delete", SpiUserDelete)
	var rc = ra.Group("/club")
	rc.POST("/rename", SpiClubRename)
	rc.POST("/cashin", SpiClubCashin)
}
