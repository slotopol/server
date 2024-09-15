package spi

import (
	"encoding/xml"

	"github.com/gin-gonic/gin"
	keno "github.com/slotopol/server/game/keno"
)

// Returns bet value.
func SpiKenoBetGet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		Bet     float64  `json:"bet" yaml:"bet" xml:"bet"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_keno_betget_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_keno_betget_nogid, ErrNoGID)
		return
	}

	var scene *Scene
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, SEC_keno_betget_notopened, ErrNotOpened)
		return
	}
	var game keno.KenoGame
	if game, ok = scene.Game.(keno.KenoGame); !ok {
		Ret403(c, SEC_keno_betget_notslot, ErrNotKeno)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALgame == 0 {
		Ret403(c, SEC_keno_betget_noaccess, ErrNoAccess)
		return
	}

	ret.Bet = game.GetBet()

	RetOk(c, ret)
}

// Set bet value.
func SpiKenoBetSet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr"`
		Bet     float64  `json:"bet" yaml:"bet" xml:"bet" binding:"required"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_keno_betset_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_keno_betset_nogid, ErrNoGID)
		return
	}

	var scene *Scene
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, SEC_keno_betset_notopened, ErrNotOpened)
		return
	}
	var game keno.KenoGame
	if game, ok = scene.Game.(keno.KenoGame); !ok {
		Ret403(c, SEC_keno_betset_notslot, ErrNotKeno)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALgame == 0 {
		Ret403(c, SEC_keno_betset_noaccess, ErrNoAccess)
		return
	}

	if err = game.SetBet(arg.Bet); err != nil {
		Ret403(c, SEC_keno_betset_badbet, err)
		return
	}

	Ret204(c)
}

// Returns selected numvers bitset.
func SpiKenoSelGet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid"`
	}
	var ret struct {
		XMLName xml.Name    `json:"-" yaml:"-" xml:"ret"`
		Sel     keno.Bitset `json:"sel" yaml:"sel" xml:"sel"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_keno_selget_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_keno_selget_nogid, ErrNoGID)
		return
	}

	var scene *Scene
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, SEC_keno_selget_notopened, ErrNotOpened)
		return
	}
	var game keno.KenoGame
	if game, ok = scene.Game.(keno.KenoGame); !ok {
		Ret403(c, SEC_keno_selget_notslot, ErrNotSlot)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALgame == 0 {
		Ret403(c, SEC_keno_selget_noaccess, ErrNoAccess)
		return
	}

	ret.Sel = game.GetSel()

	RetOk(c, ret)
}

// Set selected numbers bitset.
func SpiKenoSelSet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name    `json:"-" yaml:"-" xml:"arg"`
		GID     uint64      `json:"gid" yaml:"gid" xml:"gid,attr"`
		Sel     keno.Bitset `json:"sel" yaml:"sel" xml:"sel" binding:"required"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_keno_selset_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_keno_selset_nogid, ErrNoGID)
		return
	}

	var scene *Scene
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, SEC_keno_selset_notopened, ErrNotOpened)
		return
	}
	var game keno.KenoGame
	if game, ok = scene.Game.(keno.KenoGame); !ok {
		Ret403(c, SEC_keno_selset_notslot, ErrNotSlot)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALgame == 0 {
		Ret403(c, SEC_keno_selset_noaccess, ErrNoAccess)
		return
	}

	if err = game.SetSel(arg.Sel); err != nil {
		Ret403(c, SEC_keno_selset_badlines, err)
		return
	}

	Ret204(c)
}
