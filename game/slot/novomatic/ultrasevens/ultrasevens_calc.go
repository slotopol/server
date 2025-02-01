package ultrasevens

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
		var reshuf = float64(s.Count())
		var lrtp, srtp = s.LineRTP(g.Sel), s.ScatRTP(g.Sel)
		if srtp > 0 {
			panic("scatters have no pays")
		}
		var rtpsym = lrtp + srtp
		fmt.Fprintf(w, "reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
			len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
		fmt.Fprintf(w, "jackpots1: count %d, frequency 1/%.12g\n", s.JackCount(ssj1), reshuf/float64(s.JackCount(ssj1)))
		fmt.Fprintf(w, "jackpots2: count %d, frequency 1/%.12g\n", s.JackCount(ssj2), reshuf/float64(s.JackCount(ssj2)))
		fmt.Fprintf(w, "jackpots3: count %d, frequency 1/%.12g\n", s.JackCount(ssj3), reshuf/float64(s.JackCount(ssj3)))
		fmt.Fprintf(w, "RTP = %.6f%%\n", rtpsym)
		return rtpsym
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc,
		time.Tick(5*time.Second), time.Tick(2*time.Second))
}
