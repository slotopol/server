package slot

import (
	"encoding/json"
	"math/rand/v2"
	"sync"
)

// Screen contains symbols rectangle of the slot game.
// It can be with dimensions 3x1, 3x3, 5x3, 5x4 or others.
// (1 ,1) symbol is on left top corner.
type Screen interface {
	Dim() (Pos, Pos)                   // returns screen dimensions
	At(x, y Pos) Sym                   // returns symbol at position (x, y), starts from (1, 1)
	Pos(x Pos, line Linex) Sym         // returns symbol at position (x, line(x)), starts from (1, 1)
	Set(x, y Pos, sym Sym)             // setup symbol at given position
	SetCol(x Pos, reel []Sym, pos int) // setup column on screen with given reel at given position
	Spin(reels Reels)                  // fill the screen with random hits on those reels
	SymNum(sym Sym) (n Pos)            // returns number of symbols on the screen that can repeats on reel
	ScatNum(scat Sym) (n Pos)          // returns number of scatters on the screen
	ScatNumOdd(scat Sym) (n Pos)       // returns number of scatters on the screen on odd reels
	ScatNumCont(scat Sym) (n Pos)      // returns number of continuous scatters on the screen
	ScatPos(scat Sym) Linex            // returns line with scatters positions on the screen
	ScatPosOdd(scat Sym) Linex         // returns line with scatters positions on the screen on odd reels
	ScatPosCont(scat Sym) Linex        // returns line with continuous scatters positions on the screen
	FillSym() Sym                      // returns symbol that filled whole screen, or 0
	Free()                             // put object to pool
}

type Screenx struct {
	sx, sy Pos
	data   [40]Sym
}

var poolsx = sync.Pool{
	New: func() any {
		return &Screen3x3{}
	},
}

func NewScreenx(sx, sy Pos) *Screenx {
	var s = poolsx.Get().(*Screenx)
	s.sx, s.sy = sx, sy
	return s
}

func (s *Screenx) Free() {
	poolsx.Put(s)
}

func (s *Screenx) Len() int {
	for i := 39; i >= 0; i-- {
		if s.data[i] > 0 {
			return i + 1
		}
	}
	return 0
}

func (s *Screenx) UpdateDim() (sx, sy Pos) {
	switch s.Len() {
	case 3 * 1:
		sx, sy = 3, 1
	case 3 * 3:
		sx, sy = 3, 3
	case 5 * 3:
		sx, sy = 5, 3
	case 5 * 4:
		sx, sy = 5, 4
	case 6 * 4:
		sx, sy = 6, 4
	}
	s.sx, s.sy = sx, sy
	return
}

func (s *Screenx) Dim() (Pos, Pos) {
	return s.sx, s.sy
}

func (s *Screenx) At(x, y Pos) Sym {
	return s.data[(x-1)*s.sy+y-1]
}

func (s *Screenx) Pos(x Pos, line Linex) Sym {
	return s.data[(x-1)*s.sy+line[x-1]-1]
}

func (s *Screenx) Set(x, y Pos, sym Sym) {
	s.data[(x-1)*s.sy+y-1] = sym
}

func (s *Screenx) SetCol(x Pos, reel []Sym, pos int) {
	var i = (x - 1) * s.sy
	for y := range s.sy {
		s.data[i+y] = reel[(pos+int(y))%len(reel)]
	}
}

