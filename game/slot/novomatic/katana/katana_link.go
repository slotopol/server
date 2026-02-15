//go:build !prod || full || novomatic

package katana

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed katana_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Katana", LNum: 20, Date: game.Year(2012)},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPfgreel |
			game.GPscat |
			game.GPwild |
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
	game.DataRouter["novomatic/katana/bon"] = &ReelsBon
	game.DataRouter["novomatic/katana/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
