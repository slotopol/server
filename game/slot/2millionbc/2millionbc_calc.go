package twomillionbc

import (
	"context"
	"fmt"
	"strconv"
	"time"

	game "github.com/slotopol/server/game/slot"
)

var Eacbn float64

func ExpAcorn() float64 {
	var sum float64
	for _, v := range Acorn {
		sum += float64(v)
	}
	var E = sum / float64(len(Acorn))
	fmt.Printf("len = %d, E = %g\n", len(Acorn), E)
	return E
}

var Edlbn float64

func ExpDiamondLion() float64 {
	var sum float64
	for _, v := range DiamondLion {
		sum += float64(v)
	}
	var E = sum / float64(len(DiamondLion))
	fmt.Printf("len = %d, E = %g\n", len(DiamondLion), E)
	return E
}

func CalcStatBon(ctx context.Context) float64 {
	var reels = &ReelsBon
	var g = NewGame()
	g.SBL = game.MakeBitNum(1)
	g.FS = 4 // set free spins mode
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
	var sq = 1 / (1 - q)
	var rtp = sq * rtpsym
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/total*100, g.SBL.Num(), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCount, q, sq)
	fmt.Printf("free games frequency: 1/%.5g\n", reshuf/float64(s.FreeHits))
	fmt.Printf("RTP = sq*rtp(sym) = %.5g*%.5g = %.6f%%\n", sq, rtpsym, rtp)
	return rtp
}

func CalcStatReg(ctx context.Context, rn string) float64 {
	fmt.Printf("*bonus games calculations*\n")
	Eacbn = ExpAcorn()
	Edlbn = ExpDiamondLion()
	fmt.Printf("*bonus reels calculations*\n")
	var rtpfs = CalcStatBon(ctx)
	if ctx.Err() != nil {
		return 0
	}
	fmt.Printf("*regular reels calculations*\n")
	var reels *game.Reels5x
	if mrtp, _ := strconv.ParseFloat(rn, 64); mrtp != 0 {
		var _, r = FindReels(mrtp)
		reels = r.(*game.Reels5x)
	} else {
		reels = &ReelsReg96
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
	var sq = 1 / (1 - q)
	var qacbn = 1 / float64(len(reels.Reel(5)))
	var rtpacbn = Eacbn * qacbn * 100
	var qdlbn = float64(s.BonusCount[dlbn]) / reshuf / sbl
	var rtpdlbn = Edlbn * qdlbn * 100
	var rtp = rtpsym + rtpacbn + rtpdlbn + q*rtpfs
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/total*100, g.SBL.Num(), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCount, q, sq)
	fmt.Printf("free games frequency: 1/%.5g\n", reshuf/float64(s.FreeHits))
	fmt.Printf("acorn bonuses: frequency 1/%d, rtp = %.6f%%\n", len(reels.Reel(5)), rtpacbn)
	fmt.Printf("diamond lion bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/float64(s.BonusCount[dlbn]), rtpdlbn)
	fmt.Printf("RTP = %.5g(sym) + %.5g(acorn) + %.5g(dl) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym, rtpacbn, rtpdlbn, q, rtpfs, rtp)
	return rtp
}