func (s *Screenx) Spin(reels Reels) {
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

func (s *Screenx) ScatNum(scat Sym) (n Pos) {
	for i := range s.sx * s.sy {
		if s.data[i] == scat {
			n++
		}
	}
	return
}

func (s *Screenx) ScatNumOdd(scat Sym) (n Pos) {
	var x, y, i Pos
loopx:
	for x = 0; x < s.sx; x += 2 {
		i = x * s.sy
		for y = range s.sy {
			if s.data[i+y] == scat {
				n++
				continue loopx
			}
		}
	}
	return
}

func (s *Screenx) ScatNumCont(scat Sym) (n Pos) {
	var x, y, i Pos
loopx:
	for x = range s.sx {
		i = x * s.sy
		for y = range s.sy {
			if s.data[i+y] == scat {
				n++
				continue loopx
			}
		}
		break // scatters should be continuous
	}
	return
}

func (s *Screenx) ScatPos(scat Sym) (l Linex) {
	for i := range s.sx * s.sy {
		if s.data[i] == scat {
			l[i/s.sy+1] = i%s.sy + 1
		}
	}
	return
}

func (s *Screenx) ScatPosOdd(scat Sym) (l Linex) {
	var x, y, i Pos
loopx:
	for x = 0; x < s.sx; x += 2 {
		i = x * s.sy
		for y = range s.sy {
			if s.data[i+y] == scat {
				l[x] = y + 1
				continue loopx
			}
		}
	}
	return
}

func (s *Screenx) ScatPosCont(scat Sym) (l Linex) {
	var x, y, i Pos
loopx:
	for x = range s.sx {
		i = x * s.sy
		for y = range s.sy {
			if s.data[i+y] == scat {
				l[x] = y + 1
				continue loopx
			}
		}
		break // scatters should be continuous
	}
	return
}

func (s *Screenx) FillSym() (sym Sym) {
	sym = s.data[0]
	for i := range s.sx * s.sy {
		if s.data[i] != sym {
			return 0
		}
	}
	return
}

func (s *Screenx) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.data[:s.Len()])
}

