package jewels4all

import "github.com/slotopol/server/game"

// RTP(no eu) = 67.344781%
// RTP(eu at y=1,5) = 1706.345577%
// RTP(eu at y=2,3,4) = 7818.930041%
// euro avr: rtpeu = 5373.896256%
var Reels = game.Reels5x{
	{1, 1, 1, 2, 2, 2, 6, 6, 3, 3, 3, 7, 7, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 6, 6, 3, 3, 3, 7, 7, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 6, 6, 3, 3, 3, 7, 7, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 6, 6, 3, 3, 3, 7, 7, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 6, 6, 3, 3, 3, 7, 7, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
}

// Map with wild chances.
var ChanceMap = map[string]float64{
	// RTP = 67.345(sym) + wc*5373.9(eu) = 90.019449%
	"90": 1 / 237.,
	// RTP = 67.345(sym) + wc*5373.9(eu) = 91.995681%
	"92": 1 / 218.,
	// RTP = 67.345(sym) + wc*5373.9(eu) = 93.948228%
	"94": 1 / 202.,
	// RTP = 67.345(sym) + wc*5373.9(eu) = 96.082194%
	"96": 1 / 187.,
	// RTP = 67.345(sym) + wc*5373.9(eu) = 98.052760%
	"98": 1 / 175.,
	// RTP = 67.345(sym) + wc*5373.9(eu) = 99.913849%
	"100": 1 / 165.,
	// RTP = 67.345(sym) + wc*5373.9(eu) = 110.335951%
	"110": 1 / 125.,
}

// Lined payment.
var LinePay = [8][5]float64{
	{0, 0, 20, 100, 1000}, // 1 crown
	{0, 0, 10, 60, 500},   // 2 gold
	{0, 0, 10, 60, 500},   // 3 money
	{0, 0, 5, 40, 200},    // 4 ruby
	{0, 0, 5, 40, 200},    // 5 sapphire
	{0, 0, 5, 20, 100},    // 6 emerald
	{0, 0, 5, 20, 100},    // 7 amethyst
	{0, 0, 0, 0, 0},       // 8 euro
}

const (
	jid = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [8][5]int{
	{0, 0, 0, 0, 0}, //  1 crown
	{0, 0, 0, 0, 0}, //  2 gold
	{0, 0, 0, 0, 0}, //  3 money
	{0, 0, 0, 0, 0}, //  4 ruby
	{0, 0, 0, 0, 0}, //  5 sapphire
	{0, 0, 0, 0, 0}, //  6 emerald
	{0, 0, 0, 0, 0}, //  7 amethyst
	{0, 0, 0, 0, 0}, //  8 euro
}

type Game struct {
	game.Slot5x3 `yaml:",inline"`
}

func NewGame(rd string) *Game {
	return &Game{
		Slot5x3: game.Slot5x3{
			RD:  rd,
			SBL: game.MakeBitNum(5),
			Bet: 1,
		},
	}
}

const wild = 8

var bl = game.BetLinesNvm10

func (g *Game) Scanner(screen game.Screen, wins *game.Wins) {
	g.ScanLined(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen game.Screen, wins *game.Wins) {
	var scrnwild game.Screen5x3 = *screen.(*game.Screen5x3)
	for x := 1; x <= 5; x++ {
		for y := 1; y <= 3; y++ {
			if screen.At(x, y) == wild {
				for i := max(0, x-2); i <= min(4, x); i++ {
					for j := max(0, y-2); j <= min(2, y); j++ {
						scrnwild[i][j] = wild
					}
				}
			}
		}
	}

	for li := g.SBL.Next(0); li != 0; li = g.SBL.Next(li) {
		var line = bl.Line(li)

		var sym3 = scrnwild.Pos(3, line)
		var xy = game.NewLine5x()
		var num = 1
		xy.Set(3, line.At(3))
		if sym2 := scrnwild.Pos(2, line); sym2 == sym3 || sym2 == wild || sym3 == wild {
			if sym3 == wild {
				sym3 = sym2
			}
			xy.Set(2, line.At(2))
			num++
			if sym1 := scrnwild.Pos(1, line); sym1 == sym3 || sym1 == wild || sym3 == wild {
				if sym3 == wild {
					sym3 = sym1
				}
				xy.Set(1, line.At(1))
				num++
			}
		}
		if sym4 := scrnwild.Pos(4, line); sym4 == sym3 || sym4 == wild || sym3 == wild {
			if sym3 == wild {
				sym3 = sym4
			}
			xy.Set(4, line.At(4))
			num++
			if sym5 := scrnwild.Pos(5, line); sym5 == sym3 || sym5 == wild || sym3 == wild {
				if sym3 == wild {
					sym3 = sym5
				}
				xy.Set(5, line.At(5))
				num++
			}
		}

		if num >= 3 {
			*wins = append(*wins, game.WinItem{
				Pay:  g.Bet * LinePay[sym3-1][num-1],
				Mult: 1,
				Sym:  sym3,
				Num:  num,
				Line: li,
				XY:   xy,
			})
		} else {
			xy.Free()
		}
	}
}

func (g *Game) Spin(screen game.Screen) {
	screen.Spin(&Reels)
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

func (g *Game) SetReels(rd string) error {
	if _, ok := ChanceMap[rd]; !ok {
		return game.ErrNoReels
	}
	g.RD = rd
	return nil
}
