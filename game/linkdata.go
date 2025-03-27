package game

import (
	"context"
	"math"
	"sort"

	"github.com/slotopol/server/util"
)

type GP uint

const ( // Game properties
	GPfgno GP = 0 // free games are absent

	GPsel   GP = 1 << iota // user can change lines
	GPrline                // pays left to right and right to left
	GPcline                // pays for combination at any position
	GPjack                 // cumulative jackpot is present

	GPcsc    // cascade falls present
	GPcmult  // multipliers on cascade falls
	GPfghas  // non-retriggered free games are present
	GPretrig // free games are present and can be retriggered

	GPfgreel // separate reels on free games
	GPfgmult // any multipliers on free games
	GPrmult  // any multipliers on regular games
	GPscat   // has scatters

	GPwild  // has wild symbols
	GPrwild // has reel wild symbols
	GPbwild // has big wild (3x3)
	GPwmult // has multiplier on wilds

	GPbsym // has big symbol (usually 3x3 in the center on free games)
	GPfill // has multiplier on filled screen
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

var InfoList = []*GameInfo{}
var InfoMap = map[string]*GameInfo{}
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

func (gi *GameInfo) SetupFactory(game func() any, scan Scanner) {
	InfoList = append(InfoList, gi)
	for _, ga := range gi.Aliases {
		var aid = util.ToID(ga.Prov + "/" + ga.Name)
		InfoMap[aid] = gi
		GameFactory[aid] = game
		ScanFactory[aid] = scan // can be nil
	}
}

func (gi *GameInfo) FindClosest(mrtp float64) (rtp float64) {
	rtp = -1000 // lets to get first reels from map in any case
	for _, p := range gi.RTP {
		if math.Abs(mrtp-p) < math.Abs(mrtp-rtp) {
			rtp = p
		}
	}
	return
}
