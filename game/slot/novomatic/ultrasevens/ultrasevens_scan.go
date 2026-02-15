package ultrasevens

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) float64 {
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame(sp.Sel)
	var s = slot.NewStatGeneric(sn, 5)

	var calc = func(w io.Writer) float64 {
		var N, S, _ = s.NSQ(g.Cost())
		var rtpsym = S / N
		fmt.Fprintf(w, "jackpots1: count %g, frequency 1/%.12g\n", s.JackHitsF(ssj1), N/s.JackHitsF(ssj1))
		fmt.Fprintf(w, "jackpots2: count %g, frequency 1/%.12g\n", s.JackHitsF(ssj2), N/s.JackHitsF(ssj2))
		fmt.Fprintf(w, "jackpots3: count %g, frequency 1/%.12g\n", s.JackHitsF(ssj3), N/s.JackHitsF(ssj3))
		fmt.Fprintf(w, "RTP = %.6f%%\n", rtpsym*100)
		return rtpsym
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
}
