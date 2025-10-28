package api

import (
	"errors"
	"os"
	"runtime"
	"strings"
	"time"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/util"

	"github.com/gin-gonic/gin"
	"github.com/klauspost/cpuid/v2"
	"github.com/schwarzlichtbezirk/go-disk-usage/du"
)

var (
	ErrDiskInfo = errors.New("can not obtain disk info")
)

func isRunningInContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err != nil {
		return false // File does not exist, not running in Docker or an error occurred
	}
	return true // File exists, likely running in Docker
}

// save server start time
var starttime = time.Now()

// cached service info response
var srvinfo gin.H // lazy init

// Check service response.
func ApiPing(c *gin.Context) {
	RetOk(c, nil)
}

// Static service system information.
func ApiServInfo(c *gin.Context) {
	if srvinfo == nil {
		srvinfo = gin.H{
			// this service
			"buildvers": cfg.BuildVers,
			"buildtime": cfg.BuildTime,
			"started":   starttime.Format(time.RFC3339),
			// Go version & OS
			"govers":   runtime.Version(),
			"os":       runtime.GOOS,
			"arch":     runtime.GOARCH,
			"indocker": isRunningInContainer(),
			"maxprocs": runtime.GOMAXPROCS(0),
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
	}
	RetOk(c, srvinfo)
}

// Memory usage footprint.
func ApiMemUsage(c *gin.Context) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	var ret = gin.H{
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

func ApiDiskUsage(c *gin.Context) {
	var ds = du.NewDiskUsage(cfg.SqlPath)
	if ds == nil {
		Ret500(c, AEC_diskusage_nil, ErrDiskInfo)
		return
	}
	var ret = gin.H{
		"size":      ds.Size(),
		"used":      ds.Used(),
		"free":      ds.Free(),
		"available": ds.Available(),
		"usage":     ds.Usage(),
	}
	RetOk(c, ret)
}
