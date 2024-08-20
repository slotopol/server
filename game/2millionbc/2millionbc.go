package twomillionbc

import (
	"math"

	"github.com/slotopol/server/game"
)

// reels lengths [32, 32, 32, 32, 110], total reshuffles 115343360
// symbols: 60.055(lined) + 0(scatter) = 60.054814%
// free spins 3973320, q = 0.034448, sq = 1/(1-q) = 1.035677
// free games frequency: 1/128.39
// acorn bonuses: frequency 1/110, rtp = 10.000000%
// diamond lion bonuses: frequency 1/1213.6, rtp = 14.419556%
// RTP = 60.055(sym) + 10(acorn) + 14.42(dl) + 0.034448*124.68(fg) = 88.769193%
var ReelsReg89 = game.Reels5x{
	{10, 13, 10, 5, 1, 5, 10, 5, 2, 9, 8, 2, 6, 13, 6, 8, 1, 11, 10, 3, 4, 7, 13, 4, 9, 9, 3, 7, 4, 8, 7, 6},
	{5, 2, 10, 7, 3, 1, 11, 9, 8, 4, 1, 6, 9, 10, 6, 4, 13, 7, 6, 3, 8, 13, 5, 8, 13, 5, 10, 2, 4, 9, 10, 7},
	{3, 1, 4, 2, 7, 5, 4, 6, 13, 10, 4, 10, 13, 8, 9, 3, 8, 9, 10, 8, 9, 11, 5, 6, 2, 1, 5, 3, 6, 7, 13, 7},
	{8, 8, 3, 1, 9, 11, 5, 2, 7, 6, 13, 2, 8, 9, 13, 3, 3, 5, 10, 10, 7, 1, 10, 9, 7, 5, 13, 4, 4, 4, 6, 6},
	{2, 4, 7, 10, 4, 6, 5, 6, 10, 6, 8, 6, 10, 2, 7, 2, 7, 7, 2, 3, 13, 6, 3, 6, 9, 2, 12, 5, 6, 11, 9, 8, 13, 7, 7, 10, 8, 3, 1, 4, 6, 13, 7, 6, 8, 3, 13, 4, 10, 13, 3, 8, 10, 9, 13, 3, 3, 1, 9, 5, 11, 9, 8, 13, 2, 5, 1, 11, 1, 4, 11, 8, 1, 5, 9, 2, 8, 1, 9, 5, 13, 2, 9, 3, 13, 10, 3, 10, 4, 5, 4, 4, 9, 13, 4, 7, 4, 5, 8, 10, 5, 9, 1, 7, 5, 7, 10, 1, 8, 6},
}

// reels lengths [32, 32, 32, 32, 110], total reshuffles 115343360
// symbols: 62.674(lined) + 0(scatter) = 62.673517%
// free spins 3973320, q = 0.034448, sq = 1/(1-q) = 1.035677
// free games frequency: 1/128.39
// acorn bonuses: frequency 1/110, rtp = 10.000000%
// diamond lion bonuses: frequency 1/1213.6, rtp = 14.419556%
// RTP = 62.674(sym) + 10(acorn) + 14.42(dl) + 0.034448*124.68(fg) = 91.387896%
var ReelsReg91 = game.Reels5x{
	{10, 13, 10, 5, 1, 5, 10, 5, 2, 9, 8, 2, 6, 13, 6, 8, 1, 11, 10, 3, 4, 7, 13, 4, 9, 9, 3, 7, 4, 8, 7, 6},
	{10, 9, 4, 8, 2, 10, 7, 4, 7, 10, 6, 11, 1, 1, 7, 3, 13, 2, 5, 13, 9, 8, 5, 13, 3, 6, 8, 5, 4, 6, 9, 3},
	{3, 1, 4, 2, 7, 5, 4, 6, 13, 10, 4, 10, 13, 8, 9, 3, 8, 9, 10, 8, 9, 11, 5, 6, 2, 1, 5, 3, 6, 7, 13, 7},
	{8, 8, 3, 1, 9, 11, 5, 2, 7, 6, 13, 2, 8, 9, 13, 3, 3, 5, 10, 10, 7, 1, 10, 9, 7, 5, 13, 4, 4, 4, 6, 6},
	{2, 4, 7, 10, 4, 6, 5, 6, 10, 6, 8, 6, 10, 2, 7, 2, 7, 7, 2, 3, 13, 6, 3, 6, 9, 2, 12, 5, 6, 11, 9, 8, 13, 7, 7, 10, 8, 3, 1, 4, 6, 13, 7, 6, 8, 3, 13, 4, 10, 13, 3, 8, 10, 9, 13, 3, 3, 1, 9, 5, 11, 9, 8, 13, 2, 5, 1, 11, 1, 4, 11, 8, 1, 5, 9, 2, 8, 1, 9, 5, 13, 2, 9, 3, 13, 10, 3, 10, 4, 5, 4, 4, 9, 13, 4, 7, 4, 5, 8, 10, 5, 9, 1, 7, 5, 7, 10, 1, 8, 6},
}

