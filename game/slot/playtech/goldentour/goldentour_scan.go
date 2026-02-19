package goldentour

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

var Egolfbn float64

func ExpGolf() {
	var sum float64
	for _, v := range Golf {
		sum += float64(v)
	}
	Egolfbn = sum / float64(len(Golf))
}

func CalcStat(ctx context.Context, sp *slot.ScanPar) float64 {
	fmt.Printf("*bonus games calculations*\n")
	ExpGolf()
	fmt.Printf("len = %d, E = %g\n", len(Golf), Egolfbn)
	fmt.Printf("*reels calculations*\n")
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame(sp.Sel)
	var s = slot.NewStatGeneric(sn, 5)

	var calc = func(w io.Writer) float64 {
		var N, S, _ = s.NSQ(g.Cost())
		var µ = S / N
		var qgolfbn = s.BonusHitsF(golfbon) / N / float64(g.Sel)
		var rtpgolfbn = Egolfbn * qgolfbn
		var rtp = µ + rtpgolfbn
		fmt.Fprintf(w, "golf bonuses: hit rate 1/%.5g, rtp = %.6f%%\n", N/s.BonusHitsF(golfbon), rtpgolfbn*100)
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(golf) = %.6f%%\n", µ*100, rtpgolfbn*100, rtp*100)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
}
