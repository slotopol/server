package wildhorses

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
		g.FSR = 10 // set free spins mode
		var calc = func(w io.Writer) (float64, float64) {
			var q = 0.5 // probability of getting or not getting free spins is equal
			var ΣPL = sb.FGQ() * 10
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
			var qr = sr.FSQ()
			var qb = 0.5 // probability of getting or not getting free spins is equal
			var ΣPL = sb.FGQ() * 10
			return slot.Parsheet_fgretrig_split_custom(w, sp, sr, sb, g.Cost(), 1, qr, qb, ΣPL)
		}
		return slot.ScanReelsCommon(ctx, sp, sr, g, reels, calc)
	}
}
