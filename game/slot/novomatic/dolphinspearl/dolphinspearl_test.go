package dolphinspearl_test

import (
	"testing"

	"github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/game/slot/novomatic/dolphinspearl"
)

// go test -v -bench ^BenchmarkSpin$ -benchmem -count=5 -cover ./game/slot/dolphinspearl

func BenchmarkSpin(b *testing.B) {
	var g = dolphinspearl.NewGame()
	var wins = make(slot.Wins, 0, 10)

	b.ResetTimer()
	for range b.N {
		g.Spin(92)
		g.Scanner(&wins)
		wins.Reset()
	}
}
