package greatblue

import (
	"context"
	"fmt"
	"time"

	slot "github.com/slotopol/server/game/slot"
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

func CalcStat(ctx context.Context, mrtp float64) float64 {
	var reels, _ = slot.FindReels(ReelsMap, mrtp)
	var g = NewGame()
	var sln float64 = 5
	g.Sel.SetNum(int(sln), 1)
	var s slot.Stat

	var dur = slot.ScanReels(ctx, &s, g, reels,
		time.Tick(2*time.Second), time.Tick(2*time.Second))

	var reshuf = float64(s.Reshuffles)
	var lrtp, srtp = s.LinePay / reshuf / sln * 100, s.ScatPay / reshuf / sln * 100
	var rtpsym = lrtp + srtp
	var fghits = float64(s.FreeHits)
	var fsavr1, multavr = FirstSreespins()
	var q = fghits * fsavr1 / reshuf
	var sq = 1 / (1 - fghits*multavr*15/reshuf)
	var rtp = rtpsym + q*sq*rtpsym
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/float64(s.Planned())*100, g.Sel.Num(), dur)
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
