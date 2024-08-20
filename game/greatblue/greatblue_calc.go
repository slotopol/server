package greatblue

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/slotopol/server/game"
)

func FirstSreespins() (fsavr1 float64, multavr float64) {
	// combinations of multiplier & freespins number
	// of two shells from set [x5, x8, 7, 10, 15]
	var combs = []struct {
		mult, fsnum float64
	}{
		{2 + 5 + 8, 8},
		{2 + 5, 8 + 7},
		{2 + 5, 8 + 10},
		{2 + 5, 8 + 15},
		{2 + 8, 8 + 7},
		{2 + 8, 8 + 10},
		{2 + 8, 8 + 15},
		{2, 8 + 7 + 10},
		{2, 8 + 7 + 15},
		{2, 8 + 10 + 15},
	}
	for _, c := range combs {
		fsavr1 += c.mult * c.fsnum
		multavr += c.mult
	}
	fsavr1 /= float64(len(combs))
	multavr /= float64(len(combs))
	return
}

func CalcStat(ctx context.Context, rn string) float64 {
	var reels *game.Reels5x
	var mrtp float64
	if mrtp, _ = strconv.ParseFloat(rn, 64); mrtp != 0 {
		var _, r = FindReels(mrtp)
		reels = r.(*game.Reels5x)
	} else {
		mrtp, reels = 92, &Reels92
	}
	var g = NewGame(mrtp)
	g.SBL = game.MakeBitNum(5)
	var sbl = float64(g.SBL.Num())
	var s game.Stat

	var total = float64(reels.Reshuffles())
	var dur = func() time.Duration {
		var t0 = time.Now()
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		go s.Progress(ctx2, time.NewTicker(2*time.Second), sbl, total)
		game.BruteForce5x(ctx2, &s, g, reels)
		return time.Since(t0)
	}()

	var reshuf = float64(s.Reshuffles)
	var lrtp, srtp = s.LinePay / reshuf / sbl * 100, s.ScatPay / reshuf / sbl * 100
	var rtpsym = lrtp + srtp
	var fghits = float64(s.FreeHits)
	var fsavr1, multavr = FirstSreespins()
	var q = fghits * fsavr1 / total
	var sq = 1 / (1 - fghits*multavr*15/total)
	var rtp = rtpsym + q*sq*rtpsym
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/total*100, g.SBL.Num(), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("average plain freespins at 1st iteration: %g\n", fsavr1)
	fmt.Printf("average multiplier at free games: %g\n", multavr)
	fmt.Printf("free games %g, q = %.5g, sq = %.5g\n", fghits, q, sq)
	fmt.Printf("free games frequency: 1/%.5g\n", reshuf/fghits)
	fmt.Printf("RTP = rtpsym + q*sq*rtpsym = %.5g + %.5g = %.6f%%\n", rtpsym, q*sq*rtpsym, rtp)
	return rtp
}
