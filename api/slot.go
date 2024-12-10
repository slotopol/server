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
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, AEC_slot_betget_notopened, ErrNotOpened)
		return
	}
	var game slot.SlotGame
	if game, ok = scene.Game.(slot.SlotGame); !ok {
		Ret403(c, AEC_slot_betget_notslot, ErrNotSlot)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALgame == 0 {
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
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, AEC_slot_betset_notopened, ErrNotOpened)
		return
	}
	var game slot.SlotGame
	if game, ok = scene.Game.(slot.SlotGame); !ok {
		Ret403(c, AEC_slot_betset_notslot, ErrNotSlot)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALgame == 0 {
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
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, AEC_slot_selget_notopened, ErrNotOpened)
		return
	}
	var game slot.SlotGame
	if game, ok = scene.Game.(slot.SlotGame); !ok {
		Ret403(c, AEC_slot_selget_notslot, ErrNotSlot)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALgame == 0 {
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
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, AEC_slot_selset_notopened, ErrNotOpened)
		return
	}
	var game slot.SlotGame
	if game, ok = scene.Game.(slot.SlotGame); !ok {
		Ret403(c, AEC_slot_selset_notslot, ErrNotSlot)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALgame == 0 {
		Ret403(c, AEC_slot_selset_noaccess, ErrNoAccess)
		return
	}

	if err = game.SetSel(arg.Sel); err != nil {
		Ret403(c, AEC_slot_selset_badlines, err)
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
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, AEC_slot_modeset_notopened, ErrNotOpened)
		return
	}
	var game slot.SlotGame
	if game, ok = scene.Game.(slot.SlotGame); !ok {
		Ret403(c, AEC_slot_modeset_notslot, ErrNotSlot)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALgame == 0 {
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
	}
	var ret struct {
		XMLName xml.Name      `json:"-" yaml:"-" xml:"ret"`
		SID     uint64        `json:"sid" yaml:"sid" xml:"sid,attr"`
		Game    slot.SlotGame `json:"game" yaml:"game" xml:"game"`
		Scrn    slot.Screen   `json:"scrn" yaml:"scrn" xml:"scrn"`
		Wins    slot.Wins     `json:"wins,omitempty" yaml:"wins,omitempty" xml:"wins,omitempty"`
		Wallet  float64       `json:"wallet" yaml:"wallet" xml:"wallet"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_slot_spin_nobind, err)
		return
	}

	var scene *Scene
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, AEC_slot_spin_notopened, ErrNotOpened)
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
	if admin.UID != scene.UID && al&ALgame == 0 {
		Ret403(c, AEC_slot_spin_noaccess, ErrNoAccess)
		return
	}

	var (
		bet      = game.GetBet()
		sel      = game.GetSel()
		totalbet float64
		banksum  float64
	)
	if !game.FreeSpins() {
		totalbet = bet * float64(sel)
	}

	var props *Props
	if props, ok = user.props.Get(scene.CID); !ok {
		Ret500(c, AEC_slot_spin_noprops, ErrNoProps)
		return
	}
	if props.Wallet < totalbet {
		Ret403(c, AEC_slot_spin_nomoney, ErrNoMoney)
		return
	}

	var bank = club.Bank()
	var mrtp = GetRTP(user, club)

	// spin until gain less than bank value
	var wins slot.Wins
	defer wins.Reset()
	var scrn = game.NewScreen()
	defer scrn.Free()
	var n = 0
	game.Prepare()
	for {
		game.Spin(scrn, mrtp)
		game.Scanner(scrn, &wins)
		game.Spawn(scrn, wins)
		banksum = totalbet - wins.Gain()
		if bank+banksum >= 0 || (bank < 0 && banksum > 0) {
			break
		}
		wins.Reset()
		if n >= cfg.Cfg.MaxSpinAttempts {
			Ret500(c, AEC_slot_spin_badbank, ErrBadBank)
			return
		}
		n++
	}

	// write gain and total bet as transaction
	if Cfg.ClubUpdateBuffer > 1 {
		go BankBat[scene.CID].Put(cfg.XormStorage, scene.UID, banksum)
	} else if err = BankBat[scene.CID].Put(cfg.XormStorage, scene.UID, banksum); err != nil {
		Ret500(c, AEC_slot_spin_sqlbank, err)
		return
	}

	// make changes to memory data
	club.AddBank(banksum)
	props.Wallet -= banksum
	game.Apply(scrn, wins)

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

	if b, err = json.Marshal(scrn); err != nil {
		return
	}
	rec.Screen = util.B2S(b)

	if len(wins) > 0 {
		if b, err = json.Marshal(wins); err != nil {
			return
		}
		rec.Wins = util.B2S(b)
	}

	go func() {
		if err = SpinBuf.Put(cfg.XormSpinlog, rec); err != nil {
			log.Printf("can not write to spin log: %s", err.Error())
		}
	}()

	// prepare result
	ret.SID = sid
	ret.Game = game
	ret.Scrn = scrn
	ret.Wins = wins
	ret.Wallet = props.Wallet

	RetOk(c, ret)
}

// Double up gamble on last gain.
func ApiSlotDoubleup(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid" binding:"required"`
		Mult    int      `json:"mult" yaml:"mult" xml:"mult" form:"mult" binding:"gte=2,lte=10"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		ID      uint64   `json:"id" yaml:"id" xml:"id,attr"`
		Gain    float64  `json:"gain" yaml:"gain" xml:"gain"`
		Wallet  float64  `json:"wallet" yaml:"wallet" xml:"wallet"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_slot_doubleup_nobind, err)
		return
	}

	var scene *Scene
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, AEC_slot_doubleup_notopened, ErrNotOpened)
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
	if admin.UID != scene.UID && al&ALgame == 0 {
		Ret403(c, AEC_slot_doubleup_noaccess, ErrNoAccess)
		return
	}

	var props *Props
	if props, ok = user.props.Get(scene.CID); !ok {
		Ret500(c, AEC_slot_doubleup_noprops, ErrNoProps)
		return
	}

	var risk = game.GetGain()
	if risk == 0 {
		Ret403(c, AEC_slot_doubleup_nogain, ErrNoGain)
		return
	}

	var bank = club.Bank()
	var mrtp = GetRTP(user, club)

	var multgain float64 // new multiplied gain
	if bank >= risk*float64(arg.Mult) {
		var r = rand.Float64()
		var side = 1 / float64(arg.Mult) * mrtp / 100
		if r < side {
			multgain = risk * float64(arg.Mult)
		}
	}
	var banksum = risk - multgain

	// write gain and total bet as transaction
	if Cfg.ClubUpdateBuffer > 1 {
		go BankBat[scene.CID].Put(cfg.XormStorage, scene.UID, banksum)
	} else if err = BankBat[scene.CID].Put(cfg.XormStorage, scene.UID, banksum); err != nil {
		Ret500(c, AEC_slot_doubleup_sqlbank, err)
		return
	}

	// make changes to memory data
	club.AddBank(banksum)
	props.Wallet -= banksum

	game.SetGain(multgain)

	// write doubleup result to log and get spin ID
	var id = atomic.AddUint64(&MultCounter, 1)
	go func() {
		var rec = Multlog{
			ID:     id,
			GID:    arg.GID,
			MRTP:   mrtp,
			Mult:   arg.Mult,
			Risk:   risk,
			Gain:   multgain,
			Wallet: props.Wallet,
		}
		if err = MultBuf.Put(cfg.XormSpinlog, rec); err != nil {
			log.Printf("can not write to mult log: %s", err.Error())
		}
	}()

	// prepare result
	ret.ID = id
	ret.Gain = multgain
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
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, AEC_slot_collect_notopened, ErrNotOpened)
		return
	}
	var game slot.SlotGame
	if game, ok = scene.Game.(slot.SlotGame); !ok {
		Ret403(c, AEC_slot_collect_notslot, ErrNotSlot)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALgame == 0 {
		Ret403(c, AEC_slot_collect_noaccess, ErrNoAccess)
		return
	}

	if err = game.SetGain(0); err != nil {
		Ret403(c, AEC_slot_collect_denied, err)
		return
	}

	RetOk(c, nil)
}
