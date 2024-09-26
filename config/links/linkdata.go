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
		SX       int         `json:"sx" yaml:"sx" xml:"sx"` // screen width
		SY       int         `json:"sy" yaml:"sy" xml:"sy"` // screen height
		LN       int         `json:"ln" yaml:"ln" xml:"ln"` // number of lines
		RTP      []float64   `json:"rtp" yaml:"rtp" xml:"rtp"`
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
