package treasurehill

import (
	"context"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var g = NewGame(sp.Sel)
	var s = slot.NewStatGeneric(sn, 5)

	var calc = func(w io.Writer) (float64, float64) {
		return slot.Parsheet_fgrecur(w, sp, s, g.Cost(), 1, s.Efs)
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, calc)
}
