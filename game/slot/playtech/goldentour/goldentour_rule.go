package goldentour

// See: https://freeslotshub.com/playtech/golden-tour/

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

const golfbon = 1
const wild, scat1, scat2, scat3 = 1, 9, 10, 11

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = screen.Pos(x, line)
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

		if numl > 0 && syml > 0 {
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
		} else if numw > 0 {
			if pay := LinePay[wild-1][numw-1]; pay > 0 {
				*wins = append(*wins, slot.WinItem{
					Pay:  g.Bet * pay,
					Mult: 1,
					Sym:  wild,
					Num:  numw,
					Line: li,
					XY:   line.CopyL(numw),
				})
			}
		}

		if numl < 5 {
			var numw, numr slot.Pos = 0, 5 - numl
			var symr slot.Sym
			var x slot.Pos
			for x = 5; x > numl; x-- {
				var sx = screen.Pos(x, line)
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

			if numr > 0 && symr > 0 {
				if pay := LinePay[symr-1][numr-1]; pay > 0 {
					*wins = append(*wins, slot.WinItem{
						Pay:  g.Bet * pay,
						Mult: 1,
						Sym:  symr,
						Num:  numr,
						Line: li,
						XY:   line.CopyR5(numr),
					})
				}
			} else if numw > 0 {
				if pay := LinePay[wild-1][numw-1]; pay > 0 {
					*wins = append(*wins, slot.WinItem{
						Pay:  g.Bet * pay,
						Mult: 1,
						Sym:  wild,
						Num:  numw,
						Line: li,
						XY:   line.CopyR5(numw),
					})
				}
			}
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen slot.Screen, wins *slot.Wins) {
	if count := screen.ScatNum(scat1); count >= 3 {
		*wins = append(*wins, slot.WinItem{
			Mult: 1,
			Sym:  scat1,
			Num:  count,
			XY:   screen.ScatPos(scat1),
			BID:  golfbon,
		})
	} else if count := screen.ScatNum(scat2); count >= 3 {
		*wins = append(*wins, slot.WinItem{
			Mult: 1,
			Sym:  scat2,
			Num:  count,
			XY:   screen.ScatPos(scat2),
			BID:  golfbon,
		})
	} else if count := screen.ScatNum(scat3); count >= 3 {
		*wins = append(*wins, slot.WinItem{
			Mult: 1,
			Sym:  scat3,
			Num:  count,
			XY:   screen.ScatPos(scat3),
			BID:  golfbon,
		})
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	var reels, _ = slot.FindReels(ReelsMap, mrtp)
	screen.Spin(reels)
}

func (g *Game) Spawn(screen slot.Screen, wins slot.Wins) {
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
