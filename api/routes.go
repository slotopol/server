package api

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"xorm.io/xorm"

	cfg "github.com/slotopol/server/config"
)

type Session = xorm.Session

// "Server" field for HTTP headers.
var serverhdr = fmt.Sprintf("slotopol/%s (%s; %s)", cfg.BuildVers, runtime.GOOS, runtime.GOARCH)

var Offered = []string{
	binding.MIMEJSON,
	binding.MIMEXML,
	binding.MIMEYAML,
	binding.MIMETOML,
}

func Negotiate(c *gin.Context, code int, data any) {
	c.Writer.Header().Add("Server", serverhdr)
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

func Ret204(c *gin.Context) {
	c.Writer.Header().Add("Server", serverhdr)
	c.Status(http.StatusNoContent)
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

func SetupRouter(r *gin.Engine) {
	r.NoRoute(Handle404)
	r.NoMethod(Handle405)
	//r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.Any("/ping", ApiPing)
	r.GET("/servinfo", ApiServInfo)
	r.GET("/memusage", ApiMemUsage)

	// authorization
	r.Any("/signis", ApiSignis)
	r.GET("/sendcode", ApiSendCode)
	r.GET("/activate", Auth(false), ApiActivate)
	r.POST("/signup", Auth(false), ApiSignup)
	r.POST("/signin", ApiSignin)
	r.Any("/refresh", Auth(true), ApiRefresh)

	var ra = r.Group("/", Auth(true))

	// common game group
	r.GET("/game/algs", ApiGameAlgs)
	r.GET("/game/list", ApiGameList)
	var rg = ra.Group("/game")
	rg.POST("/new", ApiGameNew)
	rg.POST("/join", ApiGameJoin)
	rg.POST("/info", ApiGameInfo)
	rg.POST("/rtp/get", ApiGameRtpGet)

	// slot group
	var rs = ra.Group("/slot")
	rs.POST("/bet/get", ApiSlotBetGet)
	rs.POST("/bet/set", ApiSlotBetSet)
	rs.POST("/sel/get", ApiSlotSelGet)
	rs.POST("/sel/set", ApiSlotSelSet)
	rs.POST("/mode/set", ApiSlotModeSet)
	rs.POST("/spin", ApiSlotSpin)
	rs.POST("/doubleup", ApiSlotDoubleup)
	rs.POST("/collect", ApiSlotCollect)

	// keno group
	var rk = ra.Group("/keno")
	rk.POST("/bet/get", ApiKenoBetGet)
	rk.POST("/bet/set", ApiKenoBetSet)
	rk.POST("/sel/get", ApiKenoSelGet)
	rk.POST("/sel/set", ApiKenoSelSet)
	rk.POST("/sel/getslice", ApiKenoSelGetSlice)
	rk.POST("/sel/setslice", ApiKenoSelSetSlice)
	rk.POST("/spin", ApiKenoSpin)

	// properties group
	var rp = ra.Group("/prop")
	rp.POST("/get", ApiPropsGet)
	rp.POST("/wallet/get", ApiPropsWalletGet)
	rp.POST("/wallet/add", ApiPropsWalletAdd)
	rp.POST("/al/get", ApiPropsAlGet)
	rp.POST("/al/set", ApiPropsAlSet)
	rp.POST("/rtp/get", ApiPropsRtpGet)
	rp.POST("/rtp/set", ApiPropsRtpSet)

	// user group
	var ru = ra.Group("/user")
	ru.POST("/is", ApiUserIs)
	ru.POST("/rename", ApiUserRename)
	ru.POST("/secret", ApiUserSecret)
	ru.POST("/delete", ApiUserDelete)

	// club group
	var rc = ra.Group("/club")
	rc.POST("/list", ApiClubList)
	rc.POST("/is", ApiClubIs)
	rc.POST("/info", ApiClubInfo)
	rc.POST("/jpfund", ApiClubJpfund)
	rc.POST("/rename", ApiClubRename)
	rc.POST("/cashin", ApiClubCashin)
}
