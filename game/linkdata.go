package game

import (
	"context"
	"math"
	"sort"

	"github.com/slotopol/server/util"
)

type Gamble interface {
	Spin(float64)         // fill the screen with random hits on reels closest to given RTP, constant function
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
	GPjack // cumulative jackpot is present
	GPfill // has multiplier on filled screen

	GPcumul // has cumulative pays (in this case spin can fails if bank have not enough money)
	GPbmode // has non-reels bonus mode
	GPcasc  // cascade falls present
	GPcmult // multipliers on cascade avalanche

	GPfghas  // non-retriggered free games are present
	GPretrig // free games are present and can be retriggered
	GPfgreel // separate reels on free games
	GPfgmult // any multipliers on free games

	GPrmult // any multipliers on regular games
	GPscat  // has scatters
	GPwild  // has wild symbols
	GPrwild // has reel wild symbols

	GPbwild // has big (3x3) or expanding wild
	GPwturn // symbols turns to wilds
	GPwmult // has multiplier on wilds
	GPbsym  // has big symbol (usually 3x3 in the center on free games)

	GPfgno GP = 0 // free games are absent
)

type (
	// GameAlias structure describes the game target of algorithm.
	// Several games can shares single algorithm, and in this case
	// all those games presents in the list of aliases for this algorithm.
	// All game rules, paytables and lines should be equal for this games,
	// except maximum number of selected lines can differ. If maximum number
	// of lines differ, algorithm receives largest number.
	GameAlias struct {
		Prov string    `json:"prov" yaml:"prov" xml:"prov"`
		Name string    `json:"name" yaml:"name" xml:"name"`
		Date util.Unix `json:"date,omitempty" yaml:"date,omitempty" xml:"date,omitempty"`
	}

	AlgDescr struct {
		GT  GT        `json:"gt,omitempty" yaml:"gt,omitempty" xml:"gt,omitempty"` // game type
		GP  GP        `json:"gp,omitempty" yaml:"gp,omitempty" xml:"gp,omitempty"` // game properties
		SX  int       `json:"sx,omitempty" yaml:"sx,omitempty" xml:"sx,omitempty"` // screen width
		SY  int       `json:"sy,omitempty" yaml:"sy,omitempty" xml:"sy,omitempty"` // screen height
		SN  int       `json:"sn,omitempty" yaml:"sn,omitempty" xml:"sn,omitempty"` // number of symbols
		LN  int       `json:"ln,omitempty" yaml:"ln,omitempty" xml:"ln,omitempty"` // number of lines
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

func (ai *AlgInfo) SetupFactory(game func() Gamble, scan Scanner) {
	AlgList = append(AlgList, ai)
	for _, ga := range ai.Aliases {
		var aid = util.ToID(ga.Prov + "/" + ga.Name)
		if _, ok := InfoMap[aid]; ok {
			panic(ErrAidHas)
		}
		InfoMap[aid] = &GameInfo{
			GameAlias: ga,
			AlgDescr:  &ai.AlgDescr,
		}
		GameFactory[aid] = game
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
