package game

import (
	"math/rand/v2"
	"sync"
)

// Screen contains symbols rectangle of the slot game.
// It can be with dimensions 3x1, 3x3, 5x3, 5x4 or others.
// (1 ,1) symbol is on left top corner.
type Screen interface {
	Dim() (int, int)                   // returns screen dimensions
	At(x int, y int) Sym               // returns symbol at position (x, y), starts from (1, 1)
	Pos(x int, line Line) Sym          // returns symbol at position (x, line(x)), starts from (1, 1)
	Set(x int, y int, sym Sym)         // setup symbol at given position
	SetCol(x int, reel []Sym, pos int) // setup column on screen with given reel at given position
	Spin(reels Reels)                  // fill the screen with random hits on those reels
	ScatNum(scat Sym) (n int)          // returns number of scatters on the screen
	ScatNumOdd(scat Sym) (n int)       // returns number of scatters on the screen on odd reels
	ScatNumCont(scat Sym) (n int)      // returns number of continuous scatters on the screen
	ScatPos(scat Sym) Line             // returns line with scatters positions on the screen
	ScatPosOdd(scat Sym) Line          // returns line with scatters positions on the screen on odd reels
	ScatPosCont(scat Sym) Line         // returns line with continuous scatters positions on the screen
	FillSym() Sym                      // returns symbol that filled whole screen, or 0
	Free()                             // put object to pool
}

// Screen for 3x3 slots.
type Screen3x3 [3][3]Sym

var pools3x3 = sync.Pool{
	New: func() any {
		return &Screen3x3{}
	},
}

func NewScreen3x3() *Screen3x3 {
	return pools3x3.Get().(*Screen3x3)
}

func (s *Screen3x3) Free() {
	pools3x3.Put(s)
}

func (s *Screen3x3) Dim() (int, int) {
	return 3, 3
}

func (s *Screen3x3) At(x int, y int) Sym {
	return s[x-1][y-1]
}

func (s *Screen3x3) Pos(x int, line Line) Sym {
	return s[x-1][line.At(x)-1]
}

func (s *Screen3x3) Set(x int, y int, sym Sym) {
	s[x-1][y-1] = sym
}

func (s *Screen3x3) SetCol(x int, reel []Sym, pos int) {
	for y := range 3 {
		s[x-1][y] = reel[(pos+y)%len(reel)]
	}
}

func (s *Screen3x3) Spin(reels Reels) {
	for x := 1; x <= 3; x++ {
		var reel = reels.Reel(x)
		var hit = rand.N(len(reel))
		s.SetCol(x, reel, hit)
	}
}

func (s *Screen3x3) ScatNum(scat Sym) (n int) {
	for x := range 3 {
		var r = s[x]
		if r[0] == scat || r[1] == scat || r[2] == scat {
			n++
		}
	}
	return
}

func (s *Screen3x3) ScatNumOdd(scat Sym) (n int) {
	for x := 0; x < 3; x += 2 {
		var r = s[x]
		if r[0] == scat || r[1] == scat || r[2] == scat {
			n++
		}
	}
	return
}

func (s *Screen3x3) ScatNumCont(scat Sym) (n int) {
	for x := 0; x < 3; x++ {
		var r = s[x]
		if r[0] == scat || r[1] == scat || r[2] == scat {
			n++
		} else {
			break // scatters should be continuous
		}
	}
	return
}

func (s *Screen3x3) ScatPos(scat Sym) Line {
	var l = NewLine3x()
	for x := range 3 {
		var r = s[x]
		if r[0] == scat {
			l[x] = 1
		} else if r[1] == scat {
			l[x] = 2
		} else if r[2] == scat {
			l[x] = 3
		} else {
			l[x] = 0
		}
	}
	return l
}

func (s *Screen3x3) ScatPosOdd(scat Sym) Line {
	var l = NewLine3x()
	for x := 0; x < 3; x += 2 {
		var r = s[x]
		if r[0] == scat {
			l[x] = 1
		} else if r[1] == scat {
			l[x] = 2
		} else if r[2] == scat {
			l[x] = 3
		} else {
			l[x] = 0
		}
	}
	l[1] = 0
	return l
}

