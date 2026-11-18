package slot

import "slices"

type StatPage struct {
	N   uint64
	C   [][]float64
	S   [][]float64
	Cfg uint64
	Cfs uint64
	Qfs uint64
	Cb  []uint64
	Sb  []float64
	Cj  []uint64
}

func (sp *StatPage) IsZero() bool {
	return sp.N == 0 // N can not be zero for non-empty structure
}

// NewStatPage creates a new StatPage object with a specified number
// of symbols, bonuses, and jackpots. Each symbol has a specified number
// of possible occurrences in a win.
func NewStatPage(sn, rn, bn, jn int) *StatPage {
	var sp = &StatPage{
		C:  make([][]float64, sn),
		S:  make([][]float64, sn),
		Cb: make([]uint64, bn),
		Sb: make([]float64, bn),
		Cj: make([]uint64, jn),
	}
	for i := range sn {
		sp.C[i] = make([]float64, rn)
		sp.S[i] = make([]float64, rn)
	}
	return sp
}

// Reserves memory for counters and sums for given symbol by
// maximum number of this symbol in win. Symbol is 1-based.
func (sp *StatPage) MakeSym(sym Sym, rn int) *StatPage {
	sp.C[sym-1] = make([]float64, rn)
	sp.S[sym-1] = make([]float64, rn)
	return sp
}

func (sp *StatPage) Clone() *StatPage {
	var dub = &StatPage{
		N:   sp.N,
		C:   make([][]float64, len(sp.C)),
		S:   make([][]float64, len(sp.C)),
		Cfg: sp.Cfg,
		Cfs: sp.Cfs,
		Qfs: sp.Qfs,
		Cb:  slices.Clone(sp.Cb),
		Sb:  slices.Clone(sp.Sb),
		Cj:  slices.Clone(sp.Cj),
	}
	for sym, c := range sp.C {
		dub.C[sym] = slices.Clone(c)
	}
	for sym, c := range sp.S {
		dub.S[sym] = slices.Clone(c)
	}
	return dub
}

func (sp *StatPage) Clear() {
	sp.N = 0
	for _, c := range sp.C {
		clear(c)
	}
	for _, c := range sp.S {
		clear(c)
	}
	sp.Cfg = 0
	sp.Cfs = 0
	sp.Qfs = 0
	clear(sp.Cb)
	clear(sp.Sb)
	clear(sp.Cj)
}

// Add adds all counters and sums from given StatPage to this StatPage.
func (sp *StatPage) Add(ap *StatPage) {
	sp.N += ap.N
	for sym, c := range sp.C {
		for x, v := range ap.C[sym] {
			c[x] += v
		}
	}
	for sym, c := range sp.S {
		for x, v := range ap.S[sym] {
			c[x] += v
		}
	}
	sp.Cfg += ap.Cfg
	sp.Cfs += ap.Cfs
	sp.Qfs += ap.Qfs
	for id, v := range ap.Cb {
		sp.Cb[id] += v
	}
	for id, v := range ap.Sb {
		sp.Sb[id] += v
	}
	for id, v := range ap.Cj {
		sp.Cj[id] += v
	}
}

// SumPays returns sum of pays by all symbols and bonuses.
func (sp *StatPage) SumPays() (sum float64) {
	for _, pays := range sp.S {
		for i := range pays {
			sum += pays[i]
		}
	}
	for i := range sp.Sb {
		sum += sp.Sb[i]
	}
	return
}

func (sp *StatPage) Update(wins Wins) (pay float64) {
	for _, wi := range wins {
		if wi.Sym > 0 && wi.Num > 0 {
			sp.C[wi.Sym-1][wi.Num-1]++
		}
		if wi.Pay != 0 {
			var p = wi.Pay * wi.MP
			if wi.Sym > 0 { // symbols pay
				sp.S[wi.Sym-1][wi.Num-1] += p
			} else { // bonus pay, only 2 cases can be
				sp.Sb[wi.BID-1] += p
			}
			pay += p
		}
		if wi.FS != 0 {
			sp.Cfg++
			sp.Cfs += uint64(wi.FS)
			sp.Qfs += uint64(wi.FS * wi.FS)
		}
		if wi.BID != 0 {
			sp.Cb[wi.BID-1]++
		}
		if wi.JID != 0 {
			sp.Cj[wi.JID-1]++
		}
	}
	sp.N++
	return
}

type StatSnap struct {
	A  []*StatPage
	Q  float64
	EC uint64
}

// StatSnap implements Simulator interface.
var _ Simulator = (*StatSnap)(nil)

// NewStatSnap creates a new StatSnap object with a specified number
// of pages (corresponding to avalanche cascades), where each page
// contains a certain number of symbols, bonuses, and jackpots.
func NewStatSnap(cn, sn, rn, bn, jn int) *StatSnap {
	var ss = &StatSnap{
		A: make([]*StatPage, cn),
	}
	for cfn := range ss.A {
		ss.A[cfn] = NewStatPage(sn, rn, bn, jn)
	}
	return ss
}

// MakeSym reserves memory for counters and sums for given symbol by
// maximum number of this symbol in win. Symbol is 1-based.
func (ss *StatSnap) MakeSym(sym Sym, rn int) *StatSnap {
	for _, sp := range ss.A {
		sp.MakeSym(sym, rn)
	}
	return ss
}

func (ss *StatSnap) Clone() *StatSnap {
	var dub = &StatSnap{
		A:  make([]*StatPage, len(ss.A)),
		Q:  ss.Q,
		EC: ss.EC,
	}
	for cfn, sp := range ss.A {
		dub.A[cfn] = sp.Clone()
	}
	return dub
}

func (ss *StatSnap) Clear() {
	for _, sp := range ss.A {
		sp.Clear()
	}
	ss.Q = 0
	ss.EC = 0
}

// Add adds all counters and sums from given StatSnap to this StatSnap.
func (ss *StatSnap) Add(as *StatSnap) {
	for cfn, sp := range ss.A {
		sp.Add(as.A[cfn])
	}
	ss.Q += as.Q
	ss.EC += as.EC
}

// SumPays returns sum of pays by all symbols and bonuses.
func (ss *StatSnap) SumPays() (sum float64) {
	for _, sp := range ss.A {
		if sp.N == 0 {
			break
		}
		sum += sp.SumPays()
	}
	return
}

// Count returns number of processed grid reshuffles, including with no wins.
func (ss *StatSnap) Count() float64 {
	return float64(ss.A[0].N)
}

// NSQ returns number of spins, sum of pays, sum of squares of pays by spin,
// normalized by spin cost.
func (ss *StatSnap) NSQ(cost float64) (N, S, Q float64) {
	N = float64(ss.A[0].N)
	S = ss.SumPays() / cost
	Q = ss.Q / cost / cost
	return
}

// FGC returns counter of triggering free game series, total count of free spins,
// and sum of squared free spin series lengths.
func (ss *StatSnap) FGC() (Cfg, Cfs, Qfs float64) {
	for _, sp := range ss.A {
		Cfg += float64(sp.Cfg)
		Cfs += float64(sp.Cfs)
		Qfs += float64(sp.Qfs)
	}
	return
}

// Simulate performs spin simulation, calculates results, update statistics.
func (ss *StatSnap) Simulate(g SlotGame, wins *Wins, mrtp float64) {
	if g.Scanner(wins) != nil {
		ss.EC++
		return
	}
	if pay := ss.A[0].Update(*wins); pay != 0 {
		ss.Q += pay * pay
	}
}
