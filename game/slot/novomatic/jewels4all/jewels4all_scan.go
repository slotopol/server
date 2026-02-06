package jewels4all

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/slotopol/server/game/slot"
)

func BruteForceEuro(ctx context.Context, s slot.Stater, g *Game, reels slot.Reelx, x, y slot.Pos) {
	var wins slot.Wins
	var r1 = reels.Reel(1)
	var r2 = reels.Reel(2)
	var r3 = reels.Reel(3)
	var r4 = reels.Reel(4)
	var r5 = reels.Reel(5)
	var N uint64
	for i1 := range r1 {
		g.SetCol(1, r1, i1)
		for i2 := range r2 {
			g.SetCol(2, r2, i2)
			for i3 := range r3 {
				g.SetCol(3, r3, i3)
				for i4 := range r4 {
					g.SetCol(4, r4, i4)
					for i5 := range r5 {
						N++
						if N%slot.CtxGranulation == 0 {
							select {
							case <-ctx.Done():
								return
							default:
							}
						}
						g.SetCol(5, r5, i5)
						var sym slot.Sym
						if x > 0 {
							sym = g.At(x, y)
							g.SetSym(x, y, wild)
						}
						s.Simulate(g, reels, &wins)
						if x > 0 {
							g.SetSym(x, y, sym)
						}
						wins.Reset()
					}
				}
			}
		}
	}
}

func CalcStatEuro(ctx context.Context, sp *slot.ScanPar, x, y slot.Pos) float64 {
	var reels = Reels
	var g = NewGame(sp.Sel)
	var s slot.StatGeneric

	fmt.Printf("calculations of euro at [%d,%d]\n", x, y)

	var calc = func(w io.Writer) float64 {
		var N, S, _ = s.NSQ(g.Cost())
		var rtp = S / N
		fmt.Fprintf(w, "RTP[%d,%d] = %.6f%%\n", x, y, rtp*100)
		return rtp
	}

	func() time.Duration {
		var t0 = time.Now()
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		s.SetPlan(reels.Reshuffles())
		go slot.Progress(ctx2, &s, calc)
		BruteForceEuro(ctx2, &s, g, reels, x, y)
		return time.Since(t0)
	}()
	return calc(os.Stdout)
}

func CalcStat(ctx context.Context, sp *slot.ScanPar) (rtp float64) {
	var wc, _ = ChanceMap.FindClosest(sp.MRTP) // wild chance

	var rtp00 = CalcStatEuro(ctx, sp, 0, 0)
	var rtpeu float64
	var x, y slot.Pos
	for x = 1; x <= 5; x++ {
		for y = 1; y <= 3; y++ {
			rtpeu += CalcStatEuro(ctx, sp, x, y)
		}
	}
	rtpeu /= 15
	rtp = (1-wc)*rtp00 + wc*rtpeu
	fmt.Printf("euro avr: rtpeu = %.6f%%\n", rtpeu*100)
	fmt.Printf("wild chance: 1/%.5g\n", 1/wc)
	fmt.Printf("RTP = (1-wc)*%.5g(sym) + wc*%.5g(eu) = %.6f%%\n", rtp00*100, rtpeu*100, rtp*100)
	return
}
