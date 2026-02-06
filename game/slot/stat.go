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
	"github.com/slotopol/server/game"
)

type Stater interface {
	GetPlan() uint64
	SetPlan(n uint64)
	Errors() uint64
	Count() float64
	Simulate(SlotGame, Reelx, *Wins)
}

type ScanPar = game.ScanPar

// Maximum possible number of symbols at any game
const MaxSymNum = 20

type StatCounter struct {
	N   atomic.Uint64             // number of processed grid reshuffles, including with no wins
	S   [MaxSymNum]atomic.Float64 // sum of pays by symbols
	FSC atomic.Uint64             // free spins count
	FHC atomic.Uint64             // free games hits count
	BHC [8]atomic.Uint64          // bonus hits count
	JHC [4]atomic.Uint64          // jackpot hits count
}

func (c *StatCounter) Update(wins Wins) (pay float64) {
	for _, wi := range wins {
		if wi.Pay != 0 {
			var p = wi.Pay * wi.MP
			c.S[wi.Sym].Add(p)
			pay += p
		}
		if wi.FS != 0 {
			c.FSC.Add(uint64(wi.FS))
			c.FHC.Inc()
		}
		if wi.BID != 0 {
			c.BHC[wi.BID].Inc()
		}
		if wi.JID != 0 {
			c.JHC[wi.JID].Inc()
		}
	}
	c.N.Inc()
	return
}

func (s *StatCounter) SumS() (S float64) {
	for sym := range MaxSymNum {
		S += s.S[sym].Load()
	}
	return
}

type StatGeneric struct {
	Plan     atomic.Uint64
	ErrCount atomic.Uint64
	StatCounter
	Q atomic.Float64 // sum of squares of pays by symbols
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
	return float64(s.N.Load())
}

func (s *StatGeneric) RTPsym(cost float64, scat Sym) (lrtp, srtp float64) {
	var sym Sym
	for sym = range MaxSymNum {
		if sym != scat {
			lrtp += s.S[sym].Load()
		} else {
			srtp += s.S[sym].Load()
		}
	}
	var N = s.Count()
	lrtp /= N * cost
	srtp /= N * cost
	return
}

func (s *StatGeneric) RTPsym2(cost float64, scat1, scat2 Sym) (lrtp, srtp float64) {
	var sym Sym
	for sym = range MaxSymNum {
		if sym != scat1 && sym != scat2 {
			lrtp += s.S[sym].Load()
		} else {
			srtp += s.S[sym].Load()
		}
	}
	var N = s.Count()
	lrtp /= N * cost
	srtp /= N * cost
	return
}

func (s *StatGeneric) NSQ(cost float64) (N float64, S float64, Q float64) {
	N = s.Count()
	S = s.SumS() / cost
	Q = s.Q.Load() / cost / cost
	return
}

func (s *StatGeneric) SymVariance(cost float64) float64 {
	var N, S, Q = s.NSQ(cost)
	// another way: Q/N - S*S/N/N
	return N*Q - S*S/N/N
}

func (s *StatGeneric) SymSigma(cost float64) float64 {
	var N, S, Q = s.NSQ(cost)
	// another way: math.Sqrt(Q/N - S*S/N/N)
	return math.Sqrt(N*Q-S*S) / N
}

// Returns (q, sq), where q = free spins quantifier, sq = 1/(1-q)
// sum of a decreasing geometric progression for retriggered free spins.
func (s *StatGeneric) FSQ() (q float64, sq float64) {
	q = float64(s.FSC.Load()) / s.Count()
	sq = 1 / (1 - q)
	return
}

// Quantifier of free games per reshuffles.
func (s *StatGeneric) FGQ() float64 {
	return float64(s.FHC.Load()) / s.Count()
}

// Free Games Frequency: average number of reshuffles per free games hit.
func (s *StatGeneric) FGF() float64 {
	return s.Count() / float64(s.FHC.Load())
}

func (s *StatGeneric) BonusHitsF(bid int) float64 {
	return float64(s.BHC[bid].Load())
}

func (s *StatGeneric) JackHitsF(jid int) float64 {
	return float64(s.JHC[jid].Load())
}

func (s *StatGeneric) Simulate(g SlotGame, reels Reelx, wins *Wins) {
	if g.Scanner(wins) != nil {
		s.ErrCount.Inc()
		return
	}
	if pay := s.Update(*wins); pay != 0 {
		s.Q.Add(pay * pay)
	}
}

