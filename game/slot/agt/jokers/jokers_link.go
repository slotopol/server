//go:build !prod || full || agt

package jokers

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed jokers_reel.yaml
var reels []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Jokers"},
		{Prov: "AGT", Name: "Happy Santa"}, // see: https://agtsoftware.com/games/agt/happysanta
		{Prov: "AGT", Name: "Bigfoot"},     // see: https://agtsoftware.com/games/agt/bigfoot
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgno |
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
	game.DataRouter["agt/jokers/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, reels)
}
