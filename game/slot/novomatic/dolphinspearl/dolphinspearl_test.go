package dolphinspearl_test

import (
	"bytes"
	_ "embed"
	"math/rand/v2"
	"testing"

	"github.com/slotopol/server/game"
	"github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/game/slot/novomatic/dolphinspearl"
)

// go test -v -bench=^BenchmarkSpin$ -run=^$ -benchmem -count=5 ./game/slot/novomatic/dolphinspearl

//go:embed dolphinspearl_data.yaml
var data []byte

func BenchmarkSpin(b *testing.B) {
	game.MustReadChain(bytes.NewReader(data))
	var rng = rand.New(game.NewWyRandSource(0))
	var g = dolphinspearl.NewGame(10)
	var wins = make(slot.Wins, 0, 10)

	b.ResetTimer()
	for b.Loop() {
		g.SpinReels(rng, g.GetReels(92))
		g.Scanner(&wins)
		wins.Reset()
	}
}
