package jewels4all

import (
	"context"
	"fmt"
	"time"

	"github.com/slotopol/server/game"
)

func BruteForce5x(ctx context.Context, s game.Stater, g game.SlotGame, reels game.Reels, x, y int) {
	var screen = g.NewScreen()
	defer screen.Free()
	var wins game.Wins
	var r1 = reels.Reel(1)
	var r2 = reels.Reel(2)
	var r3 = reels.Reel(3)
	var r4 = reels.Reel(4)
	var r5 = reels.Reel(5)
	for i1 := range r1 {
		screen.SetCol(1, r1, i1)
		for i2 := range r2 {
			screen.SetCol(2, r2, i2)
			for i3 := range r3 {
				screen.SetCol(3, r3, i3)
				for i4 := range r4 {
					screen.SetCol(4, r4, i4)
					for i5 := range r5 {
						screen.SetCol(5, r5, i5)
						var sym game.Sym
						if x > 0 {
							sym = screen.At(x, y)
							screen.Set(x, y, wild)
						}
						g.Scanner(screen, &wins)
						if x > 0 {
							screen.Set(x, y, sym)
						}
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
		}
	}
}

func CalcStatEuro(ctx context.Context, x, y int) float64 {
	var reels = &Reels
	var g = NewGame("92")
	g.SBL = game.MakeBitNum(1)
	var sbl = float64(g.SBL.Num())
	var s game.Stat

	fmt.Printf("calculations of euro at [%d,%d]\n", x, y)

	var total = float64(reels.Reshuffles())
	var dur = func() time.Duration {
		var t0 = time.Now()
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		go s.Progress(ctx2, time.NewTicker(2*time.Second), sbl, total)
		BruteForce5x(ctx2, &s, g, reels, x, y)
		return time.Since(t0)
	}()

	var reshuf = float64(s.Reshuffles)
	var lrtp = s.LinePay / reshuf / sbl * 100
	_ = jid
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/total*100, g.SBL.Num(), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	fmt.Printf("RTP[%d,%d] = %.6f%%\n", x, y, lrtp)
	return lrtp
}

func CalcStat(ctx context.Context, rn string) (rtp float64) {
	var wc float64
	if rn != "" {
		var ok bool
		if wc, ok = ChanceMap[rn]; !ok {
			return 0
		}
	} else {
		wc, rn = ChanceMap["95"], "95"
	}

	var b = 1 / wc
	fmt.Printf("wild chance %.5g, b = %.5g\n", wc, b)
	var rtp00 = CalcStatEuro(ctx, 0, 0)
	var rtpeu float64
	for x := 1; x <= 5; x++ {
		for y := 1; y <= 3; y++ {
			rtpeu += CalcStatEuro(ctx, x, y)
		}
	}
	rtpeu /= 15
	rtp = rtp00 + wc*rtpeu
	fmt.Printf("euro avr: rtpeu = %.6f%%\n", rtpeu)
	fmt.Printf("RTP = %.5g(sym) + wc*%.5g(eu) = %.6f%%\n", rtp00, rtpeu, rtp)
	return
}
