//go:build !prod || full || megajack

package slotopol

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed slotopol_reel.yaml
var reels []byte

//go:embed slotopol_jack.yaml
var jack []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Megajack", Name: "Slotopol", Date: game.Year(1999)},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPjack |
			game.GPfgno |
			game.GPscat |
			game.GPwild |
			game.GPwmult,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 2,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["megajack/slotopol/reel"] = &ReelsMap
	game.DataRouter["megajack/slotopol/jack"] = &JackMap
	game.LoadMap = append(game.LoadMap, reels, jack)
}
