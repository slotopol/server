package diamonddogs

import (
	"math"

	slot "github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/util"
)

// reels lengths [35, 35, 35, 36, 36], total reshuffles 55566000
// symbols: 32.588(lined) + 13.69(scatter) = 46.278440%
// free spins 2967840, q = 0.053411, sq = 1/(1-q) = 1.056425
// free games frequency: 1/187.23
// ne12 bonuses: frequency 1/1786.5, rtp = 16.793003%
// RTP = 46.278(sym) + 16.793(ne12) + 0.053411*475.07(fg) = 88.445592%
var ReelsReg88 = slot.Reels5x{
	{3, 6, 5, 8, 4, 1, 11, 6, 8, 5, 7, 8, 10, 2, 1, 6, 9, 8, 2, 6, 4, 7, 9, 5, 7, 3, 5, 2, 1, 7, 8, 3, 7, 9, 4},
	{11, 6, 1, 3, 7, 6, 4, 1, 9, 7, 2, 6, 8, 10, 5, 6, 2, 8, 5, 7, 8, 5, 9, 3, 8, 9, 1, 4, 8, 3, 4, 7, 5, 2, 7},
	{5, 11, 7, 6, 3, 8, 2, 4, 7, 9, 1, 7, 2, 8, 7, 6, 9, 5, 4, 8, 1, 5, 10, 6, 8, 2, 5, 7, 1, 8, 4, 3, 6, 4, 3},
	{1, 2, 8, 6, 5, 4, 9, 8, 1, 7, 10, 1, 6, 7, 4, 9, 7, 4, 8, 2, 3, 6, 2, 11, 5, 8, 3, 7, 5, 9, 6, 3, 5, 7, 8, 4},
	{9, 2, 7, 5, 1, 8, 3, 6, 11, 7, 4, 8, 7, 6, 5, 8, 6, 7, 9, 4, 2, 7, 1, 9, 3, 5, 4, 10, 5, 3, 6, 8, 2, 1, 4, 8},
}

// reels lengths [34, 35, 35, 36, 36], total reshuffles 53978400
// symbols: 33.551(lined) + 13.852(scatter) = 47.403382%
// free spins 2931930, q = 0.054317, sq = 1/(1-q) = 1.057436
// free games frequency: 1/184.11
// ne12 bonuses: frequency 1/1735.4, rtp = 17.286915%
// RTP = 47.403(sym) + 17.287(ne12) + 0.054317*475.07(fg) = 90.494694%
var ReelsReg90 = slot.Reels5x{
	{7, 9, 8, 7, 9, 4, 6, 5, 2, 7, 6, 8, 10, 5, 1, 11, 8, 4, 3, 6, 5, 4, 8, 2, 3, 1, 2, 7, 5, 1, 9, 4, 6, 3},
	{11, 6, 1, 3, 7, 6, 4, 1, 9, 7, 2, 6, 8, 10, 5, 6, 2, 8, 5, 7, 8, 5, 9, 3, 8, 9, 1, 4, 8, 3, 4, 7, 5, 2, 7},
	{5, 11, 7, 6, 3, 8, 2, 4, 7, 9, 1, 7, 2, 8, 7, 6, 9, 5, 4, 8, 1, 5, 10, 6, 8, 2, 5, 7, 1, 8, 4, 3, 6, 4, 3},
	{1, 2, 8, 6, 5, 4, 9, 8, 1, 7, 10, 1, 6, 7, 4, 9, 7, 4, 8, 2, 3, 6, 2, 11, 5, 8, 3, 7, 5, 9, 6, 3, 5, 7, 8, 4},
	{9, 2, 7, 5, 1, 8, 3, 6, 11, 7, 4, 8, 7, 6, 5, 8, 6, 7, 9, 4, 2, 7, 1, 9, 3, 5, 4, 10, 5, 3, 6, 8, 2, 1, 4, 8},
}

