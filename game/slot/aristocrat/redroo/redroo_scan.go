package redroo

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("\n(1/2) free games calculations\n")
	var sb = slot.NewStatGeneric(sn, 5)
	{
		var reels, _ = ReelsMap.FindClosest(sp.MRTP)
		var g = NewGame()
		g.FSR = 12 // set free spins mode
		g.MW = [3]float64{2.5, 2.5, 2.5}
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_generic_fgretrig_series(w, sp, sb, g.Cost(), 1, ScatFreespinBon[:], scat)
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
		g.MW = [3]float64{1, 1, 1}
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_generic_fgretrig_split_series(w, sp, sr, sb, g.Cost(), 1, ScatFreespinReg[:], scat)
		}
		return slot.ScanReelsCommon(ctx, sp, sr, g, reels, calc)
	}
}
