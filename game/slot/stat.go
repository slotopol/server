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
	Count(cfn int) uint64
	LineRTP(cost float64) float64
	ScatRTP(cost float64) float64
	Update(wins Wins, cfn int)
}

// Stat is statistics calculation for slot reels.
type Stat struct {
	planned    uint64
	reshuffles [10]uint64
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

func (s *Stat) Count(cfn int) uint64 {
	return atomic.LoadUint64(&s.reshuffles[cfn-1])
}

func (s *Stat) LineRTP(cost float64) float64 {
	var reshuf = float64(atomic.LoadUint64(&s.reshuffles[0]))
	s.lpm.Lock()
	var lp = s.linepay
	s.lpm.Unlock()
	return lp / reshuf / cost * 100
}

func (s *Stat) ScatRTP(cost float64) float64 {
	var reshuf = float64(atomic.LoadUint64(&s.reshuffles[0]))
	s.spm.Lock()
	var sp = s.scatpay
	s.spm.Unlock()
	return sp / reshuf / cost * 100
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

func (s *Stat) Update(wins Wins, cfn int) {
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
	if cfn < len(s.reshuffles) {
		atomic.AddUint64(&s.reshuffles[cfn-1], 1)
	}
}

func Progress(ctx context.Context, s Stater, calc func(io.Writer) float64) {
	const stepdur = 1000 * time.Millisecond
	var t0 = time.Now()
	var steps = time.Tick(stepdur)
	fmt.Printf("calculation started...\r")
	for {
		select {
		case <-ctx.Done():
			return
		case <-steps:
			var reshuf = float64(s.Count(1))
			var total = float64(s.Planned())
			var rtp = calc(io.Discard)
			var dur = time.Since(t0)
			var exp = time.Duration(float64(dur) * total / reshuf)
			fmt.Printf("processed %.1fm, ready %2.2f%% (%v / %v), RTP = %2.2f%%  \r",
				reshuf/1e6, reshuf/total*100,
				dur.Truncate(stepdur), exp.Truncate(stepdur),
				rtp)
		}
	}
}

func PrintSymPays(s Stater, cost float64) func(io.Writer) float64 {
	return func(w io.Writer) float64 {
		var lrtp, srtp = s.LineRTP(cost), s.ScatRTP(cost)
		var rtpsym = lrtp + srtp
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		return rtpsym
	}
}

const CtxGranulation = 100

func CorrectThrNum() int {
	if cfg.DevMode {
		return 1
	} else if cfg.MTCount < 1 {
		return runtime.GOMAXPROCS(0)
	}
	return cfg.MTCount
}

func BruteForce3x(ctx context.Context, s Stater, g SlotGame, reels *Reels3x) {
	var tn = CorrectThrNum()
	var tn64 = uint64(tn)
	var r1 = reels.Reel(1)
	var r2 = reels.Reel(2)
	var r3 = reels.Reel(3)
	var wg sync.WaitGroup
	wg.Add(tn)
	for ti := range tn64 {
		var sg = g.Clone().(ClassicSlot)     // classic slot game
		var cs, iscascade = sg.(CascadeSlot) // cascade slot game
		var reshuf uint64
		go func() {
			defer wg.Done()

			var wins Wins
			for i1 := range r1 {
				sg.SetCol(1, r1, i1)
				for i2 := range r2 {
					sg.SetCol(2, r2, i2)
					for i3 := range r3 {
						reshuf++
						if reshuf%CtxGranulation == 0 {
							select {
							case <-ctx.Done():
								return
							default:
							}
						}
						if reshuf%tn64 != ti {
							continue
						}
						sg.SetCol(3, r3, i3)
						if iscascade {
							var cfn int
							for {
								cs.NewFall()
								cfn++
								cs.Scanner(&wins)
								s.Update(wins, cfn)
								cs.Strike(wins)
								if len(wins) == 0 {
									break
								}
								cs.NextFall(reels)
								wins.Reset()
							}
							if cfn > 1 {
								cs.SetCol(1, r1, i1)
								cs.SetCol(2, r2, i2)
							}
						} else {
							sg.Scanner(&wins)
							s.Update(wins, 1)
							wins.Reset()
						}
					}
				}
			}
		}()
	}
	wg.Wait()
}

func BruteForce4x(ctx context.Context, s Stater, g SlotGame, reels *Reels4x) {
	var tn = CorrectThrNum()
	var tn64 = uint64(tn)
	var r1 = reels.Reel(1)
	var r2 = reels.Reel(2)
	var r3 = reels.Reel(3)
	var r4 = reels.Reel(4)
	var wg sync.WaitGroup
	wg.Add(tn)
	for ti := range tn64 {
		var sg = g.Clone().(ClassicSlot)     // classic slot game
		var cs, iscascade = sg.(CascadeSlot) // cascade slot game
		var reshuf uint64
		go func() {
			defer wg.Done()

			var wins Wins
			for i1 := range r1 {
				sg.SetCol(1, r1, i1)
				for i2 := range r2 {
					sg.SetCol(2, r2, i2)
					for i3 := range r3 {
						sg.SetCol(3, r3, i3)
						for i4 := range r4 {
							reshuf++
							if reshuf%CtxGranulation == 0 {
								select {
								case <-ctx.Done():
									return
								default:
								}
							}
							if reshuf%tn64 != ti {
								continue
							}
							sg.SetCol(4, r4, i4)
							if iscascade {
								var cfn int
								for {
									cs.NewFall()
									cfn++
									cs.Scanner(&wins)
									s.Update(wins, cfn)
									cs.Strike(wins)
									if len(wins) == 0 {
										break
									}
									cs.NextFall(reels)
									wins.Reset()
								}
								if cfn > 1 {
									cs.SetCol(1, r1, i1)
									cs.SetCol(2, r2, i2)
									cs.SetCol(3, r3, i3)
								}
							} else {
								sg.Scanner(&wins)
								s.Update(wins, 1)
								wins.Reset()
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
	var tn = CorrectThrNum()
	var tn64 = uint64(tn)
	var r1 = reels.Reel(1)
	var r2 = reels.Reel(2)
	var r3 = reels.Reel(3)
	var r4 = reels.Reel(4)
	var r5 = reels.Reel(5)
	var wg sync.WaitGroup
	wg.Add(tn)
	for ti := range tn64 {
		var sg = g.Clone().(ClassicSlot)     // classic slot game
		var cs, iscascade = sg.(CascadeSlot) // cascade slot game
		var reshuf uint64
		go func() {
			defer wg.Done()

			var wins Wins
			for i1 := range r1 {
				sg.SetCol(1, r1, i1)
				for i2 := range r2 {
					sg.SetCol(2, r2, i2)
					for i3 := range r3 {
						sg.SetCol(3, r3, i3)
						for i4 := range r4 {
							sg.SetCol(4, r4, i4)
							for i5 := range r5 {
								reshuf++
								if reshuf%CtxGranulation == 0 {
									select {
									case <-ctx.Done():
										return
									default:
									}
								}
								if reshuf%tn64 != ti {
									continue
								}
								sg.SetCol(5, r5, i5)
								if iscascade {
									var cfn int
									for {
										cs.NewFall()
										cfn++
										cs.Scanner(&wins)
										s.Update(wins, cfn)
										cs.Strike(wins)
										if len(wins) == 0 {
											break
										}
										cs.NextFall(reels)
										wins.Reset()
									}
									if cfn > 1 {
										cs.SetCol(1, r1, i1)
										cs.SetCol(2, r2, i2)
										cs.SetCol(3, r3, i3)
										cs.SetCol(4, r4, i4)
									}
								} else {
									sg.Scanner(&wins)
									s.Update(wins, 1)
									wins.Reset()
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

func BruteForce5x3Big(ctx context.Context, s Stater, g SlotGame, r1, rb, r5 []Sym) {
	var tn = CorrectThrNum()
	var tn64 = uint64(tn)
	var wg sync.WaitGroup
	wg.Add(tn)
	for ti := range tn64 {
		var sg = g.Clone().(ClassicSlot)
		var cb = sg.(Bigger)
		var reshuf uint64
		go func() {
			defer wg.Done()

			var wins Wins
			for i1 := range r1 {
				sg.SetCol(1, r1, i1)
				for _, big := range rb {
					cb.SetBig(big)
					for i5 := range r5 {
						reshuf++
						if reshuf%CtxGranulation == 0 {
							select {
							case <-ctx.Done():
								return
							default:
							}
						}
						if reshuf%tn64 != ti {
							continue
						}
						sg.SetCol(5, r5, i5)
						sg.Scanner(&wins)
						s.Update(wins, 1)
						wins.Reset()
					}
				}
			}
		}()
	}
	wg.Wait()
}

func BruteForce6x(ctx context.Context, s Stater, g SlotGame, reels *Reels6x) {
	var tn = CorrectThrNum()
	var tn64 = uint64(tn)
	var r1 = reels.Reel(1)
	var r2 = reels.Reel(2)
	var r3 = reels.Reel(3)
	var r4 = reels.Reel(4)
	var r5 = reels.Reel(5)
	var r6 = reels.Reel(6)
	var wg sync.WaitGroup
	wg.Add(tn)
	for ti := range tn64 {
		var sg = g.Clone().(ClassicSlot)     // classic slot game
		var cs, iscascade = sg.(CascadeSlot) // cascade slot game
		var reshuf uint64
		go func() {
			defer wg.Done()

			var wins Wins
			for i1 := range r1 {
				sg.SetCol(1, r1, i1)
				for i2 := range r2 {
					sg.SetCol(2, r2, i2)
					for i3 := range r3 {
						sg.SetCol(3, r3, i3)
						for i4 := range r4 {
							sg.SetCol(4, r4, i4)
							for i5 := range r5 {
								sg.SetCol(5, r5, i5)
								for i6 := range r6 {
									reshuf++
									if reshuf%CtxGranulation == 0 {
										select {
										case <-ctx.Done():
											return
										default:
										}
									}
									if reshuf%tn64 != ti {
										continue
									}
									sg.SetCol(6, r6, i6)
									if iscascade {
										var cfn int
										for {
											cs.NewFall()
											cfn++
											cs.Scanner(&wins)
											s.Update(wins, cfn)
											cs.Strike(wins)
											if len(wins) == 0 {
												break
											}
											cs.NextFall(reels)
											wins.Reset()
										}
										if cfn > 1 {
											cs.SetCol(1, r1, i1)
											cs.SetCol(2, r2, i2)
											cs.SetCol(3, r3, i3)
											cs.SetCol(4, r4, i4)
											cs.SetCol(5, r5, i5)
										}
									} else {
										sg.Scanner(&wins)
										s.Update(wins, 1)
										wins.Reset()
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
	var tn = CorrectThrNum()
	var tn64 = uint64(tn)
	var n = s.Planned()
	var wg sync.WaitGroup
	wg.Add(tn)
	for range tn64 {
		var sg = g.Clone().(ClassicSlot)     // classic slot game
		var cs, iscascade = sg.(CascadeSlot) // cascade slot game
		var reshuf uint64
		go func() {
			defer wg.Done()

			var wins Wins
			for range n / tn64 {
				reshuf++
				if reshuf%CtxGranulation == 0 {
					select {
					case <-ctx.Done():
						return
					default:
					}
				}
				if iscascade {
					var cfn int
					for {
						cs.NewFall()
						cfn++
						cs.ReelSpin(reels)
						cs.Scanner(&wins)
						s.Update(wins, cfn)
						cs.Strike(wins)
						if len(wins) == 0 {
							break
						}
						wins.Reset()
					}
				} else {
					sg.ReelSpin(reels)
					sg.Scanner(&wins)
					s.Update(wins, 1)
					wins.Reset()
				}
			}
		}()
	}
	wg.Wait()
}

func ScanReels3x(ctx context.Context, s Stater, g SlotGame, reels *Reels3x,
	calc func(io.Writer) float64) float64 {
	var t0 = time.Now()
	var ctx2, cancel2 = context.WithCancel(ctx)
	defer cancel2()
	if cfg.MCCount > 0 {
		s.SetPlan(cfg.MCCount * 1e6)
		go Progress(ctx2, s, calc)
		MonteCarlo(ctx2, s, g, reels)
	} else {
		s.SetPlan(reels.Reshuffles())
		go Progress(ctx2, s, calc)
		BruteForce3x(ctx2, s, g, reels)
	}
	var dur = time.Since(t0)
	var comp = float64(s.Count(1)) / float64(s.Planned()) * 100
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v            \n", comp, g.GetSel(), dur)
	fmt.Printf("reels lengths [%d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), reels.Reshuffles())
	return calc(os.Stdout)
}

func ScanReels4x(ctx context.Context, s Stater, g SlotGame, reels *Reels4x,
	calc func(io.Writer) float64) float64 {
	var t0 = time.Now()
	var ctx2, cancel2 = context.WithCancel(ctx)
	defer cancel2()
	if cfg.MCCount > 0 {
		s.SetPlan(cfg.MCCount * 1e6)
		go Progress(ctx2, s, calc)
		MonteCarlo(ctx2, s, g, reels)
	} else {
		s.SetPlan(reels.Reshuffles())
		go Progress(ctx2, s, calc)
		BruteForce4x(ctx2, s, g, reels)
	}
	var dur = time.Since(t0)
	var comp = float64(s.Count(1)) / float64(s.Planned()) * 100
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v            \n", comp, g.GetSel(), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), reels.Reshuffles())
	return calc(os.Stdout)
}

func ScanReels5x(ctx context.Context, s Stater, g SlotGame, reels *Reels5x,
	calc func(io.Writer) float64) float64 {
	var t0 = time.Now()
	var ctx2, cancel2 = context.WithCancel(ctx)
	defer cancel2()
	if cfg.MCCount > 0 {
		s.SetPlan(cfg.MCCount * 1e6)
		go Progress(ctx2, s, calc)
		MonteCarlo(ctx2, s, g, reels)
	} else {
		s.SetPlan(reels.Reshuffles())
		go Progress(ctx2, s, calc)
		BruteForce5x(ctx2, s, g, reels)
	}
	var dur = time.Since(t0)
	var comp = float64(s.Count(1)) / float64(s.Planned()) * 100
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v            \n", comp, g.GetSel(), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	return calc(os.Stdout)
}

func ScanReels6x(ctx context.Context, s Stater, g SlotGame, reels *Reels6x,
	calc func(io.Writer) float64) float64 {
	var t0 = time.Now()
	var ctx2, cancel2 = context.WithCancel(ctx)
	defer cancel2()
	if cfg.MCCount > 0 {
		s.SetPlan(cfg.MCCount * 1e6)
		go Progress(ctx2, s, calc)
		MonteCarlo(ctx2, s, g, reels)
	} else {
		s.SetPlan(reels.Reshuffles())
		go Progress(ctx2, s, calc)
		BruteForce6x(ctx2, s, g, reels)
	}
	var dur = time.Since(t0)
	var comp = float64(s.Count(1)) / float64(s.Planned()) * 100
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v            \n", comp, g.GetSel(), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), len(reels.Reel(6)), reels.Reshuffles())
	return calc(os.Stdout)
}
