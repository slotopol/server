package firejoker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/slotopol/server/game/slot"
)

func CalcStatBon(ctx context.Context, sp *slot.ScanPar) float64 {
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame(sp.Sel)
	g.FSR = 10 // set free spins mode
	var s slot.StatGeneric

	var calc = func(w io.Writer) float64 {
		var lrtp, srtp = s.RTPsym(g.Cost(), scat)
		var rtpsym = lrtp + srtp
		fmt.Fprintf(w, "RTP = %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp*100, srtp*100, rtpsym*100)
		return rtpsym
	}

	func() {
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		slot.BruteForce5x3Big(ctx2, &s, g, reels.Reel(1), BonusReel, reels.Reel(5))
	}()
	return calc(os.Stdout)
}

func CalcStatReg(ctx context.Context, sp *slot.ScanPar) float64 {
	fmt.Printf("*free games calculations*\n")
	var rtpfs = CalcStatBon(ctx, sp)
	if ctx.Err() != nil {
		return 0
	}
	fmt.Printf("*regular games calculations*\n")
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame(sp.Sel)
	var s slot.StatGeneric

	var calc = func(w io.Writer) float64 {
		var lrtp, srtp = s.RTPsym(g.Cost(), scat)
		var rtpsym = lrtp + srtp
		var q, _ = s.FSQ()
		var rtp = rtpsym + q*rtpfs
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp*100, srtp*100, rtpsym*100)
		fmt.Fprintf(w, "free spins %d, q = %.6f\n", s.FSC.Load(), q)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", s.FGF())
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym*100, q, rtpfs*100, rtp*100)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}
