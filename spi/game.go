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

func SpiGameJoin(c *gin.Context) {
	var err error
	var arg struct {
		XMLName  xml.Name `json:"-" yaml:"-" xml:"arg"`
		UID      uint64   `json:"uid" yaml:"uid" xml:"uid" form:"uid"`
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
	if arg.UID == 0 || arg.GameName == "" {
		Ret400(c, SEC_game_join_nodata, ErrNoData)
		return
	}

	var user, has = Users[arg.UID]
	if !has {
		Ret400(c, SEC_game_join_nouser, ErrNoUser)
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
	var og = OpenGame{
		GID:   gid,
		UID:   arg.UID,
		Alias: alias,
		game:  slotgame.(game.SlotGame),
	}
	OpenGames[gid] = og
	user.games[gid] = og
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

	var og, is = OpenGames[arg.GID]
	if !is {
		Ret400(c, SEC_game_part_notopened, ErrNotOpened)
		return
	}

	var user, has = Users[og.UID]
	if !has {
		Ret500(c, SEC_game_part_nouser, ErrNoUser)
		return
	}

	delete(OpenGames, arg.GID)
	delete(user.games, arg.GID)

	c.Status(http.StatusOK)
}

func SpiGameBet(c *gin.Context) {
	var err error
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid"`
		Bet     int      `json:"bet" yaml:"bet" xml:"bet"`
	}

	if err = c.Bind(&arg); err != nil {
		Ret400(c, SEC_game_bet_nobind, err)
		return
	}
	if arg.GID == 0 || arg.Bet == 0 {
		Ret400(c, SEC_game_bet_nodata, ErrNoData)
		return
	}

	var og, is = OpenGames[arg.GID]
	if !is {
		Ret400(c, SEC_game_bet_notopened, ErrNotOpened)
		return
	}

	if err = og.game.SetBet(arg.Bet); err != nil {
		Ret400(c, SEC_game_bet_badbet, err)
		return
	}

	c.Status(http.StatusOK)
}

func SpiGameLine(c *gin.Context) {
	var err error
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid"`
		Lines   []int    `json:"lines" yaml:"lines" xml:"lines"`
	}

	if err = c.Bind(&arg); err != nil {
		Ret400(c, SEC_game_line_nobind, err)
		return
	}
	if arg.GID == 0 || len(arg.Lines) == 0 {
		Ret400(c, SEC_game_line_nodata, ErrNoData)
		return
	}

	var og, is = OpenGames[arg.GID]
	if !is {
		Ret400(c, SEC_game_line_notopened, ErrNotOpened)
		return
	}

	if err = og.game.SetLines(arg.Lines); err != nil {
		Ret400(c, SEC_game_line_badlines, err)
		return
	}

	c.Status(http.StatusOK)
}
