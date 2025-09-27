//go:build !prod || full || agt

package aislot

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed aislot_reel.yaml
var reels []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "AI"},          // see: https://agtsoftware.com/games/agt/aislot
		{Prov: "AGT", Name: "Tesla"},       // see: https://agtsoftware.com/games/agt/tesla
		{Prov: "AGT", Name: "Book of Set"}, // see: https://agtsoftware.com/games/agt/bookofset
		{Prov: "AGT", Name: "Pharaoh II"},  // see: https://agtsoftware.com/games/agt/pharaoh2
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPretrig |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["agt/aislot/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, reels)
}
