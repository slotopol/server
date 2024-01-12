package spi

import (
	"encoding/xml"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/slotopol/server/game/champagne"
	"github.com/slotopol/server/game/dolphinspearl"
	"github.com/slotopol/server/game/jewels"
	"github.com/slotopol/server/game/sizzlinghot"
	"github.com/slotopol/server/game/slotopol"
	"github.com/slotopol/server/game/slotopoldeluxe"

	"github.com/gin-gonic/gin"
	"github.com/slotopol/server/game"
)

var GameFactory = map[string]func(string) game.SlotGame{
	"slotopol": func(name string) game.SlotGame {
		if reels, ok := slotopol.ReelsMap[name]; ok {
			return slotopol.NewGame(reels)
		}
		return nil
	},
	"slotopoldeluxe": func(name string) game.SlotGame {
		if reels, ok := slotopoldeluxe.ReelsMap[name]; ok {
			return slotopoldeluxe.NewGame(reels)
		}
		return nil
	},
	"champagne": func(name string) game.SlotGame {
		if reels, ok := champagne.ReelsMap[name]; ok {
			return champagne.NewGame(reels)
		}
		return nil
	},
	"jewels": func(name string) game.SlotGame {
		if reels, ok := jewels.ReelsMap[name]; ok {
			return jewels.NewGame(reels)
		}
		return nil
	},
	"sizzlinghot": func(name string) game.SlotGame {
		if reels, ok := sizzlinghot.ReelsMap[name]; ok {
			return sizzlinghot.NewGame(reels)
		}
		return nil
	},
	"dolphinspearl": func(name string) game.SlotGame {
		if reels, ok := dolphinspearl.ReelsMap[name]; ok {
			return dolphinspearl.NewGame(reels)
		}
		return nil
	},
}

var GameAliases = map[string]string{
	// Megajack games
	"slotopol":       "slotopol",
	"slotopoldeluxe": "slotopoldeluxe",
	"champagne":      "champagne",

	// Novomatic games
	"jewels":              "jewels",
	"sizzlinghot":         "sizzlinghot",
	"sizzlinghotdeluxe":   "sizzlinghot",
	"dolphinspearl":       "dolphinspearl",
	"dolphinspearldeluxe": "dolphinspearl",
	"attila":              "dolphinspearl",
	"bananasplash":        "dolphinspearl",
	"dynastyofming":       "dolphinspearl",
	"gryphonsgold":        "dolphinspearl",
	"jokerdolphin":        "dolphinspearl",
	"pharaonsgold2":       "dolphinspearl",
	"pharaonsgold3":       "dolphinspearl",
	"polarfox":            "dolphinspearl",
	"secretforest":        "dolphinspearl",
	"themoneygame":        "dolphinspearl",
	"unicornmagic":        "dolphinspearl",
}

var GIDcounter uint64
var OpenedGames = map[uint64]game.SlotGame{}

func SpiGameJoin(c *gin.Context) {
	var err error
	var arg struct {
		XMLName  xml.Name `json:"-" yaml:"-" xml:"arg"`
		GameName string   `json:"gamename" yaml:"gamename" xml:"gamename"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid"`
	}

	if err = c.Bind(&arg); err != nil {
		Ret400(c, SEC_game_join_nobind, err)
		return
	}

	var alias, is = GameAliases[strings.ToLower(arg.GameName)]
	if !is {
		Ret400(c, SEC_game_join_noalias, ErrNoAliase)
		return
	}

	var maker = GameFactory[alias]
	var game = maker("96")
	if game == nil {
		Ret400(c, SEC_game_join_noreels, ErrNoReels)
		return
	}

	var gid = atomic.AddUint64(&GIDcounter, 1)
	OpenedGames[gid] = game
	ret.GID = gid

	RetOk(c, ret)
}

func SpiGamePart(c *gin.Context) {
	var err error
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid"`
	}

	if err = c.Bind(&arg); err != nil {
		Ret400(c, SEC_game_part_nobind, err)
		return
	}

	var _, is = OpenedGames[arg.GID]
	if !is {
		Ret400(c, SEC_game_part_notopened, ErrNotOpened)
		return
	}

	delete(OpenedGames, arg.GID)

	c.Status(http.StatusOK)
}
