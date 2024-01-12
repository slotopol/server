package spi

import (
	"encoding/xml"
	"net/http"
	"strings"
	"sync/atomic"

	cfg "github.com/slotopol/server/config"

	"github.com/gin-gonic/gin"
	"github.com/slotopol/server/game"
)

var GIDcounter uint64
var OpenedGames = map[uint64]game.SlotGame{}

func SpiGameJoin(c *gin.Context) {
	var err error
	var arg struct {
		XMLName  xml.Name `json:"-" yaml:"-" xml:"arg"`
		GameName string   `json:"gamename" yaml:"gamename" xml:"gamename" form:"gamename"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid"`
	}

	if err = c.Bind(&arg); err != nil {
		Ret400(c, SEC_game_join_nobind, err)
		return
	}
	if arg.GameName == "" {
		Ret400(c, SEC_game_join_nodata, ErrNoData)
		return
	}

	var alias, is = cfg.GameAliases[strings.ToLower(arg.GameName)]
	if !is {
		Ret400(c, SEC_game_join_noalias, ErrNoAliase)
		return
	}

	var maker = cfg.GameFactory[alias]
	var slotgame = maker("96")
	if slotgame == nil {
		Ret400(c, SEC_game_join_noreels, ErrNoReels)
		return
	}

	var gid = atomic.AddUint64(&GIDcounter, 1)
	OpenedGames[gid] = slotgame.(game.SlotGame)
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
	if arg.GID == 0 {
		Ret400(c, SEC_game_part_nodata, ErrNoData)
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
