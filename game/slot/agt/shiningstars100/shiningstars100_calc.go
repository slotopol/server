package shiningstars100

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, mrtp float64) float64 {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	var g = NewGame()
	g.Sel = 1
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var lrtp, srtp = s.LineRTP(g.Sel), s.ScatRTP(g.Sel)
		var rtpsym = lrtp + srtp
		fmt.Fprintf(w, "reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
			len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
		fmt.Fprintf(w, "RTP = %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		return rtpsym
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc,
		time.Tick(6*time.Second), time.Tick(2*time.Second))
}
