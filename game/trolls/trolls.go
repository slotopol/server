package trolls

import (
	"github.com/slotopol/server/game"
)

// *bonus reels calculations*
// reels lengths [39, 39, 40, 39, 39], total reshuffles 92537640
// symbols: 198.11(lined) + 34.997(scatter) = 233.102719%
// free spins 3838590, q = 0.041481, sq = 1/(1-q) = 1.043277
// RTP = sq*rtp(sym) = 1.0433*233.1 = 243.190604%
// *regular reels calculations*
// reels lengths [39, 39, 40, 39, 39], total reshuffles 92537640
// symbols: 66.035(lined) + 11.666(scatter) = 77.700906%
// free spins 3838590, q = 0.041481, sq = 1/(1-q) = 1.043277
// RTP = 77.701(sym) + 0.041481*243.19(fg) = 87.788791%
var Reels88 = game.Reels5x{
	{2, 11, 3, 10, 5, 8, 1, 7, 10, 9, 6, 2, 7, 8, 1, 12, 10, 4, 6, 10, 5, 7, 9, 4, 11, 3, 7, 14, 9, 4, 11, 6, 8, 2, 9, 11, 5, 3, 8},
	{6, 11, 2, 5, 10, 14, 4, 5, 3, 7, 6, 4, 9, 3, 12, 9, 11, 3, 9, 8, 2, 5, 8, 1, 7, 10, 8, 9, 11, 10, 8, 2, 11, 10, 7, 4, 6, 1, 7},
	{13, 11, 8, 3, 11, 2, 10, 4, 7, 9, 6, 2, 10, 6, 3, 7, 4, 14, 1, 8, 11, 5, 12, 10, 4, 9, 3, 11, 9, 1, 7, 2, 9, 7, 8, 5, 6, 8, 10, 5},
	{11, 1, 7, 9, 4, 5, 8, 2, 10, 6, 1, 7, 6, 8, 11, 2, 10, 8, 9, 3, 12, 2, 5, 11, 3, 8, 4, 10, 7, 9, 4, 10, 11, 14, 7, 6, 5, 3, 9},
	{11, 14, 8, 10, 4, 12, 10, 11, 7, 9, 5, 8, 2, 11, 9, 2, 8, 3, 9, 2, 11, 7, 10, 3, 8, 10, 9, 1, 6, 5, 7, 4, 6, 1, 5, 3, 6, 4, 7},
}

// *bonus reels calculations*
// reels lengths [38, 38, 39, 38, 38], total reshuffles 81320304
// symbols: 197.82(lined) + 36.892(scatter) = 234.716725%
// free spins 3637710, q = 0.044733, sq = 1/(1-q) = 1.046828
// RTP = sq*rtp(sym) = 1.0468*234.72 = 245.708008%
// *regular reels calculations*
// reels lengths [38, 38, 39, 38, 38], total reshuffles 81320304
// symbols: 65.941(lined) + 12.297(scatter) = 78.238908%
// free spins 3637710, q = 0.044733, sq = 1/(1-q) = 1.046828
// RTP = 78.239(sym) + 0.044733*245.71(fg) = 89.230191%
var Reels89 = game.Reels5x{
	{9, 8, 7, 1, 11, 8, 4, 6, 7, 4, 6, 10, 5, 3, 6, 2, 5, 11, 10, 9, 4, 12, 2, 8, 1, 7, 9, 3, 14, 10, 7, 9, 3, 5, 11, 8, 10, 11},
	{9, 10, 5, 3, 10, 11, 5, 10, 9, 11, 3, 8, 2, 6, 4, 7, 11, 2, 7, 9, 8, 12, 6, 7, 1, 11, 8, 5, 1, 6, 4, 14, 9, 3, 7, 10, 8, 4},
	{8, 2, 6, 10, 9, 7, 8, 14, 2, 7, 5, 3, 11, 8, 1, 11, 7, 13, 9, 4, 11, 5, 9, 10, 5, 11, 8, 1, 9, 10, 4, 12, 3, 6, 10, 4, 7, 3, 6},
	{10, 3, 11, 9, 7, 4, 6, 5, 7, 3, 6, 2, 8, 10, 3, 9, 8, 14, 4, 10, 1, 5, 10, 1, 5, 8, 4, 11, 7, 2, 9, 7, 8, 11, 12, 9, 11, 6},
	{8, 3, 6, 11, 1, 8, 6, 12, 11, 9, 2, 11, 5, 10, 9, 3, 11, 7, 1, 8, 4, 9, 5, 7, 3, 8, 9, 2, 10, 4, 6, 7, 4, 10, 7, 5, 14, 10},
}

