//go:build !prod || full || playngo

package copsnrobbers

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed copsnrobbers_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Play'n GO", Name: "Cops'n'Robbers", LNum: 9, Date: game.Date(2014, 12, 1)}, // see: https://www.slotsmate.com/software/play-n-go/cops-n-robbers
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgonce |
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
	game.DataRouter["playngo/copsnrobbers/bon"] = &ReelsBon
	game.DataRouter["playngo/copsnrobbers/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
