package game

import (
	"context"
	"sort"
)

type GP uint

const ( // Game properties
	GPsel    GP = 0b_0000_0000_0001 // user can change lines
	GPjack   GP = 0b_0000_0000_0010 // cumulative jackpot is present
	GPfgno   GP = 0                 // free games are absent
	GPfghas  GP = 0b_0000_0000_0100 // non-retriggered free games are present
	GPretrig GP = 0b_0000_0000_1000 // free games are present and can be retriggered
	GPfgmult GP = 0b_0000_0001_0000 // any multipliers on free games
	GPfgreel GP = 0b_0000_0010_0000 // separate reels on free games
	GPscat   GP = 0b_0000_0100_0000 // has scatters
	GPwild   GP = 0b_0000_1000_0000 // has wild symbols
	GPrwild  GP = 0b_0001_0000_0000 // has reel wild symbols
	GPbwild  GP = 0b_0010_0000_0000 // has big wild (3x3)
	GPbsym   GP = 0b_1000_0000_0000 // has big symbol (usually 3x3 in the center on free games)
)

type (
	GameAlias struct {
		Prov string `json:"prov" yaml:"prov" xml:"prov"`
		Name string `json:"name" yaml:"name" xml:"name"`
	}

	GameInfo struct {
		Aliases []GameAlias `json:"aliases" yaml:"aliases" xml:"aliases"`
		GP      GP          `json:"gp,omitempty" yaml:"gp,omitempty" xml:"gp,omitempty"` // game properties
		SX      int         `json:"sx,omitempty" yaml:"sx,omitempty" xml:"sx,omitempty"` // screen width
		SY      int         `json:"sy,omitempty" yaml:"sy,omitempty" xml:"sy,omitempty"` // screen height
		SN      int         `json:"sn,omitempty" yaml:"sn,omitempty" xml:"sn,omitempty"` // number of symbols
		LN      int         `json:"ln,omitempty" yaml:"ln,omitempty" xml:"ln,omitempty"` // number of lines
		BN      int         `json:"bn,omitempty" yaml:"bn,omitempty" xml:"bn,omitempty"` // number of bonuses
		RTP     []float64   `json:"rtp" yaml:"rtp" xml:"rtp"`                            // 'Return to Player' percents list
	}

	Scanner func(context.Context, float64) float64
)

var GameList = []*GameInfo{}
var GameFactory = map[string]func() any{}
var ScanFactory = map[string]Scanner{}

func MakeRtpList[T any](reelsmap map[float64]T) []float64 {
	var list = make([]float64, 0, len(reelsmap))
	for rtp := range reelsmap {
		list = append(list, rtp)
	}
	sort.Float64s(list)
	return list
}
