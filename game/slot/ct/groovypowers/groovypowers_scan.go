package groovypowers

import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/slotopol/server/game/slot"
)

func CalcStatBon(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame(sp.Sel)
	g.BM = true // set bonus mode
	var s = slot.NewStatGeneric(sn, 5)

	var calc = func(w io.Writer) (float64, float64) {
		return slot.Parsheet_generic_simple(w, sp, s, g.Cost())
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
}

func CalcStatReg(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("*free games calculations*\n")
	var rtpbm, _ = CalcStatBon(ctx, sp)
	if ctx.Err() != nil {
		return 0, 0
	}
	fmt.Printf("*regular games calculations*\n")
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame(sp.Sel)
	var s = slot.NewStatGeneric(sn, 5)

	var calc = func(w io.Writer) (float64, float64) {
		var lrtp, srtp = s.RTPsym(g.Cost(), scat)
		var rtpsym = lrtp + srtp
		var rtp = rtpsym*(1-Pbm) + rtpbm*Pbm
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp*100, srtp*100, rtpsym*100)
		fmt.Fprintf(w, "RTP = %.5g(reg)*%.5g + %.5g(bm)*%.5g = %.6f%%\n", rtpsym*100, 1-Pbm, rtpbm*100, Pbm, rtp*100)
		return rtp, math.NaN()
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
}
