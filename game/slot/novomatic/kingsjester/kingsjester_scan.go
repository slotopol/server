package kingsjester

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
		var reshuf = s.Count()
		var lrtp, srtp = s.SymRTP(g.Cost())
		var rtpsym = lrtp + srtp
		var q, sq = s.FSQ()
		var rtp = rtpsym + q*sq*rtpsym
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCount.Load(), q, sq)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", s.FGF())
		fmt.Fprintf(w, "jackpots1: count %g, frequency 1/%.12g\n", s.JackHitsF(kjj1), reshuf/s.JackHitsF(kjj1))
		fmt.Fprintf(w, "jackpots2: count %g, frequency 1/%.12g\n", s.JackHitsF(kjj2), reshuf/s.JackHitsF(kjj2))
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g*%.5g(sym) = %.6f%%\n", rtpsym, q, sq, rtpsym, rtp)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}
