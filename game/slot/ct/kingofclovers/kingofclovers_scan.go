package kingofclovers

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStatBon(ctx context.Context) float64 {
	var reels = ReelsBon
	var g = NewGame()
	g.FSR = 14 // set free spins mode
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var reshuf1 = s.Reshuf(1)
		var reshuf2 = s.Reshuf(2)
		var reshuf3 = s.Reshuf(3)
		var reshuf4 = s.Reshuf(4)
		var reshuf5 = s.Reshuf(5)
		var lrtp, srtp = s.SymRTP(g.Cost())
		var rtpsym = lrtp + srtp
		var q, sq = s.FSQ()
		var rtp = sq * rtpsym
		fmt.Fprintf(w, "fall[2] = %.10g, freq = 1/%.5g\n", reshuf2, reshuf1/reshuf2)
		fmt.Fprintf(w, "fall[3] = %.10g, freq = 1/%.5g, freq2 = 1/%.5g\n", reshuf3, reshuf1/reshuf3, reshuf2/reshuf3)
		fmt.Fprintf(w, "fall[4] = %.10g, freq = 1/%.5g, freq2 = 1/%.5g\n", reshuf4, reshuf1/reshuf4, reshuf3/reshuf4)
		fmt.Fprintf(w, "fall[5] = %.10g, freq = 1/%.5g, freq2 = 1/%.5g\n", reshuf5, reshuf1/reshuf5, reshuf4/reshuf5)
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCountU(), q, sq)
		fmt.Fprintf(w, "free games frequency: 1/%.5g\n", s.FGF())
		fmt.Fprintf(w, "RTP = sq*rtp(sym) = %.5g*%.5g = %.20f%%\n", sq, rtpsym, rtp)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}

func CalcStatReg(ctx context.Context, mrtp float64) float64 {
	fmt.Printf("*bonus reels calculations*\n")
	var rtpfs = CalcStatBon(ctx) // 119.47281932232095869040
	if ctx.Err() != nil {
		return 0
	}
	fmt.Printf("*regular reels calculations*\n")
	var reels, _ = ReelsMap.FindClosest(mrtp)
	var g = NewGame()
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var reshuf1 = s.Reshuf(1)
		var reshuf2 = s.Reshuf(2)
		var reshuf3 = s.Reshuf(3)
		var reshuf4 = s.Reshuf(4)
		var reshuf5 = s.Reshuf(5)
		var lrtp, srtp = s.SymRTP(g.Cost())
		var rtpsym = lrtp + srtp
		var q, sq = s.FSQ()
		var rtp = rtpsym + q*rtpfs
		fmt.Fprintf(w, "fall[2] = %.10g, freq = 1/%.5g\n", reshuf2, reshuf1/reshuf2)
		fmt.Fprintf(w, "fall[3] = %.10g, freq = 1/%.5g, freq2 = 1/%.5g\n", reshuf3, reshuf1/reshuf3, reshuf2/reshuf3)
		fmt.Fprintf(w, "fall[4] = %.10g, freq = 1/%.5g, freq2 = 1/%.5g\n", reshuf4, reshuf1/reshuf4, reshuf3/reshuf4)
		fmt.Fprintf(w, "fall[5] = %.10g, freq = 1/%.5g, freq2 = 1/%.5g\n", reshuf5, reshuf1/reshuf5, reshuf4/reshuf5)
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCountU(), q, sq)
		fmt.Fprintf(w, "free games frequency: 1/%.5g\n", s.FGF())
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym, q, rtpfs, rtp)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}
