package powerstars

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/slotopol/server/game/slot"
)

// Returns the probability of getting at least one star on the 3 reels,
// including several stars at once.
func AnyStarProb(b float64) float64 {
	return (b*b + (b-1)*b + (b-1)*(b-1)) / b / b / b
}

func BruteForceStars(ctx context.Context, s slot.Stater, g *Game, reels slot.Reelx, wc2, wc3, wc4 bool) {
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
						var sym2, sym3, sym4 = g.At(2, 1), g.At(3, 1), g.At(4, 1)
						if wc2 {
							g.SetSym(2, 1, wild)
						}
						if wc3 {
							g.SetSym(3, 1, wild)
						}
						if wc4 {
							g.SetSym(4, 1, wild)
						}
						s.Simulate(g, reels, &wins)
						g.SetSym(2, 1, sym2)
						g.SetSym(3, 1, sym3)
						g.SetSym(4, 1, sym4)
						wins.Reset()
					}
				}
			}
		}
	}
}

func CalcStatStars(ctx context.Context, wc2, wc3, wc4 bool) float64 {
	var reels = Reels
	var g = NewGame(1)
	var s slot.StatGeneric

	var wcsym = func(wc bool) byte {
		if wc {
			return '*'
		}
		return '-'
	}
	fmt.Printf("calculations of star combinations [%c%c%c]\n", wcsym(wc2), wcsym(wc3), wcsym(wc4))

	var calc = func(w io.Writer) float64 {
		var lrtp = s.LineRTP(g.Cost())
		fmt.Fprintf(w, "RTP[%c%c%c] = %.6f%%\n", wcsym(wc2), wcsym(wc3), wcsym(wc4), lrtp*100)
		return lrtp
	}

	func() {
		var t0 = time.Now()
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		s.SetPlan(reels.Reshuffles())
		go slot.Progress(ctx2, &s, calc)
		BruteForceStars(ctx2, &s, g, reels, wc2, wc3, wc4)
		var dur = time.Since(t0)
		var comp = s.Count() / float64(s.GetPlan()) * 100
		fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", comp, g.GetSel(), dur)
	}()
	return calc(os.Stdout)
}

func CalcStat(ctx context.Context, mrtp float64) (rtp float64) {
	var wc, _ = ChanceMap.FindClosest(mrtp) // wild chance

	var b = 1 / wc
	var rtp000 = CalcStatStars(ctx, false, false, false)
	var rtp100 = CalcStatStars(ctx, true, false, false)
	var rtp010 = CalcStatStars(ctx, false, true, false)
	var rtp001 = CalcStatStars(ctx, false, false, true)
	var rtp110 = CalcStatStars(ctx, true, true, false)
	var rtp011 = CalcStatStars(ctx, false, true, true)
	var rtp101 = CalcStatStars(ctx, true, false, true)
	var rtp111 = CalcStatStars(ctx, true, true, true)
	var q = AnyStarProb(b)
	var rtpfs = ((rtp100+rtp010+rtp001)*(b-1)*(b-1) + (rtp110+rtp011+rtp101)*(b-1) + rtp111) / (b*b + (b-1)*b + (b-1)*(b-1))
	rtp = (1-q)*rtp000 + q*rtpfs
	fmt.Printf("wild chance: 1/%.5g\n", 1/wc)
	fmt.Printf("free spins: q = %.5g, 1/q = %.5g, rtpfs = %.6f%%\n", q, 1/q, rtpfs*100)
	fmt.Printf("RTP = (1-q)*%.5g(sym) + q*%.5g(fg) = %.6f%%\n", rtp000*100, rtpfs*100, rtp*100)
	return
}
