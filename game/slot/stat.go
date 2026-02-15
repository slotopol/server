package slot

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"sync"

	"github.com/slotopol/server/game"

	"go.uber.org/atomic"
	"gopkg.in/yaml.v3"
)

type Simulator interface {
	// Returns number of spins, sum of pays, sum of squares of pays by spin,
	// normalized by spin cost.
	NSQ(float64) (float64, float64, float64)
	// Performs spin simulation, calculates results, update statistics.
	Simulate(SlotGame, Reelx, *Wins)
}

type ScanPar = game.ScanPar

type Uint64 struct {
	atomic.Uint64
}

func (x Uint64) String() string {
	return strconv.FormatUint(x.Load(), 10)
}

// MarshalXML encodes the wrapped uint64 into XML.
func (x *Uint64) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(x.Load(), start)
}

// UnmarshalXML decodes a uint64 from XML.
func (x *Uint64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v uint64
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	x.Store(v)
	return nil
}

// MarshalYAML encodes the wrapped uint64 into YAML.
func (x Uint64) MarshalYAML() (any, error) {
	return x.Load(), nil
}

// UnmarshalYAML decodes a uint64 from YAML.
func (x *Uint64) UnmarshalYAML(value *yaml.Node) error {
	var v uint64
	if err := value.Decode(&v); err != nil {
		return err
	}
	x.Store(v)
	return nil
}

type Float64 struct {
	atomic.Float64
}

func (x Float64) String() string {
	return strconv.FormatFloat(x.Load(), 'f', -1, 64)
}

// MarshalXML encodes the wrapped float64 into XML.
func (x *Float64) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(x.String(), start)
}

// UnmarshalXML decodes a float64 from XML.
func (x *Float64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v float64
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	x.Store(v)
	return nil
}

// MarshalYAML encodes the wrapped float64 into YAML.
func (x Float64) MarshalYAML() (any, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: x.String(),
	}, nil
}

// UnmarshalYAML decodes a float64 from YAML.
func (x *Float64) UnmarshalYAML(value *yaml.Node) error {
	var v float64
	if err := value.Decode(&v); err != nil {
		return err
	}
	x.Store(v)
	return nil
}

type Counts[T fmt.Stringer] [][]T

func (c Counts[T]) MarshalYAML() (any, error) {
	var node1 = &yaml.Node{
		Kind:  yaml.SequenceNode,
		Style: yaml.FoldedStyle,
	}
	for sym, line := range c {
		var node2 = &yaml.Node{
			Kind:  yaml.SequenceNode,
			Style: yaml.FlowStyle,
		}
		if sym > 0 {
			node2.LineComment = strconv.Itoa(sym)
		}

		for _, v := range line {
			node2.Content = append(node2.Content, &yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: v.String(),
			})
		}

		node1.Content = append(node1.Content, node2)
	}
	return node1, nil
}

type StatCounter struct {
	N   Uint64          // number of processed grid reshuffles, including with no wins
	C   Counts[Uint64]  // hit counter
	S   Counts[Float64] // sum of pays by symbols
	FSC Uint64          // free spins count
	FGH Uint64          // free games hits count
	BH  [8]Uint64       `yaml:",flow"` // bonus hits count
	JH  [4]Uint64       `yaml:",flow"` // jackpot hits count
}

func (c *StatCounter) SymDim(sym Sym, pn int) {
	c.C[sym] = make([]Uint64, pn+1)
	c.S[sym] = make([]Float64, pn+1)
}

func (c *StatCounter) CntDim(sn, pn int) {
	c.C = make([][]Uint64, sn+1)
	c.S = make([][]Float64, sn+1)
	for sym := range c.S {
		c.SymDim(Sym(sym), pn)
	}
}

func (c *StatCounter) Update(wins Wins) (pay float64) {
	for _, wi := range wins {
		c.C[wi.Sym][wi.Num].Inc()
		if wi.Pay != 0 {
			var p = wi.Pay * wi.MP
			c.S[wi.Sym][wi.Num].Add(p)
			pay += p
		}
		if wi.FS != 0 {
			c.FSC.Add(uint64(wi.FS))
			c.FGH.Inc()
		}
		if wi.BID != 0 {
			c.BH[wi.BID].Inc()
		}
		if wi.JID != 0 {
			c.JH[wi.JID].Inc()
		}
	}
	c.N.Inc()
	return
}

func (c *StatCounter) SymS(sym Sym) (S float64) {
	var pays = c.S[sym]
	for sym := range pays {
		S += pays[sym].Load()
	}
	return
}

func (c *StatCounter) SumS() (S float64) {
	for _, pays := range c.S {
		for sym := range pays {
			S += pays[sym].Load()
		}
	}
	return
}

type StatGeneric struct {
	StatCounter
	Q  Float64 // sum of squares of pays by symbols
	EC Uint64  // errors count
}

// Declare conformity with Stater interface.
var _ Simulator = (*StatGeneric)(nil)

func NewStatGeneric(sn, pn int) *StatGeneric {
	var s StatGeneric
	s.CntDim(sn, pn)
	return &s
}

func (s *StatGeneric) Errors() uint64 {
	return s.EC.Load()
}

