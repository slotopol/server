//go:build !prod || full || agt

package valkyrie

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed valkyrie_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Valkyrie", LNum: 30, Date: game.Year(2024)}, // see: https://agtsoftware.com/games/agt/valkyrie
		{Prov: "AGT", Name: "Aquaman", LNum: 10, Date: game.Year(2025)},  // see: https://agtsoftware.com/games/agt/aquaman
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgonce |
			game.GPscat |
			game.GPwild |
			game.GPbsym,
		SX: 5,
		SY: 3,
		SN: sn,
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStatReg)
	game.DataRouter["agt/valkyrie/big"] = &ReelBig
	game.DataRouter["agt/valkyrie/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
