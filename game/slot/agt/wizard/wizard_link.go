//go:build !prod || full || agt

package wizard

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed wizard_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Wizard", LNum: 50},           // see: https://agtsoftware.com/games/agt/wizard
		{Prov: "AGT", Name: "Around The World", LNum: 40}, // see: https://agtsoftware.com/games/agt/aroundtheworld
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPstscat |
			game.GPwild,
		SX: 5,
		SY: 4,
		SN: sn,
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["agt/wizard/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
