package api

import (
	"os"
	"runtime"
	"strings"
	"time"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/util"

	"github.com/gin-gonic/gin"
	"github.com/klauspost/cpuid/v2"
)

func isRunningInContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true // File exists, likely running in Docker
	}
	return false // File does not exist, not running in Docker or an error occurred
}

// save server start time
var starttime = time.Now()

// save if running in container
var indocker = isRunningInContainer()

// Check service response.
func ApiPing(c *gin.Context) {
	RetOk(c, nil)
}

// Static service system information.
func ApiServInfo(c *gin.Context) {
	var ret = gin.H{
		// this service
		"buildvers": cfg.BuildVers,
		"buildtime": cfg.BuildTime,
		"started":   starttime.Format(time.RFC3339),
		// Go version & OS
		"govers":   runtime.Version(),
		"os":       runtime.GOOS,
		"arch":     runtime.GOARCH,
		"maxprocs": runtime.GOMAXPROCS(0),
		"indocker": indocker,
		// CPU
		"cpubrand": cpuid.CPU.BrandName,
		"cpuvend":  cpuid.CPU.VendorString,
		"cpuphys":  cpuid.CPU.PhysicalCores,
		"cpulogic": cpuid.CPU.LogicalCores,
		"cpufreq":  cpuid.CPU.Hz,
		"cpufeat":  strings.Join(cpuid.CPU.FeatureSet(), ","),
		// paths
		"exepath": util.ToSlash(cfg.ExePath),
		"cfgpath": util.ToSlash(cfg.CfgPath),
		"sqlpath": util.ToSlash(cfg.SqlPath),
	}
	RetOk(c, ret)
}

// Memory usage footprint.
func ApiMemUsage(c *gin.Context) {
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
