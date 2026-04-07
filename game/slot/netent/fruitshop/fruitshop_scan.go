package fruitshop

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func ΣPL(s *slot.StatGeneric) (sum float64) {
	var N = s.Count()
	for sym, L := range LineFreespinReg {
		for i, Li := range L {
			var Pfgi = float64(s.C[sym][i].Load()) / N
			sum += Pfgi * float64(Li)
		}
	}
	return
}

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("\n(1/2) free games calculations\n")
	var sb = slot.NewStatGeneric(sn, 5)
	{
		var reels, _ = ReelsMap.FindClosest(sp.MRTP)
		var g = NewGame()
		g.FSR = 5 // set free spins mode
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_fgretrig_custom(w, sp, sb, g.Cost(), 1, sb.FSQ(), ΣPL(sb))
		}
		slot.ScanReelsCommon(ctx, sp, sb, g, reels, calc)
	}

	if ctx.Err() != nil {
		return 0, 0
	}

	fmt.Printf("\n(2/2) regular games calculations\n")
	var sr = slot.NewStatGeneric(sn, 5)
	{
		var reels, _ = ReelsMap.FindClosest(sp.MRTP)
		var g = NewGame()
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_fgretrig_split_custom(w, sp, sr, sb, g.Cost(), 1, sr.FSQ(), sb.FSQ(), ΣPL(sr))
		}
		return slot.ScanReelsCommon(ctx, sp, sr, g, reels, calc)
	}
}
