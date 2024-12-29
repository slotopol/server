package valentinesday

// See: https://demo.agtsoftware.com/games/agt/valentine

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed valentinesday_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [8][5]float64{
	{},                      // 1 scatter
	{0, 0, 100, 1000, 5000}, // 2 angel
	{0, 0, 40, 400, 1000},   // 3 nymph
	{0, 0, 24, 60, 200},     // 4 soul
	{0, 0, 20, 50, 200},     // 5 toy
	{0, 0, 20, 50, 200},     // 6 balloon
	{0, 5, 16, 40, 160},     // 7 hearts
	{0, 0, 5, 20, 100},      // 8 medallion
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 12, 60} // 1 scatter

// Bet lines
var BetLines = slot.BetLinesAgt5x3[:5]

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: len(BetLines),
			Bet: 1,
		},
	}
}

const scat = 1

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	var scrn5x3 = screen.(*slot.Screen5x3)
	g.ScanLined(scrn5x3, wins)
	g.ScanScatters(scrn5x3, wins)
}

func FillMult(screen *slot.Screen5x3) float64 {
	var sym = screen[2][1] // center symbol
	if sym < 4 || sym > 7 {
		return 1
	}
	var r *[3]slot.Sym
	if r = &screen[2]; r[0] != sym || r[2] != sym {
		return 1
	}
	var n = 1
	if r = &screen[1]; r[0] == sym && r[1] == sym && r[2] == sym {
		n++
		if r = &screen[0]; r[0] == sym && r[1] == sym && r[2] == sym {
			n++
		}
	}
	if r = &screen[3]; r[0] == sym && r[1] == sym && r[2] == sym {
		n++
		if r = &screen[4]; r[0] == sym && r[1] == sym && r[2] == sym {
			n++
		}
	}
	if n < 3 {
		return 1
	}
	return float64(n - 1)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen *slot.Screen5x3, wins *slot.Wins) {
	var fm float64 // fill mult
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		var numl slot.Pos = 5
		var syml = screen.Pos(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx != syml {
				numl = x - 1
				break
			}
		}

		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			if fm == 0 { // lazy calculation
				fm = FillMult(screen)
			}
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * pay,
				Mult: fm,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen *slot.Screen5x3, wins *slot.Wins) {
	if count := screen.ScatNum(scat); count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
		})
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	var reels, _ = slot.FindReels(ReelsMap, mrtp)
	screen.Spin(reels)
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
