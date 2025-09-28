//go:build !prod || full || netent

package arabiannights

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed arabiannights_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "NetEnt", Name: "Arabian Nights", Date: game.Date(2005, 5, 15)}, // see: https://www.slotsmate.com/software/netent/arabian-nights
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPretrig |
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
	game.DataRouter["netent/arabiannights/bon"] = &ReelsBon
	game.DataRouter["netent/arabiannights/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
