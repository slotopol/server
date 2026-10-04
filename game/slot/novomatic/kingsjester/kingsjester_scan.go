package kingsjester

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var g = NewGame(sp.Sel)
	var s = slot.NewStatGeneric(sn, 5)
	s.JackDim(kjj2)

	var calc = func(w io.Writer) (float64, float64) {
		if sp.IsJack() {
			var N = s.Count()
			var q = s.FSQ()
			var sq = 1 / (1 - q)
			for idj := range s.Cj {
				var Cj = float64(s.Cj[idj].Load()) / float64(sp.Sel)
				var HRj = N / Cj * (1 + q*sq)
				fmt.Fprintf(w, "jackpots%d: count per line %g, hit rate 1/%.12g\n", idj+1, Cj, HRj)
			}
		}
		return slot.Parsheet_fgrecur(w, sp, s, g.Cost(), 1, s.Efs)
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, calc)
}
