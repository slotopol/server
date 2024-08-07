package arabiannights

import (
	"github.com/slotopol/server/game"
)

// reels lengths [34, 34, 34, 34, 34], total reshuffles 45435424
// symbols: 44.982(lined) + 15.478(scatter) = 60.459319%
// free spins 4084020, q = 0.089886, sq = 1/(1-q) = 1.098764
// free games frequency: 1/166.88
// RTP = 60.459(sym) + 0.089886*276.85(fg) = 85.344075%
var ReelsReg85 = game.Reels5x{
	{6, 9, 5, 11, 4, 6, 8, 5, 9, 3, 7, 8, 4, 12, 5, 7, 8, 1, 10, 9, 6, 10, 8, 7, 10, 9, 5, 2, 10, 3, 9, 6, 10, 7},
	{10, 11, 9, 10, 6, 4, 5, 7, 3, 6, 8, 9, 10, 5, 1, 8, 9, 7, 12, 9, 7, 6, 5, 9, 10, 8, 3, 7, 10, 4, 5, 8, 2, 6},
	{7, 9, 5, 11, 9, 3, 6, 4, 10, 5, 7, 6, 8, 12, 9, 5, 10, 3, 8, 6, 1, 8, 5, 10, 9, 7, 10, 8, 7, 9, 4, 10, 6, 2},
	{7, 8, 5, 6, 8, 10, 7, 4, 11, 10, 3, 9, 5, 7, 6, 10, 3, 6, 1, 5, 4, 9, 8, 7, 9, 2, 6, 5, 10, 9, 8, 10, 12, 9},
	{5, 7, 4, 9, 6, 4, 8, 9, 2, 10, 6, 9, 10, 3, 7, 8, 10, 6, 5, 3, 9, 7, 5, 12, 10, 6, 8, 10, 5, 9, 8, 1, 7, 11},
}

// reels lengths [34, 34, 34, 34, 34], total reshuffles 45435424
// symbols: 47.525(lined) + 15.478(scatter) = 63.002643%
// free spins 4084020, q = 0.089886, sq = 1/(1-q) = 1.098764
// free games frequency: 1/166.88
// RTP = 63.003(sym) + 0.089886*276.85(fg) = 87.887399%
var ReelsReg88 = game.Reels5x{
	{12, 3, 5, 10, 8, 5, 10, 9, 7, 5, 8, 6, 10, 9, 7, 10, 2, 8, 7, 9, 4, 8, 3, 7, 1, 8, 6, 11, 9, 6, 7, 4, 9, 10},
	{10, 5, 9, 7, 8, 5, 9, 10, 1, 8, 6, 10, 2, 9, 4, 8, 5, 7, 8, 4, 7, 6, 3, 12, 10, 6, 8, 7, 9, 10, 11, 3, 9, 7},
	{10, 2, 8, 7, 6, 9, 7, 8, 5, 6, 8, 1, 12, 10, 3, 9, 5, 7, 10, 8, 9, 4, 10, 7, 9, 11, 6, 10, 4, 8, 3, 9, 5, 7},
	{7, 5, 6, 3, 9, 7, 12, 9, 7, 8, 2, 10, 6, 8, 7, 10, 1, 8, 5, 6, 8, 10, 3, 9, 7, 4, 9, 10, 11, 4, 8, 9, 5, 10},
	{4, 9, 7, 8, 9, 4, 10, 6, 8, 1, 9, 10, 5, 6, 9, 8, 6, 7, 9, 5, 3, 10, 7, 8, 11, 2, 7, 10, 8, 7, 3, 10, 5, 12},
}

