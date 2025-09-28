//go:build !prod || full || playtech

package goldentour

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed goldentour_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Playtech", Name: "Golden Tour", Date: game.Date(2009, 1, 25)}, // see: https://www.slotsmate.com/software/playtech/golden-tour
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPrpay |
			game.GPlsel |
			game.GPfgno |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 1,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["playtech/goldentour/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
