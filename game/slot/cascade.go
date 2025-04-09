package slot

import "math/rand/v2"

type Cascade interface {
	Screen
	Cascade() bool        // returns true on avalanche continue
	NewFall()             // set fall number before fall
	RiseFall(reels Reels) // first fall in cascade
	NextFall(reels Reels) // any next fall in cascade
	Strike(wins Wins)     // strike win symbols on the screen
}

type CascadeSlot interface {
	Cascade
	SlotGame
}

type Cascade5x3 struct {
	Sym [5][3]Sym `json:"sym" yaml:"sym" xml:"sym"` // game screen with symbols
	Hit [5][3]int `json:"hit" yaml:"hit" xml:"hit"` // hits to fall down
	Pos [5]int    `json:"pos" yaml:"pos" xml:"pos"` // reels positions
	CFN int       `json:"cfn" yaml:"cfn" xml:"cfn"` // cascade fall number
}

// Declare conformity with Cascade interface.
var _ Cascade = (*Cascade5x3)(nil)

func (s *Cascade5x3) Dim() (Pos, Pos) {
	return 5, 3
}

func (s *Cascade5x3) At(x, y Pos) Sym {
	return s.Sym[x-1][y-1]
}

func (s *Cascade5x3) LY(x Pos, line Linex) Sym {
	return s.Sym[x-1][line[x-1]-1]
}

func (s *Cascade5x3) SetSym(x, y Pos, sym Sym) {
	s.Sym[x-1][y-1] = sym
}

func (s *Cascade5x3) SetCol(x Pos, reel []Sym, pos int) {
	for y := range 3 {
		s.Sym[x-1][y] = reel[(pos+y)%len(reel)]
	}
	s.Pos[x-1] = pos
}

func (s *Cascade5x3) ReelSpin(reels Reels) {
	if s.CFN > 1 {
		s.NextFall(reels)
	} else {
		s.RiseFall(reels)
	}
}

func (s *Cascade5x3) RiseFall(reels Reels) {
	var r5x = reels.(*Reels5x)
	for x := range Pos(5) {
		var reel = r5x[x]
		var pos = rand.N(len(reel))
		s.SetCol(x+1, reel, pos)
	}
}

func (s *Cascade5x3) NextFall(reels Reels) {
	var r5x = reels.(*Reels5x)
	for x := range 5 {
		// fall old symbols
		var n = 0
		for y := range 3 {
			if s.Hit[x][y] > 0 {
				for i := range y {
					s.Sym[x][y-i] = s.Sym[x][y-i-1]
				}
				n++
			}
		}
		// fall new symbols
		s.Pos[x] = (s.Pos[x] - n + len(r5x[x])) % len(r5x[x])
		for y := range n {
			var pos = (s.Pos[x] + y) % len(r5x[x])
			s.Sym[x][y] = r5x[x][pos]
		}
	}
}

func (s *Cascade5x3) SymNum(sym Sym) (n Pos) {
	for x := range 5 {
		for y := range 3 {
			if s.Sym[x][y] == sym {
				n++
			}
		}
	}
	return
}

func (s *Cascade5x3) ScatNum(scat Sym) (n Pos) {
	for x := range 5 {
		var r = s.Sym[x]
		if r[0] == scat || r[1] == scat || r[2] == scat {
			n++
		}
	}
	return
}

func (s *Cascade5x3) ScatPos(scat Sym) (l Linex) {
	for x := range 5 {
		var r = s.Sym[x]
		if r[0] == scat {
			l[x] = 1
		} else if r[1] == scat {
			l[x] = 2
		} else if r[2] == scat {
			l[x] = 3
		}
	}
	return
}

func (s *Cascade5x3) Cascade() bool {
	for _, r := range s.Hit {
		if r[0] > 0 || r[1] > 0 || r[2] > 0 {
			return true
		}
	}
	return false
}

func (s *Cascade5x3) NewFall() {
	if s.Cascade() {
		s.CFN++
	} else {
		s.CFN = 1
	}
}

