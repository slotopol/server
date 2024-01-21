package justjewels

import "github.com/slotopol/server/game"

// reels lengths [39, 39, 39, 39, 39], total reshuffles 90224199
// RTP = 114.75(lined) + 8.0152(scatter) = 122.764204%
var Reels123 = game.Reels5x{
	{1, 1, 1, 6, 6, 6, 2, 2, 2, 2, 5, 5, 5, 3, 3, 3, 3, 7, 7, 7, 8, 4, 4, 4, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 4, 4, 4},
	{1, 1, 1, 6, 6, 6, 2, 2, 2, 2, 5, 5, 5, 3, 3, 3, 3, 7, 7, 7, 8, 4, 4, 4, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 4, 4, 4},
	{1, 1, 1, 6, 6, 6, 2, 2, 2, 2, 5, 5, 5, 3, 3, 3, 3, 7, 7, 7, 8, 4, 4, 4, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 4, 4, 4},
	{1, 1, 1, 6, 6, 6, 2, 2, 2, 2, 5, 5, 5, 3, 3, 3, 3, 7, 7, 7, 8, 4, 4, 4, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 4, 4, 4},
	{1, 1, 1, 6, 6, 6, 2, 2, 2, 2, 5, 5, 5, 3, 3, 3, 3, 7, 7, 7, 8, 4, 4, 4, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 4, 4, 4},
}

// Map with available reels.
var ReelsMap = map[string]*game.Reels5x{
	"123": &Reels123, // minimum possible percentage
}

// Lined payment.
var LinePay = [8][5]int{
	{0, 0, 50, 500, 5000}, // crown
	{0, 0, 30, 150, 500},  // gold
	{0, 0, 30, 150, 500},  // money
	{0, 0, 15, 50, 200},   // ruby
	{0, 0, 15, 50, 200},   // sapphire
	{0, 0, 10, 25, 150},   // emerald
	{0, 0, 10, 25, 150},   // amethyst
	{0, 0, 0, 0, 0},       // euro
}

// Scatters payment.
var ScatPay = [5]int{0, 0, 2, 10, 50} // 8 euro

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

func NewGame(ri string) *Game {
	return &Game{
		Slot5x3: game.Slot5x3{
			RI:  ri,
			BLI: "nvm10",
			SBL: game.MakeSBL(1),
			Bet: 1,
		},
	}
}

const scat = 8

func (g *Game) Scanner(screen game.Screen, ws *game.WinScan) {
	g.ScanLined(screen, ws)
	g.ScanScatters(screen, ws)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen game.Screen, ws *game.WinScan) {
	var bl = game.BetLines5x[g.BLI]
	for li := g.SBL.Next(0); li != 0; li = g.SBL.Next(li) {
		var line = bl.Line(li)

		var syml = screen.At(3, line.At(3))
		var xy game.Line5x
		var count = 1
		xy.Set(3, line.At(3))
		if screen.At(2, line.At(2)) == syml {
			xy.Set(2, line.At(2))
			count++
			if screen.At(1, line.At(1)) == syml {
				xy.Set(1, line.At(1))
				count++
			}
		}
		if screen.At(4, line.At(4)) == syml {
			xy.Set(4, line.At(4))
			count++
			if screen.At(5, line.At(5)) == syml {
				xy.Set(5, line.At(5))
				count++
			}
		}

		if pay := LinePay[syml-1][count-1]; pay > 0 {
			ws.Wins = append(ws.Wins, game.WinItem{
				Pay:  g.Bet * pay,
				Mult: 1,
				Sym:  syml,
				Num:  count,
				Line: li,
				XY:   &xy,
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen game.Screen, ws *game.WinScan) {
	var xy game.Line5x
	var count = 0
	for x := 1; x <= 5; x++ {
		for y := 1; y <= 3; y++ {
			if screen.At(x, y) == scat {
				xy.Set(x, y)
				count++
				break
			}
		}
	}

	if count > 0 {
		if pay := ScatPay[count-1]; pay > 0 {
			ws.Wins = append(ws.Wins, game.WinItem{
				Pay:  g.Bet * pay, // independent from selected lines
				Mult: 1,
				Sym:  scat,
				Num:  count,
				XY:   &xy,
			})
		}
	}
}

func (g *Game) Spin(screen game.Screen) {
	screen.Spin(ReelsMap[g.RI])
}
