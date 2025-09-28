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
		{Prov: "Novomatic", Name: "Katana", Date: game.Year(2012)},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPretrig |
			game.GPfgreel |
			game.GPscat |
			game.GPwild |
			game.GPrwild,
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
	game.DataRouter["novomatic/katana/bon"] = &ReelsBon
	game.DataRouter["novomatic/katana/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