// reels lengths [34, 34, 34, 34, 34], total reshuffles 45435424
// symbols: 49.584(lined) + 15.478(scatter) = 65.062023%
// free spins 4084020, q = 0.089886, sq = 1/(1-q) = 1.098764
// free games frequency: 1/166.88
// RTP = 65.062(sym) + 0.089886*276.85(fg) = 89.946779%
var ReelsReg90 = game.Reels5x{
	{3, 10, 12, 4, 10, 8, 6, 7, 10, 4, 9, 8, 6, 10, 5, 6, 4, 7, 9, 3, 8, 5, 7, 8, 9, 3, 5, 7, 2, 11, 1, 9, 5, 6},
	{9, 5, 6, 10, 3, 6, 5, 8, 7, 4, 9, 7, 11, 10, 9, 5, 3, 8, 6, 2, 10, 8, 7, 6, 5, 3, 12, 4, 7, 8, 4, 9, 1, 10},
	{8, 12, 9, 6, 7, 9, 5, 4, 9, 8, 3, 5, 10, 4, 7, 10, 3, 6, 7, 8, 3, 5, 2, 10, 11, 9, 1, 6, 7, 10, 5, 8, 6, 4},
	{7, 4, 9, 10, 3, 8, 9, 7, 6, 8, 5, 4, 10, 3, 5, 8, 9, 10, 6, 5, 10, 1, 6, 7, 4, 6, 12, 3, 7, 8, 11, 5, 9, 2},
	{1, 10, 12, 9, 4, 10, 2, 9, 5, 6, 8, 7, 9, 11, 8, 7, 3, 6, 4, 10, 9, 6, 5, 4, 8, 3, 5, 10, 7, 6, 3, 8, 7, 5},
}

// reels lengths [32, 34, 32, 34, 32], total reshuffles 37879808
// symbols: 47.004(lined) + 16.687(scatter) = 63.690550%
// free spins 3780270, q = 0.099796, sq = 1/(1-q) = 1.110860
// free games frequency: 1/150.31
// RTP = 63.691(sym) + 0.099796*276.85(fg) = 91.318913%
var ReelsReg91 = game.Reels5x{
	{6, 10, 9, 8, 5, 9, 2, 10, 5, 4, 10, 9, 7, 8, 1, 9, 8, 10, 3, 9, 7, 11, 6, 4, 10, 3, 6, 5, 7, 8, 12, 7},
	{9, 5, 6, 10, 3, 6, 5, 8, 7, 4, 9, 7, 11, 10, 9, 5, 3, 8, 6, 2, 10, 8, 7, 6, 5, 3, 12, 4, 7, 8, 4, 9, 1, 10},
	{4, 5, 8, 10, 7, 6, 3, 8, 11, 9, 7, 10, 12, 8, 6, 9, 4, 5, 9, 1, 10, 9, 2, 8, 5, 7, 3, 10, 7, 6, 9, 10},
	{7, 4, 9, 10, 3, 8, 9, 7, 6, 8, 5, 4, 10, 3, 5, 8, 9, 10, 6, 5, 10, 1, 6, 7, 4, 6, 12, 3, 7, 8, 11, 5, 9, 2},
	{7, 2, 5, 6, 10, 9, 3, 5, 8, 6, 4, 9, 7, 1, 9, 3, 8, 12, 7, 10, 9, 7, 10, 11, 8, 5, 10, 4, 6, 9, 8, 10},
}

// reels lengths [32, 32, 34, 32, 32], total reshuffles 35651584
// symbols: 46.38(lined) + 17.104(scatter) = 63.484826%
// free spins 3682260, q = 0.10328, sq = 1/(1-q) = 1.115181
// free games frequency: 1/145.23
// RTP = 63.485(sym) + 0.10328*276.85(fg) = 92.078880%
var ReelsReg92 = game.Reels5x{
	{6, 10, 9, 8, 5, 9, 2, 10, 5, 4, 10, 9, 7, 8, 1, 9, 8, 10, 3, 9, 7, 11, 6, 4, 10, 3, 6, 5, 7, 8, 12, 7},
	{8, 9, 5, 6, 10, 3, 9, 11, 5, 4, 10, 9, 7, 8, 9, 7, 2, 10, 7, 8, 4, 10, 3, 5, 7, 1, 6, 8, 9, 6, 10, 12},
	{8, 12, 9, 6, 7, 9, 5, 4, 9, 8, 3, 5, 10, 4, 7, 10, 3, 6, 7, 8, 3, 5, 2, 10, 11, 9, 1, 6, 7, 10, 5, 8, 6, 4},
	{7, 10, 9, 6, 10, 9, 7, 4, 6, 5, 2, 9, 6, 8, 3, 7, 10, 1, 8, 10, 3, 8, 5, 11, 7, 9, 4, 5, 9, 12, 10, 8},
	{7, 2, 5, 6, 10, 9, 3, 5, 8, 6, 4, 9, 7, 1, 9, 3, 8, 12, 7, 10, 9, 7, 10, 11, 8, 5, 10, 4, 6, 9, 8, 10},
}