// reels lengths [30, 32, 32, 34, 110], total reshuffles 114892800
// symbols: 63.397(lined) + 0(scatter) = 63.396906%
// free spins 3971592, q = 0.034568, sq = 1/(1-q) = 1.035806
// free games frequency: 1/127.95
// acorn bonuses: frequency 1/110, rtp = 10.000000%
// diamond lion bonuses: frequency 1/1137.8, rtp = 15.380859%
// RTP = 63.397(sym) + 10(acorn) + 15.381(dl) + 0.034568*124.68(fg) = 93.087556%
var ReelsReg93 = game.Reels5x{
	{11, 3, 9, 13, 2, 6, 8, 9, 4, 4, 8, 5, 1, 9, 5, 2, 6, 10, 10, 13, 7, 6, 13, 10, 8, 1, 7, 5, 3, 7},
	{10, 9, 4, 8, 2, 10, 7, 4, 7, 10, 6, 11, 1, 1, 7, 3, 13, 2, 5, 13, 9, 8, 5, 13, 3, 6, 8, 5, 4, 6, 9, 3},
	{3, 1, 4, 2, 7, 5, 4, 6, 13, 10, 4, 10, 13, 8, 9, 3, 8, 9, 10, 8, 9, 11, 5, 6, 2, 1, 5, 3, 6, 7, 13, 7},
	{8, 10, 11, 9, 2, 4, 5, 8, 13, 8, 9, 3, 9, 4, 2, 1, 6, 5, 1, 6, 3, 13, 7, 4, 3, 7, 13, 10, 5, 6, 1, 7, 10, 2},
	{2, 4, 7, 10, 4, 6, 5, 6, 10, 6, 8, 6, 10, 2, 7, 2, 7, 7, 2, 3, 13, 6, 3, 6, 9, 2, 12, 5, 6, 11, 9, 8, 13, 7, 7, 10, 8, 3, 1, 4, 6, 13, 7, 6, 8, 3, 13, 4, 10, 13, 3, 8, 10, 9, 13, 3, 3, 1, 9, 5, 11, 9, 8, 13, 2, 5, 1, 11, 1, 4, 11, 8, 1, 5, 9, 2, 8, 1, 9, 5, 13, 2, 9, 3, 13, 10, 3, 10, 4, 5, 4, 4, 9, 13, 4, 7, 4, 5, 8, 10, 5, 9, 1, 7, 5, 7, 10, 1, 8, 6},
}

// reels lengths [31, 32, 32, 32, 110], total reshuffles 111738880
// symbols: 64.366(lined) + 0(scatter) = 64.365814%
// free spins 3921264, q = 0.035093, sq = 1/(1-q) = 1.036369
// free games frequency: 1/126.11
// acorn bonuses: frequency 1/110, rtp = 10.000000%
// diamond lion bonuses: frequency 1/1175.7, rtp = 14.884703%
// RTP = 64.366(sym) + 10(acorn) + 14.885(dl) + 0.035093*124.68(fg) = 93.625800%
var ReelsReg94 = game.Reels5x{
	{6, 10, 13, 2, 7, 7, 11, 6, 5, 3, 9, 2, 13, 6, 9, 3, 13, 7, 9, 1, 10, 5, 5, 1, 4, 8, 8, 10, 4, 8, 4},
	{10, 9, 4, 8, 2, 10, 7, 4, 7, 10, 6, 11, 1, 1, 7, 3, 13, 2, 5, 13, 9, 8, 5, 13, 3, 6, 8, 5, 4, 6, 9, 3},
	{3, 1, 4, 2, 7, 5, 4, 6, 13, 10, 4, 10, 13, 8, 9, 3, 8, 9, 10, 8, 9, 11, 5, 6, 2, 1, 5, 3, 6, 7, 13, 7},
	{8, 8, 3, 1, 9, 11, 5, 2, 7, 6, 13, 2, 8, 9, 13, 3, 3, 5, 10, 10, 7, 1, 10, 9, 7, 5, 13, 4, 4, 4, 6, 6},
	{2, 4, 7, 10, 4, 6, 5, 6, 10, 6, 8, 6, 10, 2, 7, 2, 7, 7, 2, 3, 13, 6, 3, 6, 9, 2, 12, 5, 6, 11, 9, 8, 13, 7, 7, 10, 8, 3, 1, 4, 6, 13, 7, 6, 8, 3, 13, 4, 10, 13, 3, 8, 10, 9, 13, 3, 3, 1, 9, 5, 11, 9, 8, 13, 2, 5, 1, 11, 1, 4, 11, 8, 1, 5, 9, 2, 8, 1, 9, 5, 13, 2, 9, 3, 13, 10, 3, 10, 4, 5, 4, 4, 9, 13, 4, 7, 4, 5, 8, 10, 5, 9, 1, 7, 5, 7, 10, 1, 8, 6},
}

