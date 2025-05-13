package api

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"sync/atomic"

	"github.com/gin-gonic/gin"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/game/keno"
	"github.com/slotopol/server/util"
)

// Returns bet value.
func ApiKenoBetGet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid" binding:"required"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		Bet     float64  `json:"bet" yaml:"bet" xml:"bet"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_keno_betget_nobind, err)
		return
	}

	var scene *Scene
	if scene, err = GetScene(arg.GID); err != nil {
		Ret404(c, AEC_keno_betget_noscene, err)
		return
	}
	var game keno.KenoGame
	if game, ok = scene.Game.(keno.KenoGame); !ok {
		Ret403(c, AEC_keno_betget_notslot, ErrNotKeno)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALdealer == 0 {
		Ret403(c, AEC_keno_betget_noaccess, ErrNoAccess)
		return
	}

	ret.Bet = game.GetBet()

	RetOk(c, ret)
}

// Set bet value.
func ApiKenoBetSet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" binding:"required"`
		Bet     float64  `json:"bet" yaml:"bet" xml:"bet" binding:"required"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_keno_betset_nobind, err)
		return
	}

	var scene *Scene
	if scene, err = GetScene(arg.GID); err != nil {
		Ret404(c, AEC_keno_betset_noscene, err)
		return
	}
	var game keno.KenoGame
	if game, ok = scene.Game.(keno.KenoGame); !ok {
		Ret403(c, AEC_keno_betset_notslot, ErrNotKeno)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALdealer == 0 {
		Ret403(c, AEC_keno_betset_noaccess, ErrNoAccess)
		return
	}

	if err = game.SetBet(arg.Bet); err != nil {
		Ret403(c, AEC_keno_betset_badbet, err)
		return
	}

	Ret204(c)
}

// Returns selected numbers bitset.
func ApiKenoSelGet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid" binding:"required"`
	}
	var ret struct {
		XMLName xml.Name    `json:"-" yaml:"-" xml:"ret"`
		Sel     keno.Bitset `json:"sel" yaml:"sel" xml:"sel"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_keno_selget_nobind, err)
		return
	}

	var scene *Scene
	if scene, err = GetScene(arg.GID); err != nil {
		Ret404(c, AEC_keno_selget_noscene, err)
		return
	}
	var game keno.KenoGame
	if game, ok = scene.Game.(keno.KenoGame); !ok {
		Ret403(c, AEC_keno_selget_notslot, ErrNotSlot)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALdealer == 0 {
		Ret403(c, AEC_keno_selget_noaccess, ErrNoAccess)
		return
	}

	ret.Sel = game.GetSel()

	RetOk(c, ret)
}

// Set selected numbers bitset.
func ApiKenoSelSet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name    `json:"-" yaml:"-" xml:"arg"`
		GID     uint64      `json:"gid" yaml:"gid" xml:"gid,attr" binding:"required"`
		Sel     keno.Bitset `json:"sel" yaml:"sel" xml:"sel" binding:"required"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_keno_selset_nobind, err)
		return
	}

	var scene *Scene
	if scene, err = GetScene(arg.GID); err != nil {
		Ret404(c, AEC_keno_selset_noscene, err)
		return
	}
	var game keno.KenoGame
	if game, ok = scene.Game.(keno.KenoGame); !ok {
		Ret403(c, AEC_keno_selset_notslot, ErrNotSlot)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALdealer == 0 {
		Ret403(c, AEC_keno_selset_noaccess, ErrNoAccess)
		return
	}

	if err = game.SetSel(arg.Sel); err != nil {
		Ret403(c, AEC_keno_selset_badsel, err)
		return
	}

	Ret204(c)
}

// Returns selected numbers slice.
func ApiKenoSelGetSlice(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid" binding:"required"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		Sel     []int    `json:"sel" yaml:"sel" xml:"sel"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_keno_selgetslice_nobind, err)
		return
	}

	var scene *Scene
	if scene, err = GetScene(arg.GID); err != nil {
		Ret404(c, AEC_keno_selgetslice_noscene, err)
		return
	}
	var game keno.KenoGame
	if game, ok = scene.Game.(keno.KenoGame); !ok {
		Ret403(c, AEC_keno_selgetslice_notslot, ErrNotSlot)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALdealer == 0 {
		Ret403(c, AEC_keno_selgetslice_noaccess, ErrNoAccess)
		return
	}

	var bs = game.GetSel()
	ret.Sel = bs.Expand()

	RetOk(c, ret)
}

