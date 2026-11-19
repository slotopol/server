package ultrasevens

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var g = NewGame(sp.Sel)
	var s = slot.NewStatSnap(1, sn, 5, 0, ssj3)

	var calc = func(w io.Writer) (float64, float64) {
		if sp.IsJack() {
			var N = s.Count()
			for idj := range s.A[0].Cj {
				var Cj = float64(s.A[0].Cj[idj])
				var HRj = N / Cj
				fmt.Fprintf(w, "jackpots%d: count %g, hit rate 1/%.12g\n", idj+1, Cj, HRj)
			}
		}
		return slot.Parsheet_simple(w, sp, s, g.Cost())
	}

	return slot.ScanReelsCommon2(ctx, sp, s, g, calc)
}
