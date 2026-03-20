package lovelymermaid

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame(sp.Sel)
	var s = slot.NewStatGeneric(sn, 5)
	s.JackDim(lmj)

	var calc = func(w io.Writer) (float64, float64) {
		if sp.IsJack() {
			var N = s.Count()
			var q = s.FSQ()
			var sq = 1 / (1 - q)
			var Cj = float64(s.JH[lmj-1].Load())
			var HRj = N / Cj * (1 + q*sq)
			fmt.Fprintf(w, "jackpots: count %g, hit rate 1/%.12g\n", Cj, HRj)
		}
		return slot.Parsheet_fgretrig(w, sp, s, g.Cost(), 1, 25)
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
}
