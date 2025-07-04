package api

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"math/rand/v2"
	"sync/atomic"

	"github.com/gin-gonic/gin"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/util"
)

// Returns bet value.
func ApiSlotBetGet(c *gin.Context) {
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
		Ret400(c, AEC_slot_betget_nobind, err)
		return
	}

	var scene *Scene
	if scene, err = GetScene(arg.GID); err != nil {
		Ret404(c, AEC_slot_betget_noscene, err)
		return
	}
	var game slot.SlotGame
	if game, ok = scene.Game.(slot.SlotGame); !ok {
		Ret403(c, AEC_slot_betget_notslot, ErrNotSlot)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALdealer == 0 {
		Ret403(c, AEC_slot_betget_noaccess, ErrNoAccess)
		return
	}

	ret.Bet = game.GetBet()

	RetOk(c, ret)
}

// Set bet value.
func ApiSlotBetSet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" binding:"required"`
		Bet     float64  `json:"bet" yaml:"bet" xml:"bet" binding:"required"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_slot_betset_nobind, err)
		return
	}

	var scene *Scene
	if scene, err = GetScene(arg.GID); err != nil {
		Ret404(c, AEC_slot_betset_noscene, err)
		return
	}
	var game slot.SlotGame
	if game, ok = scene.Game.(slot.SlotGame); !ok {
		Ret403(c, AEC_slot_betset_notslot, ErrNotSlot)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALdealer == 0 {
		Ret403(c, AEC_slot_betset_noaccess, ErrNoAccess)
		return
	}

	if err = game.SetBet(arg.Bet); err != nil {
		Ret403(c, AEC_slot_betset_badbet, err)
		return
	}

	Ret204(c)
}

// Returns selected bet lines bitset.
func ApiSlotSelGet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid" binding:"required"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		Sel     int      `json:"sel" yaml:"sel" xml:"sel"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_slot_selget_nobind, err)
		return
	}

	var scene *Scene
	if scene, err = GetScene(arg.GID); err != nil {
		Ret404(c, AEC_slot_selget_noscene, err)
		return
	}
	var game slot.SlotGame
	if game, ok = scene.Game.(slot.SlotGame); !ok {
		Ret403(c, AEC_slot_selget_notslot, ErrNotSlot)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALdealer == 0 {
		Ret403(c, AEC_slot_selget_noaccess, ErrNoAccess)
		return
	}

	ret.Sel = game.GetSel()

	RetOk(c, ret)
}