// reels lengths [32, 32, 34, 32, 32], total reshuffles 35651584
// symbols: 47.784(lined) + 17.104(scatter) = 64.888253%
// free spins 3682260, q = 0.10328, sq = 1/(1-q) = 1.115181
// free games frequency: 1/145.23
// RTP = 64.888(sym) + 0.10328*276.85(fg) = 93.482307%
var ReelsReg93 = game.Reels5x{
	{10, 8, 11, 10, 7, 9, 12, 3, 9, 8, 7, 6, 5, 10, 3, 5, 9, 4, 10, 2, 8, 9, 6, 5, 7, 6, 4, 7, 1, 8, 6, 5},
	{6, 7, 9, 10, 8, 3, 12, 6, 7, 10, 5, 7, 8, 9, 3, 6, 11, 5, 4, 9, 5, 10, 4, 6, 9, 2, 7, 8, 5, 10, 1, 8},
	{8, 12, 9, 6, 7, 9, 5, 4, 9, 8, 3, 5, 10, 4, 7, 10, 3, 6, 7, 8, 3, 5, 2, 10, 11, 9, 1, 6, 7, 10, 5, 8, 6, 4},
	{9, 5, 8, 7, 11, 6, 8, 7, 4, 6, 7, 9, 3, 6, 2, 5, 10, 9, 8, 10, 9, 8, 5, 10, 6, 4, 12, 7, 10, 1, 5, 3},
	{6, 2, 9, 8, 7, 9, 3, 10, 4, 5, 7, 8, 5, 7, 10, 9, 6, 11, 5, 8, 6, 5, 4, 6, 1, 10, 9, 7, 12, 8, 3, 10},
}

// reels lengths [32, 32, 32, 32, 32], total reshuffles 33554432
// symbols: 48.099(lined) + 17.529(scatter) = 65.627718%
// free spins 3585870, q = 0.10687, sq = 1/(1-q) = 1.119654
// free games frequency: 1/140.36
// RTP = 65.628(sym) + 0.10687*276.85(fg) = 95.213616%
var ReelsReg95 = game.Reels5x{
	{6, 10, 9, 8, 5, 9, 2, 10, 5, 4, 10, 9, 7, 8, 1, 9, 8, 10, 3, 9, 7, 11, 6, 4, 10, 3, 6, 5, 7, 8, 12, 7},
	{8, 9, 5, 6, 10, 3, 9, 11, 5, 4, 10, 9, 7, 8, 9, 7, 2, 10, 7, 8, 4, 10, 3, 5, 7, 1, 6, 8, 9, 6, 10, 12},
	{4, 5, 8, 10, 7, 6, 3, 8, 11, 9, 7, 10, 12, 8, 6, 9, 4, 5, 9, 1, 10, 9, 2, 8, 5, 7, 3, 10, 7, 6, 9, 10},
	{7, 10, 9, 6, 10, 9, 7, 4, 6, 5, 2, 9, 6, 8, 3, 7, 10, 1, 8, 10, 3, 8, 5, 11, 7, 9, 4, 5, 9, 12, 10, 8},
	{7, 2, 5, 6, 10, 9, 3, 5, 8, 6, 4, 9, 7, 1, 9, 3, 8, 12, 7, 10, 9, 7, 10, 11, 8, 5, 10, 4, 6, 9, 8, 10},
}

