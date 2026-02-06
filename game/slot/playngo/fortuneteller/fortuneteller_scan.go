package fortuneteller

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

var Ecards float64

func ExpCards() {
	var sum float64
	for c1 := 1; c1 <= 4; c1++ {
		for c2 := 1; c2 <= 4; c2++ {
			for c3 := 1; c3 <= 4; c3++ {
				sum += CardsWin(c1, c2, c3)
			}
		}
	}
	Ecards = sum / 4 / 4 / 4
}

func CalcStatBon(ctx context.Context, mrtp float64) float64 {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	var g = NewGame(1)
	g.FSR = 15 // set free spins mode
	var s slot.StatGeneric

	var calc = func(w io.Writer) float64 {
		var N = s.Count()
		var lrtp, srtp = s.RTPsym(g.Cost(), scat)
		var rtpsym = lrtp + srtp
		var qcbn = s.BonusHitsF(cbn) / N / float64(g.Sel)
		var rtpcbn = Ecards * qcbn
		var rtp = rtpsym + rtpcbn
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp*100, srtp*100, rtpsym*100)
		fmt.Fprintf(w, "cards bonuses: hit rate 1/%.5g, rtp = %.6f%%\n", N/s.BonusHitsF(cbn), rtpcbn*100)
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(cards) = %.6f%%\n", rtpsym*100, rtpcbn*100, rtp*100)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}

func CalcStatReg(ctx context.Context, mrtp float64) float64 {
	fmt.Printf("*bonus games calculations*\n")
	ExpCards()
	fmt.Printf("total = %d, E = %g\n", 4*4*4, Ecards)
	fmt.Printf("*free games calculations*\n")
	var rtpfs = CalcStatBon(ctx, mrtp)
	if ctx.Err() != nil {
		return 0
	}
	fmt.Printf("*regular games calculations*\n")
	var reels, _ = ReelsMap.FindClosest(mrtp)
	var g = NewGame(1)
	var s slot.StatGeneric

	var calc = func(w io.Writer) float64 {
		var N = s.Count()
		var lrtp, srtp = s.RTPsym(g.Cost(), scat)
		var rtpsym = lrtp + srtp
		var q, _ = s.FSQ()
		var qcbn = s.BonusHitsF(cbn) / N / float64(g.Sel)
		var rtpcbn = Ecards * qcbn
		var rtp = rtpsym + rtpcbn + q*rtpfs
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp*100, srtp*100, rtpsym*100)
		fmt.Fprintf(w, "free spins %d, q = %.6f\n", s.FSC.Load(), q)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", s.FGF())
		fmt.Fprintf(w, "cards bonuses: hit rate 1/%.5g, rtp = %.6f%%\n", N/s.BonusHitsF(cbn), rtpcbn*100)
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(cards) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym*100, rtpcbn*100, q, rtpfs*100, rtp*100)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}
