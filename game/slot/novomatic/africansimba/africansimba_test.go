package africansimba_test

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/slotopol/server/game"
	"github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/game/slot/novomatic/africansimba"
)

type preset struct {
	Grid [3][5]slot.Sym
	Num  int
	Win  float64
}

func (p *preset) Setup(s *slot.Grid5x3) {
	for x := range 5 {
		for y := range 3 {
			s.Grid[x][y] = p.Grid[y][x]
		}
	}
}

var presets = []preset{
	{ // #1: way on 4 symbols
		Grid: [3][5]slot.Sym{
			{2, 7, 1, 3, 5},
			{4, 4, 2, 8, 6},
			{7, 6, 4, 4, 2},
		},
		// 4x2 sym4, 3x1 sym7, scatters
		Num: 3,
		Win: 150*2 + 10,
	},
	{ // #2: way on 5 symbols
		Grid: [3][5]slot.Sym{
			{5, 7, 1, 3, 5},
			{8, 5, 8, 5, 8},
			{7, 8, 5, 4, 2},
		},
		// 5x2 sym5, 3x1 sym7, 3x2 sym8
		Num: 3,
		Win: 250*2 + 10 + 10*2,
	},
	{ // #3: way compiled with wilds
		Grid: [3][5]slot.Sym{
			{5, 7, 2, 3, 5},
			{8, 5, 8, 5, 8},
			{7, 1, 5, 1, 2},
		},
		// 5x4 sym5, 5x1 sym8
		Num: 2,
		Win: 250*4 + 125,
	},
	{ // #4: wilded scatters (scatters by ways should be ignored)
		Grid: [3][5]slot.Sym{
			{2, 7, 6, 3, 5},
			{9, 5, 1, 5, 8},
			{7, 1, 5, 1, 2},
		},
		// 4x2 sym7, 4x1 sym9
		Num: 2,
		Win: 25*2 + 25,
	},
	{ // #5: no pays by ways
		Grid: [3][5]slot.Sym{
			{6, 7, 7, 3, 5},
			{3, 5, 4, 5, 8},
			{5, 1, 2, 1, 2},
		},
		// no wins
		Num: 0,
		Win: 0,
	},
	{ // #6: multiways for one symbol
		Grid: [3][5]slot.Sym{
			{6, 4, 4, 3, 5},
			{4, 5, 4, 4, 8},
			{5, 1, 4, 1, 2},
		},
		// 4x12 sym4
		Num: 1,
		Win: 150 * 12,
	},
	{ // #7: filled grid
		Grid: [3][5]slot.Sym{
			{3, 3, 1, 3, 3},
			{3, 3, 3, 3, 3},
			{3, 3, 3, 3, 3},
		},
		// 5x243 sym3
		Num: 1,
		Win: 2500 * 243,
	},
}

// go test -run=^TestScanner$ ./game/slot/novomatic/africansimba

func TestScanner(t *testing.T) {
	var g = africansimba.NewGame()
	var wins slot.Wins
	for i, p := range presets {
		p.Setup(&g.Grid5x3)
		g.Scanner(&wins)
		if len(wins) != p.Num {
			t.Errorf("error at %d grid, expected %d, gets %d wins", i+1, p.Num, len(wins))
		}
		var gain = wins.Gain()
		if gain != p.Win {
			t.Errorf("error at %d grid, expected %g, gets %g gain", i+1, p.Win, gain)
		}
		wins.Reset()
	}
}

// go test -v -bench=^BenchmarkSpin$ -run=^$ -benchmem -count=5 ./game/slot/novomatic/africansimba

//go:embed africansimba_data.yaml
var data []byte

func BenchmarkSpin(b *testing.B) {
	game.MustReadChain(bytes.NewReader(data))
	var g = africansimba.NewGame()
	var wins = make(slot.Wins, 0, 10)

	b.ResetTimer()
	for range b.N {
		g.Spin(96)
		g.Scanner(&wins)
		wins.Reset()
	}
}
