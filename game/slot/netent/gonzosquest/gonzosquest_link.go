//go:build !prod || full || netent

package gonzosquest

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed gonzosquest_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "NetEnt", Name: "Gonzo's Quest", Date: game.Date(2011, 5, 15)}, // see: https://www.slotsmate.com/software/netent/gonzos-quest
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPcfall |
			game.GPcmult |
			game.GPretrig |
			game.GPfgmult |
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
	game.DataRouter["netent/gonzosquest/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
