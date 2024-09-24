package links

import (
	"context"
	"sort"

	"github.com/spf13/pflag"
)

type (
	GameAlias struct {
		ID   string `json:"id" yaml:"id" xml:"id"`
		Name string `json:"name" yaml:"name" xml:"name"`
	}

	GameInfo struct {
		Aliases  []GameAlias `json:"aliases" yaml:"aliases" xml:"aliases"`
		Provider string      `json:"provider" yaml:"provider" xml:"provider"`
		ScrnX    int         `json:"scrnx" yaml:"scrnx" xml:"scrnx"`
		ScrnY    int         `json:"scrny" yaml:"scrny" xml:"scrny"`
		RtpList  []float64   `json:"rtplist" yaml:"rtplist" xml:"rtplist"`
	}
)

var GameList = []GameInfo{}

var FlagsSetters = []func(*pflag.FlagSet){}

var ScanIters = []func(*pflag.FlagSet, context.Context){}

var GameFactory = map[string]func() any{}

func MakeRtpList[T any](reelsmap map[float64]T) []float64 {
	var list = make([]float64, 0, len(reelsmap))
	for rtp := range reelsmap {
		list = append(list, rtp)
	}
	sort.Float64s(list)
	return list
}