// reels lengths [32, 32, 32, 32, 32], total reshuffles 33554432
// symbols: 48.783(lined) + 17.529(scatter) = 66.312021%
// free spins 3585870, q = 0.10687, sq = 1/(1-q) = 1.119654
// free games frequency: 1/140.36
// RTP = 66.312(sym) + 0.10687*276.85(fg) = 95.897919%
var ReelsReg96 = game.Reels5x{
	{10, 8, 11, 10, 7, 9, 12, 3, 9, 8, 7, 6, 5, 10, 3, 5, 9, 4, 10, 2, 8, 9, 6, 5, 7, 6, 4, 7, 1, 8, 6, 5},
	{6, 7, 9, 10, 8, 3, 12, 6, 7, 10, 5, 7, 8, 9, 3, 6, 11, 5, 4, 9, 5, 10, 4, 6, 9, 2, 7, 8, 5, 10, 1, 8},
	{8, 9, 11, 5, 10, 6, 8, 9, 7, 10, 5, 7, 10, 5, 1, 9, 6, 2, 5, 6, 3, 8, 10, 3, 8, 6, 4, 7, 12, 4, 9, 7},
	{9, 5, 8, 7, 11, 6, 8, 7, 4, 6, 7, 9, 3, 6, 2, 5, 10, 9, 8, 10, 9, 8, 5, 10, 6, 4, 12, 7, 10, 1, 5, 3},
	{6, 2, 9, 8, 7, 9, 3, 10, 4, 5, 7, 8, 5, 7, 10, 9, 6, 11, 5, 8, 6, 5, 4, 6, 1, 10, 9, 7, 12, 8, 3, 10},
}

// reels lengths [32, 32, 32, 32, 32], total reshuffles 33554432
// symbols: 50(lined) + 17.529(scatter) = 67.529148%
// free spins 3585870, q = 0.10687, sq = 1/(1-q) = 1.119654
// free games frequency: 1/140.36
// RTP = 67.529(sym) + 0.10687*276.85(fg) = 97.115046%
var ReelsReg97 = game.Reels5x{
	{8, 11, 3, 5, 9, 10, 7, 3, 8, 4, 9, 8, 6, 10, 1, 8, 10, 7, 12, 10, 9, 5, 6, 7, 4, 6, 2, 7, 9, 8, 5, 7},
	{4, 8, 10, 9, 3, 5, 7, 6, 4, 7, 10, 3, 6, 7, 9, 10, 6, 7, 8, 10, 5, 9, 7, 8, 11, 2, 8, 9, 1, 8, 12, 5},
	{8, 9, 11, 5, 10, 6, 8, 9, 7, 10, 5, 7, 10, 5, 1, 9, 6, 2, 5, 6, 3, 8, 10, 3, 8, 6, 4, 7, 12, 4, 9, 7},
	{4, 5, 7, 4, 5, 10, 12, 6, 3, 9, 6, 7, 8, 9, 3, 8, 10, 11, 8, 9, 10, 7, 8, 1, 5, 2, 10, 8, 7, 6, 9, 7},
	{8, 5, 7, 6, 9, 8, 11, 10, 6, 4, 8, 7, 3, 10, 7, 5, 4, 12, 5, 9, 7, 1, 10, 7, 6, 8, 9, 10, 2, 8, 9, 3},
}

// reels lengths [32, 32, 32, 32, 32], total reshuffles 33554432
// symbols: 52.077(lined) + 17.529(scatter) = 69.606400%
// free spins 3585870, q = 0.10687, sq = 1/(1-q) = 1.119654
// free games frequency: 1/140.36
// RTP = 69.606(sym) + 0.10687*276.85(fg) = 99.192298%
var ReelsReg99 = game.Reels5x{
	{8, 11, 3, 5, 9, 10, 7, 3, 8, 4, 9, 8, 6, 10, 1, 8, 10, 7, 12, 10, 9, 5, 6, 7, 4, 6, 2, 7, 9, 8, 5, 7},
	{4, 8, 10, 9, 3, 5, 7, 6, 4, 7, 10, 3, 6, 7, 9, 10, 6, 7, 8, 10, 5, 9, 7, 8, 11, 2, 8, 9, 1, 8, 12, 5},
	{8, 6, 7, 10, 6, 1, 9, 8, 7, 11, 8, 5, 9, 12, 7, 10, 3, 8, 4, 9, 3, 5, 7, 4, 10, 8, 7, 9, 2, 6, 10, 5},
	{4, 5, 7, 4, 5, 10, 12, 6, 3, 9, 6, 7, 8, 9, 3, 8, 10, 11, 8, 9, 10, 7, 8, 1, 5, 2, 10, 8, 7, 6, 9, 7},
	{8, 5, 7, 6, 9, 8, 11, 10, 6, 4, 8, 7, 3, 10, 7, 5, 4, 12, 5, 9, 7, 1, 10, 7, 6, 8, 9, 10, 2, 8, 9, 3},
}

