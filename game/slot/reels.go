package slot

import (
	"fmt"
	"math"
	"os"

	"gopkg.in/yaml.v3"
)

type Reels interface {
	Cols() int          // returns number of columns
	Reel(col Pos) []Sym // returns reel at given column, index from
	Reshuffles() uint64 // returns total number of reshuffles
	fmt.Stringer
}

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

func (r *Reels3x) Clear() {
	clear(r[:])
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

func (r *Reels4x) Clear() {
	clear(r[:])
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

func (r *Reels5x) Clear() {
	clear(r[:])
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

func (r *Reels6x) Clear() {
	clear(r[:])
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

func ReadObj[T any](b []byte) (obj T) {
	var err error
	if err = yaml.Unmarshal(b, &obj); err != nil {
		panic(err)
	}
	return
}

func LoadObj[T any](fpath string) (obj T) {
	var b, err = os.ReadFile(fpath)
	if err != nil {
		panic(err)
	}
	return ReadObj[T](b)
}

func ReadMap[T any](b []byte) (rm ReelsMap[T]) {
	var err error
	if err = yaml.Unmarshal(b, &rm); err != nil {
		panic(err)
	}
	return
}

func DataRouter[T any](fpath string) (rm ReelsMap[T]) {
	var b, err = os.ReadFile(fpath)
	if err != nil {
		panic(err)
	}
	return ReadMap[T](b)
}
