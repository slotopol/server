package kingsjester

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var g = NewGame(sp.Sel)
	var s = slot.NewStatSnap(1, sn, 5, 0, kjj2)

	var calc = func(w io.Writer) (float64, float64) {
		if sp.IsJack() {
			var N = s.Count()
			var _, q, _ = slot.Efs(s)
			var sq = 1 / (1 - q)
			for idj, cj := range s.A[0].Cj {
				var Cj = float64(cj) / float64(sp.Sel)
				var HRj = Cj * (1 + q*sq) / N
				fmt.Fprintf(w, "jackpots%d: count per line %g, hit rate 1/%.12g\n", idj+1, Cj, 1/HRj)
			}
		}
		return slot.Parsheet_solereelrecur(w, sp, s, g.Cost(), 1)
	}

	return slot.ScanReelsCommon2(ctx, sp, s, g, calc)
}
