package slot

import "math/rand/v2"

type Cascade interface {
	Screen
	UntoFall()            // set fall number before scanner call
	PushFall(reels Reels) // fill screen on fall in avalanche chain
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
	Scr [5][3]Sym `json:"scr" yaml:"scr,flow" xml:"scr"` // game screen with symbols
	Hit [5][3]Sym `json:"hit" yaml:"hit,flow" xml:"hit"` // hits to fall down
	Pos [5]int    `json:"pos" yaml:"pos,flow" xml:"pos"` // reels positions
	// cascade fall number
	CFN int `json:"cfn,omitempty" yaml:"cfn,omitempty" xml:"cfn,omitempty"`
}

// Declare conformity with Cascade interface.
var _ Cascade = (*Cascade5x3)(nil)

func (s *Cascade5x3) Dim() (Pos, Pos) {
	return 5, 3
}

func (s *Cascade5x3) At(x, y Pos) Sym {
	return s.Scr[x-1][y-1]
}

func (s *Cascade5x3) LY(x Pos, line Linex) Sym {
	return s.Scr[x-1][line[x-1]-1]
}

func (s *Cascade5x3) SetSym(x, y Pos, sym Sym) {
	s.Scr[x-1][y-1] = sym
}

func (s *Cascade5x3) SetCol(x Pos, reel []Sym, pos int) {
	for y := range 3 {
		s.Scr[x-1][y] = ReelAt(reel, pos+y)
	}
	s.Pos[x-1] = pos
}

func (s *Cascade5x3) ReelSpin(reels Reels) {
	if s.CFN > 1 {
		s.PushFall(reels)
	} else {
		s.TopFall(reels)
	}
}

func (s *Cascade5x3) TopFall(reels Reels) {
	var r5x = reels.(*Reels5x)
	for x := range Pos(5) {
		var reel = r5x[x]
		var pos = rand.N(len(reel))
		s.SetCol(x+1, reel, pos)
	}
}

func (s *Cascade5x3) PushFall(reels Reels) {
	var r5x = reels.(*Reels5x)
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
			s.Scr[x][y] = ReelAt(r5x[x], s.Pos[x]+y)
		}
	}
}

func (s *Cascade5x3) SymNum(sym Sym) (n Pos) {
	for x := range 5 {
		var r = s.Scr[x]
		for y := range 3 {
			if r[y] == sym {
				n++
			}
		}
	}
	return
}

func (s *Cascade5x3) SymPos(sym Sym) (c Hitx) {
	var x, y, i Pos
	for x = range 5 {
		var r = s.Scr[x]
		for y = range 3 {
			if r[y] == sym {
				c[i][0], c[i][1] = x+1, y+1
				i++
			}
		}
	}
	return
}

func (s *Cascade5x3) SymNum2(sym1, sym2 Sym) (n1, n2 Pos) {
	for x := range 5 {
		var r = s.Scr[x]
		for y := range 3 {
			if r[y] == sym1 {
				n1++
			} else if r[y] == sym2 {
				n2++
			}
		}
	}
	return
}

func (s *Cascade5x3) SymPos2(sym1, sym2 Sym) (c Hitx) {
	var x, y, i Pos
	for x = range 5 {
		var r = s.Scr[x]
		for y = range 3 {
			if r[y] == sym1 || r[y] == sym2 {
				c[i][0], c[i][1] = x+1, y+1
				i++
			}
		}
	}
	return
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
			s.Hit[wi.XY[i][0]-1][wi.XY[i][1]-1] = 1
		}
	}
}

type Cascade5x4 struct {
	Scr [5][4]Sym `json:"scr" yaml:"scr,flow" xml:"scr"` // game screen with symbols
	Hit [5][4]Sym `json:"hit" yaml:"hit,flow" xml:"hit"` // hits to fall down
	Pos [5]int    `json:"pos" yaml:"pos,flow" xml:"pos"` // reels positions
	// cascade fall number
	CFN int `json:"cfn,omitempty" yaml:"cfn,omitempty" xml:"cfn,omitempty"`
}

// Declare conformity with Cascade interface.
var _ Cascade = (*Cascade5x4)(nil)

func (s *Cascade5x4) Dim() (Pos, Pos) {
	return 5, 4
}

func (s *Cascade5x4) At(x, y Pos) Sym {
	return s.Scr[x-1][y-1]
}

func (s *Cascade5x4) LY(x Pos, line Linex) Sym {
	return s.Scr[x-1][line[x-1]-1]
}

func (s *Cascade5x4) SetSym(x, y Pos, sym Sym) {
	s.Scr[x-1][y-1] = sym
}

func (s *Cascade5x4) SetCol(x Pos, reel []Sym, pos int) {
	for y := range 4 {
		s.Scr[x-1][y] = ReelAt(reel, pos+y)
	}
	s.Pos[x-1] = pos
}

func (s *Cascade5x4) ReelSpin(reels Reels) {
	if s.CFN > 1 {
		s.PushFall(reels)
	} else {
		s.TopFall(reels)
	}
}

func (s *Cascade5x4) TopFall(reels Reels) {
	var r5x = reels.(*Reels5x)
	for x := range Pos(5) {
		var reel = r5x[x]
		var pos = rand.N(len(reel))
		s.SetCol(x+1, reel, pos)
	}
}

func (s *Cascade5x4) PushFall(reels Reels) {
	var r5x = reels.(*Reels5x)
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
			s.Scr[x][y] = ReelAt(r5x[x], s.Pos[x]+y)
		}
	}
}

func (s *Cascade5x4) SymNum(sym Sym) (n Pos) {
	for x := range 5 {
		var r = s.Scr[x]
		for y := range 4 {
			if r[y] == sym {
				n++
			}
		}
	}
	return
}

func (s *Cascade5x4) SymPos(sym Sym) (c Hitx) {
	var x, y, i Pos
	for x = range 5 {
		var r = s.Scr[x]
		for y = range 4 {
			if r[y] == sym {
				c[i][0], c[i][1] = x+1, y+1
				i++
			}
		}
	}
	return
}

func (s *Cascade5x4) SymNum2(sym1, sym2 Sym) (n1, n2 Pos) {
	for x := range 5 {
		var r = s.Scr[x]
		for y := range 4 {
			if r[y] == sym1 {
				n1++
			} else if r[y] == sym2 {
				n2++
			}
		}
	}
	return
}

func (s *Cascade5x4) SymPos2(sym1, sym2 Sym) (c Hitx) {
	var x, y, i Pos
	for x = range 5 {
		var r = s.Scr[x]
		for y = range 4 {
			if r[y] == sym1 || r[y] == sym2 {
				c[i][0], c[i][1] = x+1, y+1
				i++
			}
		}
	}
	return
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
			s.Hit[wi.XY[i][0]-1][wi.XY[i][1]-1] = 1
		}
	}
}

// CascadeGain calculates total gain on avalanche chain for cascading slots.
func CascadeGain(game SlotGame, wins Wins, fund, mrtp float64) (cascgain float64, err error) {
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
		cascgain += cw.Gain()
		casc.Strike(cw)
		cw.Reset()
	}
}
