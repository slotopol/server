package champagne

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

var (
	Ebot float64 // expectation of 1 bottle
	Emjc float64 // Bottle game calculated expectation
)

func ExpBottle() {
	// avr 1 bottle gain
	Ebot = 0
	for _, v := range Bottles {
		Ebot += float64(v)
	}
	Ebot /= float64(len(Bottles))

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
	Emjc = E
}

func CalcStatBon(ctx context.Context, mrtp float64) float64 {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	var g = NewGame()
	g.Sel = 1
	g.FSR = 15 // set free spins mode
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var reshuf = s.Count()
		var cost, _ = g.Cost()
		var lrtp, srtp = s.LineRTP(cost), s.ScatRTP(cost)
		var rtpsym = lrtp + srtp
		var q = s.FreeCount() / reshuf
		var sq = 1 / (1 - q)
		var qmjc = s.BonusCount(mjc) / reshuf / float64(g.Sel)
		var rtpmjc = Emjc * qmjc * 100
		var rtp = sq * (rtpsym + rtpmjc)
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCountU(), q, sq)
		fmt.Fprintf(w, "free games frequency: 1/%.5g\n", reshuf/s.FreeHits())
		fmt.Fprintf(w, "champagne bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/s.BonusCount(mjc), rtpmjc)
		if s.JackCount(mjj) > 0 {
			fmt.Fprintf(w, "jackpots: count %g, frequency 1/%.12g\n", s.JackCount(mjj), reshuf/s.JackCount(mjj))
		}
		fmt.Fprintf(w, "RTP = sq*(rtp(sym)+rtp(mjc)) = %.5g*(%.5g+%.5g) = %.6f%%\n", sq, rtpsym, rtpmjc, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc)
}

func CalcStatReg(ctx context.Context, mrtp float64) float64 {
	fmt.Printf("*bonus games calculations*\n")
	ExpBottle()
	fmt.Printf("len = %d, avr bottle gain = %.5g, E = %g\n", len(Bottles), Ebot, Emjc)
	fmt.Printf("*bonus reels calculations*\n")
	var rtpfs = CalcStatBon(ctx, mrtp)
	if ctx.Err() != nil {
		return 0
	}
	fmt.Printf("*regular reels calculations*\n")
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	var g = NewGame()
	g.Sel = 1
	g.FSR = 0 // no free spins
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var reshuf = s.Count()
		var cost, _ = g.Cost()
		var lrtp, srtp = s.LineRTP(cost), s.ScatRTP(cost)
		var rtpsym = lrtp + srtp
		var q = s.FreeCount() / reshuf
		var qmjc = s.BonusCount(mjc) / reshuf / float64(g.Sel)
		var rtpmjc = Emjc * qmjc * 100
		var rtp = rtpsym + rtpmjc + q*rtpfs
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.6f\n", s.FreeCountU(), q)
		fmt.Fprintf(w, "free games frequency: 1/%.5g\n", reshuf/s.FreeHits())
		fmt.Fprintf(w, "champagne bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/s.BonusCount(mjc), rtpmjc)
		if s.JackCount(mjj) > 0 {
			fmt.Fprintf(w, "jackpots: count %g, frequency 1/%.12g\n", s.JackCount(mjj), reshuf/s.JackCount(mjj))
		}
		fmt.Fprintf(w, "RTP = rtp(sym) + rtp(mjc) + q*rtp(fg) = %.5g + %.5g + %.5g*%.5g = %.6f%%\n", rtpsym, rtpmjc, q, rtpfs, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc)
}
