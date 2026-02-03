package game

import (
	"context"
	"fmt"
	"log"
	"math"
	"sort"

	"github.com/slotopol/server/util"
)

type Gamble interface {
	Spin(float64)         // fill the grid with random hits on reels closest to given RTP, constant function
	GetBet() float64      // returns current bet, constant function
	SetBet(float64) error // set bet to given value
}

type GT uint // Game type

const (
	GTslot GT = 1 + iota
	GTkeno
)

type GP uint // Game properties

const (
	GPlpay GP = 1<<(iota+1) - 1 // pays left to right
	GPrpay                      // pays left to right and right to left
	GPapay                      // pays for combination at any position
	GPcpay GP = 1 << iota       // pays by clusters

	GPlsel // user can select bet lines
	GPwsel // user can change ways set
	_
	_

	GPjack  // cumulative jackpot is present
	GPfill  // has multiplier on filled grid
	GPmix   // has pays by combinations with mixed symbols (non-wilds)
	GPcumul // has cumulative pays (in this case spin can fails if bank have not enough money)

	GPbmode // has non-reels bonus mode
	GPpick  // has game mode depending on the user's choice
	GPcasc  // cascade falls present
	GPcmult // multipliers on cascade avalanche

	GPfgonce // non-retriggered free games are present
	GPfgtwic // free games that can be retriggered only once
	GPfgseq  // free games that can be retriggered in a sequence
	GPfgreel // separate reels on free games

	GPrgmult // any multipliers on regular games
	GPfgmult // any multipliers on free games
	GPscat   // has scatters
	GPstscat // has stacked scatters (several may appear on the reel)

	GPwsc   // has wild/scatters symbols
	GPwild  // has wild symbols
	GPrwild // has reel wild symbols
	GPbwild // has big (3x3) wild symbols

	GPewild // has expanding wilds
	GPwturn // symbols turns to wilds
	GPwmult // has multiplier on wilds
	GPbsym  // has big symbol reel (usually with 3x3 in the center on free games)

	// free games are absent
	GPfgno GP = 0
	// any free games
	GPfgany GP = GPfgonce | GPfgtwic | GPfgseq
	// any wild symbols
	GPwany GP = GPwsc | GPwild | GPrwild | GPbwild | GPewild
	// any scatter symbols
	GPscany GP = GPwsc | GPscat | GPstscat
)

type (
	// GameAlias structure describes the game target of algorithm.
	// Several games can shares single algorithm, and in this case
	// all those games presents in the list of aliases for this algorithm.
	// All game rules, paytables and lines should be equal for this games,
	// except maximum number of selected lines can differ. If maximum number
	// of lines differ, algorithm receives largest number.
	GameAlias struct {
		Prov string    `json:"prov" yaml:"prov" xml:"prov"`                               // game provider
		Name string    `json:"name" yaml:"name" xml:"name"`                               // game name
		LNum int       `json:"lnum,omitempty" yaml:"lnum,omitempty" xml:"lnum,omitempty"` // maximum number of selectable lines
		Date util.Unix `json:"date,omitempty" yaml:"date,omitempty" xml:"date,omitempty"` // game release date
	}

	AlgDescr struct {
		GT  GT        `json:"gt,omitempty" yaml:"gt,omitempty" xml:"gt,omitempty"` // game type
		GP  GP        `json:"gp,omitempty" yaml:"gp,omitempty" xml:"gp,omitempty"` // game properties
		SX  int       `json:"sx,omitempty" yaml:"sx,omitempty" xml:"sx,omitempty"` // grid width
		SY  int       `json:"sy,omitempty" yaml:"sy,omitempty" xml:"sy,omitempty"` // grid height
		SN  int       `json:"sn,omitempty" yaml:"sn,omitempty" xml:"sn,omitempty"` // number of symbols
		LN  int       `json:"ln,omitempty" yaml:"ln,omitempty" xml:"ln,omitempty"` // number of lines in bet lines set
		WN  int       `json:"wn,omitempty" yaml:"wn,omitempty" xml:"wn,omitempty"` // number of ways
		BN  int       `json:"bn,omitempty" yaml:"bn,omitempty" xml:"bn,omitempty"` // number of bonuses
		RTP []float64 `json:"rtp" yaml:"rtp,flow" xml:"rtp"`                       // 'Return to Player' percents list
	}

	AlgInfo struct {
		Aliases  []GameAlias `json:"aliases" yaml:"aliases" xml:"aliases"`
		AlgDescr `yaml:",inline"`
		Update   func(ai *AlgInfo) `json:"-" yaml:"-" xml:"-"` // closure to update fields
	}

	GameInfo struct {
		GameAlias `yaml:",inline"`
		*AlgDescr `yaml:",inline"`
	}

	Scanner func(context.Context, float64) float64
)

var (
	AlgList     = []*AlgInfo{}
	InfoMap     = map[string]*GameInfo{}
	GameFactory = map[string]func() Gamble{}
	ScanFactory = map[string]Scanner{}
)

var (
	Year = util.Year
	Date = util.Date
)

func MakeRtpList[T any](reelsmap map[float64]T) []float64 {
	var list = make([]float64, 0, len(reelsmap))
	for rtp := range reelsmap {
		list = append(list, rtp)
	}
	sort.Float64s(list)
	return list
}

func (ai *AlgInfo) SetupFactory(game func(int) Gamble, scan Scanner) {
	AlgList = append(AlgList, ai)
	for _, ga := range ai.Aliases {
		var aid = util.ToID(ga.Prov + "/" + ga.Name)
		if _, ok := InfoMap[aid]; ok {
			panic(fmt.Errorf("%s: %w", aid, ErrAidHas))
		}
		if ai.GT == GTslot {
			if ga.LNum > ai.LN {
				panic(fmt.Errorf("%s: %w", aid, ErrLNumOut))
			}
			if ai.LN > 0 && ga.LNum == 0 {
				log.Printf("%s: LNum is not set for game with lines set of %d lines", aid, ai.LN)
			}
			if ai.LN == 0 && ai.WN == 0 && ai.GP&GPcpay == 0 {
				log.Printf("%s: both LN and WN are zero", aid)
			}
		}
		InfoMap[aid] = &GameInfo{
			GameAlias: ga,
			AlgDescr:  &ai.AlgDescr,
		}
		GameFactory[aid] = func() Gamble { return game(ga.LNum) }
		ScanFactory[aid] = scan // can be nil
	}
}

func (ad *AlgDescr) FindClosest(mrtp float64) (rtp float64) {
	rtp = -1000 // lets to get first reels from map in any case
	for _, p := range ad.RTP {
		if math.Abs(mrtp-p) < math.Abs(mrtp-rtp) {
			rtp = p
		}
	}
	return
}
