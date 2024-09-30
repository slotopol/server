package links

import (
	"context"
	"sort"

	"github.com/spf13/pflag"
)

const (
	FGno     = 0 // free games are absent
	FGhas    = 1 // free games are present
	FGretrig = 2 // free games are present and can be retriggered
)

type (
	GameAlias struct {
		ID   string `json:"id" yaml:"id" xml:"id"`
		Name string `json:"name" yaml:"name" xml:"name"`
	}

	GameInfo struct {
		Aliases  []GameAlias `json:"aliases" yaml:"aliases" xml:"aliases"`
		Provider string      `json:"provider" yaml:"provider" xml:"provider"`
		SX       int         `json:"sx" yaml:"sx" xml:"sx"`                               // screen width
		SY       int         `json:"sy" yaml:"sy" xml:"sy"`                               // screen height
		LN       int         `json:"ln,omitempty" yaml:"ln,omitempty" xml:"ln,omitempty"` // number of lines
		FG       int         `json:"fg,omitempty" yaml:"fg,omitempty" xml:"fg,omitempty"` // free games type
		BN       int         `json:"bn,omitempty" yaml:"bn,omitempty" xml:"bn,omitempty"` // number of bonuses
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
