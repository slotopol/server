package chicago

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
		var cost = g.Cost()
		var N, S, Q = s.NSQ(cost)
		var lrtp, srtp = s.RTPsym(cost, scat)
		var rtpsym = lrtp + srtp
		var q, sq = s.FSQ()
		var rtpfs = EVmc * sq * rtpsym
		var rtp = rtpsym + q*rtpfs
		var sigma = math.Sqrt(sq*(Q/N-S*S/N/N) + q*(S/N*sq)*(S/N*sq))
		var vi95 = slot.GetZ(0.90) * sigma
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp*100, srtp*100, rtpsym*100)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FSC.Load(), q, sq)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", s.FGF())
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym*100, q, rtpfs*100, rtp*100)
		fmt.Fprintf(w, "sigma = %.6g, VI[95%%] = %.6g (%s)\n", sigma, vi95, slot.VIname6[slot.VIclass6(vi95)])
		fmt.Fprintf(w, "CI[90%%] = %d, CI[95%%] = %d, CI[99%%] = %d\n",
			int(slot.CI(0.90, rtp, sigma)), int(slot.CI(0.95, rtp, sigma)), int(slot.CI(0.99, rtp, sigma)))
		return rtp
	}

	return slot.ScanReelsCommon(ctx, sp, &s, g, reels, calc)
}
