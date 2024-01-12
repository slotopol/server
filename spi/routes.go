package spi

import (
	"encoding/xml"
	"net/http"

	"github.com/gin-contrib/gzip"
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
}

func Ret400(c *gin.Context, code int, err error) {
	Negotiate(c, http.StatusBadRequest, ajaxerr{
		What: err.Error(),
		Code: code,
	})
}

func Ret500(c *gin.Context, code int, err error) {
	Negotiate(c, http.StatusInternalServerError, ajaxerr{
		What: err.Error(),
		Code: code,
	})
}

func Router() *gin.Engine {
	var r = gin.New()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.GET("/ping", SpiPing)
	r.GET("/info", SpiInfo)
	var rg = r.Group("/game")
	rg.POST("/join", SpiGameJoin)
	rg.POST("/part", SpiGamePart)
	return r
}
