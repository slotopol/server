package trolls

import (
	"context"
	"fmt"
	"time"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, mrtp float64) float64 {
	var reels, _ = slot.FindReels(ReelsMap, mrtp)
	var g = NewGame()
	var sln float64 = 1
	g.Sel = int(sln)
	var s slot.Stat

	var dur = slot.ScanReels5x(ctx, &s, g, reels,
		time.Tick(2*time.Second), time.Tick(2*time.Second))

	var reshuf = float64(s.Reshuffles)
	var lrtp, srtp = s.LinePay / reshuf / sln * 100, s.ScatPay / reshuf / sln * 100
	var rtpsym = lrtp + srtp
	var q = float64(s.FreeCount) / reshuf
	var sq = 1 / (1 - q)
	var rtpfs = sq * 3 * rtpsym
	var rtp = rtpsym + q*rtpfs
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/float64(s.Planned())*100, g.Sel, dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCount, q, sq)
	fmt.Printf("free games frequency: 1/%.5g\n", reshuf/float64(s.FreeHits))
	if s.JackCount[jid] > 0 {
		fmt.Printf("jackpots: count %d, frequency 1/%d\n", s.JackCount[jid], int(reshuf/float64(s.JackCount[jid])))
	}
	fmt.Printf("RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym, q, rtpfs, rtp)
	return rtp
}
