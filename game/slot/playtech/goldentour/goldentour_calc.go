package goldentour

import (
	"context"
	"fmt"
	"io"
	"time"

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
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	var g = NewGame()
	g.Sel = 1
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var reshuf = float64(s.Count())
		var lrtp, srtp = s.LineRTP(g.Sel), s.ScatRTP(g.Sel)
		var rtpsym = lrtp + srtp
		var qgolfbn = float64(s.BonusCount(golfbon)) / reshuf / float64(g.Sel)
		var rtpgolfbn = Egolfbn * qgolfbn * 100
		var rtp = rtpsym + rtpgolfbn
		fmt.Fprintf(w, "reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
			len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
		fmt.Fprintf(w, "golf bonuses: count %d, rtp = %.6f%%\n", s.BonusCount(golfbon), rtpgolfbn)
		fmt.Fprintf(w, "golf bonuses frequency: 1/%.5g\n", reshuf/float64(s.BonusCount(golfbon)))
		fmt.Fprintf(w, "RTP = %.5g(lined) + %.5g(scatter) + %.5g(golf) = %.6f%%\n", lrtp, srtp, rtpgolfbn, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc,
		time.Tick(2*time.Second), time.Tick(2*time.Second))
}
