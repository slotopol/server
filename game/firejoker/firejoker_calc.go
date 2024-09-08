package firejoker

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/slotopol/server/game"
)

func BruteForce5x(ctx context.Context, s game.Stater, g game.SlotGame, reels game.Reels, gs game.Sym) {
	var screen = g.NewScreen()
	defer screen.Free()
	var wins game.Wins
	for x := 2; x <= 4; x++ {
		screen.Set(x, 1, gs)
		screen.Set(x, 2, gs)
		screen.Set(x, 3, gs)
	}
	var r1 = reels.Reel(1)
	var r5 = reels.Reel(5)
	for i1 := range r1 {
		screen.SetCol(1, r1, i1)
		for i5 := range r5 {
			screen.SetCol(5, r5, i5)
			g.Scanner(screen, &wins)
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

func CalcStatSym(ctx context.Context, g *Game, reels game.Reels, gs game.Sym) float64 {
	var sbl = float64(g.SBL.Num())
	var s game.Stat

	var ctx2, cancel2 = context.WithCancel(ctx)
	defer cancel2()
	BruteForce5x(ctx2, &s, g, reels, gs)

	var reshuf = float64(s.Reshuffles)
	var lrtp, srtp = s.LinePay / reshuf / sbl * 100, s.ScatPay / reshuf / sbl * 100
	var rtpsym = lrtp + srtp
	fmt.Printf("RTP[%d] = %.5g(lined) + %.5g(scatter) = %.6f%%\n", gs, lrtp, srtp, rtpsym)
	return rtpsym
}

func CalcStatBon(ctx context.Context, rn string) (rtp float64) {
	var reels *game.Reels5x
	if mrtp, _ := strconv.ParseFloat(rn, 64); mrtp != 0 {
		var _, r = FindReels(mrtp)
		reels = r.(*game.Reels5x)
	} else {
		reels = &Reels92
	}
	var g = NewGame()
	g.SBL = game.MakeBitNum(1)

	for gs := game.Sym(1); gs <= 7; gs++ {
		rtp += CalcStatSym(ctx, g, reels, gs)
	}
	rtp /= 7
	fmt.Printf("average freespins RTP = %.6f%%\n", rtp)
	return
}

func CalcStatReg(ctx context.Context, rn string) float64 {
	fmt.Printf("*bonus reels calculations*\n")
	var rtpfs = CalcStatBon(ctx, rn)
	if ctx.Err() != nil {
		return 0
	}
	fmt.Printf("*regular reels calculations*\n")
	var reels *game.Reels5x
	if mrtp, _ := strconv.ParseFloat(rn, 64); mrtp != 0 {
		var _, r = FindReels(mrtp)
		reels = r.(*game.Reels5x)
	} else {
		reels = &Reels92
	}
	var g = NewGame()
	g.SBL = game.MakeBitNum(1)
	var sbl = float64(g.SBL.Num())
	var s game.Stat

	var total = float64(reels.Reshuffles())
	var dur = func() time.Duration {
		var t0 = time.Now()
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		go s.Progress(ctx2, time.Tick(2*time.Second), sbl, total)
		game.BruteForce5x(ctx2, &s, g, reels)
		return time.Since(t0)
	}()

	var reshuf = float64(s.Reshuffles)
	var lrtp, srtp = s.LinePay / reshuf / sbl * 100, s.ScatPay / reshuf / sbl * 100
	var rtpsym = lrtp + srtp
	var q = float64(s.FreeCount) / reshuf
	var rtp = rtpsym + q*rtpfs
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/total*100, g.SBL.Num(), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("free spins %d, q = %.6f\n", s.FreeCount, q)
	fmt.Printf("free games frequency: 1/%.5g\n", reshuf/float64(s.FreeHits))
	fmt.Printf("RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym, q, rtpfs, rtp)
	return rtp
}
