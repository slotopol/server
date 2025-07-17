package shiningstars

// See: https://demo.agtsoftware.com/games/agt/shiningstars

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed shiningstars_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [11][5]float64{
	{},                     //  1 wild (on 2, 3, 4 reels)
	{},                     //  2 scatter1 (on all reels)
	{},                     //  3 scatter2 (on 1, 3, 5 reels)
	{0, 10, 50, 250, 5000}, //  4 seven
	{0, 0, 40, 120, 700},   //  5 grape
	{0, 0, 40, 120, 700},   //  6 watermelon
	{0, 0, 20, 40, 200},    //  7 avocado
	{0, 0, 10, 30, 150},    //  8 pomegranate
	{0, 0, 10, 30, 150},    //  9 carambola
	{0, 0, 10, 30, 150},    // 10 maracuya
	{0, 0, 10, 30, 150},    // 11 orange
}

// Scatters payment.
var ScatPay1 = [5]float64{0, 0, 5, 20, 100} // 2 scatter1
var ScatPay2 = [5]float64{0, 0, 20}         // 3 scatter2

// Bet lines
var BetLines = slot.BetLinesAgt5x3[:20]

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
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

const wild, scat1, scat2 = 1, 2, 3

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	var reelwild [5]bool
	for x := 1; x < 4; x++ { // 2, 3, 4 reel only
		for y := 0; y < 3; y++ {
			if g.Scr[x][y] == wild {
				reelwild[x] = true
				break
			}
		}
	}

	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml = g.LY(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LY(x, line)
			if reelwild[x-1] {
				continue
			} else if sx != syml {
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
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.ScatNum(scat1); count >= 3 {
		var pay = ScatPay1[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  scat1,
			Num:  count,
			XY:   g.ScatPos(scat1),
		})
	} else if count := g.ScatNum(scat2); count >= 3 {
		var pay = ScatPay2[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  scat2,
			Num:  count,
			XY:   g.ScatPos(scat2),
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	g.ReelSpin(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
