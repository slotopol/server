//go:build !prod || full || playtech

package panthermoon

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed panthermoon_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Playtech", Name: "Panther Moon", LNum: 15, Date: game.Date(2009, 2, 28)}, // see: https://www.slotsmate.com/software/playtech/panther-moon
		{Prov: "Playtech", Name: "Safari Heat", LNum: 15, Date: game.Date(2009, 3, 1)},   // see: https://www.slotsmate.com/software/playtech/safari-heat
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
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
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStatReg)
	game.DataRouter["playtech/panthermoon/bon"] = &ReelsBon
	game.DataRouter["playtech/panthermoon/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
