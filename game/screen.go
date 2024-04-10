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
	SetCol(x int, reel []Sym, pos int) // setup column on screen with given reel at given position
	Spin(reels Reels)                  // fill the screen with random hits on those reels
	ScatNum(scat Sym) (n int)          // returns number of scatters on the screen
	ScatNumOdd(scat Sym) (n int)       // returns number of scatters on the screen on odd reels
	ScatNumCont(scat Sym) (n int)      // returns number of continuous scatters on the screen
	ScatPos(scat Sym) Line             // returns line with scatters positions on the screen
	ScatPosOdd(scat Sym) Line          // returns line with scatters positions on the screen on odd reels
	ScatPosCont(scat Sym) Line         // returns line with continuous scatters positions on the screen
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
