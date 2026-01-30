package champagne

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

var (
	EVbot float64 // expectation of 1 bottle
	EVmjc float64 // Bottle game calculated expectation
)

func ExpBottle() {
	// avr 1 bottle gain
	EVbot = 0
	for _, v := range Bottles {
		EVbot += float64(v)
	}
	EVbot /= float64(len(Bottles))

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
	EVmjc = E
}

func CalcStatBon(ctx context.Context, mrtp float64) float64 {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	var g = NewGame(1)
	g.FSR = 15 // set free spins mode
	var s slot.StatGeneric

	var calc = func(w io.Writer) float64 {
		var reshuf = s.Count()
		var lrtp, srtp = s.SymRTP(g.Cost())
		var rtpsym = lrtp + srtp
		var q, sq = s.FSQ()
		var qmjc = s.BonusHitsF(mjc) / reshuf / float64(g.Sel)
		var rtpmjc = EVmjc * qmjc * 100
		var rtp = sq * (rtpsym + rtpmjc)
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCount.Load(), q, sq)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", s.FGF())
		fmt.Fprintf(w, "bottle bonuses: hit rate 1/%.5g, rtp = %.6f%%\n", reshuf/s.BonusHitsF(mjc), rtpmjc)
		if s.JackHitsF(mjj) > 0 {
			fmt.Fprintf(w, "jackpots: count %g, frequency 1/%.12g\n", s.JackHitsF(mjj), reshuf/s.JackHitsF(mjj))
		}
		fmt.Fprintf(w, "RTP = sq*(rtp(sym)+rtp(mjc)) = %.5g*(%.5g+%.5g) = %.6f%%\n", sq, rtpsym, rtpmjc, rtp)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}

func CalcStatReg(ctx context.Context, mrtp float64) float64 {
	fmt.Printf("*bonus games calculations*\n")
	ExpBottle()
	fmt.Printf("len = %d, avr bottle gain = %.5g, EV = %g\n", len(Bottles), EVbot, EVmjc)
	fmt.Printf("*free games calculations*\n")
	var rtpfs = CalcStatBon(ctx, mrtp)
	if ctx.Err() != nil {
		return 0
	}
	fmt.Printf("*regular games calculations*\n")
	var reels, _ = ReelsMap.FindClosest(mrtp)
	var g = NewGame(1)
	g.FSR = 0 // no free spins
	var s slot.StatGeneric

	var calc = func(w io.Writer) float64 {
		var reshuf = s.Count()
		var lrtp, srtp = s.SymRTP(g.Cost())
		var rtpsym = lrtp + srtp
		var q, _ = s.FSQ()
		var qmjc = s.BonusHitsF(mjc) / reshuf / float64(g.Sel)
		var rtpmjc = EVmjc * qmjc * 100
		var rtp = rtpsym + rtpmjc + q*rtpfs
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.6f\n", s.FreeCount.Load(), q)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", s.FGF())
		fmt.Fprintf(w, "champagne bonuses: hit rate 1/%.5g, rtp = %.6f%%\n", reshuf/s.BonusHitsF(mjc), rtpmjc)
		if s.JackHitsF(mjj) > 0 {
			fmt.Fprintf(w, "jackpots: count %g, frequency 1/%.12g\n", s.JackHitsF(mjj), reshuf/s.JackHitsF(mjj))
		}
		fmt.Fprintf(w, "RTP = rtp(sym) + rtp(mjc) + q*rtp(fg) = %.5g + %.5g + %.5g*%.5g = %.6f%%\n", rtpsym, rtpmjc, q, rtpfs, rtp)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}
