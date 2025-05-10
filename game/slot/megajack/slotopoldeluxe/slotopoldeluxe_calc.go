package slotopoldeluxe

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/game/slot/megajack/slotopol"
)

func CalcStat(ctx context.Context, mrtp float64) float64 {
	fmt.Printf("*bonus games calculations*\n")
	slotopol.ExpEldorado()
	slotopol.ExpMonopoly()
	fmt.Printf("*reels calculations*\n")
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	var g = NewGame()
	g.Sel = 1
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var reshuf = s.Count()
		var cost, _ = g.Cost()
		var lrtp, srtp = s.LineRTP(cost), s.ScatRTP(cost)
		var rtpsym = lrtp + srtp
		var qmje1 = s.BonusCount(mje1) / reshuf / float64(g.Sel)
		var rtpmje1 = slotopol.Emje * 1 * qmje1 * 100
		var qmje3 = s.BonusCount(mje3) / reshuf / float64(g.Sel)
		var rtpmje3 = slotopol.Emje * 3 * qmje3 * 100
		var qmje6 = s.BonusCount(mje6) / reshuf / float64(g.Sel)
		var rtpmje6 = slotopol.Emje * 6 * qmje6 * 100
		var qmjm = s.BonusCount(mjm) / reshuf / float64(g.Sel)
		var rtpmjm = slotopol.Emjm * qmjm * 100
		var rtp = rtpsym + rtpmje1 + rtpmje3 + rtpmje6 + rtpmjm
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "spin1 bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/s.BonusCount(mje1), rtpmje1)
		fmt.Fprintf(w, "spin3 bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/s.BonusCount(mje3), rtpmje3)
		fmt.Fprintf(w, "spin6 bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/s.BonusCount(mje6), rtpmje6)
		fmt.Fprintf(w, "monopoly bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/s.BonusCount(mjm), rtpmjm)
		if s.JackCount(mjj) > 0 {
			fmt.Fprintf(w, "jackpots: count %g, frequency 1/%.12g\n", s.JackCount(mjj), reshuf/s.JackCount(mjj))
		}
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(mje) + %.5g(mjm) = %.6f%%\n", rtpsym, rtpmje1+rtpmje3+rtpmje6, rtpmjm, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc)
}
