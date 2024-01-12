package spi

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	cfg "github.com/slotopol/server/config"
)

// save server start time
var starttime = time.Now()

func SpiPing(c *gin.Context) {
	var ret = gin.H{
		"message": "pong",
	}
	Negotiate(c, http.StatusOK, ret)
}

func SpiInfo(c *gin.Context) {
	var ret = gin.H{
		"buildvers": cfg.BuildVers,
		"buildtime": cfg.BuildTime,
		"started":   starttime,
		"govers":    runtime.Version(),
		"os":        runtime.GOOS,
		"numcpu":    runtime.NumCPU(),
		"maxprocs":  runtime.GOMAXPROCS(0),
		"exepath":   cfg.ExePath,
		"cfgpath":   cfg.CfgPath,
	}
	Negotiate(c, http.StatusOK, ret)
}
