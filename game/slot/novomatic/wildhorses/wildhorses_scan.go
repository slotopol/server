package wildhorses

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

type StatBon struct {
	slot.StatGeneric
}

func NewStatBon() *StatBon {
	var s StatBon
	s.CntDim(sn, 5)
	return &s
}

func (s *StatBon) Efs() (HRfg, EK, EK2 float64) {
	const p = 0.5 // probability of getting or not getting free spins is equal
	const L = 10  // number of free spins
	var N = s.Count()
	var c = p * N / L
	HRfg = c / N
	EK = c * L / N
	EK2 = c * L * L / N
	return
}

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("\n(1/2) bonus reels calculations\n")
	var sb = NewStatBon()
	{
		var g = NewGame(sp.Sel)
		g.FSR = 10 // set free spins mode
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_solereelrecur(w, sp, sb, g.Cost(), 1)
		}
		slot.ScanReelsCommon(ctx, sp, sb, g, calc)
	}

	if ctx.Err() != nil {
		return 0, 0
	}

	fmt.Printf("\n(2/2) regular reels calculations\n")
	var sr = slot.NewStatGeneric(sn, 5)
	{
		var g = NewGame(sp.Sel)
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_dualreelrecur(w, sp, sr, sb, g.Cost())
		}
		return slot.ScanReelsCommon(ctx, sp, sr, g, calc)
	}
}
