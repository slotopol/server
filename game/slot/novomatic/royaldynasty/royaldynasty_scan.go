package royaldynasty

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("\n(1/2) bonus reels calculations\n")
	var sb = slot.NewStatGeneric(sn, 5)
	{
		var reels = ReelsBon
		var g = NewGame(sp.Sel)
		g.FSR = 35 // set free spins mode
		g.TS = scat1
		var calc = func(w io.Writer) (float64, float64) {
			var N = sb.Count()
			var q = float64(sb.FGH.Load()*35) / N
			var Pfgi = float64(sb.FGH.Load()) / float64(len(Freegames)) / N
			var ΣPL float64
			for _, Li := range Freegames {
				ΣPL += Pfgi * float64(Li)
			}
			return slot.Parsheet_fgretrig_custom(w, sp, sb, g.Cost(), 1, q, ΣPL)
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
		var calc = func(w io.Writer) (float64, float64) {
			// bonus reels parameters
			var Nb = sb.Count()
			var qb = float64(sb.FGH.Load()*35) / Nb
			// regular reels parameters
			var Nr = sr.Count()
			var qr = float64(sr.FGH.Load()*35) / Nr
			var Pfgi = float64(sr.FGH.Load()) / float64(len(Freegames)) / Nr
			var ΣPL float64
			for _, Li := range Freegames {
				ΣPL += Pfgi * float64(Li)
			}
			return slot.Parsheet_fgretrig_split_custom(w, sp, sr, sb, g.Cost(), 1, qr, qb, ΣPL)
		}
		return slot.ScanReelsCommon(ctx, sp, sr, g, reels, calc)
	}
}
