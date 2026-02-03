package slot

import (
	"encoding/json"
	"math/rand/v2"
)

// Grid contains symbols rectangle of the slot game.
// It can be with dimensions 3x1, 3x3, 4x4, 5x3, 5x4 or others.
// (1 ,1) symbol is on left top corner.
type Grider interface {
	Dim() (Pos, Pos)                   // returns grid dimensions
	At(x, y Pos) Sym                   // returns symbol at position (x, y), starts from (1, 1)
	LX(x Pos, line Linex) Sym          // returns symbol at position (x, line(x)), starts from (1, 1)
	SetSym(x, y Pos, sym Sym)          // setup symbol at given position
	SetCol(x Pos, reel []Sym, pos int) // setup column on grid with given reel at given position
	SpinReels(reels Reelx)             // fill the grid with random hits on those reels
	SymNum(sym Sym) (n Pos)            // returns number of symbols on the grid
	SymPos(sym Sym) Hitx               // returns symbols positions on the grid
}

type Bigger interface {
	SetBig(big Sym)
}

// Gridx is a grid with dimensions defined during construction.
// Gridx has fixed size to avoid extra memory allocations.
type Gridx struct {
	sx, sy Pos
	data   [GridSize]Sym
}

// Declare conformity with Grider interface.
var _ Grider = (*Gridx)(nil)

// Construct grid with given dimensions. Maximum possible size is cx*cy=GridSize.
func GridDim(sx, sy Pos) Gridx {
	return Gridx{
		sx: sx, sy: sy,
	}
}

func (s *Gridx) SetDim(sx, sy Pos) {
	s.sx, s.sy = sx, sy
}

func (s *Gridx) Dim() (Pos, Pos) {
	return s.sx, s.sy
}

func (s *Gridx) At(x, y Pos) Sym {
	return s.data[(x-1)*s.sy+y-1]
}

func (s *Gridx) LX(x Pos, line Linex) Sym {
	return s.data[(x-1)*s.sy+line[x-1]-1]
}

func (s *Gridx) SetSym(x, y Pos, sym Sym) {
	s.data[(x-1)*s.sy+y-1] = sym
}

func (s *Gridx) SetCol(x Pos, reel []Sym, pos int) {
	var sr = s.data[(x-1)*s.sy : x*s.sy]
	var n = len(reel)
	pos = (n + pos%n) % n // correct position
	for y := range sr {
		sr[y] = reel[(pos+y)%n]
	}
}

func (s *Gridx) SpinReels(reels Reelx) {
	for x, reel := range reels {
		var hit = rand.N(len(reel))
		s.SetCol(Pos(x+1), reel, hit)
	}
}

func (s *Gridx) SymNum(sym Sym) (n Pos) {
	for _, si := range s.data[:s.sx*s.sy] {
		if si == sym {
			n++
		}
	}
	return
}

func (s *Gridx) SymPos(sym Sym) (c Hitx) {
	var j Pos
	for i, si := range s.data[:s.sx*s.sy] {
		if si == sym {
			c[j][0], c[j][1] = Pos(i)/s.sy+1, Pos(i)%s.sy+1
			j++
		}
	}
	return
}

func (s *Gridx) SymPos2(sym1, sym2 Sym) (c Hitx) {
	var j Pos
	for i, si := range s.data[:s.sx*s.sy] {
		if si == sym1 || si == sym2 {
			c[j][0], c[j][1] = Pos(i)/s.sy+1, Pos(i)%s.sy+1
			j++
		}
	}
	return
}

func (s *Gridx) SymPosL(n Pos, sym Sym) (c Hitx) {
	var j Pos
	for i, si := range s.data[:n*s.sy] {
		if si == sym {
			c[j][0], c[j][1] = Pos(i)/s.sy+1, Pos(i)%s.sy+1
			j++
		}
	}
	return
}

func (s *Gridx) SymPosL2(n Pos, sym1, sym2 Sym) (c Hitx) {
	var j Pos
	for i, si := range s.data[:n*s.sy] {
		if si == sym1 || si == sym2 {
			c[j][0], c[j][1] = Pos(i)/s.sy+1, Pos(i)%s.sy+1
			j++
		}
	}
	return
}

type gridx struct {
	Grid [][]Sym `json:"grid" yaml:"grid,flow" xml:"grid"`
}

func (s *Gridx) MarshalJSON() ([]byte, error) {
	var tmp gridx
	tmp.Grid = make([][]Sym, s.sx)
	for x := range s.sx {
		tmp.Grid[x] = s.data[x*s.sy : (x+1)*s.sy]
	}
	return json.Marshal(tmp)
}

