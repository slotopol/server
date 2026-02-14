//go:build !prod || full || ct

package amazonsspear

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed amazonsspear_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Amazons Spear", LNum: 25, Date: game.Date(2019, 5, 1)}, // see: https://www.livebet2.com/casino/slots/ct-interactive/amazons-spear
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgseq |
			game.GPfgreel |
			game.GPscat |
			game.GPrwild,
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
	game.DataRouter["ctinteractive/amazonsspear/bon"] = &ReelsBon
	game.DataRouter["ctinteractive/amazonsspear/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
