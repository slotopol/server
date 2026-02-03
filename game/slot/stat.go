package slot

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"sync"
	"time"

	"go.uber.org/atomic"

	cfg "github.com/slotopol/server/config"
)

type Stater interface {
	GetPlan() uint64
	SetPlan(n uint64)
	Errors() uint64
	Count() float64
	Simulate(SlotGame, Reelx, *Wins)
}

type StatCounter struct {
	Reshuf    atomic.Uint64    // number of processed grid reshuffles, including with no wins
	LinePay   atomic.Float64   // sum of pays by all bet lines
	ScatPay   atomic.Float64   // sum of pays by scatters
	FreeCount atomic.Uint64    // count of free spins
	FreeHits  atomic.Uint64    // count of free games hits
	BonusHits [8]atomic.Uint64 // count of bonus hits
	JackHits  [4]atomic.Uint64
}

type StatGeneric struct {
	Plan     atomic.Uint64
	ErrCount atomic.Uint64
	StatCounter
}

// Declare conformity with Stater interface.
var _ Stater = (*StatGeneric)(nil)

func (s *StatGeneric) GetPlan() uint64 {
	return s.Plan.Load()
}

func (s *StatGeneric) SetPlan(n uint64) {
	s.Plan.Store(n)
}

func (s *StatGeneric) Errors() uint64 {
	return s.ErrCount.Load()
}

func (s *StatGeneric) Count() float64 {
	return float64(s.Reshuf.Load())
}

func (s *StatGeneric) LineRTP(cost float64) float64 {
	return s.LinePay.Load() / s.Count() / cost * 100
}

func (s *StatGeneric) ScatRTP(cost float64) float64 {
	return s.ScatPay.Load() / s.Count() / cost * 100
}

func (s *StatGeneric) SymRTP(cost float64) (lrtp, srtp float64) {
	var reshuf = s.Count()
	lrtp = s.LinePay.Load() / reshuf / cost * 100
	srtp = s.ScatPay.Load() / reshuf / cost * 100
	return
}

// Returns (q, sq), where q = free spins quantifier, sq = 1/(1-q)
// sum of a decreasing geometric progression for retriggered free spins.
func (s *StatGeneric) FSQ() (q float64, sq float64) {
	q = float64(s.FreeCount.Load()) / s.Count()
	sq = 1 / (1 - q)
	return
}

// Quantifier of free games per reshuffles.
func (s *StatGeneric) FGQ() float64 {
	return float64(s.FreeHits.Load()) / s.Count()
}

// Free Games Frequency: average number of reshuffles per free games hit.
func (s *StatGeneric) FGF() float64 {
	return s.Count() / float64(s.FreeHits.Load())
}

func (s *StatGeneric) BonusHitsF(bid int) float64 {
	return float64(s.BonusHits[bid].Load())
}

func (s *StatGeneric) JackHitsF(jid int) float64 {
	return float64(s.JackHits[jid].Load())
}

func (s *StatGeneric) Update(wins Wins) {
	var lpay, spay float64
	for _, wi := range wins {
		if wi.Pay != 0 {
			if wi.LI != 0 { // line win
				lpay += wi.Pay * wi.MP
			} else { // scatter win
				spay += wi.Pay * wi.MP
			}
		}
		if wi.FS != 0 {
			s.FreeCount.Add(uint64(wi.FS))
			s.FreeHits.Inc()
		}
		if wi.BID != 0 {
			s.BonusHits[wi.BID].Inc()
		}
		if wi.JID != 0 {
			s.JackHits[wi.JID].Inc()
		}
	}
	if lpay != 0 {
		s.LinePay.Add(lpay)
	}
	if spay != 0 {
		s.ScatPay.Add(spay)
	}
	s.Reshuf.Inc()
}

func (s *StatGeneric) Simulate(g SlotGame, reels Reelx, wins *Wins) {
	if g.Scanner(wins) != nil {
		s.ErrCount.Inc()
		return
	}
	s.Update(*wins)
}

type StatCascade struct {
	Plan     atomic.Uint64
	ErrCount atomic.Uint64
	Casc     [FallLimit]StatCounter
}

// Declare conformity with Stater interface.
var _ Stater = (*StatCascade)(nil)

func (s *StatCascade) GetPlan() uint64 {
	return s.Plan.Load()
}

