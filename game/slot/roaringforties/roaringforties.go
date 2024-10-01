package roaringforties

// See: https://freeslotshub.com/novomatic/roaring-forties/

import (
	slot "github.com/slotopol/server/game/slot"
)

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
	{0, 0, 0, 0, 0},       //  9 wild
	{0, 0, 0, 0, 0},       // 10 star
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

type Game struct {
	slot.Slot5x4 `yaml:",inline"`
}

func NewGame() *Game {
	return &Game{
		Slot5x4: slot.Slot5x4{
			Sel: slot.MakeBitNum(40, 1),
			Bet: 1,
		},
	}
}

const wild, scat = 9, 10

var bl = slot.BetLinesNvm40

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	for li := g.Sel.Next(0); li != -1; li = g.Sel.Next(li) {
		var line = bl[li-1]

		var numl slot.Pos = 1
		var syml = screen.Pos(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx != syml && sx != wild {
				break
			}
			numl++
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
func (g *Game) ScanScatters(screen slot.Screen, wins *slot.Wins) {
	if count := screen.ScatNum(scat); count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel.Num()) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
		})
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	var _, reels = slot.FindReels(ReelsMap, mrtp)
	screen.Spin(reels)
}

func (g *Game) SetSel(sel slot.Bitset) error {
	if sel.IsZero() {
		return slot.ErrNoLineset
	}
	if bs := sel; !bs.AndNot(slot.MakeBitNum(len(bl), 1)).IsZero() {
		return slot.ErrLinesetOut
	}
	if g.FreeSpins() > 0 {
		return slot.ErrDisabled
	}
	g.Sel = sel
	return nil
}
