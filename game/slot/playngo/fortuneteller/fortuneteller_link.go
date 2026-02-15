//go:build !prod || full || playngo

package fortuneteller

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed fortuneteller_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Play'n GO", Name: "Fortune Teller", LNum: 20, Date: game.Date(2012, 5, 30)}, // see: https://www.slotsmate.com/software/play-n-go/play-n-go-fortune-teller
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgonce |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 3,
		SN: sn,
		LN: len(BetLines),
		BN: 1,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStatReg)
	game.DataRouter["playngo/fortuneteller/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
