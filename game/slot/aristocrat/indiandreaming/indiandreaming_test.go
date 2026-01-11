package indiandreaming_test

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/slotopol/server/game"
	"github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/game/slot/aristocrat/indiandreaming"
)

type preset struct {
	Scr    [3][5]slot.Sym
	Num    int
	WinReg float64
	WinBon float64
}

func (p *preset) Setup(s *slot.Screen5x3) {
	for x := range 5 {
		for y := range 3 {
			s.Scr[x][y] = p.Scr[y][x]
		}
	}
}

var presets = []preset{
	{ // #1: way on 4 symbols
		Scr: [3][5]slot.Sym{
			{3, 7, 1, 3, 5},
			{4, 4, 3, 8, 6},
			{7, 6, 4, 4, 2},
		},
		// 4x2 sym4, 3x1 sym7
		Num:    2,
		WinReg: 100*2 + 6,
		WinBon: 100*5 + 100 + 6*5,
	},
	{ // #2: way on 5 symbols
		Scr: [3][5]slot.Sym{
			{5, 7, 1, 3, 5},
			{8, 5, 8, 5, 8},
			{7, 8, 5, 4, 2},
		},
		// 5x2 sym5, 3x1 sym7, 3x2 sym8
		Num:    3,
		WinReg: 1000*2 + 6 + 6*2,
		WinBon: 1000 + 1000*5 + 6*5 + 6 + 6*5,
	},
	{ // #3: way compiled with wilds
		Scr: [3][5]slot.Sym{
			{5, 7, 2, 3, 5},
			{8, 5, 8, 5, 8},
			{7, 1, 5, 1, 2},
		},
		// 5x4 sym5, 5x1 sym8, 4 scatters
		Num:    3,
		WinReg: 1000*4 + 150 + 15,
		WinBon: 1000 + 1000*3*5 + 150*5 + 15*5,
	},
	{ // #4: wilded scatters (scatters by ways should be ignored)
		Scr: [3][5]slot.Sym{
			{2, 7, 6, 3, 5},
			{9, 5, 1, 5, 8},
			{7, 1, 5, 1, 2},
		},
		// 4x2 sym7, 4x1 sym9, 5 scatters
		Num:    3,
		WinReg: 25*2 + 15 + 100,
		WinBon: 25*2*5 + 15*5 + 100*5,
	},
	{ // #5: no pays by ways
		Scr: [3][5]slot.Sym{
			{6, 7, 7, 3, 5},
			{3, 5, 4, 5, 8},
			{5, 1, 2, 1, 2},
		},
		// 4 scatters
		Num:    1,
		WinReg: 15,
		WinBon: 15 * 5,
	},
	{ // #6: multiways for one symbol
		Scr: [3][5]slot.Sym{
			{6, 4, 4, 3, 5},
			{4, 5, 4, 4, 8},
			{5, 1, 4, 1, 9},
		},
		// 4x12 sym4
		Num:    1,
		WinReg: 100 * 12,
		WinBon: 100*3 + 100*9*5,
	},
	{ // #7: filled screen
		Scr: [3][5]slot.Sym{
			{3, 3, 1, 3, 3},
			{3, 3, 3, 3, 3},
			{3, 3, 3, 3, 3},
		},
		// 5x243 sym3
		Num:    1,
		WinReg: 5000 * 243,
		WinBon: 5000*162 + 5000*81*5,
	},
	{ // #8: wild on 1, 3 reels
		Scr: [3][5]slot.Sym{
			{1, 5, 3, 6, 3},
			{5, 7, 5, 9, 3},
			{7, 4, 1, 8, 3},
		},
		// 3x1 sym4, 3x4 sym5, 3x2 sym7
		Num:    3,
		WinReg: 50 + 50*4 + 6*2,
		WinBon: 50*5 + 50 + 50*3*5 + 6*2*5,
	},
	{ // #9: continuous wilds
		Scr: [3][5]slot.Sym{
			{1, 5, 7, 6, 4},
			{5, 1, 5, 9, 7},
			{7, 7, 1, 8, 5},
		},
		// 3x8 sym5, 3x8 sym7, 4x1 sym6, 4x1 sym8, 4x1 sym9, 3 scatters
		Num:    6,
		WinReg: 50*7 + 6*7 + 40 + 25 + 15 + 2,
		WinBon: 50 + 50*6*5 + 6 + 6*6*5 + 40*5 + 25*5 + 15*5 + 2*5,
	},
}

// go test -run=^TestScanner$ ./game/slot/aristocrat/indiandreaming

func TestScanner(t *testing.T) {
	var g = indiandreaming.NewGame()
	var wins slot.Wins
	var gain float64
	for i, p := range presets {
		p.Setup(&g.Screen5x3)
		g.FSR = 0 // set regular spins mode
		g.Scanner(&wins)
		if len(wins) != p.Num {
			t.Errorf("error at %d screen on regular spins, expected %d, gets %d wins", i+1, p.Num, len(wins))
		}
		gain = wins.Gain()
		if gain != p.WinReg {
			t.Errorf("error at %d screen on regular spins, expected %g, gets %g gain", i+1, p.WinReg, gain)
		}
		wins.Reset()
		g.FSR = 12 // set free spins mode
		g.Scanner(&wins)
		gain = wins.Gain()
		if len(wins) != p.Num {
			t.Errorf("error at %d screen on free spins, expected %d, gets %d wins", i+1, p.Num, len(wins))
		}
		if gain != p.WinBon {
			t.Errorf("error at %d screen on free spins, expected %g, gets %g gain", i+1, p.WinBon, gain)
		}
		wins.Reset()
	}
}

// go test -v -bench=^BenchmarkSpin$ -run=^$ -benchmem -count=5 ./game/slot/aristocrat/indiandreaming

//go:embed indiandreaming_data.yaml
var data []byte

func BenchmarkSpin(b *testing.B) {
	game.MustReadChain(bytes.NewReader(data))
	var g = indiandreaming.NewGame()
	var wins = make(slot.Wins, 0, 10)

	b.ResetTimer()
	for range b.N {
		g.Spin(95)
		g.Scanner(&wins)
		wins.Reset()
	}
}
