//go:build !prod || full || playtech

package captainstreasure

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed captainstreasure_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Playtech", Name: "Captain's Treasure", Date: game.Date(2009, 1, 1)}, // see: https://www.slotsmate.com/software/playtech/captain-treasure
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPrpay |
			game.GPlsel |
			game.GPfgno |
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
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["playtech/captainstreasure/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
