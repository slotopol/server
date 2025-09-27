//go:build !prod || full || netent

package diamonddogs

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed diamonddogs_bon.yaml
var rbon []byte

//go:embed diamonddogs_reel.yaml
var reels []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "NetEnt", Name: "Diamond Dogs", Date: game.Year(2013)},
		{Prov: "NetEnt", Name: "Voodoo Vibes", Date: game.Year(2009)},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPretrig |
			game.GPfgreel |
			game.GPfgmult |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 1,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["netent/diamonddogs/bon"] = &ReelsBon
	game.DataRouter["netent/diamonddogs/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, rbon, reels)
}