func (s *StatCascade) SetPlan(n uint64) {
	s.Plan.Store(n)
}

func (s *StatCascade) Errors() uint64 {
	return s.ErrCount.Load()
}

func (s *StatCascade) Count() float64 {
	return float64(s.Casc[0].Reshuf.Load())
}

func (s *StatCascade) SumLinePay() float64 {
	var sum float64
	for i := range FallLimit {
		sum += s.Casc[i].LinePay.Load()
	}
	return sum
}

func (s *StatCascade) SumScatPay() float64 {
	var sum float64
	for i := range FallLimit {
		sum += s.Casc[i].ScatPay.Load()
	}
	return sum
}

func (s *StatCascade) SumFreeCount() uint64 {
	var sum uint64
	for i := range FallLimit {
		sum += s.Casc[i].FreeCount.Load()
	}
	return sum
}

func (s *StatCascade) SumFreeHits() uint64 {
	var sum uint64
	for i := range FallLimit {
		sum += s.Casc[i].FreeHits.Load()
	}
	return sum
}

func (s *StatCascade) SumBonusHits(bid int) uint64 {
	var sum uint64
	for i := range FallLimit {
		sum += s.Casc[i].BonusHits[bid].Load()
	}
	return sum
}

func (s *StatCascade) SumJackHits(jid int) uint64 {
	var sum uint64
	for i := range FallLimit {
		sum += s.Casc[i].JackHits[jid].Load()
	}
	return sum
}

func (s *StatCascade) LineRTP(cost float64) float64 {
	return s.SumLinePay() / s.Count() / cost * 100
}

func (s *StatCascade) ScatRTP(cost float64) float64 {
	return s.SumScatPay() / s.Count() / cost * 100
}

func (s *StatCascade) SymRTP(cost float64) (lrtp, srtp float64) {
	var lpay, spay float64
	for i := range FallLimit {
		lpay += s.Casc[i].LinePay.Load()
		spay += s.Casc[i].ScatPay.Load()
	}
	var reshuf = s.Count()
	lrtp = lpay / reshuf / cost * 100
	srtp = spay / reshuf / cost * 100
	return
}

// Returns (q, sq), where q = free spins quantifier, sq = 1/(1-q)
// sum of a decreasing geometric progression for retriggered free spins.
func (s *StatCascade) FSQ() (q float64, sq float64) {
	q = float64(s.SumFreeCount()) / s.Count()
	sq = 1 / (1 - q)
	return
}

// Quantifier of free games per reshuffles.
func (s *StatCascade) FGQ() float64 {
	return float64(s.SumFreeHits()) / s.Count()
}

// Free Games Frequency: average number of reshuffles per free games hit.
func (s *StatCascade) FGF() float64 {
	return s.Count() / float64(s.SumFreeHits())
}

// Cascade multiplier.
func (s *StatCascade) Mcascade() float64 {
	var pay1 = s.Casc[0].LinePay.Load() + s.Casc[0].ScatPay.Load()
	var pays float64
	for i := range FallLimit {
		var payi = s.Casc[i].LinePay.Load() + s.Casc[i].ScatPay.Load()
		pays += payi
		if payi == 0 {
			break
		}
	}
	return pays / pay1
}

// Average Cascade Length.
func (s *StatCascade) ACL() float64 {
	var sum uint64
	for i := 1; i < FallLimit; i++ {
		sum += s.Casc[i].Reshuf.Load()
	}
	return float64(sum) / float64(s.Casc[1].Reshuf.Load())
}

// Inverse coefficient of fading (C1/Cn)^(1/(n-1)).
func (s *StatCascade) Kfading() float64 {
	var reshuf1 = s.Casc[0].Reshuf.Load()
	var reshufn uint64
	var i int
	for i = range FallLimit {
		var reshuf = s.Casc[i].Reshuf.Load()
		if reshuf == 0 {
			break
		}
		reshufn = reshuf
	}
	return math.Pow(float64(reshuf1)/float64(reshufn), 1/(float64(i)-1))
}

// Maximum number of cascades in avalanche.
func (s *StatCascade) Ncascmax() int {
	for i := range FallLimit {
		if s.Casc[i].Reshuf.Load() == 0 {
			return i
		}
	}
	return FallLimit
}

