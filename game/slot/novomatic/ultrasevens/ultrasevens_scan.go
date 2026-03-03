package ultrasevens

import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame(sp.Sel)
	var s = slot.NewStatGeneric(sn, 5)
	s.JackDim(ssj3)

	var calc = func(w io.Writer) (float64, float64) {
		var N, S, Q = s.NSQ(g.Cost())
		var µ = S / N
		var sigma = math.Sqrt(Q/N - µ*µ)
		fmt.Fprintf(w, "jackpots1: count %g, frequency 1/%.12g\n", s.JackHits(ssj1), N/s.JackHits(ssj1))
		fmt.Fprintf(w, "jackpots2: count %g, frequency 1/%.12g\n", s.JackHits(ssj2), N/s.JackHits(ssj2))
		fmt.Fprintf(w, "jackpots3: count %g, frequency 1/%.12g\n", s.JackHits(ssj3), N/s.JackHits(ssj3))
		fmt.Fprintf(w, "RTP = %.6f%%\n", µ*100)
		return µ, sigma
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
}