// reels lengths [34, 34, 35, 36, 36], total reshuffles 52436160
// symbols: 34.915(lined) + 14.015(scatter) = 48.930028%
// free spins 2896290, q = 0.055235, sq = 1/(1-q) = 1.058464
// free games frequency: 1/181.05
// ne12 bonuses: frequency 1/1685.8, rtp = 17.795353%
// RTP = 48.93(sym) + 17.795(ne12) + 0.055235*475.07(fg) = 92.965833%
var ReelsReg93 = slot.Reels5x{
	{7, 9, 8, 7, 9, 4, 6, 5, 2, 7, 6, 8, 10, 5, 1, 11, 8, 4, 3, 6, 5, 4, 8, 2, 3, 1, 2, 7, 5, 1, 9, 4, 6, 3},
	{5, 6, 11, 4, 7, 1, 8, 7, 6, 9, 7, 2, 3, 1, 9, 8, 4, 2, 5, 8, 4, 2, 5, 3, 7, 1, 10, 3, 6, 4, 5, 6, 9, 8},
	{2, 5, 7, 6, 1, 7, 8, 4, 3, 1, 2, 9, 8, 5, 4, 6, 10, 4, 2, 8, 6, 7, 1, 9, 6, 8, 7, 5, 3, 8, 5, 3, 11, 7, 4},
	{9, 7, 6, 9, 4, 5, 8, 1, 4, 11, 7, 3, 5, 7, 8, 6, 10, 1, 4, 2, 7, 8, 2, 4, 8, 1, 9, 8, 3, 2, 6, 5, 3, 7, 5, 6},
	{9, 3, 6, 9, 8, 2, 11, 3, 7, 5, 2, 8, 4, 5, 7, 4, 1, 7, 8, 10, 5, 8, 1, 6, 7, 2, 6, 1, 4, 5, 9, 7, 4, 6, 3, 8},
}

// reels lengths [34, 34, 34, 36, 36], total reshuffles 50937984
// symbols: 34.791(lined) + 14.179(scatter) = 48.969592%
// free spins 2860920, q = 0.056165, sq = 1/(1-q) = 1.059507
// free games frequency: 1/178.05
// ne12 bonuses: frequency 1/1637.7, rtp = 18.318746%
// RTP = 48.97(sym) + 18.319(ne12) + 0.056165*475.07(fg) = 93.970689%
var ReelsReg94 = slot.Reels5x{
	{7, 9, 8, 7, 9, 4, 6, 5, 2, 7, 6, 8, 10, 5, 1, 11, 8, 4, 3, 6, 5, 4, 8, 2, 3, 1, 2, 7, 5, 1, 9, 4, 6, 3},
	{5, 6, 11, 4, 7, 1, 8, 7, 6, 9, 7, 2, 3, 1, 9, 8, 4, 2, 5, 8, 4, 2, 5, 3, 7, 1, 10, 3, 6, 4, 5, 6, 9, 8},
	{4, 7, 5, 1, 8, 11, 6, 4, 5, 9, 2, 7, 3, 6, 2, 1, 8, 5, 7, 3, 6, 4, 9, 1, 8, 10, 2, 7, 6, 8, 7, 5, 3, 8},
	{9, 7, 6, 9, 4, 5, 8, 1, 4, 11, 7, 3, 5, 7, 8, 6, 10, 1, 4, 2, 7, 8, 2, 4, 8, 1, 9, 8, 3, 2, 6, 5, 3, 7, 5, 6},
	{9, 3, 6, 9, 8, 2, 11, 3, 7, 5, 2, 8, 4, 5, 7, 4, 1, 7, 8, 10, 5, 8, 1, 6, 7, 2, 6, 1, 4, 5, 9, 7, 4, 6, 3, 8},
}

// reels lengths [34, 34, 34, 35, 35], total reshuffles 48147400
// symbols: 34.83(lined) + 14.494(scatter) = 49.324545%
// free spins 2791530, q = 0.057979, sq = 1/(1-q) = 1.061547
// free games frequency: 1/172.48
// ne12 bonuses: frequency 1/1637.7, rtp = 18.318746%
// RTP = 49.325(sym) + 18.319(ne12) + 0.057979*475.07(fg) = 95.187455%
var ReelsReg95 = slot.Reels5x{
	{7, 9, 8, 7, 9, 4, 6, 5, 2, 7, 6, 8, 10, 5, 1, 11, 8, 4, 3, 6, 5, 4, 8, 2, 3, 1, 2, 7, 5, 1, 9, 4, 6, 3},
	{5, 6, 11, 4, 7, 1, 8, 7, 6, 9, 7, 2, 3, 1, 9, 8, 4, 2, 5, 8, 4, 2, 5, 3, 7, 1, 10, 3, 6, 4, 5, 6, 9, 8},
	{4, 7, 5, 1, 8, 11, 6, 4, 5, 9, 2, 7, 3, 6, 2, 1, 8, 5, 7, 3, 6, 4, 9, 1, 8, 10, 2, 7, 6, 8, 7, 5, 3, 8},
	{8, 7, 2, 5, 7, 3, 6, 7, 1, 5, 8, 3, 6, 5, 9, 3, 4, 6, 2, 7, 8, 1, 7, 4, 8, 9, 6, 4, 9, 8, 1, 10, 2, 5, 11},
	{3, 8, 5, 10, 2, 6, 8, 3, 1, 7, 8, 6, 2, 11, 7, 6, 3, 9, 8, 1, 4, 5, 9, 4, 6, 9, 7, 5, 8, 4, 7, 2, 5, 7, 1},
}

