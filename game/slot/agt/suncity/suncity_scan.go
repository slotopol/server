package suncity

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
		g.FSR = -1 // set free spins mode
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_simple(w, sp, sb, g.Cost())
		}
		slot.ScanReelsCommon(ctx, sp, sb, g, reels, calc)
	}

	if ctx.Err() != nil {
		return 0, 0
	}

	fmt.Printf("\n(2/2) regular reels calculations\n")
	var sr = slot.NewStatGeneric(sn, 5)
	{
		var reels, _ = ReelsMap.FindClosest(sp.MRTP)
		var g = NewGame(sp.Sel)
		// custom parsheet
		var cost = g.Cost()
		var m = 1.0
		var calc = func(w io.Writer) (float64, float64) {
			// bonus reels parameters
			var Nb, Sb, Qb = sb.NSQ(cost)
			var µb = Sb / Nb
			var Dsymb = Qb/Nb - µb*µb
			var Pfgb = sb.FGQ()
			// regular reels parameters
			var Nr, Sr, Qr = sr.NSQ(cost)
			var µr = Sr / Nr
			var Dsymr = Qr/Nr - µr*µr
			var Pfgr = sr.FGQ()
			// calculation
			var q = Pfgr / Pfgb
			var rtp = µr + m*q*µb
			var EL = 1 / Pfgb
			var VL = (1 - Pfgb) / Pfgb / Pfgb
			var Vbon = m * m * (EL*Dsymb + µb*µb*VL)
			var D = Dsymr + Pfgr*Vbon + Pfgr*(1-Pfgr)*(EL*m*µb-µr)*(EL*m*µb-µr)
			if sp.IsFG() {
				fmt.Fprintf(w, "*bonus reels*\n")
				fmt.Fprintf(w, "RTP(fg) = %.8g%%\n", µb*100)
			}
			if sp.IsMain() {
				fmt.Fprintf(w, "*regular reels*\n")
				fmt.Fprintf(w, "symbols: µ = %.8g%%, sigma(sym) = %.6g\n", µr*100, math.Sqrt(Dsymr))
				fmt.Fprintf(w, "free games: HRfg = 1/%.5g, EL = %.5g, q = %.5g\n", 1/Pfgr, EL, q)
				fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.8g%%\n", µr*100, q, µb*100, rtp*100)
			}
			slot.Print_all(w, sp, sr, rtp, D)
			return rtp, D
		}
		return slot.ScanReelsCommon(ctx, sp, sr, g, reels, calc)
	}
}