// *bonus reels calculations*
// reels lengths [37, 37, 38, 37, 37], total reshuffles 71218118
// symbols: 200.36(lined) + 38.949(scatter) = 239.312352%
// free spins 3442230, q = 0.048334, sq = 1/(1-q) = 1.050788
// RTP = sq*rtp(sym) = 1.0508*239.31 = 251.466647%
// *regular reels calculations*
// reels lengths [37, 37, 38, 37, 37], total reshuffles 71218118
// symbols: 66.788(lined) + 12.983(scatter) = 79.770784%
// free spins 3442230, q = 0.048334, sq = 1/(1-q) = 1.050788
// RTP = 79.771(sym) + 0.048334*251.47(fg) = 91.925079%
var Reels92 = game.Reels5x{
	{4, 7, 5, 11, 6, 10, 11, 4, 7, 9, 3, 5, 11, 9, 2, 10, 7, 1, 8, 7, 5, 2, 8, 10, 3, 8, 4, 14, 9, 1, 6, 9, 12, 11, 6, 8, 10},
	{4, 5, 10, 1, 7, 5, 2, 10, 14, 3, 8, 11, 7, 9, 8, 3, 11, 6, 10, 4, 9, 11, 12, 8, 9, 7, 2, 11, 4, 9, 6, 5, 1, 10, 7, 8, 6},
	{7, 8, 5, 6, 8, 1, 9, 3, 14, 4, 11, 3, 13, 5, 11, 6, 2, 7, 10, 8, 9, 7, 10, 11, 4, 5, 1, 9, 10, 12, 2, 7, 4, 6, 11, 10, 8, 9},
	{1, 9, 8, 2, 5, 7, 10, 6, 7, 3, 11, 6, 10, 4, 11, 2, 9, 10, 8, 1, 9, 7, 5, 4, 8, 9, 14, 7, 11, 4, 10, 8, 3, 6, 5, 12, 11},
	{7, 11, 2, 8, 9, 3, 10, 6, 14, 2, 9, 4, 12, 3, 6, 7, 9, 1, 5, 11, 9, 5, 7, 11, 10, 8, 4, 7, 8, 10, 11, 1, 10, 5, 4, 6, 8},
}