// reels lengths [34, 34, 35, 35, 35], total reshuffles 49563500
// symbols: 36.801(lined) + 14.328(scatter) = 51.128942%
// free spins 2826360, q = 0.057025, sq = 1/(1-q) = 1.060474
// free games frequency: 1/175.36
// ne12 bonuses: frequency 1/1685.8, rtp = 17.795353%
// RTP = 51.129(sym) + 17.795(ne12) + 0.057025*475.07(fg) = 96.015333%
var ReelsReg96 = slot.Reels5x{
	{7, 9, 8, 7, 9, 4, 6, 5, 2, 7, 6, 8, 10, 5, 1, 11, 8, 4, 3, 6, 5, 4, 8, 2, 3, 1, 2, 7, 5, 1, 9, 4, 6, 3},
	{5, 6, 11, 4, 7, 1, 8, 7, 6, 9, 7, 2, 3, 1, 9, 8, 4, 2, 5, 8, 4, 2, 5, 3, 7, 1, 10, 3, 6, 4, 5, 6, 9, 8},
	{2, 7, 1, 3, 9, 8, 1, 10, 2, 4, 5, 3, 9, 7, 2, 3, 5, 4, 8, 7, 6, 4, 8, 6, 2, 1, 6, 5, 7, 4, 8, 11, 5, 3, 6},
	{2, 7, 3, 4, 6, 10, 5, 6, 2, 5, 8, 4, 7, 1, 5, 3, 11, 4, 1, 7, 3, 6, 8, 1, 6, 8, 2, 7, 9, 4, 8, 7, 9, 8, 5},
	{5, 4, 8, 6, 7, 1, 4, 11, 5, 7, 8, 2, 6, 9, 7, 1, 8, 2, 7, 5, 4, 8, 2, 1, 7, 3, 8, 6, 3, 5, 10, 3, 4, 6, 9},
}

// reels lengths [34, 34, 33, 36, 36], total reshuffles 49439808
// symbols: 36.953(lined) + 14.353(scatter) = 51.305454%
// free spins 2825550, q = 0.057151, sq = 1/(1-q) = 1.060616
// free games frequency: 1/174.97
// ne12 bonuses: frequency 1/1589.5, rtp = 18.873860%
// RTP = 51.305(sym) + 18.874(ne12) + 0.057151*475.07(fg) = 97.330346%
var ReelsReg97 = slot.Reels5x{
	{7, 9, 8, 7, 9, 4, 6, 5, 2, 7, 6, 8, 10, 5, 1, 11, 8, 4, 3, 6, 5, 4, 8, 2, 3, 1, 2, 7, 5, 1, 9, 4, 6, 3},
	{5, 6, 11, 4, 7, 1, 8, 7, 6, 9, 7, 2, 3, 1, 9, 8, 4, 2, 5, 8, 4, 2, 5, 3, 7, 1, 10, 3, 6, 4, 5, 6, 9, 8},
	{4, 9, 7, 2, 1, 4, 8, 7, 6, 1, 8, 5, 6, 2, 8, 10, 4, 7, 3, 11, 2, 5, 3, 6, 8, 9, 4, 7, 5, 6, 1, 5, 3},
	{8, 9, 3, 6, 2, 7, 4, 6, 5, 3, 8, 5, 1, 8, 3, 4, 9, 6, 4, 5, 11, 1, 7, 2, 9, 7, 4, 10, 2, 5, 6, 8, 7, 1, 3, 2},
	{1, 11, 2, 5, 6, 9, 4, 8, 2, 6, 3, 4, 5, 8, 3, 5, 1, 7, 8, 6, 10, 3, 7, 5, 2, 9, 4, 1, 3, 9, 7, 4, 6, 8, 7, 2},
}

