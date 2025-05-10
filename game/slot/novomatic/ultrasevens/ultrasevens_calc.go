package ultrasevens

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, mrtp float64) float64 {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	var g = NewGame()
	g.Sel = 1
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var reshuf = s.Count()
		var cost, _ = g.Cost()
		var lrtp, srtp = s.LineRTP(cost), s.ScatRTP(cost)
		if srtp > 0 {
			panic("scatters have no pays")
		}
		var rtpsym = lrtp + srtp
		fmt.Fprintf(w, "jackpots1: count %g, frequency 1/%.12g\n", s.JackCount(ssj1), reshuf/s.JackCount(ssj1))
		fmt.Fprintf(w, "jackpots2: count %g, frequency 1/%.12g\n", s.JackCount(ssj2), reshuf/s.JackCount(ssj2))
		fmt.Fprintf(w, "jackpots3: count %g, frequency 1/%.12g\n", s.JackCount(ssj3), reshuf/s.JackCount(ssj3))
		fmt.Fprintf(w, "RTP = %.6f%%\n", rtpsym)
		return rtpsym
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc)
}
