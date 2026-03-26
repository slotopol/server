package jewels4all

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/slotopol/server/game/slot"
)

func BruteForceEuro(ctx context.Context, s slot.Simulator, g *Game, reels slot.Reelx, x, y slot.Pos) {
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

func CalcStatEuro(ctx context.Context, sp *slot.ScanPar, s *slot.StatGeneric, x, y slot.Pos) (float64, float64) {
	var reels = Reels
	var g = NewGame(sp.Sel)

	var calc = func(w io.Writer) (float64, float64) {
		return slot.Parsheet_simple(w, sp, s, g.Cost())
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		go func() {
			defer wg.Done()
			slot.ProgressBF(ctx2, sp, s, calc, float64(reels.Reshuffles()))
		}()
		BruteForceEuro(ctx2, s, g, reels, x, y)
	}()
	wg.Wait()
	return calc(os.Stdout)
}

// custom parsheet
func CalcStat(ctx context.Context, sp *slot.ScanPar) (rtp, D float64) {
	var Pw, _ = ChanceMap.FindClosest(sp.MRTP) // wild chance

	var µw, µsum2, Dsum float64
	for i := range 15 {
		var x = slot.Pos(i/3 + 1)
		var y = slot.Pos(i%3 + 1)
		fmt.Printf("\n(%d/16) calculations of euro at [%d,%d]\n", i+1, x, y)
		var s = slot.NewStatGeneric(sn, 5)
		var µ, D = CalcStatEuro(ctx, sp, s, x, y)
		µw += µ
		µsum2 += µ * µ
		Dsum += D
	}
	µw /= 15
	var Dw = Dsum/15 + µsum2/15 - µw*µw
	fmt.Printf("\n(16/16) regular games calculations\n")
	var sr = slot.NewStatGeneric(sn, 5)
	var µr, Dr = CalcStatEuro(ctx, sp, sr, 0, 0)
	rtp = (1-Pw)*µr + Pw*µw
	D = (1-Pw)*Dr + Pw*Dw + Pw*(1-Pw)*(µw-µr)*(µw-µr)
	if sp.IsMain() {
		fmt.Printf("*final calculations*\n")
		fmt.Printf("euro avr: rtp(eu) = %.6f%%\n", µw*100)
		fmt.Printf("wild chance: 1/%.5g\n", 1/Pw)
		fmt.Printf("RTP = (1-Pw)*%.5g(sym) + Pw*%.5g(eu) = %.6f%%\n", µr*100, µw*100, rtp*100)
	}
	var w = os.Stdout
	slot.Print_all(w, sp, sr, rtp, D)
	return
}
