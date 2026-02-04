package chicago

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
		var cost = g.Cost()
		var N, S, Q = s.NSQ(cost)
		var lrtp, srtp = s.SymRTP(cost)
		var rtpsym = lrtp + srtp
		var q, sq = s.FSQ()
		var rtpfs = EVmc * sq * rtpsym
		var rtp = rtpsym + q*rtpfs
		var sigma = math.Sqrt(sq*(Q/N-S*S/N/N) + q*(S/N*sq)*(S/N*sq))
		var vi90 = slot.GetZ(0.90) * sigma
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCount.Load(), q, sq)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", s.FGF())
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym, q, rtpfs, rtp)
		fmt.Fprintf(w, "sigma = %.6g, VI[90%%] = %.6g (%s)\n", sigma, vi90, slot.VIname6[slot.VIclass6(vi90)])
		fmt.Fprintf(w, "CI[90%%] = %d, CI[68.27%%] = %d, CI[95.45%%] = %d, CI[99.73%%] = %d\n", int(slot.CI(0.90, rtp/100, sigma)), int(slot.CI(slot.CP(1), rtp/100, sigma)), int(slot.CI(slot.CP(2), rtp/100, sigma)), int(slot.CI(slot.CP(3), rtp/100, sigma)))
		return rtp
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}
