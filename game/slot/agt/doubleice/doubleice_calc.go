package doubleice

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
	var sln float64 = 1
	g.Sel = int(sln)
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var reshuf = float64(s.Reshuffles)
		var lrtp = s.LinePay / reshuf / sln * 100
		fmt.Printf("reels lengths [%d, %d, %d], total reshuffles %d\n",
			len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), reels.Reshuffles())
		fmt.Printf("RTP = %.6f%%\n", lrtp)
		return lrtp
	}

	return slot.ScanReels3x(ctx, &s, g, reels, calc,
		time.Tick(2*time.Second), time.Tick(2*time.Second))
}
