//go:build !prod || full || ct

package purplehot2

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed purplehot2_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Purple Hot 2", Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/purple-hot-2
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPfgno |
			game.GPscat,
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
	game.DataRouter["ctinteractive/purplehot2/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