// *bonus reels calculations*
// reels lengths [39, 39, 40, 39, 39], total reshuffles 92537640
// symbols: 212.11(lined) + 34.997(scatter) = 247.103094%
// free spins 3838590, q = 0.041481, sq = 1/(1-q) = 1.043277
// RTP = sq*rtp(sym) = 1.0433*247.1 = 257.796867%
// *regular reels calculations*
// reels lengths [39, 39, 40, 39, 39], total reshuffles 92537640
// symbols: 70.702(lined) + 11.666(scatter) = 82.367698%
// free spins 3838590, q = 0.041481, sq = 1/(1-q) = 1.043277
// RTP = 82.368(sym) + 0.041481*257.8(fg) = 93.061471%
var Reels93 = game.Reels5x{
	{5, 6, 3, 7, 11, 4, 9, 7, 1, 9, 10, 5, 8, 4, 10, 2, 8, 11, 5, 3, 11, 7, 2, 11, 1, 10, 4, 8, 6, 12, 2, 6, 9, 14, 3, 8, 9, 10, 1},
	{6, 3, 9, 12, 2, 7, 8, 1, 14, 11, 5, 9, 1, 7, 10, 11, 8, 4, 5, 7, 8, 4, 9, 3, 10, 2, 5, 1, 6, 9, 10, 4, 6, 3, 11, 8, 10, 11, 2},
	{8, 5, 10, 11, 4, 6, 8, 4, 10, 6, 2, 13, 11, 2, 9, 10, 7, 6, 3, 8, 1, 11, 9, 2, 10, 3, 8, 1, 5, 7, 9, 3, 5, 7, 4, 9, 12, 1, 11, 14},
	{1, 8, 12, 11, 3, 6, 5, 9, 1, 5, 7, 8, 2, 9, 1, 8, 7, 4, 6, 3, 5, 9, 2, 10, 14, 2, 10, 11, 6, 8, 11, 7, 10, 4, 9, 3, 11, 4, 10},
	{3, 9, 5, 3, 8, 4, 10, 11, 2, 10, 1, 8, 4, 9, 1, 7, 11, 6, 1, 9, 11, 2, 10, 5, 11, 4, 9, 6, 7, 8, 6, 12, 3, 7, 8, 5, 10, 2, 14},
}

// *bonus reels calculations*
// reels lengths [37, 37, 38, 37, 37], total reshuffles 71218118
// symbols: 205.51(lined) + 38.949(scatter) = 244.462488%
// free spins 3442230, q = 0.048334, sq = 1/(1-q) = 1.050788
// RTP = sq*rtp(sym) = 1.0508*244.46 = 256.878350%
// *regular reels calculations*
// reels lengths [37, 37, 38, 37, 37], total reshuffles 71218118
// symbols: 68.504(lined) + 12.983(scatter) = 81.487496%
// free spins 3442230, q = 0.048334, sq = 1/(1-q) = 1.050788
// RTP = 81.487(sym) + 0.048334*256.88(fg) = 93.903358%
var Reels94 = game.Reels5x{
	{2, 11, 5, 8, 3, 11, 4, 7, 1, 6, 10, 8, 6, 3, 9, 11, 7, 10, 5, 8, 4, 7, 2, 12, 3, 11, 8, 9, 4, 5, 10, 1, 9, 14, 6, 10, 9},
	{1, 7, 3, 6, 9, 1, 7, 8, 11, 4, 6, 9, 10, 5, 3, 8, 9, 2, 14, 8, 5, 2, 10, 6, 8, 4, 11, 5, 4, 10, 12, 3, 11, 10, 7, 9, 11},
	{8, 1, 11, 10, 5, 4, 6, 10, 9, 4, 7, 13, 11, 6, 1, 10, 5, 3, 9, 7, 2, 10, 11, 9, 3, 8, 11, 2, 8, 14, 4, 7, 8, 12, 9, 5, 6, 3},
	{10, 11, 1, 7, 11, 12, 4, 5, 9, 7, 3, 8, 10, 4, 6, 11, 8, 10, 3, 5, 9, 1, 5, 11, 4, 8, 2, 14, 6, 10, 9, 8, 2, 9, 6, 3, 7},
	{10, 11, 4, 14, 9, 8, 6, 9, 1, 10, 5, 3, 7, 10, 2, 6, 5, 4, 6, 7, 11, 9, 1, 11, 8, 10, 2, 8, 5, 3, 7, 9, 11, 3, 8, 12, 4},
}

