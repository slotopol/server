package guardianofasgard

// See: https://www.slotsmate.com/software/ct-interactive/guardian-of-asgard

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[*slot.Reels5x]

// Lined payment.
var LinePay = [10][7]float64{
	{0, 0, 80, 200, 2000},        //  1 scatter
	{0, 0, 0, 0, 200, 400, 1500}, //  2 wild
	{0, 0, 0, 0, 14, 20, 150},    //  3 bowl
	{0, 0, 0, 0, 2, 5, 40},       //  4 wolf
	{0, 0, 0, 0, 2, 5, 40},       //  5 helm
	{0, 0, 0, 0, 2, 5, 40},       //  6 axe
	{0, 0, 0, 0, 1, 2, 10},       //  7 ace
	{0, 0, 0, 0, 1, 2, 10},       //  8 king
	{0, 0, 0, 0, 1, 2, 10},       //  9 queen
	{0, 0, 0, 0, 1, 2, 10},       // 10 jack
}

type Game struct {
	slot.Cascade5x4 `yaml:",inline"`
	slot.Slotx      `yaml:",inline"`
}

// Declare conformity with CascadeSlot interface.
var _ slot.CascadeSlot = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx{
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

func (g *Game) Free() bool {
	return g.FSR != 0 || g.Cascade()
}

const wild, scat = 2, 1

func (g *Game) Scanner(wins *slot.Wins) error {
	var sn [10]slot.Pos
	var x, y slot.Pos
	for x = range 5 {
		for y = range 3 {
			sn[g.Scr[x][y]-1]++
		}
	}

	if count := sn[scat-1]; count >= 3 {
		var pay = LinePay[scat-1][count-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * pay,
			MP:  1,
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
		})
	}
	if count := sn[wild-1]; count >= 5 {
		var pay = LinePay[wild-1][min(count-1, 6)]
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
		if count := sn[sym-1] + sn[wild-1]; count >= 5 {
			var pay = LinePay[sym-1][min(count-1, 6)]
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
	return g.Bet * 4
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.ReelSpin(reels)
}

func (g *Game) Prepare() {
	g.UntoFall()
}

func (g *Game) Apply(wins slot.Wins) {
	g.Slotx.Apply(wins)
	g.Strike(wins)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
