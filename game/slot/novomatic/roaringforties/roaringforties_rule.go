package roaringforties

// See: https://freeslotshub.com/novomatic/roaring-forties/

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed roaringforties_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [10][5]float64{
	{0, 4, 60, 200, 1000}, //  1 seven
	{0, 0, 40, 100, 300},  //  2 bell
	{0, 0, 20, 80, 200},   //  3 melon
	{0, 0, 20, 80, 200},   //  4 grapes
	{0, 0, 8, 40, 100},    //  5 plum
	{0, 0, 8, 40, 100},    //  6 orange
	{0, 0, 8, 40, 100},    //  7 lemon
	{0, 0, 8, 40, 100},    //  8 cherry
	{},                    //  9 wild
	{},                    // 10 star
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 20, 500} // star

const (
	jid = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [10][5]int{
	{0, 0, 0, 0, 0}, //  1 seven
	{0, 0, 0, 0, 0}, //  2 bell
	{0, 0, 0, 0, 0}, //  3 melon
	{0, 0, 0, 0, 0}, //  4 grapes
	{0, 0, 0, 0, 0}, //  5 plum
	{0, 0, 0, 0, 0}, //  6 orange
	{0, 0, 0, 0, 0}, //  7 lemon
	{0, 0, 0, 0, 0}, //  8 cherry
	{0, 0, 0, 0, 0}, //  9 wild
	{0, 0, 0, 0, 0}, // 10 star
}

// Bet lines
var BetLines = slot.BetLinesNvm5x4[:40]

type Game struct {
	slot.Slotx[slot.Screen5x4] `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx[slot.Screen5x4]{
			Sel: len(BetLines),
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

const wild, scat = 9, 10

func (g *Game) Scanner(wins *slot.Wins) {
	g.ScanLined(wins)
	g.ScanScatters(wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		var numl slot.Pos = 5
		var syml = g.Scrn.Pos(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.Scrn.Pos(x, line)
			if sx != syml && sx != wild {
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
				Line: li,
				XY:   line.CopyL(numl),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.Scrn.ScatNum(scat); count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   g.Scrn.ScatPos(scat),
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindReels(ReelsMap, mrtp)
	g.Scrn.Spin(reels)
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
