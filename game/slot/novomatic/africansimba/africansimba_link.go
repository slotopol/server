//go:build !prod || full || novomatic

package africansimba

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed africansimba_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "African Simba", Date: game.Date(2012, 10, 16)}, // see: https://www.slotsmate.com/software/novomatic/african-simba
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgseq |
			game.GPfgmult |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		WN: 243,
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["novomatic/africansimba/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
