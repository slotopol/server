package dolphinspearl_test

import (
	"testing"

	"github.com/slotopol/server/game"
	"github.com/slotopol/server/game/dolphinspearl"
)

// go test -v -bench ^BenchmarkSpin$ -benchmem -count=5 -cover ./slot/dolphinspearl

func BenchmarkSpin(b *testing.B) {
	var g = dolphinspearl.NewGame()
	var screen = g.NewScreen()
	defer screen.Free()
	var wins = make(game.Wins, 0, 10)

	b.ResetTimer()
	for range b.N {
		g.Spin(screen, 92)
		g.Scanner(screen, &wins)
		wins.Reset()
	}
}
