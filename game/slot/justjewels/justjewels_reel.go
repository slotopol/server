package justjewels

import (
	"math"

	slot "github.com/slotopol/server/game/slot"
)

// reels lengths [39, 39, 39, 39, 39], total reshuffles 90224199
// RTP = 114.75(lined) + 8.0152(scatter) = 122.764204%
var Reels123 = slot.Reels5x{
	{1, 1, 1, 6, 6, 6, 2, 2, 2, 2, 5, 5, 5, 3, 3, 3, 3, 7, 7, 7, 8, 4, 4, 4, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 4, 4, 4},
	{1, 1, 1, 6, 6, 6, 2, 2, 2, 2, 5, 5, 5, 3, 3, 3, 3, 7, 7, 7, 8, 4, 4, 4, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 4, 4, 4},
	{1, 1, 1, 6, 6, 6, 2, 2, 2, 2, 5, 5, 5, 3, 3, 3, 3, 7, 7, 7, 8, 4, 4, 4, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 4, 4, 4},
	{1, 1, 1, 6, 6, 6, 2, 2, 2, 2, 5, 5, 5, 3, 3, 3, 3, 7, 7, 7, 8, 4, 4, 4, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 4, 4, 4},
	{1, 1, 1, 6, 6, 6, 2, 2, 2, 2, 5, 5, 5, 3, 3, 3, 3, 7, 7, 7, 8, 4, 4, 4, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 4, 4, 4},
}

// Map with available reels.
var ReelsMap = map[float64]*slot.Reels5x{
	122.764204: &Reels123, // minimum possible percentage
}

func FindReels(mrtp float64) (rtp float64, reels slot.Reels) {
	for p, r := range ReelsMap {
		if math.Abs(mrtp-p) < math.Abs(mrtp-rtp) {
			rtp, reels = p, r
		}
	}
	return
}
