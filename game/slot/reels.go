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
func (rx Reelx) Reel(col Pos) []Sym {
	return rx[col-1]
}

// Returns total number of reshuffles.
func (rx Reelx) Reshuffles() uint64 {
	var res uint64 = 1
	for _, reel := range rx {
		res *= uint64(len(reel))
	}
	return res
}

// Returns seed position on reels from array of positions on each reel.
func (rx Reelx) PosToSeed(pos []int) (seed uint64) {
	var m uint64 = 1
	for i := len(rx) - 1; i >= 0; i-- {
		seed += uint64(pos[i]) * m
		m *= uint64(len(rx[i]))
	}
	return
}

// Fills the array of positions on each reel from given seed.
func (rx Reelx) SeedToPos(seed uint64, pos []int) {
	var cur = seed
	var rl uint64
	for i := len(rx) - 1; i >= 0; i-- {
		rl = uint64(len(rx[i]))
		pos[i] = int(cur % rl)
		cur /= rl
	}
}

func (rx Reelx) String() string {
	var b strings.Builder
	b.WriteString("[")
	for i, reel := range rx {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(strconv.Itoa(len(reel)))
	}
	b.WriteString("]")
	return b.String()
}

func (rx Reelx) Clear() {
	clear(rx)
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
