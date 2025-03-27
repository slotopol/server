package slot

import "math/rand/v2"

type Cascade interface {
	Screen
	Cascade() bool
	Strike(wins Wins)
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
}

func (s *Cascade5x3) ReelSpin(reels Reels) {
	var r5x = reels.(*Reels5x)
	if s.Cascade() {
		s.Fall(r5x)
	} else {
		s.NewSpin(r5x)
	}
}

func (s *Cascade5x3) NewSpin(reels *Reels5x) {
	for x := range Pos(5) {
		var reel = reels[x]
		var pos = rand.N(len(reel))
		s.SetCol(x+1, reel, pos)
		s.Pos[x] = pos
	}
}

func (s *Cascade5x3) Fall(reels *Reels5x) {
	for x := range 5 {
		// fall old symbols
		var n = 0
		for y := range 3 {
			if s.Hit[x][y] > 0 {
				for i := range y {
					s.Hit[x][y-i] = s.Hit[x][y-i-1]
				}
				n++
			}
		}
		// fall new symbols
		s.Pos[x] -= n
		for y := range n {
			s.Sym[x][y] = reels[x][(s.Pos[x]+y)%len(reels[x])]
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

func (s *Cascade5x3) Strike(wins Wins) {
	clear(s.Hit[:])
	for _, wi := range wins {
		for x := range wi.Num {
			var y = wi.XY[x-1]
			s.Hit[x][y] = 1
		}
	}
}
