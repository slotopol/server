package beetlemania

import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/slotopol/server/game/slot"
)

// Attention! On freespins can be calculated median only, not expectation.

func CalcStatBon(ctx context.Context) float64 {
	var reels = ReelsBon
	var g = NewGame(1)
	g.FSR = 10 // set free spins mode
	var s slot.StatGeneric

	var calc = func(w io.Writer) float64 {
		var N, S, _ = s.NSQ(g.Cost())
		var rtpsym = S / N
		var qjazz = s.BonusHitsF(jbonus) / N
		var jpow = math.Pow(2, 10*qjazz) // jazz power
		var rtpjazz = rtpsym*jpow - rtpsym
		var rtp = rtpsym * jpow
		fmt.Fprintf(w, "rtp(sym) = %.6f%%\n", rtpsym*100)
		fmt.Fprintf(w, "jazzbee bonuses: hit rate 1/%.5g, pow = %.5g, rtp = %.6f%%\n", N/s.BonusHitsF(jbonus), jpow, rtpjazz)
		fmt.Fprintf(w, "RTP = rtp(sym) + rtp(jazz) = %.5g + %.5g = %.6f%%\n", rtpsym*100, rtpjazz*100, rtp*100)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
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
	var s slot.StatGeneric

	var calc = func(w io.Writer) float64 {
		var lrtp, srtp = s.RTPsym(g.Cost(), scat)
		var rtpsym = lrtp + srtp
		var q, _ = s.FSQ()
		var rtp = rtpsym + q*rtpfs
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp*100, srtp*100, rtpsym*100)
		fmt.Fprintf(w, "free spins %d, q = %.5g\n", s.FSC.Load(), q)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", s.FGF())
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym*100, q, rtpfs*100, rtp*100)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}
