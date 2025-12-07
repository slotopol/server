package slot

import (
	"math"
	"strconv"
	"strings"
)

// ReelAt returns symbol on the reel at true cyclic position.
// Incoming "pos" can be greater than reel length, or can be negative.
func ReelAt(reel []Sym, pos int) Sym {
	var n = len(reel)
	return reel[(n+pos%n)%n]
}

type Reelx [][]Sym

// Returns reel at given column, index is 1-based.
func (r Reelx) Reel(col Pos) []Sym {
	return r[col-1]
}

// Returns total number of reshuffles.
func (r Reelx) Reshuffles() uint64 {
	var res uint64 = 1
	for _, reel := range r {
		res *= uint64(len(reel))
	}
	return res
}

func (r Reelx) String() string {
	var b strings.Builder
	b.WriteString("[")
	for i, reel := range r {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(strconv.Itoa(len(reel)))
	}
	b.WriteString("]")
	return b.String()
}

func (r Reelx) Clear() {
	clear(r)
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
