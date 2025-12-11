package slot

import (
	"encoding/json"
	"math/rand/v2"
)

// Screen contains symbols rectangle of the slot game.
// It can be with dimensions 3x1, 3x3, 4x4, 5x3, 5x4 or others.
// (1 ,1) symbol is on left top corner.
type Screen interface {
	Dim() (Pos, Pos)                   // returns screen dimensions
	At(x, y Pos) Sym                   // returns symbol at position (x, y), starts from (1, 1)
	LY(x Pos, line Linex) Sym          // returns symbol at position (x, line(x)), starts from (1, 1)
	SetSym(x, y Pos, sym Sym)          // setup symbol at given position
	SetCol(x Pos, reel []Sym, pos int) // setup column on screen with given reel at given position
	ReelSpin(reels Reelx)              // fill the screen with random hits on those reels
	SymNum(sym Sym) (n Pos)            // returns number of symbols on the screen
	SymPos(sym Sym) Hitx               // returns symbols positions on the screen
}

type Bigger interface {
	SetBig(big Sym)
}

// Screenx is a screen with dimensions defined during construction.
// Screenx has fixed size to avoid extra memory allocations.
type Screenx struct {
	sx, sy Pos
	data   [ScrxSize]Sym
}

// Declare conformity with Screen interface.
var _ Screen = (*Screenx)(nil)

// Construct screen with given dimensions. Maximum possible size is 8x5.
func ScreenDim(sx, sy Pos) Screenx {
	return Screenx{
		sx: sx, sy: sy,
	}
}

func (s *Screenx) SetDim(sx, sy Pos) {
	s.sx, s.sy = sx, sy
}

func (s *Screenx) Dim() (Pos, Pos) {
	return s.sx, s.sy
}

func (s *Screenx) At(x, y Pos) Sym {
	return s.data[(x-1)*s.sy+y-1]
}

func (s *Screenx) LY(x Pos, line Linex) Sym {
	return s.data[(x-1)*s.sy+line[x-1]-1]
}

func (s *Screenx) SetSym(x, y Pos, sym Sym) {
	s.data[(x-1)*s.sy+y-1] = sym
}

func (s *Screenx) SetCol(x Pos, reel []Sym, pos int) {
	var d = s.data[(x-1)*s.sy:]
	for y := range int(s.sy) {
		d[y] = ReelAt(reel, pos+y)
	}
}

func (s *Screenx) ReelSpin(reels Reelx) {
	var x Pos
	for x = 1; x <= s.sx; x++ {
		var reel = reels.Reel(x)
		var hit = rand.N(len(reel))
		s.SetCol(x, reel, hit)
	}
}

func (s *Screenx) SymNum(sym Sym) (n Pos) {
	for i := range s.sx * s.sy {
		if s.data[i] == sym {
			n++
		}
	}
	return
}

func (s *Screenx) SymPos(sym Sym) (c Hitx) {
	var i, n Pos
	for i = range s.sx * s.sy {
		if s.data[i] == sym {
			c[n][0], c[n][1] = i/s.sy+1, i%s.sy+1
			n++
		}
	}
	return
}

type scrx struct {
	Scr [][]Sym `json:"scr" yaml:"scr,flow" xml:"scr"`
}

func (s *Screenx) MarshalJSON() ([]byte, error) {
	var tmp scrx
	tmp.Scr = make([][]Sym, s.sx)
	for x := range s.sx {
		tmp.Scr[x] = s.data[x*s.sy : (x+1)*s.sy]
	}
	return json.Marshal(tmp)
}

func (s *Screenx) UnmarshalJSON(b []byte) (err error) {
	var tmp scrx
	if err = json.Unmarshal(b, &tmp); err != nil {
		return
	}
	s.sx, s.sy = Pos(len(tmp.Scr)), Pos(len(tmp.Scr[0]))
	for x := range s.sx {
		copy(s.data[x*s.sy:], tmp.Scr[x])
	}
	return
}

// Screen for 3x3 slots.
type Screen3x3 struct {
	Scr [3][3]Sym `json:"scr" yaml:"scr,flow" xml:"scr"`
}

// Declare conformity with Screen interface.
var _ Screen = (*Screen3x3)(nil)

func (s *Screen3x3) Dim() (Pos, Pos) {
	return 3, 3
}

func (s *Screen3x3) At(x, y Pos) Sym {
	return s.Scr[x-1][y-1]
}

func (s *Screen3x3) LY(x Pos, line Linex) Sym {
	return s.Scr[x-1][line[x-1]-1]
}

func (s *Screen3x3) SetSym(x, y Pos, sym Sym) {
	s.Scr[x-1][y-1] = sym
}