func (s *Screen3x3) ScatPosCont(scat Sym) Line {
	var l = NewLine3x()
	var x int
	for x = 0; x < 3; x++ {
		var r = s[x]
		if r[0] == scat {
			l[x] = 1
		} else if r[1] == scat {
			l[x] = 2
		} else if r[2] == scat {
			l[x] = 3
		} else {
			break // scatters should be continuous
		}
	}
	for ; x < 3; x++ {
		l[x] = 0
	}
	return l
}

func (s *Screen3x3) FillSym() Sym {
	var sym = s[0][0]
	if s[1][0] == sym && s[2][0] == sym &&
		s[0][1] == sym && s[1][1] == sym && s[2][1] == sym &&
		s[0][2] == sym && s[1][2] == sym && s[2][2] == sym {
		return sym
	}
	return 0
}

// Screen for 5x3 slots.
type Screen5x3 [5][3]Sym

var pools5x3 = sync.Pool{
	New: func() any {
		return &Screen5x3{}
	},
}

func NewScreen5x3() *Screen5x3 {
	return pools5x3.Get().(*Screen5x3)
}

func (s *Screen5x3) Free() {
	pools5x3.Put(s)
}

func (s *Screen5x3) Dim() (int, int) {
	return 5, 3
}

func (s *Screen5x3) At(x int, y int) Sym {
	return s[x-1][y-1]
}

func (s *Screen5x3) Pos(x int, line Line) Sym {
	return s[x-1][line.At(x)-1]
}

func (s *Screen5x3) Set(x int, y int, sym Sym) {
	s[x-1][y-1] = sym
}

func (s *Screen5x3) SetCol(x int, reel []Sym, pos int) {
	for y := range 3 {
		s[x-1][y] = reel[(pos+y)%len(reel)]
	}
}

func (s *Screen5x3) Spin(reels Reels) {
	for x := 1; x <= 5; x++ {
		var reel = reels.Reel(x)
		var hit = rand.N(len(reel))
		s.SetCol(x, reel, hit)
	}
}

func (s *Screen5x3) ScatNum(scat Sym) (n int) {
	for x := range 5 {
		var r = s[x]
		if r[0] == scat || r[1] == scat || r[2] == scat {
			n++
		}
	}
	return
}

func (s *Screen5x3) ScatNumOdd(scat Sym) (n int) {
	for x := 0; x < 5; x += 2 {
		var r = s[x]
		if r[0] == scat || r[1] == scat || r[2] == scat {
			n++
		}
	}
	return
}

func (s *Screen5x3) ScatNumCont(scat Sym) (n int) {
	for x := 0; x < 5; x++ {
		var r = s[x]
		if r[0] == scat || r[1] == scat || r[2] == scat {
			n++
		} else {
			break // scatters should be continuous
		}
	}
	return
}

func (s *Screen5x3) ScatPos(scat Sym) Line {
	var l = NewLine5x()
	for x := range 5 {
		var r = s[x]
		if r[0] == scat {
			l[x] = 1
		} else if r[1] == scat {
			l[x] = 2
		} else if r[2] == scat {
			l[x] = 3
		} else {
			l[x] = 0
		}
	}
	return l
}

func (s *Screen5x3) ScatPosOdd(scat Sym) Line {
	var l = NewLine5x()
	for x := 0; x < 5; x += 2 {
		var r = s[x]
		if r[0] == scat {
			l[x] = 1
		} else if r[1] == scat {
			l[x] = 2
		} else if r[2] == scat {
			l[x] = 3
		} else {
			l[x] = 0
		}
	}
	l[1], l[3] = 0, 0
	return l
}

func (s *Screen5x3) ScatPosCont(scat Sym) Line {
	var l = NewLine5x()
	var x int
	for x = 0; x < 5; x++ {
		var r = s[x]
		if r[0] == scat {
			l[x] = 1
		} else if r[1] == scat {
			l[x] = 2
		} else if r[2] == scat {
			l[x] = 3
		} else {
			break // scatters should be continuous
		}
	}
	for ; x < 5; x++ {
		l[x] = 0
	}
	return l
}

func (s *Screen5x3) FillSym() Sym {
	var sym = s[0][0]
	if s[1][0] == sym && s[2][0] == sym && s[3][0] == sym && s[4][0] == sym &&
		s[0][1] == sym && s[1][1] == sym && s[2][1] == sym && s[3][1] == sym && s[4][1] == sym &&
		s[0][2] == sym && s[1][2] == sym && s[2][2] == sym && s[3][2] == sym && s[4][2] == sym {
		return sym
	}
	return 0
}

