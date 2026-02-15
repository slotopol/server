//go:build !prod || full || novomatic

package bananasgobahamas

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed bananasgobahamas_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Bananas Go Bahamas", LNum: 9, Date: game.Year(2004)}, // see: https://www.slotsmate.com/software/novomatic/bananas-go-bahamas
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPfgreel |
			game.GPscat |
			game.GPwild |
			game.GPwmult,
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
	game.DataRouter["novomatic/bananasgobahamas/bon"] = &ReelsBon
	game.DataRouter["novomatic/bananasgobahamas/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