// reels lengths [34, 34, 33, 35, 35], total reshuffles 46731300
// symbols: 36.426(lined) + 14.671(scatter) = 51.096406%
// free spins 2756700, q = 0.05899, sq = 1/(1-q) = 1.062688
// free games frequency: 1/169.52
// ne12 bonuses: frequency 1/1589.5, rtp = 18.873860%
// RTP = 51.096(sym) + 18.874(ne12) + 0.05899*475.07(fg) = 97.995018%
var ReelsReg98 = slot.Reels5x{
	{7, 9, 8, 7, 9, 4, 6, 5, 2, 7, 6, 8, 10, 5, 1, 11, 8, 4, 3, 6, 5, 4, 8, 2, 3, 1, 2, 7, 5, 1, 9, 4, 6, 3},
	{5, 6, 11, 4, 7, 1, 8, 7, 6, 9, 7, 2, 3, 1, 9, 8, 4, 2, 5, 8, 4, 2, 5, 3, 7, 1, 10, 3, 6, 4, 5, 6, 9, 8},
	{4, 9, 7, 2, 1, 4, 8, 7, 6, 1, 8, 5, 6, 2, 8, 10, 4, 7, 3, 11, 2, 5, 3, 6, 8, 9, 4, 7, 5, 6, 1, 5, 3},
	{1, 3, 2, 5, 6, 9, 5, 3, 7, 2, 11, 4, 6, 8, 7, 1, 9, 8, 6, 7, 9, 4, 6, 1, 3, 9, 5, 2, 4, 5, 8, 10, 7, 8, 4},
	{2, 9, 7, 8, 2, 4, 9, 1, 2, 10, 6, 3, 5, 7, 1, 5, 8, 6, 7, 4, 5, 6, 9, 4, 5, 8, 3, 9, 4, 7, 8, 6, 11, 1, 3},
}

// reels lengths [34, 34, 35, 33, 33], total reshuffles 44060940
// symbols: 37.483(lined) + 15.019(scatter) = 52.502377%
// free spins 2688120, q = 0.061009, sq = 1/(1-q) = 1.064973
// free games frequency: 1/163.91
// ne12 bonuses: frequency 1/1685.8, rtp = 17.795353%
// RTP = 52.502(sym) + 17.795(ne12) + 0.061009*475.07(fg) = 99.281511%
var ReelsReg99 = slot.Reels5x{
	{7, 9, 8, 7, 9, 4, 6, 5, 2, 7, 6, 8, 10, 5, 1, 11, 8, 4, 3, 6, 5, 4, 8, 2, 3, 1, 2, 7, 5, 1, 9, 4, 6, 3},
	{5, 6, 11, 4, 7, 1, 8, 7, 6, 9, 7, 2, 3, 1, 9, 8, 4, 2, 5, 8, 4, 2, 5, 3, 7, 1, 10, 3, 6, 4, 5, 6, 9, 8},
	{2, 6, 9, 3, 4, 8, 2, 5, 8, 11, 7, 1, 6, 7, 9, 3, 5, 4, 1, 6, 4, 2, 7, 3, 8, 1, 2, 7, 5, 8, 6, 4, 10, 5, 3},
	{10, 4, 8, 1, 7, 5, 8, 6, 1, 5, 8, 4, 7, 3, 2, 5, 6, 2, 8, 9, 6, 2, 9, 7, 4, 11, 6, 3, 5, 4, 7, 3, 1},
	{4, 7, 10, 8, 2, 4, 9, 7, 4, 3, 8, 7, 5, 1, 3, 2, 6, 5, 8, 6, 4, 5, 3, 7, 6, 1, 9, 8, 6, 11, 1, 5, 2},
}

