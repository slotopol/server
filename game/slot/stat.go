package slot

import (
	"context"
	"errors"
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
	Planned() float64
	Reshuf(cfn int) float64
	Errors() float64
	IncErr()
	Update(wins Wins, cfn int)
}

// Stat is statistics calculation for slot reels.
type Stat struct {
	planned    uint64
	reshuffles [FallLimit]uint64
	errcount   uint64
	linepay    float64
	scatpay    float64
	freecount  uint64
	freehits   uint64
	bonuscount [8]uint64
	jackcount  [4]uint64
	lpm, spm   sync.Mutex
}

// Declare conformity with Stater interface.
var _ Stater = (*Stat)(nil)

func (s *Stat) SetPlan(n uint64) {
	atomic.StoreUint64(&s.planned, n)
}

func (s *Stat) Planned() float64 {
	return float64(atomic.LoadUint64(&s.planned))
}

func (s *Stat) Count() float64 {
	var n uint64
	for i := range FallLimit {
		n += atomic.LoadUint64(&s.reshuffles[i])
	}
	return float64(n)
}

func (s *Stat) Reshuf(cfn int) float64 {
	var n uint64
	for i := cfn - 1; i < FallLimit; i++ {
		n += atomic.LoadUint64(&s.reshuffles[i])
	}
	return float64(n)
}

func (s *Stat) Errors() float64 {
	return float64(atomic.LoadUint64(&s.errcount))
}

func (s *Stat) IncErr() {
	atomic.AddUint64(&s.errcount, 1)
}

func (s *Stat) LineRTP(cost float64) float64 {
	s.lpm.Lock()
	var lp = s.linepay
	s.lpm.Unlock()
	return lp / s.Count() / cost * 100
}

func (s *Stat) ScatRTP(cost float64) float64 {
	s.spm.Lock()
	var sp = s.scatpay
	s.spm.Unlock()
	return sp / s.Count() / cost * 100
}

func (s *Stat) SymRTP(cost float64) (lrtp, srtp float64) {
	s.lpm.Lock()
	var lp = s.linepay
	s.lpm.Unlock()
	s.spm.Lock()
	var sp = s.scatpay
	s.spm.Unlock()
	var reshuf = s.Count()
	lrtp = lp / reshuf / cost * 100
	srtp = sp / reshuf / cost * 100
	return
}

func (s *Stat) FreeCountU() uint64 {
	return atomic.LoadUint64(&s.freecount)
}

func (s *Stat) FreeCount() float64 {
	return float64(atomic.LoadUint64(&s.freecount))
}

// Returns (q, sq), where q = free spins quantifier, sq = 1/(1-q)
// sum of a decreasing geometric progression for retriggered free spins.
func (s *Stat) FSQ() (q float64, sq float64) {
	q = s.FreeCount() / s.Count()
	sq = 1 / (1 - q)
	return
}

func (s *Stat) FreeHits() float64 {
	return float64(atomic.LoadUint64(&s.freehits))
}

// Quantifier of free games per reshuffles.
func (s *Stat) FGQ() float64 {
	return s.FreeHits() / s.Count()
}

// Free Games Frequency: average number of reshuffles per free games hit.
func (s *Stat) FGF() float64 {
	return s.Count() / s.FreeHits()
}

func (s *Stat) BonusCount(bid int) float64 {
	return float64(atomic.LoadUint64(&s.bonuscount[bid]))
}

func (s *Stat) JackCount(jid int) float64 {
	return float64(atomic.LoadUint64(&s.jackcount[jid]))
}

func (s *Stat) Update(wins Wins, cfn int) {
	for _, wi := range wins {
		if wi.Pay != 0 {
			if wi.LI != 0 { // line win
				s.lpm.Lock()
				s.linepay += wi.Pay * wi.MP
				s.lpm.Unlock()
			} else { // scatter win
				s.spm.Lock()
				s.scatpay += wi.Pay * wi.MP
				s.spm.Unlock()
			}
		}
		if wi.FS != 0 {
			atomic.AddUint64(&s.freecount, uint64(wi.FS))
			atomic.AddUint64(&s.freehits, 1)
		}
		if wi.BID != 0 {
			atomic.AddUint64(&s.bonuscount[wi.BID], 1)
		}
		if wi.JID != 0 {
			atomic.AddUint64(&s.jackcount[wi.JID], 1)
		}
	}
	if cfn <= FallLimit {
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
			var reshuf = s.Reshuf(1)
			var total = s.Planned()
			var rtp = calc(io.Discard)
			var dur = time.Since(t0)
			if total > 0 {
				var exp = time.Duration(float64(dur) * total / reshuf)
				fmt.Printf("processed %.1fm, ready %2.2f%% (%v / %v), RTP = %2.2f%%  \r",
					reshuf/1e6, reshuf/total*100,
					dur.Truncate(stepdur), exp.Truncate(stepdur),
					rtp)
			} else {
				fmt.Printf("processed %.1fm, spent %v, RTP = %2.2f%%  \r",
					reshuf/1e6, dur.Truncate(stepdur), rtp)
			}
		}
	}
}

