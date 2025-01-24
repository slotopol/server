package aztecgold

import (
	"context"
	"fmt"
	"time"

	"github.com/slotopol/server/game/slot"
)

var (
	Epyr  float64 // Expectation of pyramid choose by 4 tries
	Eroom float64 // Expectation of room episode
	Ebon  float64 // Expectation of whole bonus game
)

func ExpBonus() {
	Epyr = 0
	for i, p := range app {
		Epyr += p * apm[i]
	}

	var sum1 float64
	for _, cell := range Room[0] {
		if cell.Sym != 14 {
			sum1 += cell.Mult
		}
	}
	sum1 += 4 * sum1 / float64(len(Room[0])-1)
	var Erow1 = sum1 / float64(len(Room[0]))

	var sum2 float64
	for _, cell := range Room[1] {
		if cell.Sym != 14 {
			sum2 += cell.Mult
		}
	}
	sum2 += 4 * sum2 / float64(len(Room[1])-1)
	var Erow2 = sum2 / float64(len(Room[1]))

	var sum3 float64
	for _, cell := range Room[2] {
		if cell.Sym != 14 {
			sum3 += cell.Mult
		}
	}
	sum3 += 4 * sum3 / float64(len(Room[2])-1)
	var Erow3 = sum3 / float64(len(Room[2]))

	var Erow4 = (50 + 50 + 100 + 0) / 4.
	var p4 = 4 / float64(len(Room[3]))

	var Erow5 = 250.
	var p5 = 1 / float64(len(Room[4]))

	Eroom = Erow1 + Erow2 + Erow3 + Erow4*p4 + Erow5*p5*p4
	Ebon = Epyr + Eroom*app[5]
}

func CalcStat(ctx context.Context, mrtp float64) float64 {
	fmt.Printf("*bonus games calculations*\n")
	ExpBonus()
	fmt.Printf("Ebon = Epyr + Eroom*app[6] = %.5g + %.5g * %.5g = %g\n", Epyr, Eroom, app[5], Ebon)
	fmt.Printf("*reels calculations*\n")
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
	var qmjap = float64(s.BonusCount[mjap]) / reshuf
	var rtpmjap = Ebon * qmjap * 100
	var rtp = rtpsym + rtpmjap
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/float64(s.Planned())*100, g.Sel, dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("pyramid bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/float64(s.BonusCount[mjap]), rtpmjap)
	if s.JackCount[jid] > 0 {
		fmt.Printf("jackpots: count %d, frequency 1/%.12g\n", s.JackCount[jid], reshuf/float64(s.JackCount[jid]))
	}
	fmt.Printf("RTP = %.5g(sym) + %.5g(mjap) = %.6f%%\n", rtpsym, rtpmjap, rtp)
	return rtp
}
