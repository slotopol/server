package caribbeanadventure

// See: https://www.slotsmate.com/software/ct-interactive/caribbean-adventure

import (
	"github.com/slotopol/server/game/slot"
)

const (
	sn         = 13   // number of symbols
	wild, scat = 1, 2 // wild & scatter symbol IDs
	wildmin    = 2    // minimum wild symbols to win
	linemin    = 2    // minimum line symbols to win
	scatmin    = 2    // minimum scatters to win
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [sn][5]float64{
	{0, 10, 100, 500, 10000}, //  1 wild
	{},                       //  2 scatter
	{0, 2, 25, 100, 1000},    //  3 pirate
	{0, 2, 25, 100, 1000},    //  4 lady
	{0, 2, 25, 100, 1000},    //  5 spyglass
	{0, 0, 15, 75, 500},      //  6 island
	{0, 0, 15, 50, 400},      //  7 parrot
	{0, 0, 15, 50, 400},      //  8 monkey
	{0, 0, 10, 30, 200},      //  9 ace
	{0, 0, 10, 30, 200},      // 10 king
	{0, 0, 5, 20, 150},       // 11 queen
	{0, 0, 5, 20, 150},       // 12 jack
	{0, 0, 5, 20, 150},       // 13 ten
}

// Scatters payment.
var ScatPay = [5]float64{0, 1, 10, 50, 1000} // 2 scatter

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 10, 10, 10} // 2 scatter

// Bet lines
var BetLines = slot.BetLinesCT5x3[:]

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame(sel int) *Game {
	return &Game{
		Slotx: slot.Slotx{
			Sel: sel,
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

func (g *Game) Scanner(wins *slot.Wins) error {
	if g.FSR == 0 {
		g.ScanLined(wins)
	} else {
		g.ScanBonScat(wins)
	}
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.LX(x, line)
			if sx == wild {
				if syml == 0 {
					numw = x
				}
			} else if syml == 0 {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		var payw, payl float64
		if numw >= wildmin {
			payw = LinePay[wild-1][numw-1]
		}
		if numl >= linemin && syml > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl > payw {
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * payl,
				MP:  1,
				Sym: syml,
				Num: numl,
				LI:  li + 1,
				XY:  line.HitxL(numl),
			})
		} else if payw > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * payw,
				MP:  1,
				Sym: wild,
				Num: numw,
				LI:  li + 1,
				XY:  line.HitxL(numw),
			})
		}
	}
}

func (g *Game) ScanBonScat(wins *slot.Wins) {
	// Count symbols
	var counts [sn + 1]slot.Pos // symbol counts per reel
	for x := range 5 {
		for _, sym := range g.Scr[x] {
			counts[sym]++
		}
	}
	// Scatters calculation
	for sym, num := range counts {
		if num >= linemin {
			if pay := LinePay[sym-1][num-1]; pay > 0 {
				*wins = append(*wins, slot.WinItem{
					Pay: g.Bet * pay,
					MP:  1,
					Sym: slot.Sym(sym),
					Num: num,
					XY:  g.SymPos(slot.Sym(sym)),
				})
			}
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.SymNum(scat); count >= scatmin {
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  1,
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
			FS:  fs,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.SpinReels(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
