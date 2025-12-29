//go:build !prod || full || playtech

package greatblue

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed greatblue_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Playtech", Name: "Great Blue", LNum: 25, Date: game.Date(2013, 1, 1)}, // see: https://www.slotsmate.com/software/playtech/great-blue
		{Prov: "Playtech", Name: "Irish Luck", LNum: 30, Date: game.Date(2009, 1, 1)}, // see: https://www.slotsmate.com/software/playtech/irish-luck
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPpick |
			game.GPfgseq |
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
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["playtech/greatblue/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