// reels lengths [30, 30, 30, 30, 30], total reshuffles 24300000
// symbols: 51.358(lined) + 20.03(scatter) = 71.387942%
// free spins 3120120, q = 0.1284, sq = 1/(1-q) = 1.147315
// free games frequency: 1/116.82
// RTP = 71.388(sym) + 0.1284*276.85(fg) = 106.935121%
var ReelsReg107 = game.Reels5x{
	{5, 8, 10, 2, 5, 6, 10, 7, 3, 9, 12, 4, 10, 5, 6, 9, 8, 7, 3, 8, 7, 9, 1, 8, 10, 9, 6, 11, 7, 4},
	{10, 1, 8, 7, 3, 8, 6, 2, 5, 4, 10, 9, 7, 11, 9, 10, 6, 7, 9, 5, 8, 10, 9, 3, 5, 6, 4, 8, 7, 12},
	{1, 6, 4, 12, 5, 2, 10, 8, 9, 7, 4, 6, 5, 10, 8, 5, 7, 9, 8, 10, 9, 3, 8, 7, 3, 6, 9, 7, 10, 11},
	{8, 9, 2, 11, 5, 4, 7, 9, 10, 6, 5, 12, 1, 9, 7, 5, 10, 8, 3, 7, 4, 10, 7, 8, 10, 6, 3, 8, 9, 6},
	{11, 7, 10, 8, 1, 6, 3, 10, 8, 6, 10, 7, 9, 3, 6, 4, 9, 8, 2, 12, 5, 8, 9, 5, 7, 10, 5, 9, 7, 4},
}

// reels lengths [28, 28, 28, 28, 28], total reshuffles 17210368
// symbols: 164.25(lined) + 69.381(scatter) = 233.627613%
// free spins 2686770, q = 0.15611, sq = 1/(1-q) = 1.184993
// free games frequency: 1/96.084
// RTP = sq*rtp(sym) = 1.185*233.63 = 276.847183%
var ReelsBon = game.Reels5x{
	{8, 7, 5, 11, 7, 10, 4, 12, 8, 6, 3, 10, 1, 9, 3, 5, 10, 9, 6, 7, 9, 2, 5, 10, 8, 9, 6, 4},
	{3, 8, 9, 2, 10, 9, 7, 6, 11, 4, 10, 5, 8, 3, 10, 7, 9, 5, 6, 1, 9, 10, 4, 12, 7, 5, 8, 6},
	{10, 9, 6, 4, 10, 8, 1, 7, 5, 9, 12, 3, 9, 5, 8, 10, 7, 4, 10, 9, 11, 3, 6, 7, 8, 6, 2, 5},
	{10, 6, 9, 7, 5, 6, 9, 2, 5, 9, 6, 3, 8, 1, 10, 4, 9, 3, 12, 10, 7, 5, 10, 8, 7, 4, 11, 8},
	{5, 9, 1, 10, 8, 7, 3, 5, 2, 11, 8, 9, 5, 4, 10, 9, 6, 3, 12, 4, 7, 6, 10, 9, 8, 7, 6, 10},
}

// Map with available reels.
var ReelsMap = map[string]*game.Reels5x{
	"85":  &ReelsReg85,
	"88":  &ReelsReg88,
	"90":  &ReelsReg90,
	"91":  &ReelsReg91,
	"92":  &ReelsReg92,
	"93":  &ReelsReg93,
	"95":  &ReelsReg95,
	"96":  &ReelsReg96,
	"97":  &ReelsReg97,
	"99":  &ReelsReg99,
	"107": &ReelsReg107,
	"bon": &ReelsBon,
}