// *bonus reels calculations*
// reels lengths [37, 36, 37, 36, 37], total reshuffles 65646288
// symbols: 205.85(lined) + 40.273(scatter) = 246.125367%
// free spins 3327480, q = 0.050688, sq = 1/(1-q) = 1.053394
// RTP = sq*rtp(sym) = 1.0534*246.13 = 259.267101%
// *regular reels calculations*
// reels lengths [37, 36, 37, 36, 37], total reshuffles 65646288
// symbols: 68.617(lined) + 13.424(scatter) = 82.041789%
// free spins 3327480, q = 0.050688, sq = 1/(1-q) = 1.053394
// RTP = 82.042(sym) + 0.050688*259.27(fg) = 95.183523%
var Reels95 = game.Reels5x{
	{2, 11, 5, 8, 3, 11, 4, 7, 1, 6, 10, 8, 6, 3, 9, 11, 7, 10, 5, 8, 4, 7, 2, 12, 3, 11, 8, 9, 4, 5, 10, 1, 9, 14, 6, 10, 9},
	{9, 2, 7, 5, 12, 6, 7, 9, 11, 8, 6, 14, 11, 4, 5, 1, 10, 9, 11, 3, 7, 10, 3, 9, 10, 8, 2, 10, 4, 8, 6, 4, 8, 11, 5, 1},
	{12, 2, 6, 1, 10, 4, 13, 8, 10, 9, 3, 5, 11, 1, 9, 6, 11, 4, 5, 11, 8, 7, 2, 6, 7, 8, 5, 10, 9, 4, 11, 10, 3, 14, 9, 8, 7},
	{6, 5, 14, 9, 3, 7, 11, 4, 8, 2, 7, 4, 8, 2, 9, 8, 5, 4, 6, 9, 1, 10, 9, 11, 6, 7, 5, 10, 11, 1, 10, 12, 8, 3, 10, 11},
	{10, 11, 4, 14, 9, 8, 6, 9, 1, 10, 5, 3, 7, 10, 2, 6, 5, 4, 6, 7, 11, 9, 1, 11, 8, 10, 2, 8, 5, 3, 7, 9, 11, 3, 8, 12, 4},
}

// *bonus reels calculations*
// reels lengths [36, 36, 37, 36, 36], total reshuffles 62145792
// symbols: 207.55(lined) + 41.186(scatter) = 248.738872%
// free spins 3252150, q = 0.052331, sq = 1/(1-q) = 1.055221
// RTP = sq*rtp(sym) = 1.0552*248.74 = 262.474414%
// *regular reels calculations*
// reels lengths [36, 36, 37, 36, 36], total reshuffles 62145792
// symbols: 69.184(lined) + 13.729(scatter) = 82.912957%
// free spins 3252150, q = 0.052331, sq = 1/(1-q) = 1.055221
// RTP = 82.913(sym) + 0.052331*262.47(fg) = 96.648500%
var Reels97 = game.Reels5x{
	{11, 9, 3, 12, 2, 11, 4, 8, 2, 10, 5, 11, 10, 9, 8, 5, 6, 3, 9, 4, 7, 11, 10, 5, 6, 14, 9, 7, 1, 8, 7, 1, 6, 8, 10, 4},
	{9, 2, 7, 5, 12, 6, 7, 9, 11, 8, 6, 14, 11, 4, 5, 1, 10, 9, 11, 3, 7, 10, 3, 9, 10, 8, 2, 10, 4, 8, 6, 4, 8, 11, 5, 1},
	{12, 2, 6, 1, 10, 4, 13, 8, 10, 9, 3, 5, 11, 1, 9, 6, 11, 4, 5, 11, 8, 7, 2, 6, 7, 8, 5, 10, 9, 4, 11, 10, 3, 14, 9, 8, 7},
	{6, 5, 14, 9, 3, 7, 11, 4, 8, 2, 7, 4, 8, 2, 9, 8, 5, 4, 6, 9, 1, 10, 9, 11, 6, 7, 5, 10, 11, 1, 10, 12, 8, 3, 10, 11},
	{11, 5, 8, 1, 7, 10, 8, 11, 10, 9, 5, 14, 2, 7, 4, 9, 8, 10, 7, 4, 11, 3, 8, 6, 9, 10, 11, 4, 6, 9, 3, 6, 1, 5, 2, 12},
}

