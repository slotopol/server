package fruitsensation

import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, mrtp float64) float64 {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	var g = NewGame(1)
	var s slot.StatGeneric

	var calc = func(w io.Writer) float64 {
		var N, S, Q = s.NSQ(g.Cost())
		var rtp = S / N
		var sigma = math.Sqrt(Q/N - rtp*rtp)
		var vi90 = slot.GetZ(0.90) * sigma
		fmt.Fprintf(w, "RTP = %.6f%%\n", rtp*100)
		fmt.Fprintf(w, "sigma = %.6g, VI[90%%] = %.6g (%s)\n", sigma, vi90, slot.VIname6[slot.VIclass6(vi90)])
		fmt.Fprintf(w, "CI[90%%] = %d, CI[68.27%%] = %d, CI[95.45%%] = %d, CI[99.73%%] = %d\n",
			int(slot.CI(0.90, rtp, sigma)), int(slot.CI(slot.CP(1), rtp, sigma)), int(slot.CI(slot.CP(2), rtp, sigma)), int(slot.CI(slot.CP(3), rtp, sigma)))
		return rtp
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}
