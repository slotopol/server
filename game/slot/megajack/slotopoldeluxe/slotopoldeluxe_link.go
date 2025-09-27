//go:build !prod || full || megajack

package slotopoldeluxe

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed slotopoldeluxe_reel.yaml
var reels []byte

//go:embed slotopoldeluxe_jack.yaml
var jack []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Megajack", Name: "Slotopol Deluxe", Date: game.Year(1999)},
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
		BN: 4,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["megajack/slotopoldeluxe/reel"] = &ReelsMap
	game.DataRouter["megajack/slotopoldeluxe/jack"] = &JackMap
	game.LoadMap = append(game.LoadMap, reels, jack)
}