func (s *Screenx) UnmarshalJSON(b []byte) (err error) {
	if err = json.Unmarshal(b, &s.data); err != nil {
		return
	}
	s.UpdateDim()
	return
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

func (s *Screen3x3) Dim() (Pos, Pos) {
	return 3, 3
}

func (s *Screen3x3) At(x, y Pos) Sym {
	return s[x-1][y-1]
}

func (s *Screen3x3) Pos(x Pos, line Linex) Sym {
	return s[x-1][line[x-1]-1]
}

func (s *Screen3x3) Set(x, y Pos, sym Sym) {
	s[x-1][y-1] = sym
}

func (s *Screen3x3) SetCol(x Pos, reel []Sym, pos int) {
	for y := range 3 {
		s[x-1][y] = reel[(pos+y)%len(reel)]
	}
}

func (s *Screen3x3) Spin(reels Reels) {
	var x Pos
	for x = 1; x <= 3; x++ {
		var reel = reels.Reel(x)
		var hit = rand.N(len(reel))
		s.SetCol(x, reel, hit)
	}
}

func (s *Screen3x3) SymNum(sym Sym) (n Pos) {
	for x := range 3 {
		for y := range 3 {
			if s[x][y] == sym {
				n++
			}
		}
	}
	return
}

func (s *Screen3x3) ScatNum(scat Sym) (n Pos) {
	var x Pos
	for x = range 3 {
		var r = s[x]
		if r[0] == scat || r[1] == scat || r[2] == scat {
			n++
		}
	}
	return
}

func (s *Screen3x3) ScatNumOdd(scat Sym) (n Pos) {
	var x Pos
	for x = 0; x < 3; x += 2 {
		var r = s[x]
		if r[0] == scat || r[1] == scat || r[2] == scat {
			n++
		}
	}
	return
}

func (s *Screen3x3) ScatNumCont(scat Sym) (n Pos) {
	var x Pos
	for x = 0; x < 3; x++ {
		var r = s[x]
		if r[0] == scat || r[1] == scat || r[2] == scat {
			n++
		} else {
			break // scatters should be continuous
		}
	}
	return
}

func (s *Screen3x3) ScatPos(scat Sym) (l Linex) {
	for x := range 3 {
		var r = s[x]
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

func (s *Screen3x3) ScatPosOdd(scat Sym) (l Linex) {
	for x := 0; x < 3; x += 2 {
		var r = s[x]
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

func (s *Screen3x3) ScatPosCont(scat Sym) (l Linex) {
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
	return
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

func (s *Screen5x3) Dim() (Pos, Pos) {
	return 5, 3
}

func (s *Screen5x3) At(x, y Pos) Sym {
	return s[x-1][y-1]
}

func (s *Screen5x3) Pos(x Pos, line Linex) Sym {
	return s[x-1][line[x-1]-1]
}

func (s *Screen5x3) Set(x, y Pos, sym Sym) {
	s[x-1][y-1] = sym
}

func (s *Screen5x3) SetCol(x Pos, reel []Sym, pos int) {
	for y := range 3 {
		s[x-1][y] = reel[(pos+y)%len(reel)]
	}
}

func (s *Screen5x3) Spin(reels Reels) {
	var x Pos
	for x = 1; x <= 5; x++ {
		var reel = reels.Reel(x)
		var hit = rand.N(len(reel))
		s.SetCol(x, reel, hit)
	}
}

func (s *Screen5x3) SymNum(sym Sym) (n Pos) {
	for x := range 5 {
		for y := range 3 {
			if s[x][y] == sym {
				n++
			}
		}
	}
	return
}

func (s *Screen5x3) ScatNum(scat Sym) (n Pos) {
	for x := range 5 {
		var r = s[x]
		if r[0] == scat || r[1] == scat || r[2] == scat {
			n++
		}
	}
	return
}

func (s *Screen5x3) ScatNumOdd(scat Sym) (n Pos) {
	for x := 0; x < 5; x += 2 {
		var r = s[x]
		if r[0] == scat || r[1] == scat || r[2] == scat {
			n++
		}
	}
	return
}

func (s *Screen5x3) ScatNumCont(scat Sym) (n Pos) {
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

func (s *Screen5x3) ScatPos(scat Sym) (l Linex) {
	for x := range 5 {
		var r = s[x]
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

func (s *Screen5x3) ScatPosOdd(scat Sym) (l Linex) {
	for x := 0; x < 5; x += 2 {
		var r = s[x]
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

func (s *Screen5x3) ScatPosCont(scat Sym) (l Linex) {
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
	return
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

func (s *Screen5x4) Dim() (Pos, Pos) {
	return 5, 4
}

func (s *Screen5x4) At(x, y Pos) Sym {
	return s[x-1][y-1]
}

func (s *Screen5x4) Pos(x Pos, line Linex) Sym {
	return s[x-1][line[x-1]-1]
}

func (s *Screen5x4) Set(x, y Pos, sym Sym) {
	s[x-1][y-1] = sym
}

func (s *Screen5x4) SetCol(x Pos, reel []Sym, pos int) {
	for y := range 4 {
		s[x-1][y] = reel[(pos+y)%len(reel)]
	}
}

func (s *Screen5x4) Spin(reels Reels) {
	var x Pos
	for x = 1; x <= 5; x++ {
		var reel = reels.Reel(x)
		var hit = rand.N(len(reel))
		s.SetCol(x, reel, hit)
	}
}

func (s *Screen5x4) SymNum(sym Sym) (n Pos) {
	for x := range 5 {
		for y := range 4 {
			if s[x][y] == sym {
				n++
			}
		}
	}
	return
}

func (s *Screen5x4) ScatNum(scat Sym) (n Pos) {
	for x := range 5 {
		var r = s[x]
		if r[0] == scat || r[1] == scat || r[2] == scat || r[3] == scat {
			n++
		}
	}
	return
}

func (s *Screen5x4) ScatNumOdd(scat Sym) (n Pos) {
	for x := 0; x < 5; x += 2 {
		var r = s[x]
		if r[0] == scat || r[1] == scat || r[2] == scat || r[3] == scat {
			n++
		}
	}
	return
}

func (s *Screen5x4) ScatNumCont(scat Sym) (n Pos) {
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

func (s *Screen5x4) ScatPos(scat Sym) (l Linex) {
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
		}
	}
	return
}

func (s *Screen5x4) ScatPosOdd(scat Sym) (l Linex) {
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
		}
	}
	return
}

func (s *Screen5x4) ScatPosCont(scat Sym) (l Linex) {
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
	return
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
