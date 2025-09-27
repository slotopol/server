//go:build !prod || full || novomatic

package secretelixir

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed secretelixir_bon.yaml
var rbon []byte

//go:embed secretelixir_reel.yaml
var reels []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Secret Elixir", Date: game.Date(2010, 1, 15)}, // see: https://www.slotsmate.com/software/novomatic/secret-elixir
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPretrig |
			game.GPfgreel |
			game.GPfgmult |
			game.GPrmult |
			game.GPscat |
			game.GPwild,
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
	game.DataRouter["novomatic/secretelixir/bon"] = &ReelsBon
	game.DataRouter["novomatic/secretelixir/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, rbon, reels)
}