// Set selected bet lines bitset.
func ApiSlotSelSet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" binding:"required"`
		Sel     int      `json:"sel" yaml:"sel" xml:"sel" binding:"required"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_slot_selset_nobind, err)
		return
	}

	var scene *Scene
	if scene, err = GetScene(arg.GID); err != nil {
		Ret404(c, AEC_slot_selset_noscene, err)
		return
	}
	var game slot.SlotGame
	if game, ok = scene.Game.(slot.SlotGame); !ok {
		Ret403(c, AEC_slot_selset_notslot, ErrNotSlot)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALdealer == 0 {
		Ret403(c, AEC_slot_selset_noaccess, ErrNoAccess)
		return
	}

	if err = game.SetSel(arg.Sel); err != nil {
		Ret403(c, AEC_slot_selset_badsel, err)
		return
	}

	Ret204(c)
}

// Change game mode depending on the user's choice.
func ApiSlotModeSet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid" binding:"required"`
		N       int      `json:"n" yaml:"n" xml:"n,attr" form:"n"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_slot_modeset_nobind, err)
		return
	}

	var scene *Scene
	if scene, err = GetScene(arg.GID); err != nil {
		Ret404(c, AEC_slot_modeset_noscene, err)
		return
	}
	var game slot.SlotGame
	if game, ok = scene.Game.(slot.SlotGame); !ok {
		Ret403(c, AEC_slot_modeset_notslot, ErrNotSlot)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALdealer == 0 {
		Ret403(c, AEC_slot_modeset_noaccess, ErrNoAccess)
		return
	}

	if err = game.SetMode(arg.N); err != nil {
		Ret403(c, AEC_slot_modeset_badmode, err)
		return
	}

	Ret204(c)
}

// Make a spin.
func ApiSlotSpin(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid" binding:"required"`
		Bet     float64  `json:"bet,omitempty" yaml:"bet,omitempty" xml:"bet,omitempty"`
		Sel     int      `json:"sel,omitempty" yaml:"sel,omitempty" xml:"sel,omitempty"`
	}
	var ret struct {
		XMLName xml.Name      `json:"-" yaml:"-" xml:"ret"`
		SID     uint64        `json:"sid" yaml:"sid" xml:"sid,attr"`
		Game    slot.SlotGame `json:"game" yaml:"game" xml:"game"`
		Wins    slot.Wins     `json:"wins,omitempty" yaml:"wins,omitempty" xml:"wins,omitempty"`
		Wallet  float64       `json:"wallet" yaml:"wallet" xml:"wallet"`
		JpFund  float64       `json:"jpfund,omitempty" yaml:"jpfund,omitempty" xml:"jpfund,omitempty"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_slot_spin_nobind, err)
		return
	}

	var scene *Scene
	if scene, err = GetScene(arg.GID); err != nil {
		Ret404(c, AEC_slot_spin_noscene, err)
		return
	}
	var game slot.SlotGame
	if game, ok = scene.Game.(slot.SlotGame); !ok {
		Ret403(c, AEC_slot_spin_notslot, ErrNotSlot)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(scene.CID); !ok {
		Ret500(c, AEC_slot_spin_noclub, ErrNoClub)
		return
	}

	var user *User
	if user, ok = Users.Get(scene.UID); !ok {
		Ret500(c, AEC_slot_spin_nouser, ErrNoUser)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALdealer == 0 {
		Ret403(c, AEC_slot_spin_noaccess, ErrNoAccess)
		return
	}

	var props *Props
	if props, ok = user.props.Get(scene.CID); !ok {
		Ret500(c, AEC_slot_spin_noprops, ErrNoProps)
		return
	}

	if arg.Bet != 0 {
		if err = game.SetBet(arg.Bet); err != nil {
			Ret403(c, AEC_slot_spin_badbet, err)
			return
		}
	}
	if arg.Sel != 0 {
		if err = game.SetSel(arg.Sel); err != nil {
			Ret403(c, AEC_slot_spin_badsel, err)
			return
		}
	}

	var cost float64
	var isjp bool
	if !game.Free() {
		cost, isjp = game.Cost()
	}

	if props.Wallet < cost {
		Ret403(c, AEC_slot_spin_nomoney, ErrNoMoney)
		return
	}

	var bank, fund, _ = club.GetCash()
	var mrtp = GetRTP(user, club)
	var jprate float64
	if isjp {
		jprate = club.Rate()
	}

	// spin until gain less than bank value
	var wins slot.Wins
	var debit, jack float64
	defer wins.Reset()
	var n = 0
	game.Prepare()
	for { // repeat until spin will fit into bank
		for { // repeat until get valid screen
			game.Spin(mrtp - jprate)
			if game.Scanner(&wins) == nil {
				break
			}
			n++
			if n > cfg.Cfg.MaxSpinAttempts {
				Ret500(c, AEC_slot_spin_badbank, ErrBadBank)
				return
			}
		}
		game.Spawn(wins, fund, mrtp-jprate)
		debit = cost*(1-jprate/100) - wins.Gain()
		jack = wins.Jackpot()
		if (bank+debit >= 0 || debit > 0) && (jack == 0 || jack > Cfg.MinJackpot) {
			break
		}
		wins.Reset()
	}

	// write gain and total bet as transaction
	if Cfg.ClubUpdateBuffer > 1 {
		go BankBat[scene.CID].Put(cfg.XormStorage, scene.UID, debit)
	} else if err = BankBat[scene.CID].Put(cfg.XormStorage, scene.UID, debit); err != nil {
		Ret500(c, AEC_slot_spin_sqlbank, err)
		return
	}

	// make changes to memory data
	var jprent = cost * jprate / 100
	club.AddCash(debit, jprent-jack, 0)
	props.Wallet += wins.Gain() - cost
	game.Apply(wins)

	// write spin result to log and get spin ID
	var sid = atomic.AddUint64(&SpinCounter, 1)
	scene.SID = sid
	var rec = Spinlog{
		SID:    sid,
		GID:    arg.GID,
		MRTP:   mrtp,
		Gain:   game.GetGain(),
		Wallet: props.Wallet,
	}
	var b []byte

	if b, err = json.Marshal(scene.Game); err != nil {
		return
	}
	rec.Game = util.B2S(b)

	if len(wins) > 0 {
		if b, err = json.Marshal(wins); err != nil {
			return
		}
		rec.Wins = util.B2S(b)
	}

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
	if isjp {
		ret.JpFund = fund
	}

	RetOk(c, ret)
}

// Double up gamble on last gain.
func ApiSlotDoubleup(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid" binding:"required"`
		Mult    float64  `json:"mult" yaml:"mult" xml:"mult" form:"mult" binding:"gt=1,lte=10"`
		Half    bool     `json:"half" yaml:"half" xml:"half" form:"half"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		ID      uint64   `json:"id" yaml:"id" xml:"id,attr"`
		Win     bool     `json:"win" yaml:"win" xml:"win"`
		Risk    float64  `json:"risk" yaml:"risk" xml:"risk"`
		Gain    float64  `json:"gain" yaml:"gain" xml:"gain"`
		Wallet  float64  `json:"wallet" yaml:"wallet" xml:"wallet"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_slot_doubleup_nobind, err)
		return
	}

	var scene *Scene
	if scene, err = GetScene(arg.GID); err != nil {
		Ret404(c, AEC_slot_doubleup_noscene, err)
		return
	}
	var game slot.SlotGame
	if game, ok = scene.Game.(slot.SlotGame); !ok {
		Ret403(c, AEC_slot_doubleup_notslot, ErrNotSlot)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(scene.CID); !ok {
		Ret500(c, AEC_slot_doubleup_noclub, ErrNoClub)
		return
	}

	var user *User
	if user, ok = Users.Get(scene.UID); !ok {
		Ret500(c, AEC_slot_doubleup_nouser, ErrNoUser)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALdealer == 0 {
		Ret403(c, AEC_slot_doubleup_noaccess, ErrNoAccess)
		return
	}

	var props *Props
	if props, ok = user.props.Get(scene.CID); !ok {
		Ret500(c, AEC_slot_doubleup_noprops, ErrNoProps)
		return
	}

	var oldgain = game.GetGain()
	var risk = oldgain
	if risk == 0 {
		Ret403(c, AEC_slot_doubleup_nogain, ErrNoGain)
		return
	}
	if arg.Half {
		risk /= 2
	}

	var bank = club.Bank()
	var mrtp = GetRTP(user, club)

	var win bool       // true on double up is win
	var upgain float64 // gain by double up
	if bank >= risk*arg.Mult {
		var r = rand.Float64()
		var side = 1 / arg.Mult * mrtp / 100
		if r < side {
			win = true
			upgain = risk * arg.Mult
		}
	}
	var debit = risk - upgain
	var newgain = oldgain - risk + upgain

	// write gain and total bet as transaction
	if Cfg.ClubUpdateBuffer > 1 {
		go BankBat[scene.CID].Put(cfg.XormStorage, scene.UID, debit)
	} else if err = BankBat[scene.CID].Put(cfg.XormStorage, scene.UID, debit); err != nil {
		Ret500(c, AEC_slot_doubleup_sqlbank, err)
		return
	}

	// make changes to memory data
	club.AddBank(debit)
	props.Wallet -= debit

	game.SetGain(newgain)

	// write doubleup result to log and get spin ID
	var id = atomic.AddUint64(&MultCounter, 1)
	if Cfg.UseSpinLog {
		go func() {
			var rec = Multlog{
				ID:     id,
				GID:    arg.GID,
				MRTP:   mrtp,
				Mult:   arg.Mult,
				Risk:   risk,
				Win:    win,
				Gain:   newgain,
				Wallet: props.Wallet,
			}
			if err = MultBuf.Put(cfg.XormSpinlog, rec); err != nil {
				log.Printf("can not write to mult log: %s", err.Error())
			}
		}()
	}

	// prepare result
	ret.ID = id
	ret.Win = win
	ret.Risk = risk
	ret.Gain = newgain
	ret.Wallet = props.Wallet

	RetOk(c, ret)
}

func ApiSlotCollect(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid" binding:"required"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_slot_collect_nobind, err)
		return
	}

	var scene *Scene
	if scene, err = GetScene(arg.GID); err != nil {
		Ret404(c, AEC_slot_collect_noscene, err)
		return
	}
	var game slot.SlotGame
	if game, ok = scene.Game.(slot.SlotGame); !ok {
		Ret403(c, AEC_slot_collect_notslot, ErrNotSlot)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALdealer == 0 {
		Ret403(c, AEC_slot_collect_noaccess, ErrNoAccess)
		return
	}

	if err = game.SetGain(0); err != nil {
		Ret403(c, AEC_slot_collect_denied, err)
		return
	}

	RetOk(c, nil)
}