type StatCascade struct {
	Plan     atomic.Uint64
	ErrCount atomic.Uint64
	Casc     [FallLimit]StatCounter
	Q        atomic.Float64 // sum of squares of pays by symbols
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
	return float64(s.Casc[0].N.Load())
}

func (s *StatCascade) SumFreeCount() uint64 {
	var sum uint64
	for i := range FallLimit {
		sum += s.Casc[i].FSC.Load()
	}
	return sum
}

func (s *StatCascade) SumFreeHits() uint64 {
	var sum uint64
	for i := range FallLimit {
		sum += s.Casc[i].FHC.Load()
	}
	return sum
}

func (s *StatCascade) SumBonusHits(bid int) uint64 {
	var sum uint64
	for i := range FallLimit {
		sum += s.Casc[i].BHC[bid].Load()
	}
	return sum
}

func (s *StatCascade) SumJackHits(jid int) uint64 {
	var sum uint64
	for i := range FallLimit {
		sum += s.Casc[i].JHC[jid].Load()
	}
	return sum
}

func (s *StatCascade) RTPsym(cost float64, scat Sym) (lrtp, srtp float64) {
	var sym Sym
	for i := range FallLimit {
		var c = &s.Casc[i]
		for sym = range MaxSymNum {
			if sym != scat {
				lrtp += c.S[sym].Load()
			} else {
				srtp += c.S[sym].Load()
			}
		}
	}
	var N = s.Count()
	lrtp /= N * cost
	srtp /= N * cost
	return
}

func (s *StatCascade) NSQ(cost float64) (N float64, S float64, Q float64) {
	N = s.Count()
	for i := range FallLimit {
		S += s.Casc[i].SumS()
	}
	S /= cost
	Q = s.Q.Load() / cost / cost
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
	var pay1 = s.Casc[0].SumS()
	var pays float64
	for i := range FallLimit {
		var payi = s.Casc[i].SumS()
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
		sum += s.Casc[i].N.Load()
	}
	return float64(sum) / float64(s.Casc[1].N.Load())
}

// Inverse coefficient of fading (C1/Cn)^(1/(n-1)).
func (s *StatCascade) Kfading() float64 {
	var N1 = s.Casc[0].N.Load()
	var Nn uint64
	var i int
	for i = range FallLimit {
		var Ni = s.Casc[i].N.Load()
		if Ni == 0 {
			break
		}
		Nn = Ni
	}
	return math.Pow(float64(N1)/float64(Nn), 1/(float64(i)-1))
}

// Maximum number of cascades in avalanche.
func (s *StatCascade) Ncascmax() int {
	for i := range FallLimit {
		if s.Casc[i].N.Load() == 0 {
			return i
		}
	}
	return FallLimit
}

func (s *StatCascade) Simulate(g SlotGame, reels Reelx, wins *Wins) {
	var sc = g.(SlotCascade)
	var err error
	var cfn int
	var pay float64
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
		pay += s.Casc[cfn-1].Update((*wins)[wp:])
		sc.Strike((*wins)[wp:])
		if len(*wins) == wp {
			break
		}
		sc.PushFall(reels)
	}
	if pay > 0 {
		s.Q.Add(pay * pay)
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
			var N = s.Count()
			var total = float64(s.GetPlan())
			var rtp = calc(io.Discard)
			var dur = time.Since(t0)
			if total > 0 {
				var exp = time.Duration(float64(dur) * total / N)
				fmt.Printf("processed %.1fm, ready %2.2f%% (%v / %v), RTP = %2.2f%%  \r",
					N/1e6, N/total*100,
					dur.Truncate(stepdur), exp.Truncate(stepdur),
					rtp*100)
			} else {
				fmt.Printf("processed %.1fm, spent %v, RTP = %2.2f%%  \r",
					N/1e6, dur.Truncate(stepdur), rtp*100)
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
		var N uint64
		go func() {
			defer wg.Done()

			var wins Wins
			for i1 := range r1 {
				gt.SetCol(1, r1, i1)
				for _, big := range rb {
					cb.SetBig(big)
					for i5 := range r5 {
						N++
						if N%CtxGranulation == 0 {
							select {
							case <-ctx.Done():
								return
							default:
							}
						}
						if N%tn64 != ti {
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
		var N uint64
		go func() {
			defer wg.Done()

			var wins Wins
			for range n / tn64 {
				N++
				if N%CtxGranulation == 0 {
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
		var N uint64
		go func() {
			defer wg.Done()

			var wins Wins
			for rtpnum.Load() < mcc {
				N++
				if N%CtxGranulation == 0 {
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
