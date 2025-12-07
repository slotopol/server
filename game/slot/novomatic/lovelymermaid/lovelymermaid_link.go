//go:build !prod || full || novomatic

package lovelymermaid

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed lovelymermaid_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Lovely Mermaid", LNum: 40, Date: game.Date(2015, 12, 1)}, // see: https://www.slotsmate.com/software/novomatic/lovely-mermaid
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPjack |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 4,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["novomatic/lovelymermaid/reel"] = &ReelsMap
	game.DataRouter["novomatic/lovelymermaid/jack"] = &JackMap
	game.LoadMap = append(game.LoadMap, data)
}
