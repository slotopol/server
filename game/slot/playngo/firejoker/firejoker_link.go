//go:build !prod || full || playngo

package firejoker

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed firejoker_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Play'n GO", Name: "Fire Joker", LNum: 5},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgonce |
			game.GPscat |
			game.GPbsym,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStatReg)
	game.DataRouter["playngo/firejoker/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