// reels lengths [32, 32, 32, 32, 110], total reshuffles 115343360
// symbols: 66.921(lined) + 0(scatter) = 66.920705%
// free spins 3973320, q = 0.034448, sq = 1/(1-q) = 1.035677
// free games frequency: 1/128.39
// acorn bonuses: frequency 1/110, rtp = 10.000000%
// diamond lion bonuses: frequency 1/1213.6, rtp = 14.419556%
// RTP = 66.921(sym) + 10(acorn) + 14.42(dl) + 0.034448*124.68(fg) = 95.635084%
var ReelsReg96 = game.Reels5x{
	{4, 2, 10, 13, 10, 10, 4, 7, 5, 2, 13, 7, 8, 9, 7, 9, 5, 3, 13, 6, 5, 8, 6, 4, 1, 1, 6, 9, 3, 11, 8, 3},
	{10, 9, 4, 8, 2, 10, 7, 4, 7, 10, 6, 11, 1, 1, 7, 3, 13, 2, 5, 13, 9, 8, 5, 13, 3, 6, 8, 5, 4, 6, 9, 3},
	{3, 1, 4, 2, 7, 5, 4, 6, 13, 10, 4, 10, 13, 8, 9, 3, 8, 9, 10, 8, 9, 11, 5, 6, 2, 1, 5, 3, 6, 7, 13, 7},
	{8, 8, 3, 1, 9, 11, 5, 2, 7, 6, 13, 2, 8, 9, 13, 3, 3, 5, 10, 10, 7, 1, 10, 9, 7, 5, 13, 4, 4, 4, 6, 6},
	{2, 4, 7, 10, 4, 6, 5, 6, 10, 6, 8, 6, 10, 2, 7, 2, 7, 7, 2, 3, 13, 6, 3, 6, 9, 2, 12, 5, 6, 11, 9, 8, 13, 7, 7, 10, 8, 3, 1, 4, 6, 13, 7, 6, 8, 3, 13, 4, 10, 13, 3, 8, 10, 9, 13, 3, 3, 1, 9, 5, 11, 9, 8, 13, 2, 5, 1, 11, 1, 4, 11, 8, 1, 5, 9, 2, 8, 1, 9, 5, 13, 2, 9, 3, 13, 10, 3, 10, 4, 5, 4, 4, 9, 13, 4, 7, 4, 5, 8, 10, 5, 9, 1, 7, 5, 7, 10, 1, 8, 6},
}

