package game

import (
	"context"
	"sort"

	"github.com/spf13/pflag"
)

const ( // Game properties
	GPsel    = 0b_0000_0000_0001 // user can change lines
	GPjack   = 0b_0000_0000_0010 // cumulative jackpot is present
	GPfgno   = 0                 // free games are absent
	GPfghas  = 0b_0000_0001_0000 // non-retriggered free games are present
	GPretrig = 0b_0000_0010_0000 // free games are present and can be retriggered
	GPfgmult = 0b_0000_0100_0000 // any multipliers on free games
	GPfgreel = 0b_0000_1000_0000 // separate reels on free games
	GPscat   = 0b_0001_0000_0000 // has scatters
	GPwild   = 0b_0010_0000_0000 // has wild symbols
	GPrwild  = 0b_0100_0000_0000 // has reel wild symbols
	GPbwild  = 0b_1000_0000_0000 // has big wild
)

type (
	GameAlias struct {
		ID   string `json:"id" yaml:"id" xml:"id"`
		Name string `json:"name" yaml:"name" xml:"name"`
	}

	GameInfo struct {
		Aliases  []GameAlias `json:"aliases" yaml:"aliases" xml:"aliases"`
		Provider string      `json:"provider" yaml:"provider" xml:"provider"`
		GP       uint        `json:"gp,omitempty" yaml:"gp,omitempty" xml:"gp,omitempty"` // game properties
		SX       int         `json:"sx" yaml:"sx" xml:"sx"`                               // screen width
		SY       int         `json:"sy" yaml:"sy" xml:"sy"`                               // screen height
		SN       int         `json:"sn,omitempty" yaml:"sn,omitempty" xml:"sn,omitempty"` // number of symbols
		LN       int         `json:"ln,omitempty" yaml:"ln,omitempty" xml:"ln,omitempty"` // number of lines
		BN       int         `json:"bn,omitempty" yaml:"bn,omitempty" xml:"bn,omitempty"` // number of bonuses
		RTP      []float64   `json:"rtp" yaml:"rtp" xml:"rtp"`
	}
)

var GameList = []*GameInfo{}

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
