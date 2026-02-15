//go:build !prod || full || novomatic

package flamedancer

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed flamedancer_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Flame Dancer", LNum: 20, Date: game.Date(2012, 11, 15)}, // see: https://casino.ru/flame-dancer-novomatic/
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPscat |
			game.GPwild |
			game.GPrwild |
			game.GPwturn,
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
	game.DataRouter["novomatic/flamedancer/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
