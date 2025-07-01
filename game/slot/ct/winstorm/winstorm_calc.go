package winstorm

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, mrtp float64) float64 {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	var g = NewGame()
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var reshuf1 = s.Reshuf(1)
		var reshuf2 = s.Reshuf(2)
		var reshuf3 = s.Reshuf(3)
		var reshuf4 = s.Reshuf(4)
		var reshuf5 = s.Reshuf(5)
		var cost, _ = g.Cost()
		var lrtp, srtp = s.LineRTP(cost), s.ScatRTP(cost)
		var rtpsym = lrtp + srtp
		fmt.Fprintf(w, "fall[2] = %.10g, freq = 1/%.5g\n", reshuf2, reshuf1/reshuf2)
		fmt.Fprintf(w, "fall[3] = %.10g, freq = 1/%.5g, freq2 = 1/%.5g\n", reshuf3, reshuf1/reshuf3, reshuf2/reshuf3)
		fmt.Fprintf(w, "fall[4] = %.10g, freq = 1/%.5g, freq2 = 1/%.5g\n", reshuf4, reshuf1/reshuf4, reshuf3/reshuf4)
		fmt.Fprintf(w, "fall[5] = %.10g, freq = 1/%.5g, freq2 = 1/%.5g\n", reshuf5, reshuf1/reshuf5, reshuf4/reshuf5)
		fmt.Fprintf(w, "RTP = %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		return rtpsym
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc)
}
