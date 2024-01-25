package game

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

type Stater interface {
	Count() uint64
	Update(sw *WinScan)
}

// Stat is statistics calculation for slot reels.
type Stat struct {
	Reshuffles uint64
	LinePay    uint64
	ScatPay    uint64
	FreeCount  uint64
	BonusCount [8]uint64
	JackCount  [4]uint64
}

func (s *Stat) Count() uint64 {
	return atomic.LoadUint64(&s.Reshuffles)
}

func (s *Stat) Update(sw *WinScan) {
	for _, wi := range sw.Wins {
		if wi.Pay > 0 {
			if wi.Line > 0 {
				atomic.AddUint64(&s.LinePay, uint64(wi.Pay*wi.Mult))
			} else {
				atomic.AddUint64(&s.ScatPay, uint64(wi.Pay*wi.Mult))
			}
		}
		if wi.Free > 0 {
			atomic.AddUint64(&s.FreeCount, uint64(wi.Free))
		}
		if wi.BID > 0 {
			atomic.AddUint64(&s.BonusCount[wi.BID], 1)
		}
		if wi.Jack > 0 {
			atomic.AddUint64(&s.JackCount[wi.Jack], 1)
		}
	}
	atomic.AddUint64(&s.Reshuffles, 1)
}

func (s *Stat) Progress(ctx context.Context, steps *time.Ticker, sel, total float64) {
	for {
		select {
		case <-ctx.Done():
			steps.Stop()
			return
		case <-steps.C:
			var n = float64(atomic.LoadUint64(&s.Reshuffles))
			var pays = (float64(atomic.LoadUint64(&s.LinePay))/sel + float64(atomic.LoadUint64(&s.ScatPay))) / n * 100
			fmt.Printf("processed %.1fm, ready %2.2f%%, symbols pays %2.2f%%\n", n/1e6, n/total*100, pays)
		}
	}
}

func BruteForce5x(ctx context.Context, s Stater, g SlotGame, reels Reels) {
	var screen = g.NewScreen()
	defer screen.Free()
	var ws WinScan
	var r1 = reels.Reel(1)
	var r2 = reels.Reel(2)
	var r3 = reels.Reel(3)
	var r4 = reels.Reel(4)
	var r5 = reels.Reel(5)
	for i1 := range r1 {
		screen.SetCol(1, r1, i1)
		for i2 := range r2 {
			screen.SetCol(2, r2, i2)
			for i3 := range r3 {
				screen.SetCol(3, r3, i3)
				for i4 := range r4 {
					screen.SetCol(4, r4, i4)
					for i5 := range r5 {
						screen.SetCol(5, r5, i5)
						g.Scanner(screen, &ws)
						s.Update(&ws)
						ws.Reset()
						if s.Count()&100 == 0 {
							select {
							case <-ctx.Done():
								return
							default:
							}
						}
					}
				}
			}
		}
	}
}

func MonteCarlo(ctx context.Context, s Stater, g SlotGame, n int) {
	var screen = g.NewScreen()
	defer screen.Free()
	var ws WinScan
	for i := 0; i < n; i++ {
		g.Spin(screen)
		g.Scanner(screen, &ws)
		s.Update(&ws)
		ws.Reset()
		if s.Count()&100 == 0 {
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}
}
