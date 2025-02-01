package firejoker

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/slotopol/server/game/slot"
)

func BruteForceFire(ctx context.Context, s slot.Stater, g slot.SlotGame, reels slot.Reels, big slot.Sym) {
	var screen = g.Screen()
	var wins slot.Wins
	var x slot.Pos
	for x = 2; x <= 4; x++ {
		screen.Set(x, 1, big)
		screen.Set(x, 2, big)
		screen.Set(x, 3, big)
	}
	var r1 = reels.Reel(1)
	var r5 = reels.Reel(5)
	for i1 := range r1 {
		screen.SetCol(1, r1, i1)
		for i5 := range r5 {
			screen.SetCol(5, r5, i5)
			g.Scanner(&wins)
			s.Update(wins)
			wins.Reset()
			if s.Count()&100 == 0 {
				select {
				case <-ctx.Done():
					return
				default:
				}
			}
		}
	}
}

func CalcStatSym(ctx context.Context, g *Game, reels slot.Reels, big slot.Sym) float64 {
	var s slot.Stat

	var ctx2, cancel2 = context.WithCancel(ctx)
	defer cancel2()
	BruteForceFire(ctx2, &s, g, reels, big)

	var lrtp, srtp = s.LineRTP(g.Sel), s.ScatRTP(g.Sel)
	var rtpsym = lrtp + srtp
	fmt.Printf("RTP[%d] = %.5g(lined) + %.5g(scatter) = %.6f%%\n", big, lrtp, srtp, rtpsym)
	return rtpsym
}

func CalcStatBon(ctx context.Context, mrtp float64) (rtp float64) {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	var g = NewGame()
	g.Sel = 1

	for big := slot.Sym(1); big <= 7; big++ {
		rtp += CalcStatSym(ctx, g, reels, big)
	}
	rtp /= 7
	fmt.Printf("average freespins RTP = %.6f%%\n", rtp)
	return
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
		fmt.Fprintf(w, "reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
			len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.6f\n", s.FreeCount(), q)
		fmt.Fprintf(w, "free games frequency: 1/%.5g\n", reshuf/float64(s.FreeHits()))
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym, q, rtpfs, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc,
		time.Tick(2*time.Second), time.Tick(2*time.Second))
}
