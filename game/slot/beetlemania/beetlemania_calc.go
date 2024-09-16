package beetlemania

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"sync/atomic"
	"time"

	slot "github.com/slotopol/server/game/slot"
)

type Stat struct {
	slot.Stat
	FGNum [5]uint64
	FGPay uint64
}

func (s *Stat) Update(wins slot.Wins) {
	s.Stat.Update(wins)
	if len(wins) > 0 {
		if wi := wins[len(wins)-1]; wi.Free > 0 {
			atomic.AddUint64(&s.FGNum[wi.Num-1], 1)
			atomic.AddUint64(&s.FGPay, uint64(wins.Gain()))
		}
	}
}

func CalcStatBon(ctx context.Context, rn string) float64 {
	var reels *slot.Reels5x
	if mrtp, _ := strconv.ParseFloat(rn, 64); mrtp != 0 {
		var _, r = FindReels(mrtp)
		reels = r.(*slot.Reels5x)
	} else {
		reels = &ReelsReg92
	}
	var g = NewGame()
	var sln float64 = 5
	g.Sel.SetNum(int(sln), 1)
	g.FS = 10 // set free spins mode
	var s Stat

	var total = float64(reels.Reshuffles())
	var dur = func() time.Duration {
		var t0 = time.Now()
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		go s.Progress(ctx2, time.Tick(2*time.Second), sln, total)
		slot.BruteForce5x(ctx2, &s, g, reels)
		return time.Since(t0)
	}()

	var reshuf = float64(s.Reshuffles)
	var lrtp, srtp = s.LinePay / reshuf / sln * 100, s.ScatPay / reshuf / sln * 100
	var rtpsym = lrtp + srtp
	var qjazz = float64(s.BonusCount[jbonus]) / reshuf / sln
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/total*100, g.Sel.Num(), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("jazzbee bonuses: count %d, q = 1/%g\n", s.BonusCount[jbonus], 1/qjazz)
	if s.JackCount[jid] > 0 {
		fmt.Printf("jackpots: count %d, frequency 1/%d\n", s.JackCount[jid], int(reshuf/float64(s.JackCount[jid])))
	}
	fmt.Printf("RTP = rtp(sym) = %.6f%%\n", rtpsym)
	return rtpsym
}

func CalcStatReg(ctx context.Context, rn string) float64 {
	fmt.Printf("*bonus reels calculations*\n")
	var rtpfs float64
	if rn != "" && rn[len(rn)-1:] == "u" {
		rtpfs = CalcStatBon(ctx, "bonu")
	} else {
		rtpfs = CalcStatBon(ctx, "bon")
	}
	if ctx.Err() != nil {
		return 0
	}
	fmt.Printf("*regular reels calculations*\n")
	var reels *slot.Reels5x
	if mrtp, _ := strconv.ParseFloat(rn, 64); mrtp != 0 {
		var _, r = FindReels(mrtp)
		reels = r.(*slot.Reels5x)
	} else {
		reels = &ReelsReg92
	}
	var g = NewGame()
	var sln float64 = 5
	g.Sel.SetNum(int(sln), 1)
	var s Stat

	var total = float64(reels.Reshuffles())
	var dur = func() time.Duration {
		var t0 = time.Now()
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		go s.Progress(ctx2, time.Tick(2*time.Second), sln, total)
		slot.BruteForce5x(ctx2, &s, g, reels)
		return time.Since(t0)
	}()

	var reshuf = float64(s.Reshuffles)
	var lrtp, srtp = s.LinePay / reshuf / sln * 100, s.ScatPay / reshuf / sln * 100
	var rtpsym = lrtp + srtp
	var fgsum = float64(s.FGNum[2] + s.FGNum[3] + s.FGNum[4])
	var fgpay = float64(s.FGPay) / fgsum
	var rtpbon = (fgpay + rtpfs*10/100) * math.Pow(2, 1.25)
	var q = fgsum / total
	var rtp = rtpsym + q*rtpbon
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/total*100, g.Sel.Num(), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("free games numbers: [0, 0, %d, %d, %d]\n", s.FGNum[2], s.FGNum[3], s.FGNum[4])
	fmt.Printf("free games %g, q = %.6f\n", fgsum, q)
	fmt.Printf("free games frequency: 1/%.5g\n", reshuf/float64(s.FreeHits))
	fmt.Printf("average pay by freespins start %.6f\n", fgpay)
	fmt.Printf("rtpbon = (fgpay+rtpfs*10)*2^10/8 = %.6f\n", rtpbon)
	if s.JackCount[jid] > 0 {
		fmt.Printf("jackpots: count %d, frequency 1/%d\n", s.JackCount[jid], int(reshuf/float64(s.JackCount[jid])))
	}
	fmt.Printf("RTP = %.5g(sym) + q*%.5g(bon) = %.6f%%\n", rtpsym, rtpbon, rtp)
	return rtp
}
