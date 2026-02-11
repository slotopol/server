package fruitilicious

import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) float64 {
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame(sp.Sel)
	var s slot.StatGeneric

	var calc = func(w io.Writer) float64 {
		var N, S, Q = s.NSQ(g.Cost())
		var rtp = S / N
		var sigma = math.Sqrt(Q/N - rtp*rtp)
		var vi95 = slot.GetZ(0.90) * sigma
		fmt.Fprintf(w, "RTP = %.6f%%\n", rtp*100)
		fmt.Fprintf(w, "sigma = %.6g, VI[95%%] = %.6g (%s)\n", sigma, vi95, slot.VIname6[slot.VIclass6(vi95)])
		fmt.Fprintf(w, "CI[90%%] = %d, CI[95%%] = %d, CI[99%%] = %d\n",
			int(slot.CI(0.90, rtp, sigma)), int(slot.CI(0.95, rtp, sigma)), int(slot.CI(0.99, rtp, sigma)))
		return rtp
	}

	return slot.ScanReelsCommon(ctx, sp, &s, g, reels, calc)
}
