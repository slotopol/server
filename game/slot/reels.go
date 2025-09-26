package slot

import (
	"fmt"
	"math"
)

// Reels for 3-reels slots.
type Reels3x [3][]Sym

// Declare conformity with Reels interface.
var _ Reels = (*Reels3x)(nil)

func (r *Reels3x) Cols() int {
	return 3
}

func (r *Reels3x) Reel(col Pos) []Sym {
	return r[col-1]
}

func (r *Reels3x) Reshuffles() uint64 {
	return uint64(len(r[0])) * uint64(len(r[1])) * uint64(len(r[2]))
}

func (r *Reels3x) String() string {
	return fmt.Sprintf("[%d, %d, %d]", len(r[0]), len(r[1]), len(r[2]))
}

// Reels for 4-reels slots.
type Reels4x [4][]Sym

// Declare conformity with Reels interface.
var _ Reels = (*Reels4x)(nil)

func (r *Reels4x) Cols() int {
	return 4
}

func (r *Reels4x) Reel(col Pos) []Sym {
	return r[col-1]
}

func (r *Reels4x) Reshuffles() uint64 {
	return uint64(len(r[0])) * uint64(len(r[1])) * uint64(len(r[2])) * uint64(len(r[3]))
}

func (r *Reels4x) String() string {
	return fmt.Sprintf("[%d, %d, %d, %d]", len(r[0]), len(r[1]), len(r[2]), len(r[3]))
}

// Reels for 5-reels slots.
type Reels5x [5][]Sym

// Declare conformity with Reels interface.
var _ Reels = (*Reels5x)(nil)

func (r *Reels5x) Cols() int {
	return 5
}

func (r *Reels5x) Reel(col Pos) []Sym {
	return r[col-1]
}

func (r *Reels5x) Reshuffles() uint64 {
	return uint64(len(r[0])) * uint64(len(r[1])) * uint64(len(r[2])) * uint64(len(r[3])) * uint64(len(r[4]))
}

func (r *Reels5x) String() string {
	return fmt.Sprintf("[%d, %d, %d, %d, %d]", len(r[0]), len(r[1]), len(r[2]), len(r[3]), len(r[4]))
}

// Reels for 6-reels slots.
type Reels6x [6][]Sym

// Declare conformity with Reels interface.
var _ Reels = (*Reels6x)(nil)

func (r *Reels6x) Cols() int {
	return 6
}

func (r *Reels6x) Reel(col Pos) []Sym {
	return r[col-1]
}

func (r *Reels6x) Reshuffles() uint64 {
	return uint64(len(r[0])) * uint64(len(r[1])) * uint64(len(r[2])) * uint64(len(r[3])) * uint64(len(r[4])) * uint64(len(r[5]))
}

func (r *Reels6x) String() string {
	return fmt.Sprintf("[%d, %d, %d, %d, %d, %d]", len(r[0]), len(r[1]), len(r[2]), len(r[3]), len(r[4]), len(r[5]))
}

type ReelsMap[T any] map[float64]T

func (m ReelsMap[T]) Clear() {
	clear(m)
}

func (m ReelsMap[T]) FindClosest(mrtp float64) (val T, rtp float64) {
	rtp = -1000 // lets to get first reels from map in any case
	for p, v := range m {
		if math.Abs(mrtp-p) < math.Abs(mrtp-rtp) {
			val, rtp = v, p
		}
	}
	return
}