// *bonus reels calculations*
// reels lengths [37, 35, 36, 35, 37], total reshuffles 60372900
// symbols: 210.36(lined) + 41.694(scatter) = 252.056176%
// free spins 3214350, q = 0.053242, sq = 1/(1-q) = 1.056236
// RTP = sq*rtp(sym) = 1.0562*252.06 = 266.230727%
// *regular reels calculations*
// reels lengths [37, 35, 36, 35, 37], total reshuffles 60372900
// symbols: 70.121(lined) + 13.898(scatter) = 84.018725%
// free spins 3214350, q = 0.053242, sq = 1/(1-q) = 1.056236
// RTP = 84.019(sym) + 0.053242*266.23(fg) = 98.193276%
var Reels98 = game.Reels5x{
	{2, 11, 5, 8, 3, 11, 4, 7, 1, 6, 10, 8, 6, 3, 9, 11, 7, 10, 5, 8, 4, 7, 2, 12, 3, 11, 8, 9, 4, 5, 10, 1, 9, 14, 6, 10, 9},
	{3, 5, 10, 7, 8, 11, 6, 9, 3, 11, 1, 6, 2, 7, 9, 10, 5, 8, 4, 12, 10, 11, 2, 8, 1, 6, 14, 9, 5, 8, 7, 10, 4, 11, 9},
	{2, 9, 8, 4, 5, 7, 9, 10, 12, 1, 7, 9, 11, 13, 6, 4, 8, 5, 11, 7, 10, 9, 14, 5, 2, 11, 8, 3, 6, 10, 1, 11, 8, 6, 3, 10},
	{8, 11, 7, 1, 11, 8, 7, 9, 10, 2, 9, 5, 4, 6, 5, 10, 3, 14, 10, 3, 9, 4, 6, 7, 2, 12, 11, 9, 8, 6, 11, 1, 8, 5, 10},
	{10, 11, 4, 14, 9, 8, 6, 9, 1, 10, 5, 3, 7, 10, 2, 6, 5, 4, 6, 7, 11, 9, 1, 11, 8, 10, 2, 8, 5, 3, 7, 9, 11, 3, 8, 12, 4},
}

// *bonus reels calculations*
// reels lengths [35, 35, 36, 35, 35], total reshuffles 54022500
// symbols: 215.38(lined) + 43.626(scatter) = 259.010895%
// free spins 3067470, q = 0.056781, sq = 1/(1-q) = 1.060200
// RTP = sq*rtp(sym) = 1.0602*259.01 = 274.603235%
// *regular reels calculations*
// reels lengths [35, 35, 36, 35, 35], total reshuffles 54022500
// symbols: 71.795(lined) + 14.542(scatter) = 86.336965%
// free spins 3067470, q = 0.056781, sq = 1/(1-q) = 1.060200
// RTP = 86.337(sym) + 0.056781*274.6(fg) = 101.929305%
var Reels102 = game.Reels5x{
	{3, 7, 9, 8, 11, 4, 10, 14, 7, 5, 10, 2, 9, 1, 8, 10, 4, 8, 10, 6, 11, 5, 1, 7, 2, 6, 11, 5, 12, 3, 6, 9, 8, 11, 9},
	{3, 5, 10, 7, 8, 11, 6, 9, 3, 11, 1, 6, 2, 7, 9, 10, 5, 8, 4, 12, 10, 11, 2, 8, 1, 6, 14, 9, 5, 8, 7, 10, 4, 11, 9},
	{2, 9, 8, 4, 5, 7, 9, 10, 12, 1, 7, 9, 11, 13, 6, 4, 8, 5, 11, 7, 10, 9, 14, 5, 2, 11, 8, 3, 6, 10, 1, 11, 8, 6, 3, 10},
	{8, 11, 7, 1, 11, 8, 7, 9, 10, 2, 9, 5, 4, 6, 5, 10, 3, 14, 10, 3, 9, 4, 6, 7, 2, 12, 11, 9, 8, 6, 11, 1, 8, 5, 10},
	{4, 10, 1, 5, 8, 6, 9, 3, 11, 8, 1, 12, 9, 4, 7, 10, 11, 2, 5, 6, 10, 7, 11, 9, 2, 5, 11, 10, 14, 6, 8, 3, 7, 9, 8},
}