func (s *Gridx) UnmarshalJSON(b []byte) (err error) {
	var tmp gridx
	if err = json.Unmarshal(b, &tmp); err != nil {
		return
	}
	s.sx, s.sy = Pos(len(tmp.Grid)), Pos(len(tmp.Grid[0]))
	for x := range s.sx {
		copy(s.data[x*s.sy:], tmp.Grid[x])
	}
	return
}

// Grid for 3x3 slots.
type Grid3x3 struct {
	Grid [3][3]Sym `json:"grid" yaml:"grid,flow" xml:"grid"`
}

// Declare conformity with Grider interface.
var _ Grider = (*Grid3x3)(nil)

func (s *Grid3x3) Dim() (Pos, Pos) {
	return 3, 3
}

func (s *Grid3x3) At(x, y Pos) Sym {
	return s.Grid[x-1][y-1]
}

func (s *Grid3x3) LX(x Pos, line Linex) Sym {
	return s.Grid[x-1][line[x-1]-1]
}

func (s *Grid3x3) SetSym(x, y Pos, sym Sym) {
	s.Grid[x-1][y-1] = sym
}

func (s *Grid3x3) SetCol(x Pos, reel []Sym, pos int) {
	var sr = &s.Grid[x-1]
	var n = len(reel)
	pos = (n + pos%n) % n // correct position
	for y := range sr {
		sr[y] = reel[(pos+y)%n]
	}
}

func (s *Grid3x3) SpinReels(reels Reelx) {
	for x, reel := range reels {
		var hit = rand.N(len(reel))
		s.SetCol(Pos(x+1), reel, hit)
	}
}

func (s *Grid3x3) SymNum(sym Sym) (n Pos) {
	for _, sr := range s.Grid {
		for _, sy := range sr {
			if sy == sym {
				n++
			}
		}
	}
	return
}

func (s *Grid3x3) SymPos(sym Sym) (c Hitx) {
	var i Pos
	for x, sr := range s.Grid {
		for y, sy := range sr {
			if sy == sym {
				c[i][0], c[i][1] = Pos(x+1), Pos(y+1)
				i++
			}
		}
	}
	return
}

// Grid for 4x4 slots.
type Grid4x4 struct {
	Grid [4][4]Sym `json:"grid" yaml:"grid,flow" xml:"grid"`
}

// Declare conformity with Grider interface.
var _ Grider = (*Grid4x4)(nil)

func (s *Grid4x4) Dim() (Pos, Pos) {
	return 4, 4
}

func (s *Grid4x4) At(x, y Pos) Sym {
	return s.Grid[x-1][y-1]
}

func (s *Grid4x4) LX(x Pos, line Linex) Sym {
	return s.Grid[x-1][line[x-1]-1]
}

func (s *Grid4x4) SetSym(x, y Pos, sym Sym) {
	s.Grid[x-1][y-1] = sym
}

func (s *Grid4x4) SetCol(x Pos, reel []Sym, pos int) {
	var sr = &s.Grid[x-1]
	var n = len(reel)
	pos = (n + pos%n) % n // correct position
	for y := range sr {
		sr[y] = reel[(pos+y)%n]
	}
}

func (s *Grid4x4) SpinReels(reels Reelx) {
	for x, reel := range reels {
		var hit = rand.N(len(reel))
		s.SetCol(Pos(x+1), reel, hit)
	}
}

func (s *Grid4x4) SymNum(sym Sym) (n Pos) {
	for _, sr := range s.Grid {
		for _, sy := range sr {
			if sy == sym {
				n++
			}
		}
	}
	return
}

func (s *Grid4x4) SymPos(sym Sym) (c Hitx) {
	var i Pos
	for x, sr := range s.Grid {
		for y, sy := range sr {
			if sy == sym {
				c[i][0], c[i][1] = Pos(x+1), Pos(y+1)
				i++
			}
		}
	}
	return
}

// Grid for 5x3 slots.
type Grid5x3 struct {
	Grid [5][3]Sym `json:"grid" yaml:"grid,flow" xml:"grid"`
}

// Declare conformity with Grider & Bigger interface.
var _ Grider = (*Grid5x3)(nil)
var _ Bigger = (*Grid5x3)(nil)

func (s *Grid5x3) Dim() (Pos, Pos) {
	return 5, 3
}

func (s *Grid5x3) At(x, y Pos) Sym {
	return s.Grid[x-1][y-1]
}

func (s *Grid5x3) LX(x Pos, line Linex) Sym {
	return s.Grid[x-1][line[x-1]-1]
}

func (s *Grid5x3) SetSym(x, y Pos, sym Sym) {
	s.Grid[x-1][y-1] = sym
}

func (s *Grid5x3) SetCol(x Pos, reel []Sym, pos int) {
	var sr = &s.Grid[x-1]
	var n = len(reel)
	pos = (n + pos%n) % n // correct position
	for y := range sr {
		sr[y] = reel[(pos+y)%n]
	}
}

