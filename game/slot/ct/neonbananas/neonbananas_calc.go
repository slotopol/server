package neonbananas

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

const Ebb = 1 // bananas bonus expectation

func CalcStat(ctx context.Context, mrtp float64) float64 {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	var g = NewGame()
	g.Sel = 1
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var reshuf = s.Count()
		var lrtp, srtp = s.SymRTP(g.Cost())
		var rtpsym = lrtp + srtp
		var qbb = s.BonusCount(bbid) / reshuf
		var rtpbb = Ebb * qbb * 100
		var qls1 = s.BonusCount(lsb1) / reshuf / float64(g.Sel)
		var rtpls1 = Els * 1 * qls1 * 100
		var qls3 = s.BonusCount(lsb3) / reshuf / float64(g.Sel)
		var rtpls3 = Els * 3 * qls3 * 100
		var qls6 = s.BonusCount(lsb6) / reshuf / float64(g.Sel)
		var rtpls6 = Els * 6 * qls6 * 100
		var rtp = rtpsym + rtpbb + rtpls1 + rtpls3 + rtpls6
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "bananas bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/s.BonusCount(bbid), rtpbb)
		fmt.Fprintf(w, "lucky spin1 bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/s.BonusCount(lsb1), rtpls1)
		fmt.Fprintf(w, "lucky spin3 bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/s.BonusCount(lsb3), rtpls3)
		fmt.Fprintf(w, "lucky spin6 bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/s.BonusCount(lsb6), rtpls6)
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(bb) + %.5g(ls1) + %.5g(ls3) + %.5g(ls6) = %.6f%%\n", rtpsym, rtpbb, rtpls1, rtpls3, rtpls6, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc)
}
