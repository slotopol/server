package helena

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStatBon(ctx context.Context) float64 {
	var reels = ReelsBon
	var g = NewGame(1)
	g.FSR = 15 // set free spins mode
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var lrtp, srtp = s.SymRTP(g.Cost())
		var rtpsym = lrtp + srtp
		var q, sq = s.FSQ()
		var rtp = sq * rtpsym
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCountU(), q, sq)
		fmt.Fprintf(w, "free games frequency: 1/%.5g\n", s.FGF())
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
	var reels, _ = ReelsMap.FindClosest(mrtp)
	var g = NewGame(1)
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var lrtp, srtp = s.SymRTP(g.Cost())
		var rtpsym = lrtp + srtp
		var q, sq = s.FSQ()
		var rtp = rtpsym + q*rtpfs
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCountU(), q, sq)
		fmt.Fprintf(w, "free games frequency: 1/%.5g\n", s.FGF())
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym, q, rtpfs, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc)
}
