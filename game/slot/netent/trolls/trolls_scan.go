package trolls

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
		var lrtp, srtp = s.SymRTP(g.Cost())
		var rtpsym = lrtp + srtp
		var q, sq = s.FSQ()
		var rtpfs = 3 * sq * rtpsym
		var rtp = rtpsym + q*rtpfs
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp*100, srtp*100, rtpsym*100)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCount.Load(), q, sq)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", s.FGF())
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym*100, q, rtpfs*100, rtp*100)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}
