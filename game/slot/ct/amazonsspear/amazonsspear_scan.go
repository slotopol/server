package amazonsspear

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var ok bool
	var sb, sr *slot.StatGeneric
	fmt.Printf("\n(1/2) bonus reels calculations\n")
	var idb = fmt.Sprintf("ctinteractive/amazonsspear/graw/bon/%d", sp.Sel)
	if sb, ok = slot.GetStatGeneric(idb, sn, 5); ok {
		var reels = ReelsBon
		var g = NewGame(sp.Sel)
		g.FSR = 15 // set free spins mode
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_generic_fgretrig(w, sp, sb, g.Cost(), 1, 15)
		}
		slot.ScanReelsCommon(ctx, sp, sb, g, reels, calc)
	}

	if ctx.Err() != nil {
		return 0, 0
	}

	fmt.Printf("\n(2/2) regular reels calculations\n")
	var idr = fmt.Sprintf("ctinteractive/amazonsspear/graw/reg/%d", sp.Sel)
	if sr, ok = slot.FindStatGeneric(idr+"/%g", sp.MRTP, sn, 5); ok {
		var reels, _ = ReelsMap.FindClosest(sp.MRTP)
		var g = NewGame(sp.Sel)
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_generic_fgretrig_split(w, sp, sr, sb, g.Cost(), 1, 15)
		}
		return slot.ScanReelsCommon(ctx, sp, sr, g, reels, calc)
	}
	return 0, 0
}
