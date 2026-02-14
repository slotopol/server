package lucky3penguins

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) float64 {
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame(sp.Sel)
	var s = slot.NewStatGeneric(sn, 5)

	var calc = func(w io.Writer) float64 {
		var N, S, _ = s.NSQ(g.Cost())
		var rtpsym = S / N
		var q, sq = s.FSQ()
		var rtpfs = sq * rtpsym
		var rtp = rtpsym + q*rtpfs
		fmt.Fprintf(w, "symbols: rtp(sym) = %.6f%%\n", rtpsym*100)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FSC.Load(), q, sq)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", s.FGF())
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym*100, q, rtpfs*100, rtp*100)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
}
