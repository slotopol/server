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

func CalcStat(ctx context.Context, mrtp float64) float64 {
	fmt.Printf("*bonus games calculations*\n")
	ExpGolf()
	fmt.Printf("len = %d, E = %g\n", len(Golf), Egolfbn)
	fmt.Printf("*reels calculations*\n")
	var reels, _ = ReelsMap.FindClosest(mrtp)
	var g = NewGame(1)
	var s slot.StatGeneric

	var calc = func(w io.Writer) float64 {
		var N = s.Count()
		var lrtp, srtp = s.SymRTP(g.Cost())
		var rtpsym = lrtp + srtp
		var qgolfbn = s.BonusHitsF(golfbon) / N / float64(g.Sel)
		var rtpgolfbn = Egolfbn * qgolfbn * 100
		var rtp = rtpsym + rtpgolfbn
		fmt.Fprintf(w, "golf bonuses: hit rate 1/%.5g, rtp = %.6f%%\n", N/s.BonusHitsF(golfbon), rtpgolfbn)
		fmt.Fprintf(w, "RTP = %.5g(lined) + %.5g(scatter) + %.5g(golf) = %.6f%%\n", lrtp, srtp, rtpgolfbn, rtp)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}
