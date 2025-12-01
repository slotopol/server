package rainbowcharm

// See: https://www.slotsmate.com/software/ct-interactive/rainbow-charm

// Remark: bonus symbol turns to another symbol on client side
// for some random winning combinations. This is just an
// animation only and does not affect the calculation.

import (
	"math/rand/v2"

	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[*slot.Reels5x]

// Symbols payment.
var SymPay = [6][12]float64{
	{}, // 1 bonus
	{0, 0, 0, 0, 0, 0, 30, 50, 100, 500, 1000, 5000}, // 2 leprechaun
	{0, 0, 0, 0, 0, 0, 12, 15, 80, 100, 300, 1000},   // 3 clover
	{0, 0, 0, 0, 0, 0, 9, 10, 20, 80, 200, 500},      // 4 pot
	{0, 0, 0, 0, 0, 0, 8, 9, 15, 30, 80, 250},        // 5 horseshoe
	{0, 0, 0, 0, 0, 0, 7, 8, 12, 20, 50, 150},        // 6 bell
}

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
	M              [5]float64 `json:"m" yaml:"m" xml:"m"` // multipliers for bonus symbol filled at reels
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

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

var MBon = [6]float64{2, 2, 4, 8} // multipliers for bonus symbol filled at reels, E = 4

const bon = 1

func (g *Game) Scanner(wins *slot.Wins) error {
	var mb float64 = 1.0 // multiplier for bonus symbol filled at reels
	var sn [6]slot.Pos
	var x slot.Pos
	for x = range 5 {
		var r = g.Scr[x]
		if r[0] == bon && r[1] == bon && r[2] == bon {
			mb *= g.M[x]
		}
		sn[r[0]-1]++
		sn[r[1]-1]++
		sn[r[2]-1]++
	}
	if mb > 1.0 {
		var has bool
		var sym slot.Sym
		for sym = 2; sym <= 6; sym++ {
			if count := sn[sym-1]; count >= 7 {
				has = true
				break
			}
		}
		if has {
			for x = range 5 {
				if g.M[x] > 1.0 {
					*wins = append(*wins, slot.WinItem{
						MP:  g.M[x],
						Sym: bon,
						Num: 3,
						XY:  slot.Hitx{{x + 1, 1}, {x + 1, 2}, {x + 1, 3}},
					})
				}
			}
		}
	}

	var sym slot.Sym
	for sym = 2; sym <= 6; sym++ {
		if count := sn[sym-1]; count >= 7 {
			var pay = SymPay[sym-1][min(count, 12)-1]
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * pay,
				MP:  mb,
				Sym: sym,
				Num: count,
				XY:  g.SymPos(sym),
			})
		}
	}
	return nil
}

func (g *Game) Cost() float64 {
	return g.Bet * 12
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.ReelSpin(reels)
}

func (g *Game) Prepare() {
	for x := range 5 {
		g.M[x] = MBon[rand.N(len(MBon))]
	}
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
