package diamonddogs

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

const Ene12 = 3 * 100

func CalcStatBon(ctx context.Context) float64 {
	var reels = ReelsBon
	var g = NewGame()
	g.Sel = 1
	g.FSR = 10 // set free spins mode
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var reshuf = s.Count()
		var cost, _ = g.Cost()
		var lrtp, srtp = s.LineRTP(cost), s.ScatRTP(cost)
		var rtpsym = lrtp + srtp
		var q = s.FreeCount() / reshuf
		var sq = 1 / (1 - q)
		var rtp = sq * rtpsym
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCountU(), q, sq)
		fmt.Fprintf(w, "free games frequency: 1/%.5g\n", reshuf/s.FreeHits())
		fmt.Fprintf(w, "RTP = sq*rtp(sym) = %.5g*%.5g = %.6f%%\n", sq, rtpsym, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc)
}

func CalcStatReg(ctx context.Context, mrtp float64) float64 {
	fmt.Printf("*bonus reels calculations*\n")
	var rtpfs = CalcStatBon(ctx)
	if ctx.Err() != nil {
		return 0
	}
	fmt.Printf("*regular reels calculations*\n")
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
		var qne12 = s.BonusCount(ne12) / reshuf / float64(g.Sel)
		var rtpne12 = Ene12 * qne12 * 100
		var rtp = rtpsym + rtpne12 + q*rtpfs
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCountU(), q, sq)
		fmt.Fprintf(w, "free games frequency: 1/%.5g\n", reshuf/s.FreeHits())
		fmt.Fprintf(w, "ne12 bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/s.BonusCount(ne12), rtpne12)
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(ne12) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym, rtpne12, q, rtpfs, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc)
}
