package greatblue

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/slotopol/server/game"
)

type Stat struct {
	game.Stat
	FGNum [5]uint64
}

func (s *Stat) Update(sw *game.WinScan) {
	s.Stat.Update(sw)
	if len(sw.Wins) > 0 {
		if wi := sw.Wins[len(sw.Wins)-1]; wi.Sym == scat {
			atomic.AddUint64(&s.FGNum[wi.Num-1], 1)
		}
	}
}

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
	if rn != "" {
		var ok bool
		if reels, ok = ReelsMap[rn]; !ok {
			return 0
		}
	} else {
		rn, reels = "92", &Reels92
	}
	var g = NewGame(rn)
	g.SBL = game.MakeSblNum(5)
	var sbl = float64(g.SBL.Num())
	var s Stat

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
	var fgsum = float64(s.FGNum[2] + s.FGNum[3] + s.FGNum[4])
	var fsavr1, multavr = FirstSreespins()
	var q = fgsum * fsavr1 / total
	var sq = 1 / (1 - fgsum*multavr*15/total)
	var rtp = rtpsym + q*sq*rtpsym
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/total*100, g.SBL.Num(), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("free games numbers: [0, 0, %d, %d, %d]\n", s.FGNum[2], s.FGNum[3], s.FGNum[4])
	fmt.Printf("free games frequency: 1/%.5g\n", reshuf/fgsum)
	fmt.Printf("average plain freespins at 1st iteration: %g\n", fsavr1)
	fmt.Printf("average multiplier at free games: %g\n", multavr)
	fmt.Printf("free games %g, q = %.5g, sq = %.5g\n", fgsum, q, sq)
	fmt.Printf("RTP = rtpsym + q*sq*rtpsym = %.5g + %.5g = %.6f%%\n", rtpsym, q*sq*rtpsym, rtp)
	return rtp
}
