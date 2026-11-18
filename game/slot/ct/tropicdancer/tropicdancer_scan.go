package tropicdancer

import (
	"context"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var g = NewGame(sp.Sel)
	var s = slot.NewStatSnap(1, sn, 5, 0, 0).MakeSym(scat, 15)

	var calc = func(w io.Writer) (float64, float64) {
		return slot.Parsheet_solereelrecur(w, sp, s, g.Cost(), 2)
	}

	return slot.ScanReelsCommon2(ctx, sp, s, g, calc)
}
