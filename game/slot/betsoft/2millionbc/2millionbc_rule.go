package twomillionbc

// See: https://www.slotsmate.com/software/betsoft/2-million-bc

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed 2millionbc_bon.yaml
var rbon []byte

var ReelsBon = slot.ReadObj[*slot.Reels5x](rbon)

//go:embed 2millionbc_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [13][5]float64{
	{0, 30, 100, 300, 500}, //  1 girl
	{0, 15, 75, 200, 400},  //  2 lion
	{0, 10, 60, 150, 300},  //  3 bee
	{0, 5, 50, 125, 250},   //  4 stone
	{0, 5, 40, 100, 200},   //  5 wheel
	{0, 2, 30, 90, 150},    //  6 club
	{0, 0, 25, 75, 125},    //  7 chaplet
	{0, 0, 20, 60, 100},    //  8 gold
	{0, 0, 15, 50, 75},     //  9 vase
	{0, 0, 10, 25, 50},     // 10 ruby
	{},                     // 11 fire
	{},                     // 12 acorn
	{0, 0, 40, 100, 200},   // 13 diamond
}

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 4, 12, 20} // 11 fire

// Bet lines
var BetLines = slot.BetLinesNetEnt5x3[:30]

const (
	acbn = 1 // acorn bonus
	dlbn = 2 // diamond lion bonus
)

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
	// acorns number
	AN int `json:"an" yaml:"an" xml:"an"`
	// acorns bet
	AB float64 `json:"ab" yaml:"ab" xml:"ab"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx{
			Sel: len(BetLines),
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

const scat, acorn, diamond = 11, 12, 13

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml = g.LY(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LY(x, line)
			if sx != syml {
				numl = x - 1
				break
			}
		}

		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * pay,
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li + 1,
				XY:   line.CopyL(numl),
			})
		}
		if syml == diamond && numl >= 3 {
			*wins = append(*wins, slot.WinItem{
				Mult: 1,
				Sym:  diamond,
				Num:  numl,
				Line: li + 1,
				XY:   line.CopyL(numl),
				BID:  dlbn,
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.ScatNum(scat); count >= 3 {
		var fs = ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Sym:  scat,
			Num:  count,
			XY:   g.ScatPos(scat),
			Free: fs,
		})
	}

	if g.At(5, 1) == acorn || g.At(5, 2) == acorn || g.At(5, 3) == acorn {
		if (g.AN+1)%3 == 0 {
			*wins = append(*wins, slot.WinItem{
				Mult: 1,
				Sym:  acorn,
				Num:  1,
				BID:  acbn,
			})
		}
	}
}

func (g *Game) Spin(mrtp float64) {
	if g.FSR == 0 {
		var reels, _ = slot.FindClosest(ReelsMap, mrtp)
		g.ReelSpin(reels)
	} else {
		g.ReelSpin(ReelsBon)
	}
}

func (g *Game) Spawn(wins slot.Wins, fund, mrtp float64) {
	for i, wi := range wins {
		switch wi.BID {
		case acbn:
			wins[i].Pay = AcornSpawn(g.AB + g.Bet*float64(g.Sel))
		case dlbn:
			wins[i].Pay = DiamondLionSpawn(g.Bet)
		}
	}
}

func (g *Game) Apply(wins slot.Wins) {
	g.Slotx.Apply(wins)

	if g.At(5, 1) == acorn || g.At(5, 2) == acorn || g.At(5, 3) == acorn {
		g.AN++
		g.AN %= 3
		if g.AN > 0 {
			g.AB += g.Bet * float64(g.Sel)
		} else {
			g.AB = 0
		}
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
