package sugartown

// See: https://www.slotsmate.com/software/ct-interactive/sugar-town

import (
	"github.com/slotopol/server/game/slot"
)

const (
	sn         = 10   // number of symbols
	wild, scat = 2, 1 // wild & scatter symbol IDs
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Symbols payment.
var SymPay = [sn][7]float64{
	{0, 0, 800, 2000, 20000},        //  1 scatter
	{0, 0, 0, 0, 2000, 4000, 15000}, //  2 wild (2, 3, 4 reels only)
	{0, 0, 0, 0, 140, 200, 1500},    //  3 heart
	{0, 0, 0, 0, 20, 50, 400},       //  4 blue
	{0, 0, 0, 0, 20, 50, 400},       //  5 green
	{0, 0, 0, 0, 20, 50, 400},       //  6 yellow
	{0, 0, 0, 0, 10, 20, 100},       //  7 melon
	{0, 0, 0, 0, 10, 20, 100},       //  8 jujube
	{0, 0, 0, 0, 10, 20, 100},       //  9 plum
	{0, 0, 0, 0, 10, 20, 100},       // 10 cherry
}

type Game struct {
	slot.Cascade5x3 `yaml:",inline"`
	slot.Slotx      `yaml:",inline"`
}

// Declare conformity with SlotCascade interface.
var _ slot.SlotCascade = (*Game)(nil)

func NewGame() *Game {
	var g = &Game{
		Slotx: slot.Slotx{
			Bet: 1,
		},
	}
	g.SpinReels(slot.WySRNG, g.GetReels(slot.InitRTP))
	return g
}

func (g *Game) Clone() slot.SlotGeneric {
	var clone = *g
	return &clone
}

func (g *Game) FreeMode() bool {
	return g.FSR != 0 || g.Cascade()
}

func (g *Game) Scanner(wins *slot.Wins) error {
	var counts [sn + 1]slot.Pos
	for _, sr := range g.Grid {
		for _, sy := range sr {
			counts[sy]++
		}
	}

	if count := counts[scat]; count >= 3 {
		var pay = SymPay[scat-1][count-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * pay,
			MP:  1,
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
		})
	}
	if count := counts[wild]; count >= 5 {
		var pay = SymPay[wild-1][min(count, 7)-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * pay,
			MP:  1,
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
				MP:  1,
				Sym: sym,
				Num: count,
				XY:  g.SymPos2(sym, wild),
			})
		}
	}
	return nil
}

func (g *Game) Cost() float64 {
	return g.Bet * 40
}

func (g *Game) GetReels(mrtp float64) slot.Reelx {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	return reels
}

func (g *Game) Spin(rng *slot.Rand, mrtp float64) {
	g.SpinReels(rng, g.GetReels(mrtp))
}

func (g *Game) Prepare(rng *slot.Rand) {
	g.UntoFall()
}

func (g *Game) Apply(wins slot.Wins) {
	g.Slotx.Apply(wins)
	g.Strike(wins)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
