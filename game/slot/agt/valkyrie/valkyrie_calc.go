package valkyrie

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/slotopol/server/game/slot"
)

func CalcStatBon(ctx context.Context, mrtp float64) float64 {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	var g = NewGame()
	g.Sel = 1
	g.FSR = 15 // set free spins mode
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var lrtp, srtp = s.LineRTP(g.Sel), s.ScatRTP(g.Sel)
		if srtp > 0 {
			panic("scatters are presents on bonus games")
		}
		var rtp = lrtp + srtp
		fmt.Fprintf(w, "RTP = %.6f%%\n", rtp)
		return rtp
	}

	func() {
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		s.SetPlan(uint64(len(reels.Reel(1))) * uint64(len(BonusReel)) * uint64(len(reels.Reel(5))))
		slot.BruteForce5x3Big(ctx2, &s, g, reels.Reel(1), BonusReel, reels.Reel(5))
	}()
	return calc(os.Stdout)
}

func CalcStatReg(ctx context.Context, mrtp float64) float64 {
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
		var rtp = rtpsym + q*rtpfs
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.6f\n", s.FreeCount(), q)
		fmt.Fprintf(w, "free games frequency: 1/%.5g\n", reshuf/float64(s.FreeHits()))
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym, q, rtpfs, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc,
		time.Tick(2*time.Second), time.Tick(2*time.Second))
}
