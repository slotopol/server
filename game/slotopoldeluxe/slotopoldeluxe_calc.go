package slotopoldeluxe

import (
	"context"
	"fmt"
	"time"

	"github.com/slotopol/server/game"
	"github.com/slotopol/server/game/slotopol"
)

func CalcStat(ctx context.Context, rn string) float64 {
	fmt.Printf("*bonus games calculations*\n")
	slotopol.Emje = slotopol.ExpEldorado()
	slotopol.Emjm = slotopol.ExpMonopoly()
	fmt.Printf("*reels calculations*\n")
	var reels *game.Reels5x
	if rn != "" {
		var ok bool
		if reels, ok = ReelsMap[rn]; !ok {
			return 0
		}
	} else {
		rn, reels = "104", &Reels104
	}
	var g = NewGame(rn)
	g.SBL = game.MakeBitNum(1)
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
	var qmje1 = float64(s.BonusCount[mje1]) / reshuf / sbl
	var rtpmje1 = slotopol.Emje * 1 * qmje1 * 100
	var qmje3 = float64(s.BonusCount[mje3]) / reshuf / sbl
	var rtpmje3 = slotopol.Emje * 3 * qmje3 * 100
	var qmje6 = float64(s.BonusCount[mje6]) / reshuf / sbl
	var rtpmje6 = slotopol.Emje * 6 * qmje6 * 100
	var qmjm = float64(s.BonusCount[mjm]) / reshuf / sbl
	var rtpmjm = slotopol.Emjm * qmjm * 100
	var rtp = rtpsym + rtpmje1 + rtpmje3 + rtpmje6 + rtpmjm
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/total*100, g.SBL.Num(), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("spin1 bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/float64(s.BonusCount[mje1]), rtpmje1)
	fmt.Printf("spin3 bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/float64(s.BonusCount[mje3]), rtpmje3)
	fmt.Printf("spin6 bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/float64(s.BonusCount[mje6]), rtpmje6)
	fmt.Printf("monopoly bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/float64(s.BonusCount[mjm]), rtpmjm)
	if s.JackCount[jid] > 0 {
		fmt.Printf("jackpots: count %d, frequency 1/%d\n", s.JackCount[jid], int(reshuf/float64(s.JackCount[jid])))
	}
	fmt.Printf("RTP = %.5g(sym) + %.5g(mje) + %.5g(mjm) = %.6f%%\n", rtpsym, rtpmje1+rtpmje3+rtpmje6, rtpmjm, rtp)
	return rtp
}
