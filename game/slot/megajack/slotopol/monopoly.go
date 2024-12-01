package slotopol

import "math/rand/v2"

type MonCell struct {
	Mult float64 `json:"mult" yaml:"mult" xml:"mult,attr"` // bet multiplier
	Jump int     `json:"jump" yaml:"jump" xml:"jump,attr"` // jump position, or 0 if no jump
	Dice bool    `json:"dice" yaml:"dice" xml:"dice,attr"` // is here multiply on dice value
}

// count = 279936, sum = 80231330, avr = 286.60597422268, zerocount = 1, p(zero) = 0.00035722450845908%
// variance = 15529.19650266, sigma = 124.61619679103, limits = 161.98977743165...411.22217101371
var Monopoly = []MonCell{
	{5, 0, false},   //  1, x5
	{15, 0, false},  //  2, x15
	{15, 0, false},  //  3, x15
	{30, 0, false},  //  4, x30
	{30, 0, false},  //  5, x30
	{10, 0, true},   //  6, x10
	{0, 0, false},   //  7,  0
	{0, 1, false},   //  8,  0
	{50, 0, false},  //  9, x50
	{50, 0, false},  // 10, x50
	{0, 0, false},   // 11,  0
	{20, 0, true},   // 12, x20
	{0, 4, false},   // 13,  0
	{80, 0, false},  // 14, x80
	{80, 0, false},  // 15, x80
	{30, 0, true},   // 16, x30
	{0, 7, false},   // 17,  0
	{120, 0, false}, // 18, x120
	{120, 0, false}, // 19, x120
	{200, 0, false}, // 20, x200
}

type WinMonCell struct {
	MonCell `yaml:",inline"`
	Pos     int     `json:"pos" yaml:"pos" xml:"pos,attr"` // cell number, starts from 1
	Pay     float64 `json:"pay" yaml:"pay" xml:"pay,attr"` // pay by this cell
}

func MonopolySpawn(bet float64) (any, float64) {
	var res [7]WinMonCell
	var cash float64

	var pos = 0
	for i := range res {
		var dice = rand.N(6) + 1
		pos = (pos + dice) % len(Monopoly)
		var mult float64
		if Monopoly[pos].Jump > 0 {
			mult = Monopoly[Monopoly[pos].Jump-1].Mult
		} else {
			mult = Monopoly[pos].Mult
		}
		if Monopoly[pos].Dice {
			mult *= float64(dice)
		}
		res[i].Mult = mult
		res[i].Jump = Monopoly[pos].Jump
		res[i].Dice = Monopoly[pos].Dice
		res[i].Pos = pos + 1
		res[i].Pay = bet * mult
		cash += mult
		if Monopoly[pos].Jump > 0 {
			pos = Monopoly[pos].Jump - 1
		}
	}
	if cash == 0 {
		cash = 5000
	}
	return res, bet * cash
}
