package jewels4all

import (
	"math"

	slot "github.com/slotopol/server/game/slot"
)

// RTP(no eu) = 67.344781%
// RTP(eu at y=1,5) = 1706.345577%
// RTP(eu at y=2,3,4) = 7818.930041%
// euro avr: rtpeu = 5373.896256%
var Reels = slot.Reels5x{
	{1, 1, 1, 2, 2, 2, 6, 6, 3, 3, 3, 7, 7, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 6, 6, 3, 3, 3, 7, 7, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 6, 6, 3, 3, 3, 7, 7, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 6, 6, 3, 3, 3, 7, 7, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 6, 6, 3, 3, 3, 7, 7, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
}

// Map with wild chances.
var ChanceMap = map[float64]float64{
	// RTP = 67.345(sym) + wc*5373.9(eu) = 90.019449%
	90.019449: 1 / 237.,
	// RTP = 67.345(sym) + wc*5373.9(eu) = 91.995681%
	91.995681: 1 / 218.,
	// RTP = 67.345(sym) + wc*5373.9(eu) = 93.948228%
	93.948228: 1 / 202.,
	// RTP = 67.345(sym) + wc*5373.9(eu) = 96.082194%
	96.082194: 1 / 187.,
	// RTP = 67.345(sym) + wc*5373.9(eu) = 98.052760%
	98.052760: 1 / 175.,
	// RTP = 67.345(sym) + wc*5373.9(eu) = 99.913849%
	99.913849: 1 / 165.,
	// RTP = 67.345(sym) + wc*5373.9(eu) = 110.335951%
	110.335951: 1 / 125.,
}

func FindChance(mrtp float64) (rtp float64, chance float64) {
	for p, c := range ChanceMap {
		if math.Abs(mrtp-p) < math.Abs(mrtp-rtp) {
			rtp, chance = p, c
		}
	}
	return
}