// Screen for 5x4 slots.
type Screen5x4 [5][4]Sym

var pools5x4 = sync.Pool{
	New: func() any {
		return &Screen5x4{}
	},
}

func NewScreen5x4() *Screen5x4 {
	return pools5x4.Get().(*Screen5x4)
}

func (s *Screen5x4) Free() {
	pools5x4.Put(s)
}

func (s *Screen5x4) Dim() (int, int) {
	return 5, 4
}

func (s *Screen5x4) At(x int, y int) Sym {
	return s[x-1][y-1]
}

func (s *Screen5x4) Pos(x int, line Line) Sym {
	return s[x-1][line.At(x)-1]
}

func (s *Screen5x4) Set(x int, y int, sym Sym) {
	s[x-1][y-1] = sym
}

func (s *Screen5x4) SetCol(x int, reel []Sym, pos int) {
	for y := range 4 {
		s[x-1][y] = reel[(pos+y)%len(reel)]
	}
}

func (s *Screen5x4) Spin(reels Reels) {
	for x := 1; x <= 5; x++ {
		var reel = reels.Reel(x)
		var hit = rand.N(len(reel))
		s.SetCol(x, reel, hit)
	}
}

func (s *Screen5x4) ScatNum(scat Sym) (n int) {
	for x := range 5 {
		var r = s[x]
		if r[0] == scat || r[1] == scat || r[2] == scat || r[3] == scat {
			n++
		}
	}
	return
}

func (s *Screen5x4) ScatNumOdd(scat Sym) (n int) {
	for x := 0; x < 5; x += 2 {
		var r = s[x]
		if r[0] == scat || r[1] == scat || r[2] == scat || r[3] == scat {
			n++
		}
	}
	return
}

func (s *Screen5x4) ScatNumCont(scat Sym) (n int) {
	for x := 0; x < 5; x++ {
		var r = s[x]
		if r[0] == scat || r[1] == scat || r[2] == scat || r[3] == scat {
			n++
		} else {
			break // scatters should be continuous
		}
	}
	return
}

func (s *Screen5x4) ScatPos(scat Sym) Line {
	var l = NewLine5x()
	for x := range 5 {
		var r = s[x]
		if r[0] == scat {
			l[x] = 1
		} else if r[1] == scat {
			l[x] = 2
		} else if r[2] == scat {
			l[x] = 3
		} else if r[3] == scat {
			l[x] = 4
		} else {
			l[x] = 0
		}
	}
	return l
}

func (s *Screen5x4) ScatPosOdd(scat Sym) Line {
	var l = NewLine5x()
	for x := 0; x < 5; x += 2 {
		var r = s[x]
		if r[0] == scat {
			l[x] = 1
		} else if r[1] == scat {
			l[x] = 2
		} else if r[2] == scat {
			l[x] = 3
		} else if r[3] == scat {
			l[x] = 4
		} else {
			l[x] = 0
		}
	}
	l[1], l[3] = 0, 0
	return l
}

func (s *Screen5x4) ScatPosCont(scat Sym) Line {
	var l = NewLine5x()
	var x int
	for x = 0; x < 5; x++ {
		var r = s[x]
		if r[0] == scat {
			l[x] = 1
		} else if r[1] == scat {
			l[x] = 2
		} else if r[2] == scat {
			l[x] = 3
		} else if r[3] == scat {
			l[x] = 4
		} else {
			break // scatters should be continuous
		}
	}
	for ; x < 5; x++ {
		l[x] = 0
	}
	return l
}

func (s *Screen5x4) FillSym() Sym {
	var sym = s[0][0]
	if s[1][0] == sym && s[2][0] == sym && s[3][0] == sym && s[4][0] == sym &&
		s[0][1] == sym && s[1][1] == sym && s[2][1] == sym && s[3][1] == sym && s[4][1] == sym &&
		s[0][2] == sym && s[1][2] == sym && s[2][2] == sym && s[3][2] == sym && s[4][2] == sym &&
		s[0][3] == sym && s[1][3] == sym && s[2][3] == sym && s[3][3] == sym && s[4][3] == sym {
		return sym
	}
	return 0
}
