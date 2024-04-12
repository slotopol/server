package ultrahot

import (
	"github.com/slotopol/server/game"
)

// reels lengths [37, 37, 37], total reshuffles 50653
// RTP = 88.227746%
var Reels88 = game.Reels3x{
	{8, 8, 8, 1, 6, 6, 6, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 2, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 5, 5, 5},
	{8, 8, 8, 1, 6, 6, 6, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 2, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 5, 5, 5},
	{8, 8, 8, 1, 6, 6, 6, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 2, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 5, 5, 5},
}

// reels lengths [44, 44, 44], total reshuffles 85184
// RTP = 90.275169%
var Reels90 = game.Reels3x{
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 8, 5, 5, 5, 2},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 8, 5, 5, 5, 2},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 8, 5, 5, 5, 2},
}

// reels lengths [43, 43, 43], total reshuffles 79507
// RTP = 92.117675%
var Reels92 = game.Reels3x{
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 8, 5, 5, 5, 2},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 8, 5, 5, 5, 2},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 8, 5, 5, 5, 2},
}

// reels lengths [40, 40, 40], total reshuffles 64000
// RTP = 93.062500%
var Reels93 = game.Reels3x{
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 5, 5, 5},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 5, 5, 5},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 5, 5, 5},
}

// reels lengths [43, 43, 43], total reshuffles 79507
// RTP = 95.658244%
var Reels96 = game.Reels3x{
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 5, 5, 5, 2},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 5, 5, 5, 2},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 5, 5, 5, 2},
}

// reels lengths [45, 45, 45], total reshuffles 91125
// RTP = 97.816187%
var Reels98 = game.Reels3x{
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 3, 3, 5, 5, 5, 2},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 3, 3, 5, 5, 5, 2},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 3, 3, 5, 5, 5, 2},
}

// reels lengths [44, 44, 44], total reshuffles 85184
// RTP = 110.648713%
var Reels111 = game.Reels3x{
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 2, 5, 5, 5, 2},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 2, 5, 5, 5, 2},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 2, 5, 5, 5, 2},
}

// Map with available reels.
var ReelsMap = map[string]*game.Reels3x{
	"88":  &Reels88,
	"90":  &Reels90,
	"92":  &Reels92,
	"93":  &Reels93,
	"96":  &Reels96,
	"98":  &Reels98,
	"111": &Reels111,
}

// Lined payment.
var LinePay = [13][3]float64{
	{0, 0, 750}, // 1 seven 1
	{0, 0, 200}, // 2 star 4
	{0, 0, 60},  // 3 bar 3
	{0, 0, 40},  // 4 plum 6
	{0, 0, 40},  // 5 orange 6
	{0, 0, 40},  // 6 lemon 6
	{0, 0, 40},  // 7 cherry 6
	{0, 0, 5},   // 8 x 6
}

type Game struct {
	game.Slot3x3 `yaml:",inline"`
}

func NewGame(rd string) *Game {
	return &Game{
		Slot3x3: game.Slot3x3{
			RD:  rd,
			SBL: game.MakeSblNum(5),
			Bet: 1,
		},
	}
}

func (g *Game) Scanner(screen game.Screen, ws *game.WinScan) {
	var bl = game.BetLinesHot
	for li := g.SBL.Next(0); li != 0; li = g.SBL.Next(li) {
		var line = bl.Line(li)
		var fm float64 = 1 // fill mult
		if sym := screen.FillSym(); sym >= 4 && sym <= 7 {
			fm = 2
		}
		var sym1, sym2, sym3 = screen.At(1, line.At(1)), screen.At(2, line.At(2)), screen.At(3, line.At(3))
		if sym1 == sym2 && sym1 == sym3 {
			ws.Wins = append(ws.Wins, game.WinItem{
				Pay:  g.Bet * LinePay[sym1-1][2],
				Mult: fm,
				Sym:  sym1,
				Num:  3,
				Line: li,
				XY:   line.CopyN(3),
			})
		}
	}
}

func (g *Game) Spin(screen game.Screen) {
	screen.Spin(ReelsMap[g.RD])
}

func (g *Game) SetLines(sbl game.SBL) error {
	return game.ErrNoFeature
}

func (g *Game) SetReels(rd string) error {
	if _, ok := ReelsMap[rd]; !ok {
		return game.ErrNoReels
	}
	g.RD = rd
	return nil
}
