package spi

import (
	"encoding/xml"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

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
}

func RetOk(c *gin.Context, data any) {
	Negotiate(c, http.StatusOK, data)
}

type ajaxerr struct {
	XMLName xml.Name `json:"-" yaml:"-" xml:"error"`
	What    string   `json:"what" yaml:"what" xml:"what"`
	Code    int      `json:"code,omitempty" yaml:"code,omitempty" xml:"code,omitempty"`
	UID     uint64   `json:"uid,omitempty" yaml:"uid,omitempty" xml:"uid,omitempty,attr"`
}

func RetErr(c *gin.Context, status, code int, err error) {
	var claims = jwt.ExtractClaims(c)
	var uid uint64
	if v, ok := claims[identityKey]; ok {
		uid = uint64(v.(float64))
	}
	Negotiate(c, status, ajaxerr{
		What: err.Error(),
		Code: code,
		UID:  uid,
	})
}

func Ret400(c *gin.Context, code int, err error) {
	RetErr(c, http.StatusBadRequest, code, err)
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
	r.NoRoute(AuthMiddleware.MiddlewareFunc(), Handle404)
	r.GET("/ping", SpiPing)
	r.GET("/servinfo", SpiServInfo)
	r.GET("/memusage", SpiMemUsage)
	r.GET("/gamelist", SpiGameList)

	// authorization
	r.POST("/signup", SpiSignup)
	r.POST("/signin", AuthMiddleware.LoginHandler)
	r.GET("/refresh", AuthMiddleware.RefreshHandler)
	var ra = r.Group("", AuthMiddleware.MiddlewareFunc())

	//r.Use(gzip.Gzip(gzip.DefaultCompression))
	var rg = ra.Group("/game")
	rg.POST("/join", SpiGameJoin)
	rg.POST("/part", SpiGamePart)
	rg.POST("/state", SpiGameState)
	rg.POST("/bet/get", SpiGameBetGet)
	rg.POST("/bet/set", SpiGameBetSet)
	rg.POST("/sbl/get", SpiGameSblGet)
	rg.POST("/sbl/set", SpiGameSblSet)
	rg.POST("/spin", SpiGameSpin)
	rg.POST("/doubleup", SpiGameDoubleup)
	var rp = ra.Group("/prop")
	rp.POST("/wallet/get", SpiPropsWalletGet)
	rp.POST("/wallet/add", SpiPropsWalletAdd)
	var ru = ra.Group("/user")
	ru.POST("/rename", SpiUserRename)
	ru.POST("/delete", SpiUserDelete)
}
