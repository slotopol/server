package neonbananas

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

const Ebb = 1 // bananas bonus expectation

func CalcStat(ctx context.Context, sp *slot.ScanPar) float64 {
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame(sp.Sel)
	var s = slot.NewStatGeneric(sn, 5)

	var calc = func(w io.Writer) float64 {
		var N = s.Count()
		var lrtp, srtp = s.RTPsym(g.Cost(), scat1)
		var rtpsym = lrtp + srtp
		var qbb = s.BonusHitsF(bbid) / N
		var rtpbb = Ebb * qbb
		var qls1 = s.BonusHitsF(lsb1) / N / float64(g.Sel)
		var rtpls1 = Els * 1 * qls1
		var qls3 = s.BonusHitsF(lsb3) / N / float64(g.Sel)
		var rtpls3 = Els * 3 * qls3
		var qls6 = s.BonusHitsF(lsb6) / N / float64(g.Sel)
		var rtpls6 = Els * 6 * qls6
		var rtp = rtpsym + rtpbb + rtpls1 + rtpls3 + rtpls6
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp*100, srtp*100, rtpsym*100)
		fmt.Fprintf(w, "bananas bonuses: hit rate 1/%.5g, rtp = %.6f%%\n", N/s.BonusHitsF(bbid), rtpbb*100)
		fmt.Fprintf(w, "lucky spin1 bonuses: hit rate 1/%.5g, rtp = %.6f%%\n", N/s.BonusHitsF(lsb1), rtpls1*100)
		fmt.Fprintf(w, "lucky spin3 bonuses: hit rate 1/%.5g, rtp = %.6f%%\n", N/s.BonusHitsF(lsb3), rtpls3*100)
		fmt.Fprintf(w, "lucky spin6 bonuses: hit rate 1/%.5g, rtp = %.6f%%\n", N/s.BonusHitsF(lsb6), rtpls6*100)
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(bb) + %.5g(ls) = %.6f%%\n",
			rtpsym*100, rtpbb*100, (rtpls1+rtpls3+rtpls6)*100, rtp*100)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
}
