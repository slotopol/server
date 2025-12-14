package slot

import "math/rand/v2"

type Cascade interface {
	Screen
	UntoFall()            // set fall number before scanner call
	PushFall(reels Reelx) // fill screen on fall in avalanche chain
	Strike(wins Wins)     // strike win symbols on the screen
}

// Remark: On cascading slots, the reels must not match, otherwise,
// if the same positions appear on the reels during a spin, there
// will be an endless avalanche.

type CascadeSlot interface {
	Cascade
	SlotGame
}

type Cascade5x3 struct {
	Screen5x3 `yaml:",inline"`
	Hit       [5][3]Pos `json:"hit" yaml:"hit,flow" xml:"hit"` // hits to fall down
	Pos       [5]int    `json:"pos" yaml:"pos,flow" xml:"pos"` // reels positions
	// cascade fall number
	CFN int `json:"cfn,omitempty" yaml:"cfn,omitempty" xml:"cfn,omitempty"`
}

// Declare conformity with Cascade interface.
var _ Cascade = (*Cascade5x3)(nil)

func (s *Cascade5x3) SetCol(x Pos, reel []Sym, pos int) {
	var d = &s.Scr[x-1]
	var n = len(reel)
	pos = (n + pos%n) % n // correct position
	for y := range 3 {
		d[y] = reel[(pos+y)%n]
	}
	s.Pos[x-1] = pos
}

func (s *Cascade5x3) ReelSpin(reels Reelx) {
	if s.CFN > 1 {
		s.PushFall(reels)
	} else {
		s.TopFall(reels)
	}
}

func (s *Cascade5x3) TopFall(reels Reelx) {
	for x := range Pos(5) {
		var reel = reels[x]
		var pos = rand.N(len(reel))
		s.SetCol(x+1, reel, pos)
	}
}

func (s *Cascade5x3) PushFall(reels Reelx) {
	for x := range 5 {
		// fall old symbols
		var n = 0
		for y := range 3 {
			if s.Hit[x][y] > 0 {
				for i := range y {
					s.Scr[x][y-i] = s.Scr[x][y-i-1]
				}
				n++
			}
		}
		// fall new symbols
		s.Pos[x] -= n
		for y := range n {
			s.Scr[x][y] = ReelAt(reels[x], s.Pos[x]+y)
		}
	}
}

// Returns true on avalanche continue.
func (s *Cascade5x3) Cascade() bool {
	for _, r := range s.Hit {
		if r[0] > 0 || r[1] > 0 || r[2] > 0 {
			return true
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
	clear(s.Hit[:])
	for _, wi := range wins {
		for i := 0; wi.XY[i][0] > 0; i++ {
			s.Hit[wi.XY[i][0]-1][wi.XY[i][1]-1]++
		}
	}
}

type Cascade5x4 struct {
	Screen5x4 `yaml:",inline"`
	Hit       [5][4]Pos `json:"hit" yaml:"hit,flow" xml:"hit"` // hits to fall down
	Pos       [5]int    `json:"pos" yaml:"pos,flow" xml:"pos"` // reels positions
	// cascade fall number
	CFN int `json:"cfn,omitempty" yaml:"cfn,omitempty" xml:"cfn,omitempty"`
}

// Declare conformity with Cascade interface.
var _ Cascade = (*Cascade5x4)(nil)

func (s *Cascade5x4) SetCol(x Pos, reel []Sym, pos int) {
	var d = &s.Scr[x-1]
	var n = len(reel)
	pos = (n + pos%n) % n // correct position
	for y := range 4 {
		d[y] = reel[(pos+y)%n]
	}
	s.Pos[x-1] = pos
}

func (s *Cascade5x4) ReelSpin(reels Reelx) {
	if s.CFN > 1 {
		s.PushFall(reels)
	} else {
		s.TopFall(reels)
	}
}

func (s *Cascade5x4) TopFall(reels Reelx) {
	for x := range Pos(5) {
		var reel = reels[x]
		var pos = rand.N(len(reel))
		s.SetCol(x+1, reel, pos)
	}
}

func (s *Cascade5x4) PushFall(reels Reelx) {
	for x := range 5 {
		// fall old symbols
		var n = 0
		for y := range 4 {
			if s.Hit[x][y] > 0 {
				for i := range y {
					s.Scr[x][y-i] = s.Scr[x][y-i-1]
				}
				n++
			}
		}
		// fall new symbols
		s.Pos[x] -= n
		for y := range n {
			s.Scr[x][y] = ReelAt(reels[x], s.Pos[x]+y)
		}
	}
}

// Returns true on avalanche continue.
func (s *Cascade5x4) Cascade() bool {
	for _, r := range s.Hit {
		if r[0] > 0 || r[1] > 0 || r[2] > 0 || r[3] > 0 {
			return true
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
	clear(s.Hit[:])
	for _, wi := range wins {
		for i := 0; wi.XY[i][0] > 0; i++ {
			s.Hit[wi.XY[i][0]-1][wi.XY[i][1]-1]++
		}
	}
}

// CascadeGain calculates total gain on avalanche chain for cascading slots.
func CascadeGain(game SlotGame, wins Wins, fund, mrtp float64) (sumgain float64, err error) {
	if len(wins) == 0 {
		return
	}
	if _, ok := game.(CascadeSlot); !ok {
		return
	}
	var casc = game.Clone().(CascadeSlot)
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
