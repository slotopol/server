package slotopoldeluxe

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/game/slot/megajack/slotopol"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) float64 {
	fmt.Printf("*bonus games calculations*\n")
	slotopol.ExpEldorado()
	slotopol.ExpMonopoly()
	fmt.Printf("*reels calculations*\n")
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame(sp.Sel)
	var s slot.StatGeneric

	var calc = func(w io.Writer) float64 {
		var N = s.Count()
		var lrtp, srtp = s.RTPsym(g.Cost(), scat)
		var rtpsym = lrtp + srtp
		var qmje1 = s.BonusHitsF(mje1) / N / float64(g.Sel)
		var rtpmje1 = slotopol.Emje * 1 * qmje1
		var qmje3 = s.BonusHitsF(mje3) / N / float64(g.Sel)
		var rtpmje3 = slotopol.Emje * 3 * qmje3
		var qmje6 = s.BonusHitsF(mje6) / N / float64(g.Sel)
		var rtpmje6 = slotopol.Emje * 6 * qmje6
		var qmjm = s.BonusHitsF(mjm) / N / float64(g.Sel)
		var rtpmjm = slotopol.Emjm * qmjm
		var rtp = rtpsym + rtpmje1 + rtpmje3 + rtpmje6 + rtpmjm
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp*100, srtp*100, rtpsym*100)
		fmt.Fprintf(w, "spin1 bonuses: hit rate 1/%.5g, rtp = %.6f%%\n", N/s.BonusHitsF(mje1), rtpmje1*100)
		fmt.Fprintf(w, "spin3 bonuses: hit rate 1/%.5g, rtp = %.6f%%\n", N/s.BonusHitsF(mje3), rtpmje3*100)
		fmt.Fprintf(w, "spin6 bonuses: hit rate 1/%.5g, rtp = %.6f%%\n", N/s.BonusHitsF(mje6), rtpmje6*100)
		fmt.Fprintf(w, "monopoly bonuses: hit rate 1/%.5g, rtp = %.6f%%\n", N/s.BonusHitsF(mjm), rtpmjm*100)
		if s.JackHitsF(mjj) > 0 {
			fmt.Fprintf(w, "jackpots: count %g, frequency 1/%.12g\n", s.JackHitsF(mjj), N/s.JackHitsF(mjj))
		}
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(mje) + %.5g(mjm) = %.6f%%\n", rtpsym*100, (rtpmje1+rtpmje3+rtpmje6)*100, rtpmjm*100, rtp*100)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, sp, &s, g, reels, calc)
}
