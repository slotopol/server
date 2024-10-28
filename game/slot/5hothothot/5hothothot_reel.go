package hothothot

import (
	"github.com/slotopol/server/game/slot"
)

// *bonus reels calculations*
// reels lengths [20, 19, 20], total reshuffles 7600
// symbols: 228.04(lined) + 3.5526(scatter) = 231.592105%
// free spins 540, q = 0.071053, sq = 1/(1-q) = 1.076487
// free games frequency: 1/281.48
// RTP = sq*rtp(sym) = 1.0765*231.59 = 249.305949%
// *regular reels calculations*
// reels lengths [20, 19, 20], total reshuffles 7600
// symbols: 63.158(lined) + 3.5526(scatter) = 66.710526%
// free spins 540, q = 0.071053, sq = 1/(1-q) = 1.076487
// free games frequency: 1/281.48
// RTP = 66.711(sym) + 0.071053*249.31(fg) = 84.424370%
var Reels844 = slot.Reels3x{
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 1, 1, 1, 3, 3, 3, 6, 6, 6, 6, 7},
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 1, 1, 1, 3, 3, 3, 6, 6, 6, 7},
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 1, 1, 1, 3, 3, 3, 6, 6, 6, 6, 7},
}

// *bonus reels calculations*
// reels lengths [19, 21, 19], total reshuffles 7581
// symbols: 238.98(lined) + 3.5615(scatter) = 242.540562%
// free spins 540, q = 0.071231, sq = 1/(1-q) = 1.076694
// free games frequency: 1/280.78
// RTP = sq*rtp(sym) = 1.0767*242.54 = 261.141883%
// *regular reels calculations*
// reels lengths [19, 21, 19], total reshuffles 7581
// symbols: 65.888(lined) + 3.5615(scatter) = 69.449941%
// free spins 540, q = 0.071231, sq = 1/(1-q) = 1.076694
// free games frequency: 1/280.78
// RTP = 69.45(sym) + 0.071231*261.14(fg) = 88.051262%
var Reels880 = slot.Reels3x{
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 1, 1, 1, 3, 3, 3, 6, 6, 6, 7},
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 4, 1, 1, 1, 3, 3, 3, 3, 6, 6, 6, 7},
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 1, 1, 1, 3, 3, 3, 6, 6, 6, 7},
}

// *bonus reels calculations*
// reels lengths [19, 20, 19], total reshuffles 7220
// symbols: 240.96(lined) + 3.7396(scatter) = 244.695291%
// free spins 540, q = 0.074792, sq = 1/(1-q) = 1.080838
// free games frequency: 1/267.41
// RTP = sq*rtp(sym) = 1.0808*244.7 = 264.476048%
// *regular reels calculations*
// reels lengths [19, 20, 19], total reshuffles 7220
// symbols: 66.69(lined) + 3.7396(scatter) = 70.429363%
// free spins 540, q = 0.074792, sq = 1/(1-q) = 1.080838
// free games frequency: 1/267.41
// RTP = 70.429(sym) + 0.074792*264.48(fg) = 90.210120%
var Reels902 = slot.Reels3x{
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 1, 1, 1, 3, 3, 3, 6, 6, 6, 7},
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 4, 1, 1, 1, 3, 3, 3, 6, 6, 6, 7},
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 1, 1, 1, 3, 3, 3, 6, 6, 6, 7},
}

// *bonus reels calculations*
// reels lengths [19, 20, 19], total reshuffles 7220
// symbols: 248.93(lined) + 3.7396(scatter) = 252.673130%
// free spins 540, q = 0.074792, sq = 1/(1-q) = 1.080838
// free games frequency: 1/267.41
// RTP = sq*rtp(sym) = 1.0808*252.67 = 273.098802%
// *regular reels calculations*
// reels lengths [19, 20, 19], total reshuffles 7220
// symbols: 68.684(lined) + 3.7396(scatter) = 72.423823%
// free spins 540, q = 0.074792, sq = 1/(1-q) = 1.080838
// free games frequency: 1/267.41
// RTP = 72.424(sym) + 0.074792*273.1(fg) = 92.849495%
var Reels928 = slot.Reels3x{
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 1, 1, 1, 3, 3, 3, 6, 6, 6, 7},
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 1, 1, 1, 3, 3, 3, 3, 6, 6, 6, 7},
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 1, 1, 1, 3, 3, 3, 6, 6, 6, 7},
}

