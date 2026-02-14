//go:build !prod || full || agt

package zeus

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed zeus_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Zeus", LNum: 10, Date: game.Year(2025)}, // see: https://agtsoftware.org/games/agt/zeus
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPscat |
			game.GPwild,
		SX: 4,
		SY: 4,
		SN: sn,
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["agt/zeus/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