func (s *StatCascade) Update(wins Wins, cfn int) {
	var c = &s.Casc[cfn-1]
	var lpay, spay float64
	for _, wi := range wins {
		if wi.Pay != 0 {
			if wi.LI != 0 { // line win
				lpay += wi.Pay * wi.MP
			} else { // scatter win
				spay += wi.Pay * wi.MP
			}
		}
		if wi.FS != 0 {
			c.FreeCount.Add(uint64(wi.FS))
			c.FreeHits.Inc()
		}
		if wi.BID != 0 {
			c.BonusHits[wi.BID].Inc()
		}
		if wi.JID != 0 {
			c.JackHits[wi.JID].Inc()
		}
	}
	if lpay != 0 {
		c.ScatPay.Add(lpay)
	}
	if spay != 0 {
		c.ScatPay.Add(spay)
	}
	c.Reshuf.Inc()
}

func (s *StatCascade) Simulate(g SlotGame, reels Reelx, wins *Wins) {
	var sc = g.(SlotCascade)
	var err error
	var cfn int
	for {
		sc.UntoFall()
		if cfn++; cfn > FallLimit {
			err = ErrAvalanche
			break
		}
		var wp = len(*wins)
		if err = sc.Scanner(wins); err != nil {
			break
		}
		s.Update((*wins)[wp:], cfn)
		sc.Strike((*wins)[wp:])
		if len(*wins) == wp {
			break
		}
		sc.PushFall(reels)
	}
	if err != nil {
		s.ErrCount.Inc()
		return
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
			var reshuf = s.Count()
			var total = float64(s.GetPlan())
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
	FallLimit      = 20   // maximum 20 cascades for reels ~100 symbols length
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
		var gt = g.Clone().(SlotGeneric)
		var _, iscascade = gt.(SlotCascade)
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
						gt.SetCol(Pos(x+1), reels[x], pos)
						rpos[x] = pos
					} else {
						break
					}
				}
				s.Simulate(gt, reels, &wins)
				if iscascade && len(wins) > 0 {
					for x := range rn {
						rpos[x] = -1
					}
				}
				wins.Reset()
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
		var gt = g.Clone().(SlotGeneric)
		var cb = gt.(Bigger)
		var reshuf uint64
		go func() {
			defer wg.Done()

			var wins Wins
			for i1 := range r1 {
				gt.SetCol(1, r1, i1)
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
						gt.SetCol(5, r5, i5)
						s.Simulate(gt, nil, &wins)
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
	var n = s.GetPlan()
	var wg sync.WaitGroup
	wg.Add(tn)
	for range tn64 {
		var gt = g.Clone()
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
				s.Simulate(gt, reels, &wins)
				wins.Reset()
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
	var rtpnum atomic.Uint64
	var mcc = cfg.MCCount * 1e6 / CtxGranulation
	if mcc == 0 {
		mcc = 1e7 / CtxGranulation
	}
	var wg sync.WaitGroup
	wg.Add(tn)
	for range tn {
		var gt = g.Clone()
		var reshuf uint64
		go func() {
			defer wg.Done()

			var wins Wins
			for rtpnum.Load() < mcc {
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
						rtpnum.Inc()
					} else {
						rtpnum.Store(0)
						rtpcmp = rtp
					}
					rtpmux.Unlock()
				}
				s.Simulate(gt, reels, &wins)
				wins.Reset()
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

func ScanReels(ctx context.Context, s Stater, g SlotGeneric, reels Reelx,
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

	if s.GetPlan() > 0 {
		var comp = s.Count() / float64(s.GetPlan()) * 100
		fmt.Printf("completed %.5g%%, selected %d lines, time spent %v            \n", comp, g.GetSel(), dur)
	} else {
		fmt.Printf("produced %.1fm, selected %d lines, time spent %v            \n", s.Count()/1e6, g.GetSel(), dur)
	}
	if s.Errors() > 0 {
		fmt.Printf("reels lengths %s, total reshuffles %d, errors %d\n", reels.String(), reels.Reshuffles(), s.Errors())
	} else {
		fmt.Printf("reels lengths %s, total reshuffles %d\n", reels.String(), reels.Reshuffles())
	}
	return calc(os.Stdout)
}

func ScanReelsCommon(ctx context.Context, s Stater, g SlotGeneric, reels Reelx,
	calc func(io.Writer) float64) float64 {
	return ScanReels(ctx, s, g, reels, BruteForcex, MonteCarlo, MCCalcAlg(calc), calc)
}
