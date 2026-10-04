package ultrasevens

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var g = NewGame(sp.Sel)
	var s = slot.NewStatGeneric(sn, 5)
	s.JackDim(ssj3)

	var calc = func(w io.Writer) (float64, float64) {
		if sp.IsJack() {
			var N = s.Count()
			for idj := range s.Cj {
				var Cj = float64(s.Cj[idj].Load())
				var HRj = N / Cj
				fmt.Fprintf(w, "jackpots%d: count %g, hit rate 1/%.12g\n", idj+1, Cj, HRj)
			}
		}
		return slot.Parsheet_simple(w, sp, s, g.Cost())
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, calc)
}