// reels lengths [34, 34, 33, 34, 34], total reshuffles 44099088
// symbols: 36.968(lined) + 15.01(scatter) = 51.977696%
// free spins 2688390, q = 0.060962, sq = 1/(1-q) = 1.064920
// free games frequency: 1/164.04
// ne12 bonuses: frequency 1/1589.5, rtp = 18.873860%
// RTP = 51.978(sym) + 18.874(ne12) + 0.060962*475.07(fg) = 99.813172%
var ReelsReg100 = slot.Reels5x{
	{7, 9, 8, 7, 9, 4, 6, 5, 2, 7, 6, 8, 10, 5, 1, 11, 8, 4, 3, 6, 5, 4, 8, 2, 3, 1, 2, 7, 5, 1, 9, 4, 6, 3},
	{5, 6, 11, 4, 7, 1, 8, 7, 6, 9, 7, 2, 3, 1, 9, 8, 4, 2, 5, 8, 4, 2, 5, 3, 7, 1, 10, 3, 6, 4, 5, 6, 9, 8},
	{4, 9, 7, 2, 1, 4, 8, 7, 6, 1, 8, 5, 6, 2, 8, 10, 4, 7, 3, 11, 2, 5, 3, 6, 8, 9, 4, 7, 5, 6, 1, 5, 3},
	{5, 1, 9, 8, 6, 1, 7, 11, 4, 5, 6, 7, 5, 3, 9, 8, 4, 6, 9, 1, 3, 7, 8, 4, 2, 7, 6, 3, 2, 5, 8, 4, 2, 10},
	{7, 8, 6, 4, 7, 2, 5, 1, 10, 2, 8, 6, 9, 5, 4, 3, 6, 7, 3, 8, 11, 5, 6, 1, 9, 4, 1, 5, 3, 7, 4, 2, 8, 9},
}

// reels lengths [33, 33, 32, 31, 31], total reshuffles 33488928
// symbols: 38.7(lined) + 16.76(scatter) = 55.460485%
// free spins 2390040, q = 0.071368, sq = 1/(1-q) = 1.076853
// free games frequency: 1/140.12
// ne12 bonuses: frequency 1/1452, rtp = 20.661157%
// RTP = 55.46(sym) + 20.661(ne12) + 0.071368*475.07(fg) = 110.026666%
var ReelsReg110 = slot.Reels5x{
	{8, 7, 1, 6, 7, 4, 8, 5, 9, 1, 5, 2, 4, 6, 7, 5, 3, 4, 10, 3, 8, 5, 9, 1, 8, 2, 6, 11, 3, 6, 2, 7, 9},
	{2, 1, 7, 6, 4, 9, 8, 5, 7, 8, 9, 7, 6, 8, 5, 9, 2, 8, 6, 3, 1, 4, 3, 11, 5, 3, 2, 4, 6, 5, 7, 1, 10},
	{8, 6, 7, 5, 1, 2, 3, 11, 2, 1, 3, 4, 8, 10, 1, 7, 4, 8, 6, 5, 2, 3, 8, 4, 6, 5, 7, 6, 9, 5, 7, 9},
	{3, 1, 8, 6, 11, 8, 6, 5, 3, 2, 7, 8, 3, 2, 6, 7, 5, 6, 4, 7, 9, 1, 7, 2, 4, 1, 5, 8, 10, 5, 4},
	{9, 1, 3, 4, 6, 1, 2, 5, 6, 7, 2, 4, 7, 3, 8, 7, 10, 6, 5, 4, 8, 3, 11, 8, 1, 5, 6, 2, 8, 5, 7},
}

// reels lengths [40, 40, 40, 40, 40], total reshuffles 102400000
// symbols: 218.06(lined) + 130.58(scatter) = 348.646992%
// free spins 27250560, q = 0.26612, sq = 1/(1-q) = 1.362618
// free games frequency: 1/37.577
// RTP = sq*rtp(sym) = 1.3626*348.65 = 475.072762%
var ReelsBon = slot.Reels5x{
	{3, 1, 7, 3, 8, 10, 6, 7, 8, 6, 4, 7, 5, 1, 4, 3, 11, 1, 8, 5, 2, 6, 8, 5, 6, 2, 5, 7, 2, 6, 10, 2, 4, 5, 1, 7, 8, 3, 4, 11},
	{5, 2, 3, 1, 8, 7, 10, 6, 8, 5, 7, 11, 4, 8, 3, 2, 4, 5, 1, 7, 5, 3, 6, 11, 1, 8, 5, 4, 1, 6, 2, 4, 10, 7, 6, 2, 8, 3, 7, 6},
	{5, 2, 7, 8, 1, 7, 2, 3, 6, 5, 4, 3, 10, 7, 6, 5, 7, 1, 6, 11, 3, 2, 8, 3, 4, 8, 7, 4, 11, 6, 1, 5, 8, 4, 2, 1, 8, 5, 6, 10},
	{2, 1, 6, 3, 11, 8, 7, 5, 8, 10, 7, 4, 6, 7, 2, 5, 8, 3, 4, 5, 1, 11, 7, 2, 4, 7, 8, 1, 2, 3, 1, 6, 5, 8, 6, 10, 5, 6, 3, 4},
	{8, 5, 6, 8, 11, 1, 4, 3, 6, 7, 4, 10, 2, 3, 5, 1, 2, 8, 6, 5, 4, 8, 1, 7, 10, 2, 3, 7, 2, 6, 8, 5, 11, 6, 3, 4, 7, 1, 5, 7},
}