type CalcAlg = func(ctx context.Context, s Stater, g SlotGame, reels Reelx)

const (
	CtxGranulation = 1000 // check context every N reshuffles
	FallLimit      = 15
)

var (
	ErrAvalanche = errors.New("too many cascading falls")
	ErrReelCount = errors.New("unexpected number of reels")
)

func CorrectThrNum() int {
	if cfg.MTCount < 1 {
		return runtime.GOMAXPROCS(0)
	}
	return cfg.MTCount
}

func BruteForcex(ctx context.Context, s Stater, g SlotGame, reels Reelx) {
	var plan = reels.Reshuffles()
	s.SetPlan(plan)
	var tn = CorrectThrNum()
	var tn64 = uint64(tn)
	if tn%len(reels[0]) == 0 {
		panic("BruteForcex: thread number equals to 1-st reel length")
	}
	// Number of reels.
	var rn = len(reels)
	// Precompute reel lengths.
	var rlen = make([]int, rn)
	for x := range rn {
		rlen[x] = len(reels[x])
	}
	// Precompute delimiters for position calculation.
	var delim = make([]uint64, rn)
	delim[0] = 1
	for x := 1; x < rn; x++ {
		delim[x] = delim[x-1] * uint64(rlen[x-1])
	}

	var wg sync.WaitGroup
	wg.Add(tn)
	for ti := range tn64 {
		var sg = g.Clone().(ClassicSlot)     // classic slot game
		var cs, iscascade = sg.(CascadeSlot) // cascade slot game
		go func() {
			defer wg.Done()
			var wins Wins
			var rpos = make([]int, rn)
			for x := range rn {
				rpos[x] = -1
			}
			var i uint64
			var x int
			var pos int
			// Using one general loop instead of five nested loops
			// gives ~30% performance increase.
			for i = ti; i < plan; i += tn64 {
				if (i/tn64)%CtxGranulation == 0 {
					select {
					case <-ctx.Done():
						return
					default:
					}
				}
				for x = range rn {
					pos = int(i/delim[x]) % rlen[x]
					if rpos[x] != pos {
						sg.SetCol(Pos(x+1), reels[x], pos)
						rpos[x] = pos
					} else {
						break
					}
				}
				if iscascade {
					var err error
					var cfn int
					for {
						cs.UntoFall()
						if cfn++; cfn > FallLimit {
							err = ErrAvalanche
							break
						}
						var wp = len(wins)
						if err = cs.Scanner(&wins); err != nil {
							break
						}
						cs.Strike(wins[wp:])
						if len(wins) == wp {
							break
						}
						cs.PushFall(reels)
					}
					if err == nil {
						s.Update(wins, cfn)
					} else {
						s.IncErr()
					}
					wins.Reset()
					if cfn > 1 {
						for x := range rn {
							rpos[x] = -1
						}
					}
				} else {
					if sg.Scanner(&wins) == nil {
						s.Update(wins, 1)
					} else {
						s.IncErr()
					}
					wins.Reset()
				}
			}
		}()
	}
	wg.Wait()
}

func BruteForce5x3Big(ctx context.Context, s Stater, g SlotGame, r1, rb, r5 []Sym) {
	s.SetPlan(uint64(len(r1)) * uint64(len(rb)) * uint64(len(r5)))
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
						if sg.Scanner(&wins) == nil {
							s.Update(wins, 1)
						} else {
							s.IncErr()
						}
						wins.Reset()
					}
				}
			}
		}()
	}
	wg.Wait()
}

func MonteCarlo(ctx context.Context, s Stater, g SlotGame, reels Reelx) {
	s.SetPlan(cfg.MCCount * 1e6)
	var tn = CorrectThrNum()
	var tn64 = uint64(tn)
	var n = uint64(s.Planned())
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
					var err error
					var cfn int
					for {
						cs.UntoFall()
						if cfn++; cfn > FallLimit {
							err = ErrAvalanche
							break
						}
						cs.SpinReels(reels)
						var wp = len(wins)
						if err = cs.Scanner(&wins); err != nil {
							break
						}
						cs.Strike(wins[wp:])
						if len(wins) == wp {
							break
						}
					}
					if err == nil {
						s.Update(wins, cfn)
					} else {
						s.IncErr()
					}
					wins.Reset()
				} else {
					sg.SpinReels(reels)
					if sg.Scanner(&wins) == nil {
						s.Update(wins, 1)
					} else {
						s.IncErr()
					}
					wins.Reset()
				}
			}
		}()
	}
	wg.Wait()
}

