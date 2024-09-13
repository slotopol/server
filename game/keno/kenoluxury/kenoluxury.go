package kenoluxury

import (
	keno "github.com/slotopol/server/game/keno"
)

// RTP[ 2] = 92.088608%
// RTP[ 3] = 92.964946%
// RTP[ 4] = 92.027909%
// RTP[ 5] = 92.575146%
// RTP[ 6] = 92.668091%
// RTP[ 7] = 92.443234%
// RTP[ 8] = 92.751742%
// RTP[ 9] = 92.001419%
// RTP[10] = 92.551062%
// RTP[game] = 92.452462%
var Paytable = keno.KenoPaytable{
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},              //  0 sel
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},              //  1 sel
	{0, 1, 9, 0, 0, 0, 0, 0, 0, 0, 0},              //  2 sel
	{0, 0, 2, 47, 0, 0, 0, 0, 0, 0, 0},             //  3 sel
	{0, 0, 2, 5, 91, 0, 0, 0, 0, 0, 0},             //  4 sel
	{0, 0, 0, 3, 12, 820, 0, 0, 0, 0, 0},           //  5 sel
	{0, 0, 0, 3, 4, 70, 1600, 0, 0, 0, 0},          //  6 sel
	{0, 0, 0, 1, 2, 21, 400, 7000, 0, 0, 0},        //  7 sel
	{0, 0, 0, 0, 2, 12, 100, 1650, 10000, 0, 0},    //  8 sel
	{0, 0, 0, 0, 1, 6, 44, 335, 4700, 10000, 0},    //  9 sel
	{0, 0, 0, 0, 0, 5, 24, 142, 1000, 4500, 10000}, // 10 sel
}

type Game struct {
	keno.Keno80 `yaml:",inline"`
}

func NewGame() *Game {
	return &Game{
		Keno80: keno.Keno80{
			Bet: 1,
		},
	}
}

func (g *Game) Scanner(scrn *keno.Screen, wins *keno.Wins) {
	for i := range 80 {
		if scrn[i] == keno.KSselhit {
			wins.Num++
		}
	}
	wins.Pay = Paytable[len(g.Sel)][wins.Num] * g.Bet
}