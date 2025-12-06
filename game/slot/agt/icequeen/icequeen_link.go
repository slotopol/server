//go:build !prod || full || agt

package icequeen

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed icequeen_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Ice Queen", LNum: 10},       // see: https://agtsoftware.com/games/agt/iceqween
		{Prov: "AGT", Name: "STALKER", LNum: 10},         // see: https://agtsoftware.com/games/agt/stalker
		{Prov: "AGT", Name: "Big Five", LNum: 15},        // see: https://agtsoftware.com/games/agt/bigfive
		{Prov: "AGT", Name: "Arabian Nights", LNum: 10},  // see: https://agtsoftware.com/games/agt/arabiannights
		{Prov: "AGT", Name: "Anonymous", LNum: 20},       // see: https://agtsoftware.com/games/agt/anonymous
		{Prov: "AGT", Name: "Grand Theft", LNum: 20},     // see: https://agtsoftware.com/games/agt/bankofny
		{Prov: "AGT", Name: "Firefighters", LNum: 15},    // see: https://agtsoftware.com/games/agt/firefighters
		{Prov: "AGT", Name: "Time Machine II", LNum: 15}, // see: https://agtsoftware.com/games/agt/timemachine2
		{Prov: "AGT", Name: "Bitcoin", LNum: 15},         // see: https://agtsoftware.com/games/agt/bitcoin
		{Prov: "AGT", Name: "Pirates Gold", LNum: 20},    // see: https://agtsoftware.com/games/agt/piratesgold
		{Prov: "AGT", Name: "The Leprechaun", LNum: 10},  // see: https://agtsoftware.com/games/agt/leprechaun
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPfgreel |
			game.GPfgmult |
			game.GPscat |
			game.GPwild |
			game.GPwmult,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["agt/icequeen/bon"] = &ReelsBon
	game.DataRouter["agt/icequeen/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