// Set selected numbers slice.
func ApiKenoSelSetSlice(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" binding:"required"`
		Sel     []int    `json:"sel" yaml:"sel" xml:"sel" binding:"required"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_keno_selsetslice_nobind, err)
		return
	}

	var scene *Scene
	if scene, err = GetScene(arg.GID); err != nil {
		Ret404(c, AEC_keno_selsetslice_noscene, err)
		return
	}
	var game keno.KenoGame
	if game, ok = scene.Game.(keno.KenoGame); !ok {
		Ret403(c, AEC_keno_selsetslice_notslot, ErrNotSlot)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALdealer == 0 {
		Ret403(c, AEC_keno_selsetslice_noaccess, ErrNoAccess)
		return
	}

	var bs keno.Bitset
	bs.Pack(arg.Sel)
	if err = game.SetSel(bs); err != nil {
		Ret403(c, AEC_keno_selsetslice_badsel, err)
		return
	}

	Ret204(c)
}

func ApiKenoSpin(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid" binding:"required"`
		Bet     float64  `json:"bet,omitempty" yaml:"bet,omitempty" xml:"bet,omitempty"`
		Sel     []int    `json:"sel,omitempty" yaml:"sel,omitempty" xml:"sel,omitempty"`
	}
	var ret struct {
		XMLName xml.Name      `json:"-" yaml:"-" xml:"ret"`
		SID     uint64        `json:"sid" yaml:"sid" xml:"sid,attr"`
		Game    keno.KenoGame `json:"game" yaml:"game" xml:"game"`
		Wins    keno.Wins     `json:"wins" yaml:"wins" xml:"wins"`
		Wallet  float64       `json:"wallet" yaml:"wallet" xml:"wallet"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_keno_spin_nobind, err)
		return
	}

	var scene *Scene
	if scene, err = GetScene(arg.GID); err != nil {
		Ret404(c, AEC_keno_spin_noscene, err)
		return
	}
	var game keno.KenoGame
	if game, ok = scene.Game.(keno.KenoGame); !ok {
		Ret403(c, AEC_keno_spin_notslot, ErrNotKeno)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(scene.CID); !ok {
		Ret500(c, AEC_keno_spin_noclub, ErrNoClub)
		return
	}

	var user *User
	if user, ok = Users.Get(scene.UID); !ok {
		Ret500(c, AEC_keno_spin_nouser, ErrNoUser)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALdealer == 0 {
		Ret403(c, AEC_keno_spin_noaccess, ErrNoAccess)
		return
	}

	if arg.Bet != 0 {
		if err = game.SetBet(arg.Bet); err != nil {
			Ret403(c, AEC_keno_spin_badbet, err)
			return
		}
	}
	if len(arg.Sel) > 0 {
		var bs keno.Bitset
		bs.Pack(arg.Sel)
		if err = game.SetSel(bs); err != nil {
			Ret403(c, AEC_keno_spin_badsel, err)
			return
		}
	}

	var cost = game.GetBet()

	var props *Props
	if props, ok = user.props.Get(scene.CID); !ok {
		Ret500(c, AEC_keno_spin_noprops, ErrNoProps)
		return
	}
	if props.Wallet < cost {
		Ret403(c, AEC_keno_spin_nomoney, ErrNoMoney)
		return
	}

	var bank = club.Bank()
	var mrtp = GetRTP(user, club)

	// spin until gain less than bank value
	var wins keno.Wins
	var debit float64
	var n = 0
	for { // repeat until spin will fit into bank
		for { // repeat until get valid screen
			game.Spin(mrtp)
			if game.Scanner(&wins) == nil {
				break
			}
			n++
			if n > cfg.Cfg.MaxSpinAttempts {
				Ret500(c, AEC_keno_spin_badbank, ErrBadBank)
				return
			}
		}
		debit = cost - wins.Pay
		if bank+debit >= 0 || (bank < 0 && debit > 0) {
			break
		}
	}

	// write gain and total bet as transaction
	if Cfg.ClubUpdateBuffer > 1 {
		go BankBat[scene.CID].Put(cfg.XormStorage, scene.UID, debit)
	} else if err = BankBat[scene.CID].Put(cfg.XormStorage, scene.UID, debit); err != nil {
		Ret500(c, AEC_keno_spin_sqlbank, err)
		return
	}

	// make changes to memory data
	club.AddBank(debit)
	props.Wallet -= debit

	// write spin result to log and get spin ID
	var sid = atomic.AddUint64(&SpinCounter, 1)
	scene.SID = sid
	var rec = Spinlog{
		SID:    sid,
		GID:    arg.GID,
		MRTP:   mrtp,
		Gain:   wins.Pay,
		Wallet: props.Wallet,
	}
	var b []byte

	if b, err = json.Marshal(scene.Game); err != nil {
		return
	}
	rec.Game = util.B2S(b)

	if b, err = json.Marshal(wins); err != nil {
		return
	}
	rec.Wins = util.B2S(b)

	if Cfg.UseSpinLog {
		go func() {
			if err = SpinBuf.Put(cfg.XormSpinlog, rec); err != nil {
				log.Printf("can not write to spin log: %s", err.Error())
			}
		}()
	}

	// prepare result
	ret.SID = sid
	ret.Game = game
	ret.Wins = wins
	ret.Wallet = props.Wallet

	RetOk(c, ret)
}
