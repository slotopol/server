//go:build !prod || full || novomatic

package powerstars

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed powerstars_reel.yaml
var reels []byte

//go:embed powerstars_chance.yaml
var chance []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Power Stars", Date: game.Year(2013)},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPrpay |
			game.GPlsel |
			game.GPfgno |
			game.GPwild,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ChanceMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["novomatic/powerstars/reel"] = &Reels
	game.DataRouter["novomatic/powerstars/chance"] = &ChanceMap
	game.LoadMap = append(game.LoadMap, reels, chance)
}
