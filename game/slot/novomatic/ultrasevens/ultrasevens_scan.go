package ultrasevens

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, mrtp float64) float64 {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	var g = NewGame(1)
	var s slot.StatGeneric

	var calc = func(w io.Writer) float64 {
		var N = s.Count()
		var lrtp, srtp = s.SymRTP(g.Cost())
		if srtp > 0 {
			panic("scatters have no pays")
		}
		var rtpsym = lrtp + srtp
		fmt.Fprintf(w, "jackpots1: count %g, frequency 1/%.12g\n", s.JackHitsF(ssj1), N/s.JackHitsF(ssj1))
		fmt.Fprintf(w, "jackpots2: count %g, frequency 1/%.12g\n", s.JackHitsF(ssj2), N/s.JackHitsF(ssj2))
		fmt.Fprintf(w, "jackpots3: count %g, frequency 1/%.12g\n", s.JackHitsF(ssj3), N/s.JackHitsF(ssj3))
		fmt.Fprintf(w, "RTP = %.6f%%\n", rtpsym*100)
		return rtpsym
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}
