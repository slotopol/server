package aztecgold

import (
	"math"
	"math/rand/v2"
)

// aztec pyramids multipliers
var apm = [6]float64{
	10, // wooden pyramid
	12, // stone pyramid
	15, // bronze pyramid
	20, // silver pyramid
	25, // golden pyramid
	51, // diamond pyramid
}

// aztec pyramids frequency
var apf [6]float64

// aztec pyramid probability on max 3 attempts
var app [6]float64

const att = 3 // number of attempts

func init() {
	var s1 float64
	for _, m := range apm {
		s1 += m
	}
	var s2 float64
	for _, m := range apm {
		s2 += s1 / m
	}
	for i, m := range apm {
		apf[i] = s1 / m / s2
	}

	var p6 = (1 - apf[5])
	var p6p = math.Pow(p6, att-1)
	app[0] = apf[0] * p6p
	app[1] = apf[1] * p6p
	app[2] = apf[2] * p6p
	app[3] = apf[3] * p6p
	app[4] = apf[4] * p6p
	app[5] = 1 - math.Pow(p6, att)
}

func getpyramid(p float64) byte {
	var s float64
	for i, f := range apf {
		s += f
		if p <= s {
			return byte(i + 1)
		}
	}
	return 6
}

// Symbols:
//  1 - 1 green diamond
//  2 - 2 green diamonds
//  3 - 3 green diamonds
//  4 - 4 green diamonds
//  5 - 1 blue diamond
//  6 - 2 blue diamonds
//  7 - 3 blue diamonds
//  8 - 4 blue diamonds
//  9 - 1 red diamond
// 10 - 2 red diamonds
// 11 - 3 red diamonds
// 12 - 4 red diamonds
// 13 - snake
// 14 - wild that put all winnings on the line
// 15 - queen
// 16 - king
// 17 - dragon
// 18 - wildcard rescuing the snake, but not giving win

type Cell struct {
	Mult float64 `json:"mult" yaml:"mult" xml:"mult,attr"`
	Sym  byte    `json:"sym" yaml:"sym" xml:"sym,attr"`
}

type Row struct {
	Sel  Cell    `json:"sel" yaml:"sel" xml:"sel"`
	Next [4]Cell `json:"next" yaml:"next" xml:"next>cell"`
}

type Pyramid struct {
	Mult float64 `json:"mult" yaml:"mult" xml:"mult,attr"`
	Type byte    `json:"type" yaml:"type" xml:"type,attr"`
}

type Bonus struct {
	Strike [att]Pyramid `json:"strike" yaml:"strike" xml:"strike"`
	Room   [5]Row       `json:"room" yaml:"room" xml:"room>row"`
}

var Room = [5][]Cell{
	{ // row #1
		{10, 1},
		{11, 1},
		{12, 2},
		{13, 2},
		{14, 3},
		{15, 3},
		{16, 4},
		{10, 1},
		{11, 1},
		{12, 2},
		{13, 2},
		{14, 3},
		{15, 3},
		{10, 1},
		{11, 1},
		{12, 2},
		{13, 2},
		{14, 3},
		{10, 1},
		{11, 1},
		{12, 2},
		{13, 2},
		{0, 14},
	},
	{ // row #2
		{15, 5},
		{16, 5},
		{17, 6},
		{18, 6},
		{19, 7},
		{20, 8},
		{21, 8},
		{15, 5},
		{16, 5},
		{17, 6},
		{18, 6},
		{19, 7},
		{20, 8},
		{15, 5},
		{16, 5},
		{17, 6},
		{18, 6},
		{19, 7},
		{15, 5},
		{16, 5},
		{17, 6},
		{18, 6},
		{0, 14},
	},
	{ // row #3
		{20, 9},
		{21, 9},
		{22, 10},
		{23, 10},
		{24, 11},
		{25, 11},
		{26, 12},
		{20, 9},
		{21, 9},
		{22, 10},
		{23, 10},
		{24, 11},
		{25, 11},
		{20, 9},
		{21, 9},
		{22, 10},
		{23, 10},
		{24, 11},
		{20, 9},
		{21, 9},
		{22, 10},
		{23, 10},
		{0, 14},
	},
	{ // row #4
		{50, 15},
		{50, 15},
		{100, 16},
		{0, 13},
		{0, 13},
		{0, 13},
		{0, 13},
		{0, 13},
		{0, 13},
		{0, 13},
		{0, 13},
		{0, 13},
		{0, 18},
	},
	{ // row #5
		{250, 17},
		{0, 13},
		{0, 13},
		{0, 13},
		{0, 13},
		{0, 13},
		{0, 13},
		{0, 13},
		{0, 13},
		{0, 13},
		{0, 13},
		{0, 13},
		{0, 13},
		{0, 13},
		{0, 13},
		{0, 13},
		{0, 13},
		{0, 13},
		{0, 13},
	},
}

func AztecPyramidSpawn(bet float64) (any, float64) {
	var res Bonus
	var cash float64
	var enter bool
	for i := range att {
		var pyr = getpyramid(rand.Float64())
		cash = apm[pyr-1]
		res.Strike[i].Type = pyr
		res.Strike[i].Mult = cash
		if pyr == 6 {
			enter = true
			break
		}
	}
	if !enter {
		return res, bet * cash
	}

	for ri := range 5 {
		var perm = rand.Perm(len(Room[ri]))
		var row = &res.Room[ri]
		row.Sel = Room[ri][perm[0]]
		if ri != 3 {
			row.Next[0] = Room[ri][perm[1]]
			row.Next[1] = Room[ri][perm[2]]
			row.Next[2] = Room[ri][perm[3]]
			row.Next[3] = Room[ri][perm[4]]
		} else {
			if row.Sel.Sym == 15 { // selected queen
				row.Next[0] = Room[ri][3] // put snake
			} else {
				row.Next[0] = Room[ri][1] // put queen
			}
			if row.Sel.Sym == 16 { // selected king
				row.Next[1] = Room[ri][3] // put snake
			} else {
				row.Next[1] = Room[ri][2] // put king
			}
			row.Next[2] = Room[ri][perm[4]]
			row.Next[3] = Room[ri][perm[5]]
			rand.Shuffle(4, func(i int, j int) {
				row.Next[i], row.Next[j] = row.Next[j], row.Next[i]
			})
		}
		cash += row.Sel.Mult
		if row.Sel.Sym == 13 { // snake was choosed
			break
		} else if row.Sel.Sym == 14 { // wild was choosed
			for _, c := range row.Next {
				cash += c.Mult
			}
		}
	}
	return res, bet * cash
}
