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
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	var g = NewGame()
	g.Sel = 1
	g.FSR = 15 // set free spins mode
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var reshuf = s.Count()
		var cost, _ = g.Cost()
		var lrtp, srtp = s.LineRTP(cost), s.ScatRTP(cost)
		var rtpsym = lrtp + srtp
		var qcbn = s.BonusCount(cbn) / reshuf / float64(g.Sel)
		var rtpcbn = Ecards * qcbn * 100
		var rtp = rtpsym + rtpcbn
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "cards bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/s.BonusCount(cbn), rtpcbn)
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(cards) = %.6f%%\n", rtpsym, rtpcbn, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc)
}

func CalcStatReg(ctx context.Context, mrtp float64) float64 {
	fmt.Printf("*bonus games calculations*\n")
	ExpCards()
	fmt.Printf("total = %d, E = %g\n", 4*4*4, Ecards)
	fmt.Printf("*bonus reels calculations*\n")
	var rtpfs = CalcStatBon(ctx, mrtp)
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
		var qcbn = s.BonusCount(cbn) / reshuf / float64(g.Sel)
		var rtpcbn = Ecards * qcbn * 100
		var rtp = rtpsym + rtpcbn + q*rtpfs
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.6f\n", s.FreeCountU(), q)
		fmt.Fprintf(w, "free games frequency: 1/%.5g\n", reshuf/s.FreeHits())
		fmt.Fprintf(w, "cards bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/s.BonusCount(cbn), rtpcbn)
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(cards) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym, rtpcbn, q, rtpfs, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc)
}
