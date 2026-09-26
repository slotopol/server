package groovypowers

import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("\n(1/2) bonus games calculations\n")
	var sb = slot.NewStatGeneric(sn, 5)
	{
		var g = NewGame(sp.Sel)
		g.BM = true // set bonus mode
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_simple(w, sp, sb, g.Cost())
		}
		slot.ScanReelsCommon(ctx, sp, sb, g, calc)
	}

	if ctx.Err() != nil {
		return 0, 0
	}

	fmt.Printf("\n(2/2) regular games calculations\n")
	var sr = slot.NewStatGeneric(sn, 5)
	{
		var g = NewGame(sp.Sel)

		// custom parsheet
		var cost = g.Cost()
		var calc = func(w io.Writer) (float64, float64) {
			// bonus reels parameters
			var µb, Dsymb = sb.EvD(cost)
			// regular reels parameters
			var µr, Dsymr = sr.EvD(cost)
			var rtp = µr*(1-Pbm) + µb*Pbm
			var D = (1-Pbm)*Dsymr + Pbm*Dsymb + Pbm*(1-Pbm)*(µb-µr)*(µb-µr)
			if sp.IsMain() {
				fmt.Fprintf(w, "bon symbols: µb = %.8g%%, sigma(sym)b = %.6g\n", µb*100, math.Sqrt(Dsymb))
				fmt.Fprintf(w, "reg symbols: µr = %.8g%%, sigma(sym)r = %.6g\n", µr*100, math.Sqrt(Dsymr))
				fmt.Fprintf(w, "RTP = %.5g(reg)*%.5g + %.5g(bm)*%.5g = %.6f%%\n", µr*100, 1-Pbm, µb*100, Pbm, rtp*100)
			}
			slot.Print_all(w, sp, sr, rtp, D)
			return rtp, D
		}

		return slot.ScanReelsCommon(ctx, sp, sr, g, calc)
	}
}
