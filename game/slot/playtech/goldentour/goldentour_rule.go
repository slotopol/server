package goldentour

// See: https://freeslotshub.com/playtech/golden-tour/
// See: https://www.slotsmate.com/software/playtech/golden-tour

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed goldentour_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [11][5]float64{
	{0, 25, 100, 500, 2000}, //  1 two balls
	{0, 5, 50, 200, 1000},   //  2 white ball
	{0, 5, 25, 100, 500},    //  3 yellow ball
	{0, 3, 25, 100, 250},    //  4 electrocar
	{0, 2, 20, 80, 200},     //  5 golf clubs
	{0, 2, 15, 50, 150},     //  6 flag
	{0, 2, 5, 25, 100},      //  7 beer
	{0, 2, 5, 10, 50},       //  8 slippers
	{},                      //  9 fitch
	{},                      // 10 drake
	{},                      // 11 luce
}

// Bet lines
var BetLines = slot.BetLinesHot5

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

const golfbon = 1
const wild, scat1, scat2, scat3 = 1, 9, 10, 11

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
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
			var sx = g.LY(x, line)
			if sx == wild {
				if syml == 0 {
					numw = x
				} else if syml >= 4 { // wild after not ball
					numl = x - 1
					break
				}
			} else if syml == 0 && sx < scat1 { // any lined symbol
				if numw > 0 && sx >= 4 { // not ball symbol after wild
					break
				}
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		if numl >= 2 && syml > 0 {
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
		} else if numw >= 2 {
			if pay := LinePay[wild-1][numw-1]; pay > 0 {
				*wins = append(*wins, slot.WinItem{
					Pay:  g.Bet * pay,
					Mult: 1,
					Sym:  wild,
					Num:  numw,
					Line: li + 1,
					XY:   line.CopyL(numw),
				})
			}
		}

		if numl < 5 {
			var numw, numr slot.Pos = 0, 5 - numl
			var symr slot.Sym
			var x slot.Pos
			for x = 5; x > numl; x-- {
				var sx = g.LY(x, line)
				if sx == wild {
					if symr == 0 {
						numw = 6 - x
					} else if symr >= 4 { // wild after not ball
						numr = 5 - x
						break
					}
				} else if symr == 0 && sx < scat1 { // any lined symbol
					if numw > 0 && sx >= 4 { // not ball symbol after wild
						break
					}
					symr = sx
				} else if sx != symr {
					numr = 5 - x
					break
				}
			}

			if numr >= 2 && symr > 0 {
				if pay := LinePay[symr-1][numr-1]; pay > 0 {
					*wins = append(*wins, slot.WinItem{
						Pay:  g.Bet * pay,
						Mult: 1,
						Sym:  symr,
						Num:  numr,
						Line: li + 1,
						XY:   line.CopyR5(numr),
					})
				}
			} else if numw >= 2 {
				if pay := LinePay[wild-1][numw-1]; pay > 0 {
					*wins = append(*wins, slot.WinItem{
						Pay:  g.Bet * pay,
						Mult: 1,
						Sym:  wild,
						Num:  numw,
						Line: li + 1,
						XY:   line.CopyR5(numw),
					})
				}
			}
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.ScatNum(scat1); count >= 3 {
		*wins = append(*wins, slot.WinItem{
			Mult: 1,
			Sym:  scat1,
			Num:  count,
			XY:   g.ScatPos(scat1),
			BID:  golfbon,
		})
	} else if count := g.ScatNum(scat2); count >= 3 {
		*wins = append(*wins, slot.WinItem{
			Mult: 1,
			Sym:  scat2,
			Num:  count,
			XY:   g.ScatPos(scat2),
			BID:  golfbon,
		})
	} else if count := g.ScatNum(scat3); count >= 3 {
		*wins = append(*wins, slot.WinItem{
			Mult: 1,
			Sym:  scat3,
			Num:  count,
			XY:   g.ScatPos(scat3),
			BID:  golfbon,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	g.ReelSpin(reels)
}

func (g *Game) Spawn(wins slot.Wins, fund, mrtp float64) {
	for i, wi := range wins {
		switch wi.BID {
		case golfbon:
			wins[i].Pay = GolfSpawn(g.Bet * float64(g.Sel))
		}
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
