//go:build !prod || full || betsoft

package twomillionbc

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed 2millionbc_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "BetSoft", Name: "2 Million B.C.", Date: game.Year(2011)}, // see: https://www.slotsmate.com/software/betsoft/2-million-bc
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPfgreel |
			game.GPscat,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 2,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["betsoft/2millionbc/bon"] = &ReelsBon
	game.DataRouter["betsoft/2millionbc/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