func MonteCarloPrec(ctx context.Context, s Stater, g SlotGame, reels Reelx, calc func(io.Writer) float64) {
	s.SetPlan(0)
	var tn = CorrectThrNum()
	var rtpcmp float64
	var rtpmux sync.Mutex
	var rtpnum uint64
	var mcc = cfg.MCCount * 1e6 / CtxGranulation
	if mcc == 0 {
		mcc = 1e7 / CtxGranulation
	}
	var wg sync.WaitGroup
	wg.Add(tn)
	for range tn {
		var sg = g.Clone().(ClassicSlot)     // classic slot game
		var cs, iscascade = sg.(CascadeSlot) // cascade slot game
		var reshuf uint64
		go func() {
			defer wg.Done()

			var wins Wins
			for atomic.LoadUint64(&rtpnum) < mcc {
				reshuf++
				if reshuf%CtxGranulation == 0 {
					select {
					case <-ctx.Done():
						return
					default:
					}
					var rtp = calc(io.Discard)
					rtpmux.Lock()
					if diff := rtp - rtpcmp; diff < cfg.MCPrec && diff > -cfg.MCPrec {
						rtpnum++
					} else {
						rtpnum = 0
						rtpcmp = rtp
					}
					rtpmux.Unlock()
				}
				if iscascade {
					var err error
					var cfn int
					for {
						cs.UntoFall()
						if cfn++; cfn > FallLimit {
							err = ErrAvalanche
							break
						}
						cs.SpinReels(reels)
						var wp = len(wins)
						if err = cs.Scanner(&wins); err != nil {
							break
						}
						cs.Strike(wins[wp:])
						if len(wins) == wp {
							break
						}
					}
					if err == nil {
						s.Update(wins, cfn)
					} else {
						s.IncErr()
					}
					wins.Reset()
				} else {
					sg.SpinReels(reels)
					if sg.Scanner(&wins) == nil {
						s.Update(wins, 1)
					} else {
						s.IncErr()
					}
					wins.Reset()
				}
			}
		}()
	}
	wg.Wait()
}

func MCCalcAlg(calc func(io.Writer) float64) CalcAlg {
	return func(ctx context.Context, s Stater, g SlotGame, reels Reelx) {
		MonteCarloPrec(ctx, s, g, reels, calc)
	}
}

func ScanReels(ctx context.Context, s Stater, g ClassicSlot, reels Reelx,
	bruteforce, montecarlo, montecarloprec CalcAlg,
	calc func(io.Writer) float64) float64 {
	if sx, sy := g.Dim(); len(reels) != int(sx) {
		panic(fmt.Errorf("%w: %d reels provided for %dx%d slot", ErrReelCount, len(reels), sx, sy))
	}
	var t0 = time.Now()
	var wg sync.WaitGroup
	wg.Add(2)
	var ctx2, cancel2 = context.WithCancel(ctx)
	go func() {
		defer wg.Done()
		Progress(ctx2, s, calc)
	}()
	go func() {
		defer wg.Done()
		defer cancel2()
		if cfg.MCPrec > 0 {
			montecarloprec(ctx2, s, g, reels)
		} else if cfg.MCCount > 0 {
			montecarlo(ctx2, s, g, reels)
		} else {
			bruteforce(ctx2, s, g, reels)
		}
	}()
	wg.Wait()
	var dur = time.Since(t0)

	if s.Planned() > 0 {
		var comp = s.Reshuf(1) / s.Planned() * 100
		fmt.Printf("completed %.5g%%, selected %d lines, time spent %v            \n", comp, g.GetSel(), dur)
	} else {
		fmt.Printf("produced %.1fm, selected %d lines, time spent %v            \n", s.Reshuf(1)/1e6, g.GetSel(), dur)
	}
	if s.Errors() > 0 {
		fmt.Printf("reels lengths %s, total reshuffles %d, errors %g\n", reels.String(), reels.Reshuffles(), s.Errors())
	} else {
		fmt.Printf("reels lengths %s, total reshuffles %d\n", reels.String(), reels.Reshuffles())
	}
	return calc(os.Stdout)
}

func ScanReelsCommon(ctx context.Context, s Stater, g ClassicSlot, reels Reelx,
	calc func(io.Writer) float64) float64 {
	return ScanReels(ctx, s, g, reels, BruteForcex, MonteCarlo, MCCalcAlg(calc), calc)
}