// Map with available reels.
var ReelsMap = map[float64]*slot.Reels5x{
	88.445592:  &ReelsReg88,
	90.494694:  &ReelsReg90,
	92.965833:  &ReelsReg93,
	93.970689:  &ReelsReg94,
	95.187455:  &ReelsReg95,
	96.015333:  &ReelsReg96,
	97.330346:  &ReelsReg97,
	97.995018:  &ReelsReg98,
	99.281511:  &ReelsReg99,
	99.813172:  &ReelsReg100,
	110.026666: &ReelsReg110,
}

func FindReels(mrtp float64) (rtp float64, reels slot.Reels) {
	for p, r := range ReelsMap {
		if math.Abs(mrtp-p) < math.Abs(mrtp-rtp) {
			rtp, reels = p, r
		}
	}
	return
}

// Lined payment.
var LinePay = [11][5]float64{
	{0, 0, 50, 120, 600},     //  1 booth
	{0, 0, 15, 90, 240},      //  2 vip
	{0, 0, 15, 90, 240},      //  3 food
	{0, 0, 10, 60, 120},      //  4 bell
	{0, 0, 5, 60, 120},       //  5 ace
	{0, 0, 5, 30, 90},        //  6 king
	{0, 0, 2, 12, 60},        //  7 queen
	{0, 0, 2, 12, 60},        //  8 jack
	{0, 0, 0, 0, 0},          //  9 bonus
	{0, 5, 200, 2000, 10000}, // 10 wild
	{0, 0, 0, 0, 0},          // 11 scatter
}

// Scatters payment.
var ScatPay = [5]float64{0, 2, 4, 25, 100} // 11 scatter

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 10, 10, 10} // 11 scatter

const (
	ne12 = 1 // bonus ID
)

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
	// free spin number
	FS int `json:"fs,omitempty" yaml:"fs,omitempty" xml:"fs,omitempty"`
}

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: util.MakeBitNum(25, 1),
			Bet: 1,
		},
		FS: 0,
	}
}

const bon, wild, scat = 9, 10, 11

var bl = slot.BetLinesBetSoft25

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	for li := range g.Sel.Bits() {
		var line = bl.Line(li)

		var numw, numl = 0, 5
		var syml slot.Sym
		for x := 1; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx == wild {
				if syml == 0 {
					numw = x
				} else if syml == bon {
					numl = x - 1
					break
				}
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
		if payl > payw {
			var mm float64 = 1 // mult mode
			if g.FS > 0 {
				mm = 3
			}
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: mm,
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
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: mm,
				Sym:  wild,
				Num:  numw,
				Line: li,
				XY:   line.CopyL(numw),
			})
		} else if syml == bon && numl >= 3 { // appear on regular games only
			*wins = append(*wins, slot.WinItem{
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
				BID:  ne12,
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen slot.Screen, wins *slot.Wins) {
	if count := screen.ScatNum(scat); count >= 2 {
		var mm float64 = 1 // mult mode
		if g.FS > 0 {
			mm = 3
		}
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel.Num()) * pay,
			Mult: mm,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: fs,
		})
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	if g.FS == 0 {
		var _, reels = FindReels(mrtp)
		screen.Spin(reels)
	} else {
		screen.Spin(&ReelsBon)
	}
}

func (g *Game) Spawn(screen slot.Screen, wins slot.Wins) {
	for i, wi := range wins {
		switch wi.BID {
		case ne12:
			wins[i].Bon, wins[i].Pay = BonusSpawn(g.Bet)
		}
	}
}

func (g *Game) Apply(screen slot.Screen, wins slot.Wins) {
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

func (g *Game) SetSel(sel slot.Bitset) error {
	var mask slot.Bitset = (1<<len(bl) - 1) << 1
	if sel == 0 {
		return slot.ErrNoLineset
	}
	if sel&^mask != 0 {
		return slot.ErrLinesetOut
	}
	if g.FreeSpins() > 0 {
		return slot.ErrNoFeature
	}
	g.Sel = sel
	return nil
}