// Lined payment.
var LinePay = [12][5]float64{
	{0, 2, 25, 125, 2000},     //  1 knife
	{0, 2, 25, 125, 1000},     //  2 sneakers
	{0, 2, 10, 75, 500},       //  3 tent
	{0, 2, 10, 75, 300},       //  4 drum
	{0, 0, 5, 30, 150},        //  5 camel
	{0, 0, 5, 30, 150},        //  6 king
	{0, 0, 5, 20, 125},        //  7 queen
	{0, 0, 5, 20, 125},        //  8 jack
	{0, 0, 3, 15, 75},         //  9 ten
	{0, 0, 3, 15, 75},         // 10 nine
	{0, 10, 200, 2500, 10000}, // 11 wild
	{0, 0, 0, 0, 0},           // 12 scatter
}

// Scatters payment.
var ScatPay = [5]float64{0, 2, 5, 20, 500} // 12 scatter

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 15, 15, 15} // 12 scatter

const (
	jid = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [12][5]int{
	{0, 0, 0, 0, 0}, //  1 knife
	{0, 0, 0, 0, 0}, //  2 sneakers
	{0, 0, 0, 0, 0}, //  3 tent
	{0, 0, 0, 0, 0}, //  4 drum
	{0, 0, 0, 0, 0}, //  5 camel
	{0, 0, 0, 0, 0}, //  6 king
	{0, 0, 0, 0, 0}, //  7 queen
	{0, 0, 0, 0, 0}, //  8 jack
	{0, 0, 0, 0, 0}, //  9 ten
	{0, 0, 0, 0, 0}, // 10 nine
	{0, 0, 0, 0, 0}, // 11 wild
	{0, 0, 0, 0, 0}, // 12 scatter
}

type Game struct {
	game.Slot5x3 `yaml:",inline"`
	// free spin number
	FS int `json:"fs,omitempty" yaml:"fs,omitempty" xml:"fs,omitempty"`
}

func NewGame(rd string) *Game {
	return &Game{
		Slot5x3: game.Slot5x3{
			RD:  rd,
			SBL: game.MakeBitNum(10),
			Bet: 1,
		},
		FS: 0,
	}
}

const wild, scat = 11, 12

var bl = game.BetLinesNetEnt10

func (g *Game) Scanner(screen game.Screen, wins *game.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen game.Screen, wins *game.Wins) {
	for li := g.SBL.Next(0); li != 0; li = g.SBL.Next(li) {
		var line = bl.Line(li)

		var numw, numl = 0, 5
		var syml game.Sym
		var mw float64 = 1 // mult wild
		for x := 1; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx == wild {
				if syml == 0 {
					numw = x
				}
				mw = 2
			} else if syml == 0 && sx != scat {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		var payw, payl float64
		if numw > 0 {
			payw = LinePay[wild-1][numw-1]
		}
		if numl > 0 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl*mw > payw {
			var mm float64 = 1 // mult mode
			if g.FS > 0 {
				mm = 3
			}
			*wins = append(*wins, game.WinItem{
				Pay:  g.Bet * payl,
				Mult: mw * mm,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			var mm float64 = 1 // mult mode
			if g.FS > 0 {
				mm = 3
			}
			*wins = append(*wins, game.WinItem{
				Pay:  g.Bet * payw,
				Mult: mm,
				Sym:  wild,
				Num:  numw,
				Line: li,
				XY:   line.CopyL(numw),
				Jack: Jackpot[wild-1][numw-1],
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen game.Screen, wins *game.Wins) {
	if count := screen.ScatNum(scat); count >= 2 {
		var mm float64 = 1 // mult mode
		if g.FS > 0 {
			mm = 3
		}
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		*wins = append(*wins, game.WinItem{
			Pay:  g.Bet * float64(g.SBL.Num()) * pay,
			Mult: mm,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: fs,
		})
	}
}

func (g *Game) Spin(screen game.Screen) {
	if g.FS == 0 {
		screen.Spin(ReelsMap[g.RD])
	} else {
		screen.Spin(&ReelsBon)
	}
}

func (g *Game) Apply(screen game.Screen, wins game.Wins) {
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

func (g *Game) SetReels(rd string) error {
	if _, ok := ReelsMap[rd]; !ok {
		return game.ErrNoReels
	}
	g.RD = rd
	return nil
}
