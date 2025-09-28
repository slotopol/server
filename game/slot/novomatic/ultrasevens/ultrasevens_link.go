//go:build !prod || full || novomatic

package ultrasevens

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed ultrasevens_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Ultra Sevens", Date: game.Date(2015, 11, 1)}, // see: https://www.slotsmate.com/software/novomatic/ultra-sevens
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
			game.GPjack |
			game.GPscat,
		SX: 5,
		SY: 4,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["novomatic/ultrasevens/reel"] = &ReelsMap
	game.DataRouter["novomatic/ultrasevens/jack"] = &JackMap
	game.LoadMap = append(game.LoadMap, data)
}
