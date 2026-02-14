//go:build !prod || full || ct

package americangigolo

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed americangigolo_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "American Gigolo", LNum: 30, Date: game.Date(2018, 12, 31)}, // see: https://www.slotsmate.com/software/ct-interactive/american-gigolo
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgseq |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 3,
		SN: sn,
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["ctinteractive/americangigolo/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
