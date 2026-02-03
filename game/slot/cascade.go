package slot

import "math/rand/v2"

type Cascade interface {
	UntoFall()            // set fall number before scanner call
	PushFall(reels Reelx) // fill grid on fall in avalanche chain
	Strike(wins Wins)     // strike win symbols on the grid
}

// Remark: On cascading slots, the reels must not to match, otherwise,
// if the same positions appear on the reels during a spin, there
// will be an endless avalanche.

type SlotCascade interface {
	Cascade
	SlotGeneric
}

type Cascade5x3 struct {
	Grid5x3 `yaml:",inline"`
	Hits    [5][3]Pos `json:"hits" yaml:"hits,flow" xml:"hits"` // hits to fall down
	Seed    [5]int    `json:"seed" yaml:"seed,flow" xml:"seed"` // reels positions
	// cascade fall number
	CFN int `json:"cfn,omitempty" yaml:"cfn,omitempty" xml:"cfn,omitempty"`
}

// Declare conformity with Cascade interface.
var _ Cascade = (*Cascade5x3)(nil)

func (s *Cascade5x3) SetCol(x Pos, reel []Sym, pos int) {
	var sr = &s.Grid[x-1]
	var n = len(reel)
	pos = (n + pos%n) % n // correct position
	for y := range sr {
		sr[y] = reel[(pos+y)%n]
	}
	s.Seed[x-1] = pos
}

func (s *Cascade5x3) SpinReels(reels Reelx) {
	if s.CFN > 1 {
		s.PushFall(reels)
	} else {
		s.TopFall(reels)
	}
}

func (s *Cascade5x3) TopFall(reels Reelx) {
	for x, reel := range reels {
		var pos = rand.N(len(reel))
		s.SetCol(Pos(x+1), reel, pos)
	}
}

func (s *Cascade5x3) PushFall(reels Reelx) {
	for x, hr := range s.Hits {
		var sr = &s.Grid[x]
		// fall old symbols
		var n = 0
		for y, h := range hr {
			if h > 0 {
				for i := range y {
					sr[y-i] = sr[y-i-1]
				}
				n++
			}
		}
		// fall new symbols
		s.Seed[x] -= n
		for y := range n {
			sr[y] = ReelAt(reels[x], s.Seed[x]+y)
		}
	}
}

// Returns true on avalanche continue.
func (s *Cascade5x3) Cascade() bool {
	for _, hr := range s.Hits {
		for _, h := range hr {
			if h > 0 {
				return true
			}
		}
	}
	return false
}

func (s *Cascade5x3) UntoFall() {
	if s.Cascade() {
		s.CFN++
	} else {
		s.CFN = 1
	}
}

func (s *Cascade5x3) Strike(wins Wins) {
	clear(s.Hits[:])
	for _, wi := range wins {
		for i := 0; wi.XY[i][0] > 0; i++ {
			s.Hits[wi.XY[i][0]-1][wi.XY[i][1]-1]++
		}
	}
}

type Cascade5x4 struct {
	Grid5x4 `yaml:",inline"`
	Hits    [5][4]Pos `json:"hits" yaml:"hits,flow" xml:"hits"` // hits to fall down
	Seed    [5]int    `json:"seed" yaml:"seed,flow" xml:"seed"` // reels positions
	// cascade fall number
	CFN int `json:"cfn,omitempty" yaml:"cfn,omitempty" xml:"cfn,omitempty"`
}

// Declare conformity with Cascade interface.
var _ Cascade = (*Cascade5x4)(nil)

func (s *Cascade5x4) SetCol(x Pos, reel []Sym, pos int) {
	var sr = &s.Grid[x-1]
	var n = len(reel)
	pos = (n + pos%n) % n // correct position
	for y := range sr {
		sr[y] = reel[(pos+y)%n]
	}
	s.Seed[x-1] = pos
}

func (s *Cascade5x4) SpinReels(reels Reelx) {
	if s.CFN > 1 {
		s.PushFall(reels)
	} else {
		s.TopFall(reels)
	}
}

func (s *Cascade5x4) TopFall(reels Reelx) {
	for x, reel := range reels {
		var pos = rand.N(len(reel))
		s.SetCol(Pos(x+1), reel, pos)
	}
}

func (s *Cascade5x4) PushFall(reels Reelx) {
	for x, hr := range s.Hits {
		var sr = &s.Grid[x]
		// fall old symbols
		var n = 0
		for y, h := range hr {
			if h > 0 {
				for i := range y {
					sr[y-i] = sr[y-i-1]
				}
				n++
			}
		}
		// fall new symbols
		s.Seed[x] -= n
		for y := range n {
			sr[y] = ReelAt(reels[x], s.Seed[x]+y)
		}
	}
}

// Returns true on avalanche continue.
func (s *Cascade5x4) Cascade() bool {
	for _, hr := range s.Hits {
		for _, h := range hr {
			if h > 0 {
				return true
			}
		}
	}
	return false
}

func (s *Cascade5x4) UntoFall() {
	if s.Cascade() {
		s.CFN++
	} else {
		s.CFN = 1
	}
}

func (s *Cascade5x4) Strike(wins Wins) {
	clear(s.Hits[:])
	for _, wi := range wins {
		for i := 0; wi.XY[i][0] > 0; i++ {
			s.Hits[wi.XY[i][0]-1][wi.XY[i][1]-1]++
		}
	}
}

// CascadeGain calculates total gain on avalanche chain for cascading slots.
func CascadeGain(game SlotGame, wins Wins, fund, mrtp float64) (sumgain float64, err error) {
	if len(wins) == 0 {
		return
	}
	if _, ok := game.(SlotCascade); !ok {
		return
	}
	var casc = game.Clone().(SlotCascade)
	casc.Strike(wins)
	var cfn = 1
	var cw Wins
	for {
		casc.UntoFall()
		if cfn++; cfn > FallLimit {
			err = ErrAvalanche
			return
		}
		casc.Spin(mrtp)
		if err = casc.Scanner(&cw); err != nil {
			return
		}
		if len(cw) == 0 {
			return
		}
		game.Spawn(cw, fund, mrtp)
		sumgain += cw.Gain()
		casc.Strike(cw)
		cw.Reset()
	}
}