// *bonus reels calculations*
// reels lengths [19, 19, 19], total reshuffles 6859
// symbols: 251.54(lined) + 3.9364(scatter) = 255.474559%
// free spins 540, q = 0.078729, sq = 1/(1-q) = 1.085457
// free games frequency: 1/254.04
// RTP = sq*rtp(sym) = 1.0855*255.47 = 277.306536%
// *regular reels calculations*
// reels lengths [19, 19, 19], total reshuffles 6859
// symbols: 69.675(lined) + 3.9364(scatter) = 73.611314%
// free spins 540, q = 0.078729, sq = 1/(1-q) = 1.085457
// free games frequency: 1/254.04
// RTP = 73.611(sym) + 0.078729*277.31(fg) = 95.443290%
var Reels954 = slot.Reels3x{
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 1, 1, 1, 3, 3, 3, 6, 6, 6, 7},
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 1, 1, 1, 3, 3, 3, 6, 6, 6, 7},
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 1, 1, 1, 3, 3, 3, 6, 6, 6, 7},
}

// *bonus reels calculations*
// reels lengths [19, 20, 19], total reshuffles 7220
// symbols: 263.89(lined) + 3.7396(scatter) = 267.631579%
// free spins 540, q = 0.074792, sq = 1/(1-q) = 1.080838
// free games frequency: 1/267.41
// RTP = sq*rtp(sym) = 1.0808*267.63 = 289.266467%
// *regular reels calculations*
// reels lengths [19, 20, 19], total reshuffles 7220
// symbols: 72.424(lined) + 3.7396(scatter) = 76.163435%
// free spins 540, q = 0.074792, sq = 1/(1-q) = 1.080838
// free games frequency: 1/267.41
// RTP = 76.163(sym) + 0.074792*289.27(fg) = 97.798323%
var Reels977 = slot.Reels3x{
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 1, 1, 1, 3, 3, 3, 6, 6, 6, 7},
	{5, 5, 5, 2, 2, 2, 2, 4, 4, 4, 1, 1, 1, 3, 3, 3, 6, 6, 6, 7},
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 1, 1, 1, 3, 3, 3, 6, 6, 6, 7},
}

// *bonus reels calculations*
// reels lengths [19, 20, 19], total reshuffles 7220
// symbols: 280.22(lined) + 3.7396(scatter) = 283.961219%
// free spins 540, q = 0.074792, sq = 1/(1-q) = 1.080838
// free games frequency: 1/267.41
// RTP = sq*rtp(sym) = 1.0808*283.96 = 306.916168%
// *regular reels calculations*
// reels lengths [19, 20, 19], total reshuffles 7220
// symbols: 78.657(lined) + 3.7396(scatter) = 82.396122%
// free spins 540, q = 0.074792, sq = 1/(1-q) = 1.080838
// free games frequency: 1/267.41
// RTP = 82.396(sym) + 0.074792*306.92(fg) = 105.351071%
var Reels105 = slot.Reels3x{
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 1, 1, 1, 3, 3, 3, 6, 6, 6, 7},
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 1, 1, 1, 1, 3, 3, 3, 6, 6, 6, 7},
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 1, 1, 1, 3, 3, 3, 6, 6, 6, 7},
}

// *bonus reels calculations*
// reels lengths [19, 21, 19], total reshuffles 7581
// symbols: 229.48(lined) + 7.1231(scatter) = 236.604670%
// free spins 1080, q = 0.14246, sq = 1/(1-q) = 1.166128
// free games frequency: 1/140.39
// RTP = sq*rtp(sym) = 1.1661*236.6 = 275.911398%
// *regular reels calculations*
// reels lengths [19, 21, 19], total reshuffles 7581
// symbols: 63.514(lined) + 7.1231(scatter) = 70.637119%
// free spins 1080, q = 0.14246, sq = 1/(1-q) = 1.166128
// free games frequency: 1/140.39
// RTP = 70.637(sym) + 0.14246*275.91(fg) = 109.943848%
var Reels109 = slot.Reels3x{
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 1, 1, 1, 3, 3, 3, 6, 6, 6, 7},
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 4, 1, 1, 1, 7, 3, 3, 3, 6, 6, 6, 7},
	{5, 5, 5, 2, 2, 2, 4, 4, 4, 1, 1, 1, 3, 3, 3, 6, 6, 6, 7},
}

// Map with available reels.
var ReelsMap = map[float64]*slot.Reels3x{
	84.424370:  &Reels844,
	88.051262:  &Reels880,
	90.210120:  &Reels902,
	92.849495:  &Reels928,
	95.443290:  &Reels954,
	97.798323:  &Reels977,
	105.351071: &Reels105,
	109.943848: &Reels109,
}
