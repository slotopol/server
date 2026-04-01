package jaguarmoon

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
		var g = NewGame()
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
		var g = NewGame()
		sr.SymDim(scat, 6)
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
			var Pfg, PLm, PLm2, PL2m2 float64
			for i := range 6 {
				var Pfgi = float64(sr.C[scat-1][i].Load()) / Nr
				var L = float64(ScatFreespin[i])
				var m = float64(FreeMult[i])
				var plm = Pfgi * L * m
				Pfg += Pfgi
				PLm += plm
				PLm2 += plm * m
				PL2m2 += plm * L * m
			}
			var rtp = µr + PLm*µb
			var D = Dsymr + Dsymb*PLm2 + µb*µb*PL2m2
			if sp.IsFG() {
				fmt.Fprintf(w, "*bonus reels*\n")
				fmt.Fprintf(w, "RTP(fg) = %.8g%%\n", µb*100)
			}
			if sp.IsMain() {
				fmt.Fprintf(w, "*regular reels*\n")
				fmt.Fprintf(w, "symbols: µ = %.8g%%, sigma(sym) = %.6g\n", µr*100, math.Sqrt(Dsymr))
				fmt.Fprintf(w, "free games: HRfg = 1/%.5g, q = %.5g\n", 1/Pfg, PLm)
				fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.8g%%\n", µr*100, PLm, µb*100, rtp*100)
			}
			slot.Print_all(w, sp, sr, rtp, D)
			return rtp, D
		}
		return slot.ScanReelsCommon(ctx, sp, sr, g, reels, calc)
	}
}
