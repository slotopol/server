package hypercuber

// See: https://www.livebet2.com/casino/slots/ct-interactive/hyper-cuber

import (
	"math/rand/v2"

	"github.com/slotopol/server/game/slot"
)

const (
	sn         = 11   // number of symbols
	wild, scat = 2, 1 // wild & scatter symbol IDs
	bon        = 11   // cuber symbol ID
	Mavr       = 12   // average multiplier on cuber
)

var ReelsBon slot.Reelx

var ReelsMap slot.ReelsMap[slot.Reelx]

// Symbols payment.
var SymPay = [sn][7]float64{
	{0, 0, 15, 75, 300},       //  1 infinity
	{0, 0, 0, 0, 10, 20, 100}, //  2 wild (2, 3, 4 reels only)
	{0, 0, 0, 0, 10, 20, 65},  //  3 atom
	{0, 0, 0, 0, 3, 5, 20},    //  4 red
	{0, 0, 0, 0, 3, 5, 20},    //  5 yellow
	{0, 0, 0, 0, 3, 5, 20},    //  6 gold
	{0, 0, 0, 0, 1, 3, 10},    //  7 violet
	{0, 0, 0, 0, 1, 3, 10},    //  8 lilac
	{0, 0, 0, 0, 1, 3, 10},    //  9 green
	{0, 0, 0, 0, 1, 3, 10},    // 10 blue
	{},                        // 11 cuber
}

// Average multiplier = 12
var CuberMult = [...]float64{2, 2, 2, 2, 3, 3, 3, 5, 5, 5, 10, 10, 10, 15, 15, 100}

type Game struct {
	slot.Cascade5x3 `yaml:",inline"`
	slot.Slotx      `yaml:",inline"`
	M               [5]float64 `json:"m" yaml:"m" xml:"m"` // multipliers for cuber symbols
}

// Declare conformity with SlotCascade interface.
var _ slot.SlotCascade = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx{
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGeneric {
	var clone = *g
	return &clone
}

func (g *Game) FreeMode() bool {
	return g.FSR != 0 || g.Cascade()
}

func (g *Game) Scanner(wins *slot.Wins) error {
	var mc float64
	var counts [sn + 1]slot.Pos
	for x, sr := range g.Grid {
		for _, sy := range sr {
			counts[sy]++
			if sy == bon {
				mc += g.M[x]
			}
		}
	}
	if mc == 0 {
		mc = 1
	}

	if count := counts[scat]; count >= 3 {
		var pay = SymPay[scat-1][count-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * pay,
			MP:  mc,
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
			FS:  15,
		})
	}
	if count := counts[wild]; count >= 5 {
		var pay = SymPay[wild-1][min(count, 7)-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * pay,
			MP:  mc,
			Sym: wild,
			Num: count,
			XY:  g.SymPos(wild),
		})
	}

	var sym slot.Sym
	for sym = 3; sym <= 10; sym++ {
		if count := counts[sym] + counts[wild]; count >= 5 {
			var pay = SymPay[sym-1][min(count, 7)-1]
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * pay,
				MP:  mc,
				Sym: sym,
				Num: count,
				XY:  g.SymPos2(sym, wild),
			})
		}
	}
	return nil
}

func (g *Game) Cost() float64 {
	return g.Bet * 3
}

func (g *Game) Spin(mrtp float64) {
	if g.FSR == 0 {
		var reels, _ = ReelsMap.FindClosest(mrtp)
		g.SpinReels(reels)
	} else {
		g.SpinReels(ReelsBon)
	}
}

func (g *Game) Prepare() {
	g.UntoFall()
	if g.FSR != 0 {
		for x := range 5 {
			g.M[x] = CuberMult[rand.N(len(CuberMult))]
		}
	} else {
		clear(g.M[:])
	}
}

func (g *Game) Apply(wins slot.Wins) {
	g.Slotx.Apply(wins)
	g.Strike(wins)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