// *bonus reels calculations*
// reels lengths [34, 34, 35, 34, 34], total reshuffles 46771760
// symbols: 230.04(lined) + 46.295(scatter) = 276.334117%
// free spins 2888190, q = 0.061751, sq = 1/(1-q) = 1.065815
// RTP = sq*rtp(sym) = 1.0658*276.33 = 294.521002%
// *regular reels calculations*
// reels lengths [34, 34, 35, 34, 34], total reshuffles 46771760
// symbols: 76.68(lined) + 15.432(scatter) = 92.111372%
// free spins 2888190, q = 0.061751, sq = 1/(1-q) = 1.065815
// RTP = 92.111(sym) + 0.061751*294.52(fg) = 110.298257%
var Reels110 = game.Reels5x{
	{1, 10, 7, 9, 5, 7, 4, 6, 2, 8, 6, 11, 3, 5, 1, 8, 10, 7, 2, 11, 5, 4, 14, 3, 9, 6, 8, 11, 10, 4, 9, 12, 10, 11},
	{1, 11, 4, 6, 12, 9, 4, 5, 1, 8, 7, 5, 10, 6, 11, 2, 8, 6, 11, 3, 8, 5, 9, 2, 10, 4, 9, 10, 7, 14, 3, 11, 10, 7},
	{4, 10, 7, 13, 6, 8, 3, 10, 2, 7, 6, 3, 11, 12, 5, 2, 11, 10, 9, 8, 4, 9, 5, 11, 1, 8, 10, 1, 9, 6, 4, 5, 14, 11, 7},
	{9, 10, 12, 7, 4, 9, 6, 11, 2, 8, 11, 3, 5, 4, 10, 6, 9, 7, 4, 11, 1, 7, 5, 3, 10, 2, 11, 10, 8, 6, 5, 1, 8, 14},
	{1, 12, 9, 4, 5, 9, 6, 2, 11, 4, 9, 5, 7, 2, 10, 7, 5, 1, 10, 8, 3, 6, 11, 10, 8, 11, 4, 10, 8, 3, 7, 14, 6, 11},
}

// Map with available reels.
var ReelsMap = map[string]*game.Reels5x{
	"88":  &Reels88,
	"89":  &Reels89,
	"92":  &Reels92,
	"93":  &Reels93,
	"94":  &Reels94,
	"95":  &Reels95,
	"97":  &Reels97,
	"98":  &Reels98,
	"102": &Reels102,
	"110": &Reels110,
}

// Lined payment.
var LinePay = [14][5]float64{
	{0, 3, 25, 100, 750},      //  1 troll1
	{0, 0, 25, 100, 500},      //  2 troll2
	{0, 0, 15, 100, 500},      //  3 troll3
	{0, 0, 10, 75, 250},       //  4 troll4
	{0, 0, 10, 75, 250},       //  5 troll5
	{0, 0, 10, 50, 200},       //  6 troll6
	{0, 0, 5, 50, 150},        //  7 ace
	{0, 0, 5, 25, 125},        //  8 king
	{0, 0, 5, 25, 125},        //  9 queen
	{0, 0, 5, 25, 125},        // 10 jack
	{0, 2, 5, 25, 100},        // 11 ten
	{0, 10, 250, 2500, 10000}, // 12 wild
	{0, 0, 0, 0, 0},           // 13 golden
	{0, 0, 0, 0, 0},           // 14 scatter
}

// Scatters payment.
var ScatPay = [5]float64{0, 2, 5, 25, 500} // 14 scatter

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 10, 20, 30} // 14 scatter

