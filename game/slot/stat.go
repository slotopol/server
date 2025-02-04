package slot

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	cfg "github.com/slotopol/server/config"
)

type Stater interface {
	SetPlan(n uint64)
	Planned() uint64
	Count() uint64
	LineRTP(sel int) float64
	ScatRTP(sel int) float64
	Update(wins Wins)
}

// Stat is statistics calculation for slot reels.
type Stat struct {
	planned    uint64
	reshuffles uint64
	linepay    float64
	scatpay    float64
	freecount  uint64
	freehits   uint64
	bonuscount [8]uint64
	jackcount  [4]uint64
	lpm, spm   sync.Mutex
}

func (s *Stat) SetPlan(n uint64) {
	atomic.StoreUint64(&s.planned, n)
}

func (s *Stat) Planned() uint64 {
	return atomic.LoadUint64(&s.planned)
}

func (s *Stat) Count() uint64 {
	return atomic.LoadUint64(&s.reshuffles)
}

func (s *Stat) LineRTP(sel int) float64 {
	var reshuf = float64(atomic.LoadUint64(&s.reshuffles))
	s.lpm.Lock()
	var lp = s.linepay
	s.lpm.Unlock()
	return lp / reshuf / float64(sel) * 100
}

func (s *Stat) ScatRTP(sel int) float64 {
	var reshuf = float64(atomic.LoadUint64(&s.reshuffles))
	s.spm.Lock()
	var sp = s.scatpay
	s.spm.Unlock()
	return sp / reshuf / float64(sel) * 100
}

func (s *Stat) FreeCount() uint64 {
	return atomic.LoadUint64(&s.freecount)
}

func (s *Stat) FreeHits() uint64 {
	return atomic.LoadUint64(&s.freehits)
}

func (s *Stat) BonusCount(bid int) uint64 {
	return atomic.LoadUint64(&s.bonuscount[bid])
}

func (s *Stat) JackCount(jid int) uint64 {
	return atomic.LoadUint64(&s.jackcount[jid])
}

func (s *Stat) Update(wins Wins) {
	for _, wi := range wins {
		if wi.Pay != 0 {
			if wi.Line != 0 {
				s.lpm.Lock()
				s.linepay += wi.Pay * wi.Mult
				s.lpm.Unlock()
			} else {
				s.spm.Lock()
				s.scatpay += wi.Pay * wi.Mult
				s.spm.Unlock()
			}
		}
		if wi.Free != 0 {
			atomic.AddUint64(&s.freecount, uint64(wi.Free))
			atomic.AddUint64(&s.freehits, 1)
		}
		if wi.BID != 0 {
			atomic.AddUint64(&s.bonuscount[wi.BID], 1)
		}
		if wi.JID != 0 {
			atomic.AddUint64(&s.jackcount[wi.JID], 1)
		}
	}
	atomic.AddUint64(&s.reshuffles, 1)
}

func Progress(ctx context.Context, s Stater, steps <-chan time.Time, calc func(io.Writer) float64) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-steps:
			var reshuf = float64(s.Count())
			var total = float64(s.Planned())
			var rtp = calc(io.Discard)
			fmt.Printf("processed %.1fm, ready %2.2f%%, RTP = %2.2f%%\n", reshuf/1e6, reshuf/total*100, rtp)
		}
	}
}

func PrintSymPays(s Stater, sel int) func(io.Writer) float64 {
	return func(w io.Writer) float64 {
		var lrtp, srtp = s.LineRTP(sel), s.ScatRTP(sel)
		var rtpsym = lrtp + srtp
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		return rtpsym
	}
}

const ctxgranulation = 100

