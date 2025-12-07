//go:build !prod || full || ct

package moneypipe

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed moneypipe_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Money Pipe", LNum: 40, Date: game.Date(2020, 10, 1)},    // see: https://www.slotsmate.com/software/ct-interactive/money-pipe
		{Prov: "CT Interactive", Name: "More Dragons", LNum: 40, Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/more-dragons
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgseq |
			game.GPscat |
			game.GPewild,
		SX: 5,
		SY: 4,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["ctinteractive/moneypipe/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
