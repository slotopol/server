package redroo_test

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/slotopol/server/game"
	"github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/game/slot/aristocrat/redroo"
)

type preset struct {
	Grid   [4][5]slot.Sym
	MW     [3]float64
	NumReg int
	NumBon int
	WinReg float64
	WinBon float64
}

func (p *preset) Setup(s *slot.Grid5x4) {
	for x := range 5 {
		for y := range 4 {
			s.Grid[x][y] = p.Grid[y][x]
		}
	}
}

var presets = []preset{
	{ // #1: way on 4 symbols
		Grid: [4][5]slot.Sym{
			{3, 7, 1, 3, 5},
			{5, 4, 3, 8, 6},
			{4, 6, 4, 4, 2},
			{7, 8, 5, 5, 5},
		},
		MW: [3]float64{2, 2, 3},
		// 4x2 sym4, 3x1 sym7
		NumReg: 2,
		NumBon: 2,
		WinReg: 150*2 + 40,
		WinBon: 150*2 + 150 + 40*2,
	},
	{ // #2: way on 5 symbols
		Grid: [4][5]slot.Sym{
			{5, 7, 4, 3, 5},
			{8, 5, 8, 5, 8},
			{4, 8, 5, 4, 2},
			{7, 3, 1, 6, 9},
		},
		MW: [3]float64{2, 3, 3},
		// 5x2 sym5, 3x1 sym7, 3x2 sym8
		NumReg: 3,
		NumBon: 3,
		WinReg: 200*2 + 40 + 10*2,
		WinBon: 200 + 200*3 + 40*3 + 10 + 10*3,
	},
	{ // #3: way compiled with wilds
		Grid: [4][5]slot.Sym{
			{5, 7, 4, 3, 5},
			{8, 5, 8, 5, 8},
			{2, 4, 5, 1, 4},
			{9, 1, 3, 6, 9},
		},
		MW: [3]float64{3, 3, 2},
		// 5x4 sym5, 5x1 sym8
		NumReg: 2,
		NumBon: 2,
		WinReg: 200*4 + 140,
		WinBon: 200 + 200*3 + 200*2 + 200*3*2 + 140*3*2,
	},
	{ // #4: wilded scatters (scatters by ways should be ignored)
		Grid: [4][5]slot.Sym{
			{2, 7, 6, 3, 5},
			{9, 5, 1, 5, 8},
			{7, 3, 5, 9, 2},
			{4, 1, 8, 1, 6},
		},
		MW: [3]float64{2, 3, 2},
		// 4x1 sym4, 4x2 sym7, 4x2 sym9, 2 scatters
		NumReg: 3,
		NumBon: 4,
		WinReg: 150 + 100*2 + 50*2,
		WinBon: 150*2*3*2 + 100*3*2 + 100*2*3*2 + 50*2*3 + 50*2*3*2 + 0,
	},
	{ // #5: no pays by ways
		Grid: [4][5]slot.Sym{
			{6, 7, 7, 3, 5},
			{3, 2, 4, 5, 8},
			{5, 4, 1, 6, 4},
			{8, 9, 5, 1, 7},
		},
		MW: [3]float64{2, 2, 2},
		// 4 scatters
		NumReg: 0,
		NumBon: 0,
		WinReg: 0,
		WinBon: 0,
	},
	{ // #6: multiways for one symbol
		Grid: [4][5]slot.Sym{
			{8, 4, 4, 3, 5},
			{4, 5, 4, 4, 8},
			{8, 3, 4, 1, 4},
			{2, 1, 4, 7, 9},
		},
		MW: [3]float64{3, 2, 3},
		// 5x16 sym4
		NumReg: 1,
		NumBon: 1,
		WinReg: 200 * 16,
		WinBon: 200*4 + 200*4*3 + 200*4*3 + 200*4*3*3,
	},
	{ // #7: filled grid with wilds
		Grid: [4][5]slot.Sym{
			{3, 3, 3, 3, 3},
			{3, 1, 3, 3, 3},
			{3, 3, 3, 1, 3},
			{3, 3, 3, 3, 3},
		},
		MW: [3]float64{2, 2, 3},
		// 5x1024 sym3
		NumReg: 1,
		NumBon: 1,
		WinReg: 250 * 1024,
		WinBon: 250*4*3*4*3*4 + 250*4*1*4*3*4*2 + 250*4*3*4*1*4*3 + 250*4*1*4*1*4*2*3,
	},
}

// go test -run=^TestScanner$ ./game/slot/aristocrat/redroo

func TestScanner(t *testing.T) {
	var g = redroo.NewGame()
	var wins slot.Wins
	var gain float64
	for i, p := range presets {
		p.Setup(&g.Grid5x4)
		g.FSR = 0 // set regular spins mode
		g.MW = [3]float64{1, 1, 1}
		g.Scanner(&wins)
		if len(wins) != p.NumReg {
			t.Errorf("error at %d grid on regular spins, expected %d, gets %d wins", i+1, p.NumReg, len(wins))
		}
		gain = wins.Gain()
		if gain != p.WinReg {
			t.Errorf("error at %d grid on regular spins, expected %g, gets %g gain", i+1, p.WinReg, gain)
		}
		wins.Reset()
		g.FSR = 12 // set free spins mode
		g.MW = p.MW
		g.Scanner(&wins)
		gain = wins.Gain()
		if len(wins) != p.NumBon {
			t.Errorf("error at %d grid on free spins, expected %d, gets %d wins", i+1, p.NumBon, len(wins))
		}
		if gain != p.WinBon {
			t.Errorf("error at %d grid on free spins, expected %g, gets %g gain", i+1, p.WinBon, gain)
		}
		wins.Reset()
	}
}

// go test -v -bench=^BenchmarkSpin$ -run=^$ -benchmem -count=5 ./game/slot/aristocrat/redroo

//go:embed redroo_data.yaml
var data []byte

func BenchmarkSpin(b *testing.B) {
	game.MustReadChain(bytes.NewReader(data))
	var g = redroo.NewGame()
	var wins = make(slot.Wins, 0, 10)

	b.ResetTimer()
	for range b.N {
		g.Spin(95)
		g.Scanner(&wins)
		wins.Reset()
	}
}
