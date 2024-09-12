package ultrahot

import (
	"math"

	slot "github.com/slotopol/server/game/slot"
)

// reels lengths [37, 37, 37], total reshuffles 50653
// RTP = 88.227746%
var Reels88 = slot.Reels3x{
	{8, 8, 8, 1, 6, 6, 6, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 2, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 5, 5, 5},
	{8, 8, 8, 1, 6, 6, 6, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 2, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 5, 5, 5},
	{8, 8, 8, 1, 6, 6, 6, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 2, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 5, 5, 5},
}

// reels lengths [44, 44, 44], total reshuffles 85184
// RTP = 90.275169%
var Reels90 = slot.Reels3x{
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 8, 5, 5, 5, 2},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 8, 5, 5, 5, 2},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 8, 5, 5, 5, 2},
}

// reels lengths [43, 43, 43], total reshuffles 79507
// RTP = 92.117675%
var Reels92 = slot.Reels3x{
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 8, 5, 5, 5, 2},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 8, 5, 5, 5, 2},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 8, 5, 5, 5, 2},
}

// reels lengths [40, 40, 40], total reshuffles 64000
// RTP = 93.062500%
var Reels93 = slot.Reels3x{
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 5, 5, 5},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 5, 5, 5},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 5, 5, 5},
}

// reels lengths [43, 43, 43], total reshuffles 79507
// RTP = 95.658244%
var Reels96 = slot.Reels3x{
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 5, 5, 5, 2},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 5, 5, 5, 2},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 5, 5, 5, 2},
}

// reels lengths [45, 45, 45], total reshuffles 91125
// RTP = 97.816187%
var Reels98 = slot.Reels3x{
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 3, 3, 5, 5, 5, 2},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 3, 3, 5, 5, 5, 2},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 3, 3, 5, 5, 5, 2},
}

// reels lengths [44, 44, 44], total reshuffles 85184
// RTP = 110.648713%
var Reels111 = slot.Reels3x{
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 2, 5, 5, 5, 2},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 2, 5, 5, 5, 2},
	{8, 8, 8, 6, 6, 6, 1, 3, 2, 4, 4, 4, 7, 7, 7, 2, 8, 8, 8, 3, 3, 3, 6, 6, 6, 2, 5, 5, 5, 4, 4, 4, 7, 7, 7, 1, 2, 8, 3, 2, 5, 5, 5, 2},
}

// Map with available reels.
var ReelsMap = map[float64]*slot.Reels3x{
	88.227746:  &Reels88,
	90.275169:  &Reels90,
	92.117675:  &Reels92,
	93.062500:  &Reels93,
	95.658244:  &Reels96,
	97.816187:  &Reels98,
	110.648713: &Reels111,
}

func FindReels(mrtp float64) (rtp float64, reels slot.Reels) {
	for p, r := range ReelsMap {
		if math.Abs(mrtp-p) < math.Abs(mrtp-rtp) {
			rtp, reels = p, r
		}
	}
	return
}

// Lined payment.
var LinePay = [8][3]float64{
	{0, 0, 750}, // 1 seven
	{0, 0, 200}, // 2 star
	{0, 0, 60},  // 3 bar
	{0, 0, 40},  // 4 plum
	{0, 0, 40},  // 5 orange
	{0, 0, 40},  // 6 lemon
	{0, 0, 40},  // 7 cherry
	{0, 0, 5},   // 8 x
}

type Game struct {
	slot.Slot3x3 `yaml:",inline"`
}

func NewGame() *Game {
	return &Game{
		Slot3x3: slot.Slot3x3{
			SBL: slot.MakeBitNum(5),
			Bet: 1,
		},
	}
}

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	var bl = slot.BetLinesHot3
	for li := g.SBL.Next(0); li != 0; li = g.SBL.Next(li) {
		var line = bl.Line(li)
		var fm float64 = 1 // fill mult
		if sym := screen.FillSym(); sym >= 4 && sym <= 7 {
			fm = 2
		}
		var sym1, sym2, sym3 = screen.Pos(1, line), screen.Pos(2, line), screen.Pos(3, line)
		if sym1 == sym2 && sym1 == sym3 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * LinePay[sym1-1][2],
				Mult: fm,
				Sym:  sym1,
				Num:  3,
				Line: li,
				XY:   line.CopyL(3),
			})
		}
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	var _, reels = FindReels(mrtp)
	screen.Spin(reels)
}

func (g *Game) SetLines(sbl slot.Bitset) error {
	return slot.ErrNoFeature
}
