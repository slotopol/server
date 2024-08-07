package dolphinspearl_test

import (
	"testing"

	"github.com/slotopol/server/game"
	"github.com/slotopol/server/game/dolphinspearl"
)

// go test -v -bench ^BenchmarkSpin$ -benchmem -count=5 -cover ./game/dolphinspearl

func BenchmarkSpin(b *testing.B) {
	var g = dolphinspearl.NewGame("92")
	var screen = g.NewScreen()
	defer screen.Free()
	var wins = make(game.Wins, 0, 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Spin(screen)
		g.Scanner(screen, &wins)
		wins.Reset()
	}
}
