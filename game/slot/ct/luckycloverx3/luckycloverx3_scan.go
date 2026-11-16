package luckycloverx3

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

// custom parsheet
func Parsheet(w io.Writer, sp *slot.ScanPar, s *slot.StatGeneric, cost float64) (float64, float64) {
	var µ, D = slot.EvD(s, cost)
	var HRx2 = s.Count() / s.BonusHits(bonx2)
	var HRx3 = s.Count() / s.BonusHits(bonx3)
	var HRxx = s.Count() / (s.BonusHits(bonx2) + s.BonusHits(bonx3))
	if sp.IsMain() {
		fmt.Fprintf(w, "bonus mult: HRx2 = 1/%.5g, HRx3 = 1/%.5g, HRxx = 1/%.5g\n", HRx2, HRx3, HRxx)
		fmt.Fprintf(w, "RTP = %.8g%%\n", µ*100)
	}
	slot.Print_all(w, sp, s, µ, D)
	return µ, D
}

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var g = NewGame(sp.Sel)
	var s = slot.NewStatGeneric(sn, 5)
	s.BonDim(bonx3)

	var calc = func(w io.Writer) (float64, float64) {
		return Parsheet(w, sp, s, g.Cost())
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, calc)
}
