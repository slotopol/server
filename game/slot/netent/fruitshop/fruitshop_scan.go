package fruitshop

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

type StatReg struct {
	slot.StatGeneric
}

func NewStatReg() *StatReg {
	var s StatReg
	s.CntDim(sn, 5)
	return &s
}

func (s *StatReg) FGC() (Cfg, Cfs, Qfs float64) {
	for sym, L := range LineFreespinReg {
		for i, Li := range L {
			var c = float64(s.C[sym][i].Load())
			Cfg += c
			Cfs += c * float64(Li)
			Qfs += c * float64(Li*Li)
		}
	}
	return
}

type StatBon struct {
	slot.StatGeneric
}

func NewStatBon() *StatBon {
	var s StatBon
	s.CntDim(sn, 5)
	return &s
}

func (s *StatBon) FGC() (Cfg, Cfs, Qfs float64) {
	for sym, L := range LineFreespinBon {
		for i, Li := range L {
			var c = float64(s.C[sym][i].Load())
			Cfg += c
			Cfs += c * float64(Li)
			Qfs += c * float64(Li*Li)
		}
	}
	return
}

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("\n(1/2) free games calculations\n")
	var sb = NewStatBon() // slot.NewStatGeneric(sn, 5) also valid
	{
		var g = NewGame()
		g.FSR = 5 // set free spins mode
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_affixrecur(w, sp, sb, g.Cost())
		}
		slot.ScanReelsCommon(ctx, sp, sb, g, calc)
	}

	if ctx.Err() != nil {
		return 0, 0
	}

	fmt.Printf("\n(2/2) regular games calculations\n")
	var sr = NewStatReg() // slot.NewStatGeneric(sn, 5) also valid
	{
		var g = NewGame()
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_dualreelrecur(w, sp, sr, sb, g.Cost())
		}
		return slot.ScanReelsCommon(ctx, sp, sr, g, calc)
	}
}