func (s *Screen3x3) SetCol(x Pos, reel []Sym, pos int) {
	for y := range 3 {
		s.Scr[x-1][y] = ReelAt(reel, pos+y)
	}
}

func (s *Screen3x3) ReelSpin(reels Reelx) {
	var x Pos
	for x = 1; x <= 3; x++ {
		var reel = reels.Reel(x)
		var hit = rand.N(len(reel))
		s.SetCol(x, reel, hit)
	}
}

func (s *Screen3x3) SymNum(sym Sym) (n Pos) {
	for x := range 3 {
		var r = s.Scr[x]
		for y := range 3 {
			if r[y] == sym {
				n++
			}
		}
	}
	return
}

func (s *Screen3x3) SymPos(sym Sym) (c Hitx) {
	var x, y, i Pos
	for x = range 3 {
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

// Screen for 4x4 slots.
type Screen4x4 struct {
	Scr [4][4]Sym `json:"scr" yaml:"scr,flow" xml:"scr"`
}

// Declare conformity with Screen interface.
var _ Screen = (*Screen4x4)(nil)

func (s *Screen4x4) Dim() (Pos, Pos) {
	return 4, 4
}

func (s *Screen4x4) At(x, y Pos) Sym {
	return s.Scr[x-1][y-1]
}

func (s *Screen4x4) LY(x Pos, line Linex) Sym {
	return s.Scr[x-1][line[x-1]-1]
}

func (s *Screen4x4) SetSym(x, y Pos, sym Sym) {
	s.Scr[x-1][y-1] = sym
}

func (s *Screen4x4) SetCol(x Pos, reel []Sym, pos int) {
	for y := range 4 {
		s.Scr[x-1][y] = ReelAt(reel, pos+y)
	}
}

func (s *Screen4x4) ReelSpin(reels Reelx) {
	var x Pos
	for x = 1; x <= 4; x++ {
		var reel = reels.Reel(x)
		var hit = rand.N(len(reel))
		s.SetCol(x, reel, hit)
	}
}

func (s *Screen4x4) SymNum(sym Sym) (n Pos) {
	for x := range 4 {
		var r = s.Scr[x]
		for y := range 4 {
			if r[y] == sym {
				n++
			}
		}
	}
	return
}

func (s *Screen4x4) SymPos(sym Sym) (c Hitx) {
	var x, y, i Pos
	for x = range 4 {
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

// Screen for 5x3 slots.
type Screen5x3 struct {
	Scr [5][3]Sym `json:"scr" yaml:"scr,flow" xml:"scr"`
}

// Declare conformity with Screen & Bigger interface.
var _ Screen = (*Screen5x3)(nil)
var _ Bigger = (*Screen5x3)(nil)

func (s *Screen5x3) Dim() (Pos, Pos) {
	return 5, 3
}

func (s *Screen5x3) At(x, y Pos) Sym {
	return s.Scr[x-1][y-1]
}

func (s *Screen5x3) LY(x Pos, line Linex) Sym {
	return s.Scr[x-1][line[x-1]-1]
}

func (s *Screen5x3) SetSym(x, y Pos, sym Sym) {
	s.Scr[x-1][y-1] = sym
}

func (s *Screen5x3) SetCol(x Pos, reel []Sym, pos int) {
	for y := range 3 {
		s.Scr[x-1][y] = ReelAt(reel, pos+y)
	}
}

func (s *Screen5x3) SetBig(big Sym) {
	var x Pos
	for x = 1; x <= 3; x++ {
		s.Scr[x][0] = big
		s.Scr[x][1] = big
		s.Scr[x][2] = big
	}
}

func (s *Screen5x3) ReelSpin(reels Reelx) {
	var x Pos
	for x = 1; x <= 5; x++ {
		var reel = reels.Reel(x)
		var hit = rand.N(len(reel))
		s.SetCol(x, reel, hit)
	}
}

func (s *Screen5x3) SpinBig(r1, rb, r5 []Sym) {
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

func (s *Screen5x3) SymNum(sym Sym) (n Pos) {
	for x := range 5 {
		var r = s.Scr[x]
		for y := range 3 {
			if r[y] == sym {
				n++
			}
		}
	}
	return
	// Other way:
	// var b = unsafe.Slice((*byte)(unsafe.Pointer(&g.Scr[0][0])), 15)
	// return Pos(bytes.Count(b, []byte{sym}))
}

func (s *Screen5x3) SymPos(sym Sym) (c Hitx) {
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

func (s *Screen5x3) SymNum2(sym1, sym2 Sym) (n1, n2 Pos) {
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

func (s *Screen5x3) SymPos2(sym1, sym2 Sym) (c Hitx) {
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

// Screen for 5x4 slots.
type Screen5x4 struct {
	Scr [5][4]Sym `json:"scr" yaml:"scr,flow" xml:"scr"`
}

// Declare conformity with Screen interface.
var _ Screen = (*Screen5x4)(nil)

func (s *Screen5x4) Dim() (Pos, Pos) {
	return 5, 4
}

func (s *Screen5x4) At(x, y Pos) Sym {
	return s.Scr[x-1][y-1]
}

func (s *Screen5x4) LY(x Pos, line Linex) Sym {
	return s.Scr[x-1][line[x-1]-1]
}

func (s *Screen5x4) SetSym(x, y Pos, sym Sym) {
	s.Scr[x-1][y-1] = sym
}

func (s *Screen5x4) SetCol(x Pos, reel []Sym, pos int) {
	for y := range 4 {
		s.Scr[x-1][y] = ReelAt(reel, pos+y)
	}
}

func (s *Screen5x4) ReelSpin(reels Reelx) {
	var x Pos
	for x = 1; x <= 5; x++ {
		var reel = reels.Reel(x)
		var hit = rand.N(len(reel))
		s.SetCol(x, reel, hit)
	}
}

func (s *Screen5x4) SymNum(sym Sym) (n Pos) {
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

func (s *Screen5x4) SymPos(sym Sym) (c Hitx) {
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

func (s *Screen5x4) SymNum2(sym1, sym2 Sym) (n1, n2 Pos) {
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

func (s *Screen5x4) SymPos2(sym1, sym2 Sym) (c Hitx) {
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

// Screen for 6x3 slots.
type Screen6x3 struct {
	Scr [6][3]Sym `json:"scr" yaml:"scr,flow" xml:"scr"`
}

// Declare conformity with Screen interface.
var _ Screen = (*Screen6x3)(nil)

func (s *Screen6x3) Dim() (Pos, Pos) {
	return 6, 3
}

func (s *Screen6x3) At(x, y Pos) Sym {
	return s.Scr[x-1][y-1]
}

func (s *Screen6x3) LY(x Pos, line Linex) Sym {
	return s.Scr[x-1][line[x-1]-1]
}

func (s *Screen6x3) SetSym(x, y Pos, sym Sym) {
	s.Scr[x-1][y-1] = sym
}

func (s *Screen6x3) SetCol(x Pos, reel []Sym, pos int) {
	for y := range 3 {
		s.Scr[x-1][y] = ReelAt(reel, pos+y)
	}
}

func (s *Screen6x3) ReelSpin(reels Reelx) {
	var x Pos
	for x = 1; x <= 6; x++ {
		var reel = reels.Reel(x)
		var hit = rand.N(len(reel))
		s.SetCol(x, reel, hit)
	}
}

func (s *Screen6x3) SymNum(sym Sym) (n Pos) {
	for x := range 6 {
		var r = s.Scr[x]
		for y := range 3 {
			if r[y] == sym {
				n++
			}
		}
	}
	return
}

func (s *Screen6x3) SymPos(sym Sym) (c Hitx) {
	var x, y, i Pos
	for x = range 6 {
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

// Screen for 6x4 slots.
type Screen6x4 struct {
	Scr [6][4]Sym `json:"scr" yaml:"scr,flow" xml:"scr"`
}

// Declare conformity with Screen interface.
var _ Screen = (*Screen6x4)(nil)

func (s *Screen6x4) Dim() (Pos, Pos) {
	return 6, 4
}

func (s *Screen6x4) At(x, y Pos) Sym {
	return s.Scr[x-1][y-1]
}

func (s *Screen6x4) LY(x Pos, line Linex) Sym {
	return s.Scr[x-1][line[x-1]-1]
}

func (s *Screen6x4) SetSym(x, y Pos, sym Sym) {
	s.Scr[x-1][y-1] = sym
}

func (s *Screen6x4) SetCol(x Pos, reel []Sym, pos int) {
	for y := range 4 {
		s.Scr[x-1][y] = ReelAt(reel, pos+y)
	}
}

func (s *Screen6x4) ReelSpin(reels Reelx) {
	var x Pos
	for x = 1; x <= 6; x++ {
		var reel = reels.Reel(x)
		var hit = rand.N(len(reel))
		s.SetCol(x, reel, hit)
	}
}

func (s *Screen6x4) SymNum(sym Sym) (n Pos) {
	for x := range 6 {
		var r = s.Scr[x]
		for y := range 4 {
			if r[y] == sym {
				n++
			}
		}
	}
	return
}

func (s *Screen6x4) SymPos(sym Sym) (c Hitx) {
	var x, y, i Pos
	for x = range 6 {
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
