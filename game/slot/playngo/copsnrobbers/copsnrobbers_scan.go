package copsnrobbers

import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("\n(1/2) bonus reels calculations\n")
	var sb = slot.NewStatGeneric(sn, 5)
	{
		var reels = ReelsBon
		var g = NewGame(sp.Sel)
		g.FSR = Efs // set free spins mode
		g.M = 1
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_simple(w, sp, sb, g.Cost())
		}
		slot.ScanReelsCommon(ctx, sp, sb, g, reels, calc)
	}

	if ctx.Err() != nil {
		return 0, 0
	}

	fmt.Printf("\n(2/2) regular reels calculations\n")
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame(sp.Sel)
	var sr = slot.NewStatGeneric(sn, 5)

	// custom parsheet
	var cost = g.Cost()
	var calc = func(w io.Writer) (float64, float64) {
		// bonus reels parameters
		var Nb, Sb, Qb = sb.NSQ(cost)
		var µb = Sb / Nb
		var Dsymb = Qb/Nb - µb*µb
		// regular reels parameters
		var Nr, Sr, Qr = sr.NSQ(cost)
		var µr = Sr / Nr
		var Dsymr = Qr/Nr - µr*µr
		// calculation
		var EL, EL2 float64
		for _, m := range Freegames {
			var m = float64(m)
			EL += m
			EL2 += m * m
		}
		EL /= float64(len(Freegames))
		EL2 /= float64(len(Freegames))
		var Em = 2*Pmfs + 1*(1-Pmfs)
		var Em2 = 2*2*Pmfs + 1*1*(1-Pmfs)
		var Pfg = float64(sr.FGH.Load()) / Nr
		var q = Pfg * EL
		var rtp = µr + q*Em*µb
		var Vbon = Em2*EL*Dsymb + µb*µb*(Em2*EL2-Em*EL*Em*EL)
		var D = Dsymr + Pfg*Vbon + Pfg*(1-Pfg)*(µb-µr)*(µb-µr)

		if sp.IsFG() {
			fmt.Fprintf(w, "*bonus reels*\n")
			fmt.Fprintf(w, "RTP(fg) = %.8g%%\n", Em*µb*100)
		}
		if sp.IsMain() {
			fmt.Fprintf(w, "*regular reels*\n")
			fmt.Fprintf(w, "symbols: µ = %.8g%%, sigma(sym) = %.6g\n", µr*100, math.Sqrt(Dsymr))
			fmt.Fprintf(w, "free spins: q = %.5g\n", q)
			fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", 1/Pfg)
			fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.8g%%\n", µr*100, q, Em*µb*100, rtp*100)
		}
		slot.Print_all(w, sp, sr, rtp, D)
		return rtp, D
	}

	return slot.ScanReelsCommon(ctx, sp, sr, g, reels, calc)
}
