package royaldynasty

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

type Stat struct {
	slot.StatGeneric
}

func NewStat() *Stat {
	var s Stat
	s.CntDim(sn, 5)
	return &s
}

func (s *Stat) Efs() (HRfg, EK, EK2 float64) {
	var N = s.Count()
	var sum, sum2 float64
	for _, Li := range Freegames {
		sum += float64(Li)
		sum2 += float64(Li * Li)
	}
	HRfg = float64(s.Cfg.Load()) / N
	EK = HRfg * sum / float64(len(Freegames))
	EK2 = HRfg * sum2 / float64(len(Freegames))
	return
}

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("\n(1/2) bonus reels calculations\n")
	var sb = NewStat()
	{
		var g = NewGame(sp.Sel)
		g.FSR = 35 // set free spins mode
		g.TS = scat1
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_fgrecur(w, sp, sb, g.Cost(), 1, sb.Efs)
		}
		slot.ScanReelsCommon(ctx, sp, sb, g, calc)
	}

	if ctx.Err() != nil {
		return 0, 0
	}

	fmt.Printf("\n(2/2) regular reels calculations\n")
	var sr = NewStat()
	{
		var g = NewGame(sp.Sel)
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_dualreelset(w, sp, sr, sb, g.Cost(), sr.Efs, sb.Efs)
		}
		return slot.ScanReelsCommon(ctx, sp, sr, g, calc)
	}
}
