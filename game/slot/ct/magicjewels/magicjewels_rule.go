package magicjewels

// See: https://www.livebet2.com/casino/slots/ct-interactive/magic-jewels

import (
	"github.com/slotopol/server/game/slot"
)

const (
	sn   = 8 // number of symbols
	wild = 1 // wild symbol IDs
)

var ReelsCas slot.Reelx

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [sn][5]float64{
	{},                    //  1 wild (2, 3, 4 reels only)
	{0, 0, 50, 300, 5000}, //  2 crown
	{0, 0, 25, 100, 1500}, //  3 ruby
	{0, 0, 15, 40, 500},   //  4 diamond
	{0, 0, 10, 30, 200},   //  5 emerald
	{0, 0, 5, 15, 50},     //  6 amber
	{0, 0, 5, 15, 50},     //  7 sapphire
	{0, 0, 5, 10, 40},     //  8 amethyst
}

// Bet lines
var BetLines = slot.BetLinesCT5x3[:]

type Game struct {
	slot.Cascade5x3 `yaml:",inline"`
	slot.Slotx      `yaml:",inline"`
}

// Declare conformity with SlotCascade interface.
var _ slot.SlotCascade = (*Game)(nil)

func NewGame() *Game {
	var g = &Game{
		Slotx: slot.Slotx{
			Sel: 25,
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
	g.ScanLined(wins)
	if len(*wins) == 0 {
		switch {
		case g.CFN == 4:
			*wins = append(*wins, slot.WinItem{
				FS: 12,
			})
		case g.CFN == 5:
			*wins = append(*wins, slot.WinItem{
				FS: 20,
			})
		case g.CFN >= 6:
			*wins = append(*wins, slot.WinItem{
				FS: 38,
			})
		}
	}
	return nil
}

func (g *Game) GetMult() float64 {
	var mp float64 = 1
	switch g.CFN {
	case 1:
		mp = 2
	case 2:
		mp = 5
	case 3:
		mp = 10
	}
	if g.FSR != 0 {
		mp *= 2
	}
	return mp
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	var mp float64
	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml = g.LX(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LX(x, line)
			if sx != syml && sx != wild {
				numl = x - 1
				break
			}
		}
		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			if mp == 0 { // lazy calculation
				mp = g.GetMult()
			}
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * pay,
				MP:  mp,
				Sym: syml,
				Num: numl,
				LI:  li + 1,
				XY:  line.HitxL(numl),
			})
		}
	}
}

func (g *Game) GetReels(mrtp float64) slot.Reelx {
	if g.CFN == 0 {
		var reels, _ = ReelsMap.FindClosest(mrtp)
		return reels
	} else {
		return ReelsCas
	}
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