func (s *Grid5x3) SetBig(big Sym) {
	var x Pos
	for x = 1; x <= 3; x++ {
		var sr = &s.Grid[x]
		for y := range sr {
			sr[y] = big
		}
	}
}

func (s *Grid5x3) SpinReels(reels Reelx) {
	for x, reel := range reels {
		var hit = rand.N(len(reel))
		s.SetCol(Pos(x+1), reel, hit)
	}
}

func (s *Grid5x3) SpinBig(r1, rb, r5 []Sym) {
	var hit int
	// set 1 reel
	hit = rand.N(len(r1))
	s.SetCol(1, r1, hit)
	// set center
	var big = rb[rand.N(len(rb))]
	s.SetBig(big)
	// set 5 reel
	hit = rand.N(len(r5))
	s.SetCol(5, r5, hit)
}

func (s *Grid5x3) SymNum(sym Sym) (n Pos) {
	for _, sr := range s.Grid {
		for _, sy := range sr {
			if sy == sym {
				n++
			}
		}
	}
	return
	// Other way:
	// var b = unsafe.Slice((*byte)(unsafe.Pointer(&g.Grid[0][0])), 15)
	// return Pos(bytes.Count(b, []byte{sym}))
}

func (s *Grid5x3) SymPos(sym Sym) (c Hitx) {
	var i Pos
	for x, sr := range s.Grid {
		for y, sy := range sr {
			if sy == sym {
				c[i][0], c[i][1] = Pos(x+1), Pos(y+1)
				i++
			}
		}
	}
	return
}

func (s *Grid5x3) SymNum2(sym1, sym2 Sym) (n1, n2 Pos) {
	for _, sr := range s.Grid {
		for _, sy := range sr {
			if sy == sym1 {
				n1++
			} else if sy == sym2 {
				n2++
			}
		}
	}
	return
}

func (s *Grid5x3) SymPos2(sym1, sym2 Sym) (c Hitx) {
	var i Pos
	for x, sr := range s.Grid {
		for y, sy := range sr {
			if sy == sym1 || sy == sym2 {
				c[i][0], c[i][1] = Pos(x+1), Pos(y+1)
				i++
			}
		}
	}
	return
}

func (s *Grid5x3) SymPosL(n Pos, sym Sym) (c Hitx) {
	var i Pos
	for x := range n {
		var sr = &s.Grid[x]
		for y, sy := range sr {
			if sy == sym {
				c[i][0], c[i][1] = x+1, Pos(y+1)
				i++
			}
		}
	}
	return
}

func (s *Grid5x3) SymPosL2(n Pos, sym1, sym2 Sym) (c Hitx) {
	var i Pos
	for x := range n {
		var sr = &s.Grid[x]
		for y, sy := range sr {
			if sy == sym1 || sy == sym2 {
				c[i][0], c[i][1] = x+1, Pos(y+1)
				i++
			}
		}
	}
	return
}

// Grid for 5x4 slots.
type Grid5x4 struct {
	Grid [5][4]Sym `json:"grid" yaml:"grid,flow" xml:"grid"`
}

// Declare conformity with Grider interface.
var _ Grider = (*Grid5x4)(nil)

func (s *Grid5x4) Dim() (Pos, Pos) {
	return 5, 4
}

func (s *Grid5x4) At(x, y Pos) Sym {
	return s.Grid[x-1][y-1]
}

func (s *Grid5x4) LX(x Pos, line Linex) Sym {
	return s.Grid[x-1][line[x-1]-1]
}

func (s *Grid5x4) SetSym(x, y Pos, sym Sym) {
	s.Grid[x-1][y-1] = sym
}

func (s *Grid5x4) SetCol(x Pos, reel []Sym, pos int) {
	var sr = &s.Grid[x-1]
	var n = len(reel)
	pos = (n + pos%n) % n // correct position
	for y := range sr {
		sr[y] = reel[(pos+y)%n]
	}
}

func (s *Grid5x4) SpinReels(reels Reelx) {
	for x, reel := range reels {
		var hit = rand.N(len(reel))
		s.SetCol(Pos(x+1), reel, hit)
	}
}

func (s *Grid5x4) SymNum(sym Sym) (n Pos) {
	for _, sr := range s.Grid {
		for _, sy := range sr {
			if sy == sym {
				n++
			}
		}
	}
	return
}

func (s *Grid5x4) SymPos(sym Sym) (c Hitx) {
	var i Pos
	for x, sr := range s.Grid {
		for y, sy := range sr {
			if sy == sym {
				c[i][0], c[i][1] = Pos(x+1), Pos(y+1)
				i++
			}
		}
	}
	return
}