// reels lengths [32, 32, 32, 32, 110], total reshuffles 115343360
// symbols: 68.262(lined) + 0(scatter) = 68.262265%
// free spins 3973320, q = 0.034448, sq = 1/(1-q) = 1.035677
// free games frequency: 1/128.39
// acorn bonuses: frequency 1/110, rtp = 10.000000%
// diamond lion bonuses: frequency 1/1213.6, rtp = 14.419556%
// RTP = 68.262(sym) + 10(acorn) + 14.42(dl) + 0.034448*124.68(fg) = 96.976644%
var ReelsReg97 = game.Reels5x{
	{4, 2, 10, 13, 10, 10, 4, 7, 5, 2, 13, 7, 8, 9, 7, 9, 5, 3, 13, 6, 5, 8, 6, 4, 1, 1, 6, 9, 3, 11, 8, 3},
	{10, 9, 4, 8, 2, 10, 7, 4, 7, 10, 6, 11, 1, 1, 7, 3, 13, 2, 5, 13, 9, 8, 5, 13, 3, 6, 8, 5, 4, 6, 9, 3},
	{8, 3, 4, 10, 13, 6, 5, 2, 13, 8, 9, 5, 7, 10, 1, 4, 7, 13, 3, 8, 11, 2, 3, 1, 7, 1, 5, 9, 2, 6, 6, 4},
	{8, 13, 7, 5, 2, 5, 13, 4, 7, 1, 10, 3, 13, 5, 4, 9, 1, 2, 11, 1, 6, 4, 8, 3, 9, 6, 6, 3, 8, 7, 2, 10},
	{2, 4, 7, 10, 4, 6, 5, 6, 10, 6, 8, 6, 10, 2, 7, 2, 7, 7, 2, 3, 13, 6, 3, 6, 9, 2, 12, 5, 6, 11, 9, 8, 13, 7, 7, 10, 8, 3, 1, 4, 6, 13, 7, 6, 8, 3, 13, 4, 10, 13, 3, 8, 10, 9, 13, 3, 3, 1, 9, 5, 11, 9, 8, 13, 2, 5, 1, 11, 1, 4, 11, 8, 1, 5, 9, 2, 8, 1, 9, 5, 13, 2, 9, 3, 13, 10, 3, 10, 4, 5, 4, 4, 9, 13, 4, 7, 4, 5, 8, 10, 5, 9, 1, 7, 5, 7, 10, 1, 8, 6},
}

// reels lengths [32, 32, 32, 32, 110], total reshuffles 115343360
// symbols: 71.205(lined) + 0(scatter) = 71.205347%
// free spins 3973320, q = 0.034448, sq = 1/(1-q) = 1.035677
// free games frequency: 1/128.39
// acorn bonuses: frequency 1/110, rtp = 10.000000%
// diamond lion bonuses: frequency 1/1213.6, rtp = 14.419556%
// RTP = 71.205(sym) + 10(acorn) + 14.42(dl) + 0.034448*124.68(fg) = 99.919726%
var ReelsReg100 = game.Reels5x{
	{4, 2, 10, 13, 10, 10, 4, 7, 5, 2, 13, 7, 8, 9, 7, 9, 5, 3, 13, 6, 5, 8, 6, 4, 1, 1, 6, 9, 3, 11, 8, 3},
	{9, 10, 1, 11, 1, 10, 6, 13, 5, 4, 13, 10, 6, 4, 8, 3, 7, 1, 6, 8, 3, 9, 5, 13, 2, 7, 9, 8, 4, 2, 7, 5},
	{8, 3, 4, 10, 13, 6, 5, 2, 13, 8, 9, 5, 7, 10, 1, 4, 7, 13, 3, 8, 11, 2, 3, 1, 7, 1, 5, 9, 2, 6, 6, 4},
	{8, 13, 7, 5, 2, 5, 13, 4, 7, 1, 10, 3, 13, 5, 4, 9, 1, 2, 11, 1, 6, 4, 8, 3, 9, 6, 6, 3, 8, 7, 2, 10},
	{2, 4, 7, 10, 4, 6, 5, 6, 10, 6, 8, 6, 10, 2, 7, 2, 7, 7, 2, 3, 13, 6, 3, 6, 9, 2, 12, 5, 6, 11, 9, 8, 13, 7, 7, 10, 8, 3, 1, 4, 6, 13, 7, 6, 8, 3, 13, 4, 10, 13, 3, 8, 10, 9, 13, 3, 3, 1, 9, 5, 11, 9, 8, 13, 2, 5, 1, 11, 1, 4, 11, 8, 1, 5, 9, 2, 8, 1, 9, 5, 13, 2, 9, 3, 13, 10, 3, 10, 4, 5, 4, 4, 9, 13, 4, 7, 4, 5, 8, 10, 5, 9, 1, 7, 5, 7, 10, 1, 8, 6},
}