func (s *Cascade5x3) Strike(wins Wins) {
	clear(s.Hit[:])
	for _, wi := range wins {
		for x := range 5 {
			if y := wi.XY[x]; y > 0 {
				s.Hit[x][y-1] = 1
			}
		}
	}
}

type Cascade5x4 struct {
	Sym [5][4]Sym `json:"sym" yaml:"sym" xml:"sym"` // game screen with symbols
	Hit [5][4]int `json:"hit" yaml:"hit" xml:"hit"` // hits to fall down
	Pos [5]int    `json:"pos" yaml:"pos" xml:"pos"` // reels positions
	CFN int       `json:"cfn" yaml:"cfn" xml:"cfn"` // cascade fall number
}

// Declare conformity with Cascade interface.
var _ Cascade = (*Cascade5x4)(nil)

func (s *Cascade5x4) Dim() (Pos, Pos) {
	return 5, 4
}

func (s *Cascade5x4) At(x, y Pos) Sym {
	return s.Sym[x-1][y-1]
}

func (s *Cascade5x4) LY(x Pos, line Linex) Sym {
	return s.Sym[x-1][line[x-1]-1]
}

func (s *Cascade5x4) SetSym(x, y Pos, sym Sym) {
	s.Sym[x-1][y-1] = sym
}

func (s *Cascade5x4) SetCol(x Pos, reel []Sym, pos int) {
	for y := range 4 {
		s.Sym[x-1][y] = reel[(pos+y)%len(reel)]
	}
	s.Pos[x-1] = pos
}

func (s *Cascade5x4) ReelSpin(reels Reels) {
	if s.CFN > 1 {
		s.NextFall(reels)
	} else {
		s.RiseFall(reels)
	}
}

func (s *Cascade5x4) RiseFall(reels Reels) {
	var r5x = reels.(*Reels5x)
	for x := range Pos(5) {
		var reel = r5x[x]
		var pos = rand.N(len(reel))
		s.SetCol(x+1, reel, pos)
	}
}

func (s *Cascade5x4) NextFall(reels Reels) {
	var r5x = reels.(*Reels5x)
	for x := range 5 {
		// fall old symbols
		var n = 0
		for y := range 4 {
			if s.Hit[x][y] > 0 {
				for i := range y {
					s.Sym[x][y-i] = s.Sym[x][y-i-1]
				}
				n++
			}
		}
		// fall new symbols
		if s.Pos[x] -= n; s.Pos[x] < 0 {
			s.Pos[x] += len(r5x[x])
		}
		for y := range n {
			var pos = (s.Pos[x] + y) % len(r5x[x])
			s.Sym[x][y] = r5x[x][pos]
		}
	}
}

func (s *Cascade5x4) SymNum(sym Sym) (n Pos) {
	for x := range 5 {
		for y := range 4 {
			if s.Sym[x][y] == sym {
				n++
			}
		}
	}
	return
}

func (s *Cascade5x4) ScatNum(scat Sym) (n Pos) {
	for x := range 5 {
		var r = s.Sym[x]
		if r[0] == scat || r[1] == scat || r[2] == scat || r[3] == scat {
			n++
		}
	}
	return
}

func (s *Cascade5x4) ScatPos(scat Sym) (l Linex) {
	for x := range 5 {
		var r = s.Sym[x]
		if r[0] == scat {
			l[x] = 1
		} else if r[1] == scat {
			l[x] = 2
		} else if r[2] == scat {
			l[x] = 3
		} else if r[3] == scat {
			l[x] = 4
		}
	}
	return
}

func (s *Cascade5x4) Cascade() bool {
	for _, r := range s.Hit {
		if r[0] > 0 || r[1] > 0 || r[2] > 0 || r[3] > 0 {
			return true
		}
	}
	return false
}

func (s *Cascade5x4) NewFall() {
	if s.Cascade() {
		s.CFN++
	} else {
		s.CFN = 1
	}
}

func (s *Cascade5x4) Strike(wins Wins) {
	clear(s.Hit[:])
	for _, wi := range wins {
		for x := range 5 {
			if y := wi.XY[x]; y > 0 {
				s.Hit[x][y-1] = 1
			}
		}
	}
}