const (
	jid = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [14][5]int{
	{0, 0, 0, 0, 0}, //  1 troll1
	{0, 0, 0, 0, 0}, //  2 troll2
	{0, 0, 0, 0, 0}, //  3 troll3
	{0, 0, 0, 0, 0}, //  4 troll4
	{0, 0, 0, 0, 0}, //  5 troll5
	{0, 0, 0, 0, 0}, //  6 troll6
	{0, 0, 0, 0, 0}, //  7 ace
	{0, 0, 0, 0, 0}, //  8 king
	{0, 0, 0, 0, 0}, //  9 queen
	{0, 0, 0, 0, 0}, // 10 jack
	{0, 0, 0, 0, 0}, // 11 ten
	{0, 0, 0, 0, 0}, // 12 wild
	{0, 0, 0, 0, 0}, // 13 golden
	{0, 0, 0, 0, 0}, // 14 scatter
}

type Game struct {
	game.Slot5x3 `yaml:",inline"`
	FS           int `json:"fs" yaml:"fs" xml:"fs"` // free spin number
}

func NewGame(rd string) *Game {
	return &Game{
		Slot5x3: game.Slot5x3{
			RD:  rd,
			BLI: "ne20",
			SBL: game.MakeSblNum(20),
			Bet: 1,
		},
		FS: 0,
	}
}

const wild1, wild2, scat = 12, 13, 14

func (g *Game) Scanner(screen game.Screen, ws *game.WinScan) {
	g.ScanLined(screen, ws)
	g.ScanScatters(screen, ws)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen game.Screen, ws *game.WinScan) {
	var bl = game.BetLines5x[g.BLI]
	for li := g.SBL.Next(0); li != 0; li = g.SBL.Next(li) {
		var line = bl.Line(li)

		var numw, numl = 0, 5
		var syml game.Sym
		var mw float64 = 1 // mult wild
		for x := 1; x <= 5; x++ {
			var sx = screen.At(x, line.At(x))
			if sx == wild1 {
				if syml == 0 {
					numw = x
				}
				if mw < 4 {
					mw = 2
				}
			} else if sx == wild2 {
				if syml == 0 {
					numw = x
				}
				mw = 4
			} else if syml == 0 && sx != scat {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		var payw, payl float64
		if numw > 0 {
			payw = LinePay[wild1-1][numw-1]
		}
		if numl > 0 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl*mw > payw {
			var mm float64 = 1 // mult mode
			if g.FS > 0 {
				mm = 3
			}
			ws.Wins = append(ws.Wins, game.WinItem{
				Pay:  g.Bet * payl,
				Mult: mw * mm,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyN(numl),
			})
		} else if payw > 0 {
			var mm float64 = 1 // mult mode
			if g.FS > 0 {
				mm = 3
			}
			ws.Wins = append(ws.Wins, game.WinItem{
				Pay:  g.Bet * payw,
				Mult: mm, // no multiplyer on line by double symbol
				Sym:  wild1,
				Num:  numw,
				Line: li,
				XY:   line.CopyN(numw),
				Jack: Jackpot[wild1-1][numw-1],
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen game.Screen, ws *game.WinScan) {
	if count := screen.ScatNum(scat); count >= 2 {
		var mm float64 = 1 // mult mode
		if g.FS > 0 {
			mm = 3
		}
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		ws.Wins = append(ws.Wins, game.WinItem{
			Pay:  g.Bet * pay, // independent from selected lines
			Mult: mm,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: fs,
		})
	}
}

func (g *Game) Spin(screen game.Screen) {
	screen.Spin(ReelsMap[g.RD])
}

func (g *Game) Apply(screen game.Screen, sw *game.WinScan) {
	if g.FS > 0 {
		g.Gain += sw.Gain()
	} else {
		g.Gain = sw.Gain()
	}

	if g.FS > 0 {
		g.FS--
	}
	for _, wi := range sw.Wins {
		if wi.Free > 0 {
			g.FS += wi.Free
		}
	}
}

func (g *Game) FreeSpins() int {
	return g.FS
}

func (g *Game) SetReels(rd string) error {
	if _, ok := ReelsMap[rd]; !ok {
		return game.ErrNoReels
	}
	g.RD = rd
	return nil
}