// reels lengths [23, 23, 23, 23, 110], total reshuffles 30782510
// symbols: 81.964(lined) + 0(scatter) = 81.963557%
// free spins 2622240, q = 0.085186, sq = 1/(1-q) = 1.093118
// free games frequency: 1/53.956
// acorn bonuses: frequency 1/110, rtp = 10.000000%
// diamond lion bonuses: frequency 1/1520.9, rtp = 11.506534%
// RTP = 81.964(sym) + 10(acorn) + 11.507(dl) + 0.085186*124.68(fg) = 114.090782%
var ReelsReg114 = game.Reels5x{
	{13, 2, 8, 4, 13, 7, 7, 11, 8, 1, 10, 4, 2, 6, 5, 6, 5, 3, 9, 1, 9, 3, 10},
	{9, 5, 2, 1, 8, 10, 6, 13, 3, 3, 9, 1, 8, 4, 13, 4, 7, 5, 11, 6, 10, 7, 2},
	{6, 2, 10, 4, 13, 8, 4, 7, 7, 3, 3, 9, 9, 13, 10, 2, 1, 11, 5, 8, 1, 5, 6},
	{11, 10, 9, 4, 4, 2, 6, 13, 5, 8, 3, 1, 10, 3, 9, 1, 5, 7, 2, 7, 13, 8, 6},
	{4, 9, 4, 13, 8, 9, 6, 8, 13, 6, 10, 13, 6, 8, 3, 4, 11, 8, 10, 12, 1, 1, 11, 5, 4, 13, 4, 4, 9, 2, 7, 1, 9, 4, 13, 3, 3, 5, 3, 6, 6, 3, 2, 7, 7, 1, 6, 2, 5, 7, 4, 2, 13, 6, 5, 10, 6, 10, 3, 8, 11, 7, 9, 2, 3, 8, 6, 5, 2, 10, 13, 10, 5, 7, 13, 4, 7, 8, 13, 10, 7, 8, 10, 9, 6, 9, 13, 9, 5, 1, 1, 11, 7, 10, 10, 8, 1, 9, 5, 11, 3, 1, 5, 2, 7, 4, 9, 2, 8, 5},
}

// reels lengths [32, 32, 32, 32, 32], total reshuffles 33554432
// symbols: 94.885(lined) + 0(scatter) = 94.884673%
// free spins 8017920, q = 0.23895, sq = 1/(1-q) = 1.313979
// free games frequency: 1/20.505
// RTP = sq*rtp(sym) = 1.314*94.885 = 124.676436%
var ReelsBon = game.Reels5x{
	{6, 4, 5, 5, 8, 7, 10, 4, 8, 2, 6, 11, 4, 9, 7, 3, 2, 8, 11, 3, 1, 5, 1, 1, 10, 9, 9, 6, 10, 3, 2, 7},
	{7, 1, 10, 7, 1, 6, 7, 2, 5, 8, 10, 8, 4, 11, 9, 3, 5, 2, 9, 10, 8, 5, 1, 4, 11, 6, 3, 3, 4, 6, 9, 2},
	{3, 4, 10, 6, 8, 9, 1, 7, 4, 11, 3, 5, 10, 1, 6, 1, 7, 2, 11, 7, 9, 10, 8, 8, 3, 9, 2, 6, 4, 5, 2, 5},
	{2, 2, 4, 3, 3, 8, 3, 1, 9, 9, 5, 2, 5, 11, 7, 10, 10, 6, 1, 10, 1, 6, 9, 5, 7, 4, 8, 4, 6, 7, 11, 8},
	{3, 3, 8, 7, 7, 1, 10, 1, 10, 11, 2, 1, 9, 4, 6, 5, 5, 2, 10, 4, 2, 6, 11, 9, 8, 3, 9, 8, 4, 5, 6, 7},
}

// Map with available reels.
var reelsmap = map[float64]*game.Reels5x{
	88.769193:  &ReelsReg89,
	91.387896:  &ReelsReg91,
	93.087556:  &ReelsReg93,
	93.625800:  &ReelsReg94,
	95.635084:  &ReelsReg96,
	96.976644:  &ReelsReg97,
	99.919726:  &ReelsReg100,
	114.090782: &ReelsReg114,
}

func FindReels(mrtp float64) (rtp float64, reels game.Reels) {
	for p, r := range reelsmap {
		if math.Abs(mrtp-p) < math.Abs(mrtp-rtp) {
			rtp, reels = p, r
		}
	}
	return
}

