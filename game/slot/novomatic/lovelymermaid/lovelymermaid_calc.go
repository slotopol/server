package lovelymermaid

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
		var rtpsym = lrtp + srtp
		var q = s.FreeCount() / reshuf
		var sq = 1 / (1 - q)
		var rtp = rtpsym + q*sq*rtpsym
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCountU(), q, sq)
		fmt.Fprintf(w, "free games frequency: 1/%.5g\n", reshuf/s.FreeHits())
		fmt.Fprintf(w, "jackpots: count %g, frequency 1/%.12g\n", s.JackCount(lmj), reshuf/s.JackCount(lmj))
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g*%.5g(sym) = %.6f%%\n", rtpsym, q, sq, rtpsym, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc)
}