func (s *StatGeneric) Count() float64 {
	return float64(s.N.Load())
}

func (s *StatGeneric) RTPsym(cost float64, scat Sym) (lrtp, srtp float64) {
	for _, pays := range s.S {
		for sym := range pays {
			if Sym(sym) != scat {
				lrtp += pays[sym].Load()
			} else {
				srtp += pays[sym].Load()
			}
		}
	}
	var N = s.Count()
	lrtp /= N * cost
	srtp /= N * cost
	return
}

func (s *StatGeneric) RTPsym2(cost float64, scat1, scat2 Sym) (lrtp, srtp float64) {
	for _, pays := range s.S {
		for sym := range pays {
			if Sym(sym) != scat1 && Sym(sym) != scat2 {
				lrtp += pays[sym].Load()
			} else {
				srtp += pays[sym].Load()
			}
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
	return float64(s.FGH.Load()) / s.Count()
}

// Free Games Frequency: average number of reshuffles per free games hits.
func (s *StatGeneric) FGF() float64 {
	return s.Count() / float64(s.FGH.Load())
}

func (s *StatGeneric) BonusHitsF(bid int) float64 {
	return float64(s.BH[bid].Load())
}

func (s *StatGeneric) JackHitsF(jid int) float64 {
	return float64(s.JH[jid].Load())
}

func (s *StatGeneric) Simulate(g SlotGame, reels Reelx, wins *Wins) {
	if g.Scanner(wins) != nil {
		s.EC.Inc()
		return
	}
	if pay := s.Update(*wins); pay != 0 {
		s.Q.Add(pay * pay)
	}
}

type StatCascade struct {
	Casc [FallLimit]StatCounter
	Q    Float64 // sum of squares of pays by symbols
	EC   Uint64  // errors count
}

// Declare conformity with Stater interface.
var _ Simulator = (*StatCascade)(nil)

func NewStatCascade(sn, pn int) *StatCascade {
	var s StatCascade
	s.CntDim(sn, pn)
	return &s
}

func (s *StatCascade) CntDim(sn, pn int) {
	for i := range FallLimit {
		s.Casc[i].CntDim(sn, pn)
	}
}

func (s *StatCascade) Errors() uint64 {
	return s.EC.Load()
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
		sum += s.Casc[i].FGH.Load()
	}
	return sum
}

func (s *StatCascade) SumBonusHits(bid int) uint64 {
	var sum uint64
	for i := range FallLimit {
		sum += s.Casc[i].BH[bid].Load()
	}
	return sum
}

func (s *StatCascade) SumJackHits(jid int) uint64 {
	var sum uint64
	for i := range FallLimit {
		sum += s.Casc[i].JH[jid].Load()
	}
	return sum
}

func (s *StatCascade) RTPsym(cost float64, scat Sym) (lrtp, srtp float64) {
	for i := range FallLimit {
		var c = &s.Casc[i]
		for _, pays := range c.S {
			for sym := range pays {
				if Sym(sym) != scat {
					lrtp += pays[sym].Load()
				} else {
					srtp += pays[sym].Load()
				}
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
	var sc = g.(Cascade)
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
		if err = g.Scanner(wins); err != nil {
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
		s.EC.Inc()
		return
	}
}

type CalcAlg = func(ctx context.Context, sp *ScanPar, s Simulator, g SlotGeneric, reels Reelx)

const (
	CtxGranulation = 1000 // check context every N reshuffles
	FallLimit      = 20   // maximum 20 cascades for reels ~100 symbols length
)

var (
	ErrAvalanche = errors.New("too many cascading falls")
	ErrReelCount = errors.New("unexpected number of reels")
)

func ScanReels(ctx context.Context, sp *ScanPar, s Simulator, g SlotGeneric, reels Reelx,
	bruteforce, montecarlo CalcAlg,
	calc func(io.Writer) float64) float64 {
	if sx, sy := g.Dim(); len(reels) != int(sx) {
		panic(fmt.Errorf("%w: %d reels provided for %dx%d slot", ErrReelCount, len(reels), sx, sy))
	}
	if sel := g.GetSel(); sel > 0 {
		fmt.Printf("selected %d lines\n", sel)
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		switch sp.Method {
		case game.CMbruteforce:
			go func() {
				defer wg.Done()
				ProgressBF(ctx2, sp, s, calc, float64(reels.Reshuffles()))
			}()
			bruteforce(ctx2, sp, s, g, reels)
		case game.CMmontecarlo:
			go func() {
				defer wg.Done()
				ProgressMC(ctx2, sp, s, calc, g.Cost())
			}()
			montecarlo(ctx2, sp, s, g, reels)
		}
	}()
	wg.Wait()

	fmt.Printf("reels lengths %s, total reshuffles %d\n", reels.String(), reels.Reshuffles())
	return calc(os.Stdout)
}

func ScanReelsCommon(ctx context.Context, sp *ScanPar, s Simulator, g SlotGeneric, reels Reelx,
	calc func(io.Writer) float64) float64 {
	return ScanReels(ctx, sp, s, g, reels, BruteForcex, MonteCarlo, calc)
}