func BruteForce3x(ctx context.Context, s Stater, g SlotGame, reels *Reels3x) {
	var r1 = reels.Reel(1)
	var r2 = reels.Reel(2)
	var r3 = reels.Reel(3)
	var screen = g.Screen()
	var wins Wins
	for i1 := range r1 {
		screen.SetCol(1, r1, i1)
		for i2 := range r2 {
			screen.SetCol(2, r2, i2)
			for i3 := range r3 {
				screen.SetCol(3, r3, i3)
				g.Scanner(&wins)
				s.Update(wins)
				wins.Reset()
				if s.Count()%ctxgranulation == 0 {
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

func BruteForce3xGo(ctx context.Context, s Stater, g SlotGame, reels *Reels3x) {
	var r1 = reels.Reel(1)
	var r2 = reels.Reel(2)
	var r3 = reels.Reel(3)
	var wg sync.WaitGroup
	wg.Add(len(r1))
	for i1 := range r1 {
		var c = g.Clone()
		go func() {
			defer wg.Done()

			var screen = c.Screen()
			var wins Wins

			screen.SetCol(1, r1, i1)
			for i2 := range r2 {
				screen.SetCol(2, r2, i2)
				for i3 := range r3 {
					screen.SetCol(3, r3, i3)
					c.Scanner(&wins)
					s.Update(wins)
					wins.Reset()
					if s.Count()%ctxgranulation == 0 {
						select {
						case <-ctx.Done():
							return
						default:
						}
					}
				}
			}
		}()
	}
	wg.Wait()
}

func BruteForce4x(ctx context.Context, s Stater, g SlotGame, reels *Reels4x) {
	var r1 = reels.Reel(1)
	var r2 = reels.Reel(2)
	var r3 = reels.Reel(3)
	var r4 = reels.Reel(4)
	var screen = g.Screen()
	var wins Wins
	for i1 := range r1 {
		screen.SetCol(1, r1, i1)
		for i2 := range r2 {
			screen.SetCol(2, r2, i2)
			for i3 := range r3 {
				screen.SetCol(3, r3, i3)
				for i4 := range r4 {
					screen.SetCol(4, r4, i4)
					g.Scanner(&wins)
					s.Update(wins)
					wins.Reset()
					if s.Count()%ctxgranulation == 0 {
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

func BruteForce4xGo(ctx context.Context, s Stater, g SlotGame, reels *Reels4x) {
	var r1 = reels.Reel(1)
	var r2 = reels.Reel(2)
	var r3 = reels.Reel(3)
	var r4 = reels.Reel(4)
	var wg sync.WaitGroup
	wg.Add(len(r1))
	for i1 := range r1 {
		var c = g.Clone()
		go func() {
			defer wg.Done()

			var screen = c.Screen()
			var wins Wins

			screen.SetCol(1, r1, i1)
			for i2 := range r2 {
				screen.SetCol(2, r2, i2)
				for i3 := range r3 {
					screen.SetCol(3, r3, i3)
					for i4 := range r4 {
						screen.SetCol(4, r4, i4)
						c.Scanner(&wins)
						s.Update(wins)
						wins.Reset()
						if s.Count()%ctxgranulation == 0 {
							select {
							case <-ctx.Done():
								return
							default:
							}
						}
					}
				}
			}
		}()
	}
	wg.Wait()
}

func BruteForce5x(ctx context.Context, s Stater, g SlotGame, reels *Reels5x) {
	var r1 = reels.Reel(1)
	var r2 = reels.Reel(2)
	var r3 = reels.Reel(3)
	var r4 = reels.Reel(4)
	var r5 = reels.Reel(5)
	var screen = g.Screen()
	var wins Wins
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
						g.Scanner(&wins)
						s.Update(wins)
						wins.Reset()
						if s.Count()%ctxgranulation == 0 {
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

func BruteForce5xGo(ctx context.Context, s Stater, g SlotGame, reels *Reels5x) {
	var r1 = reels.Reel(1)
	var r2 = reels.Reel(2)
	var r3 = reels.Reel(3)
	var r4 = reels.Reel(4)
	var r5 = reels.Reel(5)
	var wg sync.WaitGroup
	wg.Add(len(r1))
	for i1 := range r1 {
		var c = g.Clone()
		go func() {
			defer wg.Done()

			var screen = c.Screen()
			var wins Wins

			screen.SetCol(1, r1, i1)
			for i2 := range r2 {
				screen.SetCol(2, r2, i2)
				for i3 := range r3 {
					screen.SetCol(3, r3, i3)
					for i4 := range r4 {
						screen.SetCol(4, r4, i4)
						for i5 := range r5 {
							screen.SetCol(5, r5, i5)
							c.Scanner(&wins)
							s.Update(wins)
							wins.Reset()
							if s.Count()%ctxgranulation == 0 {
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
		}()
	}
	wg.Wait()
}

func BruteForce6x(ctx context.Context, s Stater, g SlotGame, reels *Reels6x) {
	var r1 = reels.Reel(1)
	var r2 = reels.Reel(2)
	var r3 = reels.Reel(3)
	var r4 = reels.Reel(4)
	var r5 = reels.Reel(5)
	var r6 = reels.Reel(6)
	var screen = g.Screen()
	var wins Wins
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
						for i6 := range r6 {
							screen.SetCol(6, r6, i6)
							g.Scanner(&wins)
							s.Update(wins)
							wins.Reset()
							if s.Count()%ctxgranulation == 0 {
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
}

func BruteForce6xGo(ctx context.Context, s Stater, g SlotGame, reels *Reels6x) {
	var r1 = reels.Reel(1)
	var r2 = reels.Reel(2)
	var r3 = reels.Reel(3)
	var r4 = reels.Reel(4)
	var r5 = reels.Reel(5)
	var r6 = reels.Reel(6)
	var wg sync.WaitGroup
	wg.Add(len(r1))
	for i1 := range r1 {
		var c = g.Clone()
		go func() {
			defer wg.Done()

			var screen = c.Screen()
			var wins Wins

			screen.SetCol(1, r1, i1)
			for i2 := range r2 {
				screen.SetCol(2, r2, i2)
				for i3 := range r3 {
					screen.SetCol(3, r3, i3)
					for i4 := range r4 {
						screen.SetCol(4, r4, i4)
						for i5 := range r5 {
							screen.SetCol(5, r5, i5)
							for i6 := range r6 {
								screen.SetCol(6, r6, i6)
								c.Scanner(&wins)
								s.Update(wins)
								wins.Reset()
								if s.Count()%ctxgranulation == 0 {
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
		}()
	}
	wg.Wait()
}

func MonteCarlo(ctx context.Context, s Stater, g SlotGame, reels Reels) {
	var n = s.Planned()
	var screen = g.Screen()
	var wins Wins
	for range n {
		screen.Spin(reels)
		g.Scanner(&wins)
		s.Update(wins)
		wins.Reset()
		if s.Count()%ctxgranulation == 0 {
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}
}

func MonteCarloGo(ctx context.Context, s Stater, g SlotGame, reels Reels) {
	var n = s.Planned()
	var ncpu = runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(ncpu)
	for range ncpu {
		var c = g.Clone()
		go func() {
			defer wg.Done()

			var screen = c.Screen()
			var wins Wins

			for range n / uint64(ncpu) {
				screen.Spin(reels)
				c.Scanner(&wins)
				s.Update(wins)
				wins.Reset()
				if s.Count()%ctxgranulation == 0 {
					select {
					case <-ctx.Done():
						return
					default:
					}
				}
			}
		}()
	}
	wg.Wait()
}

func ScanReels3x(ctx context.Context, s Stater, g SlotGame, reels *Reels3x,
	calc func(io.Writer) float64,
	bftick, mctick <-chan time.Time) float64 {
	var t0 = time.Now()
	var ctx2, cancel2 = context.WithCancel(ctx)
	defer cancel2()
	if cfg.MCCount > 0 {
		s.SetPlan(cfg.MCCount * 1e6)
		go Progress(ctx2, s, mctick, calc)
		if cfg.MTScan && !cfg.DevMode {
			MonteCarloGo(ctx2, s, g, reels)
		} else {
			MonteCarlo(ctx2, s, g, reels)
		}
	} else {
		s.SetPlan(reels.Reshuffles())
		go Progress(ctx2, s, bftick, calc)
		if cfg.MTScan && !cfg.DevMode {
			BruteForce3xGo(ctx2, s, g, reels)
		} else {
			BruteForce3x(ctx2, s, g, reels)
		}
	}
	var dur = time.Since(t0)
	var comp = float64(s.Count()) / float64(s.Planned()) * 100
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", comp, g.GetSel(), dur)
	return calc(os.Stdout)
}

func ScanReels4x(ctx context.Context, s Stater, g SlotGame, reels *Reels4x,
	calc func(io.Writer) float64,
	bftick, mctick <-chan time.Time) float64 {
	var t0 = time.Now()
	var ctx2, cancel2 = context.WithCancel(ctx)
	defer cancel2()
	if cfg.MCCount > 0 {
		s.SetPlan(cfg.MCCount * 1e6)
		go Progress(ctx2, s, mctick, calc)
		if cfg.MTScan && !cfg.DevMode {
			MonteCarloGo(ctx2, s, g, reels)
		} else {
			MonteCarlo(ctx2, s, g, reels)
		}
	} else {
		s.SetPlan(reels.Reshuffles())
		go Progress(ctx2, s, bftick, calc)
		if cfg.MTScan && !cfg.DevMode {
			BruteForce4xGo(ctx2, s, g, reels)
		} else {
			BruteForce4x(ctx2, s, g, reels)
		}
	}
	var dur = time.Since(t0)
	var comp = float64(s.Count()) / float64(s.Planned()) * 100
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", comp, g.GetSel(), dur)
	return calc(os.Stdout)
}

func ScanReels5x(ctx context.Context, s Stater, g SlotGame, reels *Reels5x,
	calc func(io.Writer) float64,
	bftick, mctick <-chan time.Time) float64 {
	var t0 = time.Now()
	var ctx2, cancel2 = context.WithCancel(ctx)
	defer cancel2()
	if cfg.MCCount > 0 {
		s.SetPlan(cfg.MCCount * 1e6)
		go Progress(ctx2, s, mctick, calc)
		if cfg.MTScan && !cfg.DevMode {
			MonteCarloGo(ctx2, s, g, reels)
		} else {
			MonteCarlo(ctx2, s, g, reels)
		}
	} else {
		s.SetPlan(reels.Reshuffles())
		go Progress(ctx2, s, bftick, calc)
		if cfg.MTScan && !cfg.DevMode {
			BruteForce5xGo(ctx2, s, g, reels)
		} else {
			BruteForce5x(ctx2, s, g, reels)
		}
	}
	var dur = time.Since(t0)
	var comp = float64(s.Count()) / float64(s.Planned()) * 100
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", comp, g.GetSel(), dur)
	return calc(os.Stdout)
}
