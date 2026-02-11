package justjewels

// See: https://www.slotsmate.com/software/novomatic/just-jewels-deluxe

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [8][5]float64{
	{0, 0, 50, 500, 5000}, // 1 crown
	{0, 0, 30, 150, 500},  // 2 gold
	{0, 0, 30, 150, 500},  // 3 money
	{0, 0, 15, 50, 200},   // 4 ruby
	{0, 0, 15, 50, 200},   // 5 sapphire
	{0, 0, 10, 25, 150},   // 6 emerald
	{0, 0, 10, 25, 150},   // 7 amethyst
	{},                    // 8 euro
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 10, 50} // 8 euro

// Bet lines
var BetLines = slot.BetLinesNvm10[:]

type Game struct {
	slot.Grid5x3 `yaml:",inline"`
	slot.Slotx   `yaml:",inline"`
}

// Declare conformity with SlotGeneric interface.
var _ slot.SlotGeneric = (*Game)(nil)

func NewGame(sel int) *Game {
	return &Game{
		Slotx: slot.Slotx{
			Sel: sel,
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGeneric {
	var clone = *g
	return &clone
}

const scat = 8

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 1
		var syml = g.LX(3, line)
		var xy slot.Linex
		xy.Set(3, line.At(3))
		if g.LX(2, line) == syml {
			xy.Set(2, line.At(2))
			numl++
			if g.LX(1, line) == syml {
				xy.Set(1, line.At(1))
				numl++
			}
		}
		if g.LX(4, line) == syml {
			xy.Set(4, line.At(4))
			numl++
			if g.LX(5, line) == syml {
				xy.Set(5, line.At(5))
				numl++
			}
		}

		if numl >= 3 {
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * LinePay[syml-1][numl-1],
				MP:  1,
				Sym: syml,
				Num: numl,
				LI:  li + 1,
				XY:  slot.L2H(xy),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.SymNum(scat); count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  1,
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.SpinReels(reels)
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
