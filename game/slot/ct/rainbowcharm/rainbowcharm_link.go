//go:build !prod || full || ct

package rainbowcharm

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed rainbowcharm_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Rainbow Charm", Date: game.Date(2024, 11, 30)},   // see: https://www.slotsmate.com/software/ct-interactive/rainbow-charm
		{Prov: "CT Interactive", Name: "The Magic Goblet", Date: game.Date(2024, 7, 14)}, // see: https://www.slotsmate.com/software/ct-interactive/the-magic-goblet
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPcpay |
			game.GPrgmult |
			game.GPscat,
		SX: 5,
		SY: 3,
		SN: len(SymPay),
		LN: 0,
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["ctinteractive/rainbowcharm/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
