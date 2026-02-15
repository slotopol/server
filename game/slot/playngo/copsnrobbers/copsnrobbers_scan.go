package copsnrobbers

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStatBon(ctx context.Context, sp *slot.ScanPar) float64 {
	var reels = ReelsBon
	var g = NewGame(sp.Sel)
	g.FSR = Efs // set free spins mode
	g.M = 1
	var s = slot.NewStatGeneric(sn, 5)

	var calc = func(w io.Writer) float64 {
		var lrtp, srtp = s.RTPsym(g.Cost(), scat)
		var rtpsym = lrtp + srtp
		var rtp = rtpsym
		fmt.Fprintf(w, "RTP = %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp*100, srtp*100, rtpsym*100)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
}

func CalcStatReg(ctx context.Context, sp *slot.ScanPar) float64 {
	fmt.Printf("*bonus reels calculations*\n")
	var rtpfs = CalcStatBon(ctx, sp)
	if ctx.Err() != nil {
		return 0
	}
	fmt.Printf("*regular reels calculations*\n")
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame(sp.Sel)
	var s = slot.NewStatGeneric(sn, 5)

	var calc = func(w io.Writer) float64 {
		// Correct free spins count with math expectation value
		var fc = float64(s.FGH.Load()) * Efs

		var N = s.Count()
		var lrtp, srtp = s.RTPsym(g.Cost(), scat)
		var rtpsym = lrtp + srtp
		var q = fc / N
		var rtp = rtpsym + (1+Pfs)*q*rtpfs
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp*100, srtp*100, rtpsym*100)
		fmt.Fprintf(w, "free spins %g, q = %.6f\n", fc, q)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", s.FGF())
		fmt.Fprintf(w, "RTP = %.5g(sym) + %g*%.5g*%.5g(fg) = %.6f%%\n", rtpsym*100, 1+Pfs, q, rtpfs*100, rtp*100)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
}
