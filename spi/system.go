package spi

import (
	"net/http"
	"runtime"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	cfg "github.com/slotopol/server/config"
)

// save server start time
var starttime = time.Now()

// Check service response.
func SpiPing(c *gin.Context) {
	var ret = gin.H{
		"message": "pong",
	}
	Negotiate(c, http.StatusOK, ret)
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
	}
	Negotiate(c, http.StatusOK, ret)
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
	Negotiate(c, http.StatusOK, ret)
}

// Returns full list of all available games by game type IDs.
func SpiGameList(c *gin.Context) {
	var list = make([]string, len(cfg.GameAliases))
	var i int
	for alias := range cfg.GameAliases {
		list[i] = alias
		i++
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})
	Negotiate(c, http.StatusOK, list)
}
