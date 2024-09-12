package goldentour

// See: https://freeslotshub.com/playtech/golden-tour/

import (
	"math"

	"github.com/slotopol/server/game"
)

// reels lengths [32, 46, 54, 46, 32], total reshuffles 117006336
// golf bonuses: count 1264896, rtp = 17.296786%
// golf bonuses frequency: 1/92.503
// RTP = 127.57(lined) + 0(scatter) + 17.297(golf) = 144.868028%
var Reels145 = game.Reels5x{
	{4, 6, 3, 7, 8, 2, 11, 5, 6, 4, 2, 10, 5, 8, 7, 4, 1, 7, 5, 4, 8, 7, 6, 8, 5, 6, 3, 7, 6, 5, 8, 9},
	{11, 5, 3, 8, 7, 6, 2, 8, 7, 5, 8, 4, 7, 10, 8, 4, 7, 6, 5, 7, 8, 5, 6, 8, 7, 6, 3, 7, 5, 6, 8, 7, 6, 5, 8, 7, 4, 6, 8, 2, 6, 9, 8, 1, 4, 7},
	{8, 6, 7, 8, 6, 7, 2, 6, 7, 8, 2, 7, 4, 8, 6, 4, 8, 7, 5, 6, 7, 5, 8, 1, 6, 8, 7, 5, 8, 4, 7, 6, 9, 7, 3, 4, 8, 7, 6, 8, 7, 5, 8, 6, 5, 8, 7, 5, 8, 6, 11, 8, 3, 10},
	{7, 6, 8, 7, 6, 4, 8, 3, 11, 7, 5, 6, 4, 5, 7, 6, 10, 8, 2, 5, 6, 1, 8, 6, 5, 7, 6, 8, 7, 5, 8, 7, 6, 4, 8, 7, 9, 8, 7, 3, 4, 8, 7, 5, 2, 8},
	{5, 8, 4, 5, 2, 10, 6, 4, 11, 7, 2, 6, 7, 5, 8, 4, 1, 8, 9, 6, 5, 4, 6, 3, 7, 8, 3, 7, 8, 5, 7, 6},
}

// Map with available reels.
var ReelsMap = map[float64]*game.Reels5x{
	144.868028: &Reels145,
}

func FindReels(mrtp float64) (rtp float64, reels game.Reels) {
	for p, r := range ReelsMap {
		if math.Abs(mrtp-p) < math.Abs(mrtp-rtp) {
			rtp, reels = p, r
		}
	}
	return
}

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
	{0, 0, 0, 0, 0},         //  9 fitch
	{0, 0, 0, 0, 0},         // 10 drake
	{0, 0, 0, 0, 0},         // 11 luce
}

type Game struct {
	game.Slot5x3 `yaml:",inline"`
}

func NewGame() *Game {
	return &Game{
		Slot5x3: game.Slot5x3{
			SBL: game.MakeBitNum(5),
			Bet: 1,
		},
	}
}

var bl = game.Lineset5x{
	{2, 2, 2, 2, 2}, // 1
	{1, 1, 1, 1, 1}, // 2
	{3, 3, 3, 3, 3}, // 3
	{1, 2, 3, 2, 1}, // 4
	{3, 2, 1, 2, 3}, // 5
}

const golfbon = 1
const wild, scat1, scat2, scat3 = 1, 9, 10, 11

func (g *Game) Scanner(screen game.Screen, wins *game.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen game.Screen, wins *game.Wins) {
	for li := g.SBL.Next(0); li != 0; li = g.SBL.Next(li) {
		var line = bl.Line(li)

		var numw, numl = 0, 5
		var syml game.Sym
		for x := 1; x <= 5; x++ {
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
				*wins = append(*wins, game.WinItem{
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
				*wins = append(*wins, game.WinItem{
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
			var numw, numr = 0, 5 - numl
			var symr game.Sym
			for x := 5; x > numl; x-- {
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
					*wins = append(*wins, game.WinItem{
						Pay:  g.Bet * pay,
						Mult: 1,
						Sym:  symr,
						Num:  numr,
						Line: li,
						XY:   line.CopyR(numr),
					})
				}
			} else if numw > 0 {
				if pay := LinePay[wild-1][numw-1]; pay > 0 {
					*wins = append(*wins, game.WinItem{
						Pay:  g.Bet * pay,
						Mult: 1,
						Sym:  wild,
						Num:  numw,
						Line: li,
						XY:   line.CopyR(numw),
					})
				}
			}
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen game.Screen, wins *game.Wins) {
	if count := screen.ScatNum(scat1); count >= 3 {
		*wins = append(*wins, game.WinItem{
			Mult: 1,
			Sym:  scat1,
			Num:  count,
			XY:   screen.ScatPos(scat1),
			BID:  golfbon,
		})
	} else if count := screen.ScatNum(scat2); count >= 3 {
		*wins = append(*wins, game.WinItem{
			Mult: 1,
			Sym:  scat2,
			Num:  count,
			XY:   screen.ScatPos(scat2),
			BID:  golfbon,
		})
	} else if count := screen.ScatNum(scat3); count >= 3 {
		*wins = append(*wins, game.WinItem{
			Mult: 1,
			Sym:  scat3,
			Num:  count,
			XY:   screen.ScatPos(scat3),
			BID:  golfbon,
		})
	}
}

func (g *Game) Spin(screen game.Screen, mrtp float64) {
	var _, reels = FindReels(mrtp)
	screen.Spin(reels)
}

func (g *Game) Spawn(screen game.Screen, wins game.Wins) {
	for i, wi := range wins {
		switch wi.BID {
		case golfbon:
			wins[i].Pay = GolfSpawn(g.Bet * float64(g.SBL.Num()))
		}
	}
}

func (g *Game) SetLines(sbl game.Bitset) error {
	var mask game.Bitset = (1<<len(bl) - 1) << 1
	if sbl == 0 {
		return game.ErrNoLineset
	}
	if sbl&^mask != 0 {
		return game.ErrLinesetOut
	}
	if g.FreeSpins() > 0 {
		return game.ErrNoFeature
	}
	g.SBL = sbl
	return nil
}
