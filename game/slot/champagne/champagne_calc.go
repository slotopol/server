package champagne

import (
	"context"
	"fmt"
	"time"

	slot "github.com/slotopol/server/game/slot"
)

var Emjc float64 // Bottle game calculated expectation

func ExpBottle() float64 {
	// avr 1 bottle gain
	var m float64
	for _, v := range Bottles {
		m += float64(v)
	}
	m /= float64(len(Bottles))

	// expectation
	var E float64
	var n = 0
	for i := 0; i < len(Bottles); i++ {
		for j := i + 1; j < len(Bottles); j++ {
			if Bottles[i] == Bottles[j] {
				E += float64(Bottles[i]) * 4
			} else {
				E += float64(Bottles[i] + Bottles[j])
			}
			n++
		}
	}
	E /= float64(n)
	fmt.Printf("len = %d, avr bottle gain = %.5g, E = %g\n", len(Bottles), m, E)
	return E
}

func CalcStatBon(ctx context.Context, mrtp float64) float64 {
	var reels, _ = slot.FindReels(ReelsMap, mrtp)
	var g = NewGame()
	var sln float64 = 1
	g.Sel.SetNum(int(sln), 1)
	g.FS = 15 // set free spins mode
	var s slot.Stat

	var dur = slot.ScanReels(ctx, &s, g, reels,
		time.Tick(2*time.Second), time.Tick(2*time.Second))

	var reshuf = float64(s.Reshuffles)
	var lrtp, srtp = s.LinePay / reshuf / sln * 100, s.ScatPay / reshuf / sln * 100
	var rtpsym = lrtp + srtp
	var q = float64(s.FreeCount) / reshuf
	var sq = 1 / (1 - q)
	var qmjc = float64(s.BonusCount[mjc]) / reshuf / sln
	var rtpmjc = Emjc * qmjc * 100
	var rtp = sq * (rtpsym + rtpmjc)
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/float64(s.Planned())*100, g.Sel.Num(), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCount, q, sq)
	fmt.Printf("free games frequency: 1/%.5g\n", reshuf/float64(s.FreeHits))
	fmt.Printf("champagne bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/float64(s.BonusCount[mjc]), rtpmjc)
	if s.JackCount[jid] > 0 {
		fmt.Printf("jackpots: count %d, frequency 1/%d\n", s.JackCount[jid], int(reshuf/float64(s.JackCount[jid])))
	}
	fmt.Printf("RTP = sq*(rtp(sym)+rtp(mjc)) = %.5g*(%.5g+%.5g) = %.6f%%\n", sq, rtpsym, rtpmjc, rtp)
	return rtp
}

func CalcStatReg(ctx context.Context, mrtp float64) float64 {
	fmt.Printf("*bonus games calculations*\n")
	Emjc = ExpBottle()
	fmt.Printf("*bonus reels calculations*\n")
	var rtpfs = CalcStatBon(ctx, mrtp)
	if ctx.Err() != nil {
		return 0
	}
	fmt.Printf("*regular reels calculations*\n")
	var reels, _ = slot.FindReels(ReelsMap, mrtp)
	var g = NewGame()
	var sln float64 = 1
	g.Sel.SetNum(int(sln), 1)
	g.FS = 0 // no free spins
	var s slot.Stat

	var dur = slot.ScanReels(ctx, &s, g, reels,
		time.Tick(2*time.Second), time.Tick(2*time.Second))

	var reshuf = float64(s.Reshuffles)
	var lrtp, srtp = s.LinePay / reshuf / sln * 100, s.ScatPay / reshuf / sln * 100
	var rtpsym = lrtp + srtp
	var q = float64(s.FreeCount) / reshuf
	var qmjc = float64(s.BonusCount[mjc]) / reshuf / sln
	var rtpmjc = Emjc * qmjc * 100
	var rtp = rtpsym + rtpmjc + q*rtpfs
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/float64(s.Planned())*100, g.Sel.Num(), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("free spins %d, q = %.6f\n", s.FreeCount, q)
	fmt.Printf("free games frequency: 1/%.5g\n", reshuf/float64(s.FreeHits))
	fmt.Printf("champagne bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/float64(s.BonusCount[mjc]), rtpmjc)
	if s.JackCount[jid] > 0 {
		fmt.Printf("jackpots: count %d, frequency 1/%d\n", s.JackCount[jid], int(reshuf/float64(s.JackCount[jid])))
	}
	fmt.Printf("RTP = rtp(sym) + rtp(mjc) + q*rtp(fg) = %.5g + %.5g + %.5g*%.5g = %.6f%%\n", rtpsym, rtpmjc, q, rtpfs, rtp)
	return rtp
}
