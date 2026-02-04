package suncity

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStatBon(ctx context.Context) (rtp, num float64) {
	var reels = ReelsBon
	var g = NewGame(1)
	g.FSR = -1 // set free spins mode
	var s slot.StatGeneric

	var fgf float64
	var calc = func(w io.Writer) float64 {
		var N = s.Count()
		var lrtp, srtp = s.SymRTP(g.Cost())
		if srtp > 0 {
			panic("scatters have no pays")
		}
		var rtpsym = lrtp + srtp
		fgf = N / float64(s.FreeHits.Load())
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp*100, srtp*100, rtpsym*100)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", fgf)
		fmt.Fprintf(w, "RTP = rtp(sym) = %.6f%%\n", rtpsym*100)
		return rtpsym
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc), fgf
}

func CalcStatReg(ctx context.Context, mrtp float64) float64 {
	fmt.Printf("*bonus reels calculations*\n")
	var rtpfs, numfs = CalcStatBon(ctx)
	if ctx.Err() != nil {
		return 0
	}
	fmt.Printf("*regular reels calculations*\n")
	var reels, _ = ReelsMap.FindClosest(mrtp)
	var g = NewGame(1)
	var s slot.StatGeneric

	var calc = func(w io.Writer) float64 {
		var N = s.Count()
		var lrtp, srtp = s.SymRTP(g.Cost())
		if srtp > 0 {
			panic("scatters have no pays")
		}
		var rtpsym = lrtp + srtp
		var fgf = N / float64(s.FreeHits.Load())
		var rtp = rtpsym + rtpfs*numfs/fgf
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp*100, srtp*100, rtpsym*100)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", fgf)
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(fg)*%.5g/%.5g = %.6f%%\n", rtpsym*100, rtpfs*100, numfs, fgf, rtp*100)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}