func (s *Grid5x4) SymNum2(sym1, sym2 Sym) (n1, n2 Pos) {
	for _, sr := range s.Grid {
		for _, sy := range sr {
			if sy == sym1 {
				n1++
			} else if sy == sym2 {
				n2++
			}
		}
	}
	return
}

func (s *Grid5x4) SymPos2(sym1, sym2 Sym) (c Hitx) {
	var i Pos
	for x, sr := range s.Grid {
		for y, sy := range sr {
			if sy == sym1 || sy == sym2 {
				c[i][0], c[i][1] = Pos(x+1), Pos(y+1)
				i++
			}
		}
	}
	return
}

func (s *Grid5x4) SymPosL(n Pos, sym Sym) (c Hitx) {
	var i Pos
	for x := range n {
		var sr = &s.Grid[x]
		for y, sy := range sr {
			if sy == sym {
				c[i][0], c[i][1] = x+1, Pos(y+1)
				i++
			}
		}
	}
	return
}

func (s *Grid5x4) SymPosL2(n Pos, sym1, sym2 Sym) (c Hitx) {
	var i Pos
	for x := range n {
		var sr = &s.Grid[x]
		for y, sy := range sr {
			if sy == sym1 || sy == sym2 {
				c[i][0], c[i][1] = x+1, Pos(y+1)
				i++
			}
		}
	}
	return
}

// Grid for 6x3 slots.
type Grid6x3 struct {
	Grid [6][3]Sym `json:"grid" yaml:"grid,flow" xml:"grid"`
}

// Declare conformity with Grider interface.
var _ Grider = (*Grid6x3)(nil)

func (s *Grid6x3) Dim() (Pos, Pos) {
	return 6, 3
}

func (s *Grid6x3) At(x, y Pos) Sym {
	return s.Grid[x-1][y-1]
}

func (s *Grid6x3) LX(x Pos, line Linex) Sym {
	return s.Grid[x-1][line[x-1]-1]
}

func (s *Grid6x3) SetSym(x, y Pos, sym Sym) {
	s.Grid[x-1][y-1] = sym
}

func (s *Grid6x3) SetCol(x Pos, reel []Sym, pos int) {
	var sr = &s.Grid[x-1]
	var n = len(reel)
	pos = (n + pos%n) % n // correct position
	for y := range sr {
		sr[y] = reel[(pos+y)%n]
	}
}

func (s *Grid6x3) SpinReels(reels Reelx) {
	for x, reel := range reels {
		var hit = rand.N(len(reel))
		s.SetCol(Pos(x+1), reel, hit)
	}
}

func (s *Grid6x3) SymNum(sym Sym) (n Pos) {
	for _, sr := range s.Grid {
		for _, sy := range sr {
			if sy == sym {
				n++
			}
		}
	}
	return
}

func (s *Grid6x3) SymPos(sym Sym) (c Hitx) {
	var i Pos
	for x, sr := range s.Grid {
		for y, sy := range sr {
			if sy == sym {
				c[i][0], c[i][1] = Pos(x+1), Pos(y+1)
				i++
			}
		}
	}
	return
}

// Grid for 6x4 slots.
type Grid6x4 struct {
	Grid [6][4]Sym `json:"grid" yaml:"grid,flow" xml:"grid"`
}

// Declare conformity with Grider interface.
var _ Grider = (*Grid6x4)(nil)

func (s *Grid6x4) Dim() (Pos, Pos) {
	return 6, 4
}

func (s *Grid6x4) At(x, y Pos) Sym {
	return s.Grid[x-1][y-1]
}

func (s *Grid6x4) LX(x Pos, line Linex) Sym {
	return s.Grid[x-1][line[x-1]-1]
}

func (s *Grid6x4) SetSym(x, y Pos, sym Sym) {
	s.Grid[x-1][y-1] = sym
}

func (s *Grid6x4) SetCol(x Pos, reel []Sym, pos int) {
	var sr = &s.Grid[x-1]
	var n = len(reel)
	pos = (n + pos%n) % n // correct position
	for y := range sr {
		sr[y] = reel[(pos+y)%n]
	}
}

func (s *Grid6x4) SpinReels(reels Reelx) {
	for x, reel := range reels {
		var hit = rand.N(len(reel))
		s.SetCol(Pos(x+1), reel, hit)
	}
}

func (s *Grid6x4) SymNum(sym Sym) (n Pos) {
	for _, sr := range s.Grid {
		for _, sy := range sr {
			if sy == sym {
				n++
			}
		}
	}
	return
}

func (s *Grid6x4) SymPos(sym Sym) (c Hitx) {
	var i Pos
	for x, sr := range s.Grid {
		for y, sy := range sr {
			if sy == sym {
				c[i][0], c[i][1] = Pos(x+1), Pos(y+1)
				i++
			}
		}
	}
	return
}
