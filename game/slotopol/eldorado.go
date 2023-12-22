package slotopol

import "math/rand"

// Average gain = 106
var Eldorado = []int{
	1000, // 1
	10,   // 2
	25,   // 3
	100,  // 4
	10,   // 5
	50,   // 6
	300,  // 7
	10,   // 8
	50,   // 9
	200,  // 10
	10,   // 11
	25,   // 12
	100,  // 13
	25,   // 14
	10,   // 15
	100,  // 16
	10,   // 17
	50,   // 18
	10,   // 19
	25,   // 20
}

type WinElSeg struct {
	Mult int `json:"mult" yaml:"mult" xml:"mult,attr"` // bet multiplier
	Pos  int `json:"pos" yaml:"pos" xml:"pos,attr"`    // segment number, starts from 1
	Pay  int `json:"pay" yaml:"pay" xml:"pay,attr"`    // pay by this segment
}

func Eldorado1Spawn(bet int) (any, int) {
	var res = make([]WinElSeg, 1)
	var cash int
	for i := range res {
		var n = rand.Intn(len(Eldorado))
		res[i].Mult = Eldorado[n]
		res[i].Pos = n + 1
		res[i].Pay = bet * Eldorado[n]
		cash += Eldorado[n]
	}
	return res, bet * cash
}

func Eldorado3Spawn(bet int) (any, int) {
	var res = make([]WinElSeg, 3)
	var pay int
	for i := range res {
		var n = rand.Intn(len(Eldorado))
		res[i].Mult = Eldorado[n]
		res[i].Pos = n + 1
		res[i].Pay = bet * Eldorado[n]
		pay += res[i].Pay
	}
	return res, pay
}

func Eldorado6Spawn(bet int) (any, int) {
	var res = make([]WinElSeg, 6)
	var pay int
	for i := range res {
		var n = rand.Intn(len(Eldorado))
		res[i].Mult = Eldorado[n]
		res[i].Pos = n + 1
		res[i].Pay = bet * Eldorado[n]
		pay += res[i].Pay
	}
	return res, pay
}

func Eldorado9Spawn(bet int) (any, int) {
	var res = make([]WinElSeg, 9)
	var pay int
	for i := range res {
		var n = rand.Intn(len(Eldorado))
		res[i].Mult = Eldorado[n]
		res[i].Pos = n + 1
		res[i].Pay = bet * Eldorado[n]
		pay += res[i].Pay
	}
	return res, pay
}