// Lined payment.
var LinePay = [13][5]float64{
	{0, 30, 100, 300, 500}, //  1 girl
	{0, 15, 75, 200, 400},  //  2 lion
	{0, 10, 60, 150, 300},  //  3 bee
	{0, 5, 50, 125, 250},   //  4 stone
	{0, 5, 40, 100, 200},   //  5 wheel
	{0, 2, 30, 90, 150},    //  6 club
	{0, 0, 25, 75, 125},    //  7 chaplet
	{0, 0, 20, 60, 100},    //  8 gold
	{0, 0, 15, 50, 75},     //  9 vase
	{0, 0, 10, 25, 50},     // 10 ruby
	{0, 0, 0, 0, 0},        // 11 fire
	{0, 0, 0, 0, 0},        // 12 acorn
	{0, 0, 40, 100, 200},   // 13 diamond
}

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 4, 12, 20} // 11 fire

const (
	acbn = 1 // acorn bonus
	dlbn = 2 // diamond lion bonus
)

type Game struct {
	game.Slot5x3 `yaml:",inline"`
	// free spin number
	FS int `json:"fs,omitempty" yaml:"fs,omitempty" xml:"fs,omitempty"`
	// acorns number
	AN int `json:"an" yaml:"an" xml:"an"`
	// acorns bet
	AB float64 `json:"ab" yaml:"ab" xml:"ab"`
}

func NewGame(rtp float64) *Game {
	return &Game{
		Slot5x3: game.Slot5x3{
			RTP: rtp,
			SBL: game.MakeBitNum(30),
			Bet: 1,
		},
		FS: 0,
	}
}

const scat, acorn, diamond = 11, 12, 13

var bl = game.BetLinesBetSoft30

func (g *Game) Scanner(screen game.Screen, wins *game.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen game.Screen, wins *game.Wins) {
	for li := g.SBL.Next(0); li != 0; li = g.SBL.Next(li) {
		var line = bl.Line(li)

		var syml, numl = screen.Pos(1, line), 1
		for x := 2; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx != syml {
				break
			}
			numl++
		}

		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			*wins = append(*wins, game.WinItem{
				Pay:  g.Bet * pay,
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		}
		if syml == diamond && numl >= 3 {
			*wins = append(*wins, game.WinItem{
				Mult: 1,
				Sym:  diamond,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
				BID:  dlbn,
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen game.Screen, wins *game.Wins) {
	if count := screen.ScatNum(scat); count >= 3 {
		var fs = ScatFreespin[count-1]
		*wins = append(*wins, game.WinItem{
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: fs,
		})
	}

	if screen.At(5, 1) == acorn || screen.At(5, 2) == acorn || screen.At(5, 3) == acorn {
		if (g.AN+1)%3 == 0 {
			*wins = append(*wins, game.WinItem{
				Mult: 1,
				Sym:  acorn,
				Num:  1,
				BID:  acbn,
			})
		}
	}
}

func (g *Game) Spin(screen game.Screen) {
	if g.FS == 0 {
		var _, reels = FindReels(g.RTP)
		screen.Spin(reels)
	} else {
		screen.Spin(&ReelsBon)
	}
}

func (g *Game) Spawn(screen game.Screen, wins game.Wins) {
	for i, wi := range wins {
		switch wi.BID {
		case acbn:
			wins[i].Pay = AcornSpawn(g.AB + g.Bet*float64(g.SBL.Num()))
		case dlbn:
			wins[i].Pay = DiamondLionSpawn(g.Bet)
		}
	}
}

func (g *Game) Apply(screen game.Screen, wins game.Wins) {
	if screen.At(5, 1) == acorn || screen.At(5, 2) == acorn || screen.At(5, 3) == acorn {
		g.AN++
		g.AN %= 3
		if g.AN > 0 {
			g.AB += g.Bet * float64(g.SBL.Num())
		} else {
			g.AB = 0
		}
	}

	if g.FS > 0 {
		g.Gain += wins.Gain()
	} else {
		g.Gain = wins.Gain()
	}

	if g.FS > 0 {
		g.FS--
	}
	for _, wi := range wins {
		if wi.Free > 0 {
			g.FS += wi.Free
		}
	}
}

func (g *Game) FreeSpins() int {
	return g.FS
}

func (g *Game) SetLines(sbl game.Bitset) error {
	var mask game.Bitset = (1<<len(bl) - 1) << 1
	if sbl == 0 {
		return game.ErrNoLineset
	}
	if sbl&^mask != 0 {
		return game.ErrLinesetOut
	}
	if g.FreeSpins() > 0 {
		return game.ErrNoFeature
	}
	g.SBL = sbl
	return nil
}
