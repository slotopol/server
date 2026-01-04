package africansimba_test

import (
	"testing"

	"github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/game/slot/novomatic/africansimba"
)

type preset struct {
	Scr [3][5]slot.Sym
	Win float64
	Num int
}

func (p *preset) Setup(s *slot.Screen5x3) {
	for x := range 5 {
		for y := range 3 {
			s.Scr[x][y] = p.Scr[y][x]
		}
	}
}

var presets = []preset{
	{ // #1
		Scr: [3][5]slot.Sym{ // 4x2 buffalo, 3x1 ace, scatters
			{2, 7, 1, 3, 5},
			{4, 4, 2, 8, 6},
			{7, 6, 4, 4, 2},
		},
		Win: 150*2 + 10,
		Num: 3,
	},
	{ // #2
		Scr: [3][5]slot.Sym{ // 5x2 lemur, 3x1 ace, 3x2 king
			{5, 7, 1, 3, 5},
			{8, 5, 8, 5, 8},
			{7, 8, 5, 4, 2},
		},
		Win: 250*2 + 10 + 10*2,
		Num: 3,
	},
	{ // #3
		Scr: [3][5]slot.Sym{ // 5x4 lemur, 5x1 king
			{5, 7, 2, 3, 5},
			{8, 5, 8, 5, 8},
			{7, 1, 5, 1, 2},
		},
		Win: 250*4 + 125,
		Num: 2,
	},
	{ // #4
		Scr: [3][5]slot.Sym{ // 4x2 ace, 2x1 queen
			{2, 7, 6, 3, 5},
			{9, 5, 1, 5, 8},
			{7, 1, 5, 1, 2},
		},
		Win: 25*2 + 25,
		Num: 2, // scatters by ways should be ignored
	},
	{ // #5
		Scr: [3][5]slot.Sym{ // no wins
			{6, 7, 7, 3, 5},
			{3, 5, 4, 5, 8},
			{5, 1, 2, 1, 2},
		},
		Win: 0,
		Num: 0,
	},
	{ // #6
		Scr: [3][5]slot.Sym{ // 4x12 buffalo
			{6, 4, 4, 3, 5},
			{4, 5, 4, 4, 8},
			{5, 1, 4, 1, 2},
		},
		Win: 150 * 12,
		Num: 1,
	},
	{ // #7
		Scr: [3][5]slot.Sym{ // 5x243 giraffe
			{3, 3, 1, 3, 3},
			{3, 3, 3, 3, 3},
			{3, 3, 3, 3, 3},
		},
		Win: 2500 * 243,
		Num: 1,
	},
}

func TestWays(t *testing.T) {
	var g = africansimba.NewGame()
	var wins slot.Wins
	for i, p := range presets {
		p.Setup(&g.Screen5x3)
		g.Scanner(&wins)
		var num = len(wins)
		if num != p.Num {
			t.Errorf("error at %d screen, expected %d, gets %d wins", i+1, p.Num, num)
		}
		var gain = wins.Gain()
		if gain != p.Win {
			t.Errorf("error at %d screen, expected %g, gets %g gain", i+1, p.Win, gain)
		}
		wins.Reset()
	}
}
