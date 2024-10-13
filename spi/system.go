package spi

import (
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	cfg "github.com/slotopol/server/config"
	game "github.com/slotopol/server/game"
)

// save server start time
var starttime = time.Now()

// Check service response.
func SpiPing(c *gin.Context) {
	var ret = gin.H{
		"message": "pong",
	}
	RetOk(c, ret)
}

// Static service system information.
func SpiServInfo(c *gin.Context) {
	var ret = gin.H{
		"buildvers": cfg.BuildVers,
		"buildtime": cfg.BuildTime,
		"started":   starttime.Format(time.RFC3339),
		"govers":    runtime.Version(),
		"os":        runtime.GOOS,
		"numcpu":    runtime.NumCPU(),
		"maxprocs":  runtime.GOMAXPROCS(0),
		"exepath":   cfg.ExePath,
		"cfgpath":   cfg.CfgPath,
		"sqlpath":   cfg.SqlPath,
	}
	RetOk(c, ret)
}

// Memory usage footprint.
func SpiMemUsage(c *gin.Context) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	var ret = gin.H{
		"buildvers":     cfg.BuildVers,
		"buildtime":     cfg.BuildTime,
		"running":       time.Since(starttime) / time.Millisecond,
		"heapalloc":     mem.HeapAlloc,
		"heapsys":       mem.HeapSys,
		"totalalloc":    mem.TotalAlloc,
		"nextgc":        mem.NextGC,
		"numgc":         mem.NumGC,
		"pausetotalns":  mem.PauseTotalNs,
		"gccpufraction": mem.GCCPUFraction,
	}
	RetOk(c, ret)
}

// Returns full list of all available games by game type IDs.
func SpiGameList(c *gin.Context) {
	RetOk(c, game.GameList)
}
