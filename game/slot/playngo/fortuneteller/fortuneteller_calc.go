package fortuneteller

import (
	"context"
	"fmt"
	"io"
	"time"

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
		var reshuf = float64(s.Count())
		var lrtp, srtp = s.LineRTP(g.Sel), s.ScatRTP(g.Sel)
		var rtpsym = lrtp + srtp
		var qcbn = float64(s.BonusCount(cbn)) / reshuf / float64(g.Sel)
		var rtpcbn = Ecards * qcbn * 100
		var rtp = rtpsym + rtpcbn
		fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
			len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
		fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Printf("cards bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/float64(s.BonusCount(cbn)), rtpcbn)
		fmt.Printf("RTP = %.5g(sym) + %.5g(cards) = %.6f%%\n", rtpsym, rtpcbn, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc,
		time.Tick(2*time.Second), time.Tick(2*time.Second))
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
		var reshuf = float64(s.Count())
		var lrtp, srtp = s.LineRTP(g.Sel), s.ScatRTP(g.Sel)
		var rtpsym = lrtp + srtp
		var q = float64(s.FreeCount()) / reshuf
		var qcbn = float64(s.BonusCount(cbn)) / reshuf / float64(g.Sel)
		var rtpcbn = Ecards * qcbn * 100
		var rtp = rtpsym + rtpcbn + q*rtpfs
		fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
			len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
		fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Printf("free spins %d, q = %.6f\n", s.FreeCount(), q)
		fmt.Printf("free games frequency: 1/%.5g\n", reshuf/float64(s.FreeHits()))
		fmt.Printf("cards bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/float64(s.BonusCount(cbn)), rtpcbn)
		fmt.Printf("RTP = %.5g(sym) + %.5g(cards) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym, rtpcbn, q, rtpfs, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc,
		time.Tick(2*time.Second), time.Tick(2*time.Second))
}
