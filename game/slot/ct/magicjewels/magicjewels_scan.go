package magicjewels

import (
	"context"
	"io"

	"github.com/slotopol/server/game/slot"
)

func Efs(s *slot.StatCascade) (HRfg, EK, EK2 float64) {
	var N = s.Count()
	var Cfg5 = float64(s.CascNum(4))
	var Cfg6 = float64(s.CascNum(5))
	var Cfg7 float64
	const L5 = 12 // number of free spins for 5 cascades
	const L6 = 20 // number of free spins for 6 cascades
	const L7 = 38 // number of free spins for 7 or more cascades
	for i := 6; i < slot.FallLimit; i++ {
		Cfg7 += float64(s.CascNum(i))
	}
	HRfg = (Cfg5 + Cfg6 + Cfg7) / N
	EK = (Cfg5*L5 + Cfg6*L6 + Cfg7*L7) / N
	EK2 = (Cfg5*L5*L5 + Cfg6*L6*L6 + Cfg7*L7*L7) / N
	return
}

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var g = NewGame()
	var s = slot.NewStatCascade(sn, 5)

	var calc = func(w io.Writer) (float64, float64) {
		return slot.Parsheet_fgrecur(w, sp, s, g.Cost(), 2, func() (float64, float64, float64) {
			return Efs(s)
		})
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, calc)
}
