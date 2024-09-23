package indiandreaming

import (
	"math"

	slot "github.com/slotopol/server/game/slot"
)

var Reels92 = slot.Reels5x{
	{},
}

// Map with available reels.
var ReelsMap = map[float64]*slot.Reels5x{
	92: &Reels92,
}

func FindReels(mrtp float64) (rtp float64, reels slot.Reels) {
	for p, r := range ReelsMap {
		if math.Abs(mrtp-p) < math.Abs(mrtp-rtp) {
			rtp, reels = p, r
		}
	}
	return
}
