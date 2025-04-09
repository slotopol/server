package cashfarm

import (
	"math/rand/v2"
)

type CT int

const (
	pay CT = iota + 1
	mult
	buck
)

type Cell int

func (c Cell) Type() CT {
	if c == 0 {
		return buck
	} else if c == 1 || c == -1 {
		return mult
	} else {
		return pay
	}
}

func (c Cell) IsUp() bool {
	return c < 0
}

func (c Cell) Mult() float64 {
	if c != 1 && c != -1 {
		return 0
	}
	return 1
}

func (c Cell) Pay() float64 {
	if c == -1 || c == 0 || c == 1 {
		return 0
	}
	if c < 0 {
		return float64(-c)
	}
	return float64(c)
}

type Bonus [5][]Cell

var bonreel1 = []Cell{2, 2, 3, 4, 5, 6, 8, 10, 14}
var bonreel4 = []Cell{14, 14, 20, 24}
var bonreel5 = []Cell{20, 24, 40, 50}

func br1() Cell {
	return bonreel1[rand.N(Cell(len(bonreel1)))]
}

func br4() Cell {
	return bonreel4[rand.N(Cell(len(bonreel4)))]
}

func br5() Cell {
	return bonreel5[rand.N(Cell(len(bonreel5)))]
}

func NewBonus() (bon *Bonus) {
	return &Bonus{
		[]Cell{br1(), br1(), br1(), -br1(), -br1(), -br1()},
		[]Cell{br1(), br1(), -br1(), -1, 0},
		[]Cell{br1(), -br1(), -1, 0},
		[]Cell{br4(), -1, 0},
		[]Cell{br5(), 0},
	}
}

func (b *Bonus) Shuffle() {
	for lev := range b {
		rand.Shuffle(len(b[lev]), func(i int, j int) {
			b[lev][i], b[lev][j] = b[lev][j], b[lev][i]
		})
	}
}

func (b *Bonus) Calc() (p, m float64) {
	p, m = 0, 1
	for lev := range b {
		for _, c := range b[lev] {
			if c == 0 {
				return
			} else if c == 1 || c == -1 {
				m++
			} else if c < 0 {
				p -= float64(-c)
			} else {
				p += float64(c)
			}
			if c < 0 {
				break
			}
		}
	}
	return
}

func FarmSpawn(bet float64) (any, float64, float64) {
	var res = NewBonus()
	res.Shuffle()
	var p, m = res.Calc()
	return res, bet * p, m
}
