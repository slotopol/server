package spi

import (
	"encoding/xml"
	"log"
	"math/rand/v2"
	"net/http"
	"sync/atomic"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/config/links"
	"github.com/slotopol/server/util"

	"github.com/gin-gonic/gin"
	"github.com/slotopol/server/game"
)

var (
	SpinBuf util.SqlBuf[Spinlog]
	MultBuf util.SqlBuf[Multlog]
	BankBat = map[uint64]*SqlBank{}
)

// Joins to game and creates new instance of game.
func SpiGameJoin(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr" form:"uid"`
		Alias   string   `json:"alias" yaml:"alias" xml:"alias" form:"alias" binding:"required"`
	}
	var ret struct {
		XMLName xml.Name      `json:"-" yaml:"-" xml:"ret"`
		GID     uint64        `json:"gid" yaml:"gid" xml:"gid,attr"`
		Game    game.SlotGame `json:"game" yaml:"game" xml:"game"`
		Scrn    game.Screen   `json:"scrn" yaml:"scrn" xml:"scrn"`
		Wallet  float64       `json:"wallet" yaml:"wallet" xml:"wallet"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_game_join_nobind, err)
		return
	}
	if arg.CID == 0 {
		Ret400(c, SEC_game_join_norid, ErrNoCID)
		return
	}
	if arg.UID == 0 {
		Ret400(c, SEC_game_join_nouid, ErrNoUID)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(arg.CID); !ok {
		Ret404(c, SEC_game_join_noclub, ErrNoClub)
		return
	}
	_ = club

	var user *User
	if user, ok = Users.Get(arg.UID); !ok {
		Ret404(c, SEC_game_join_nouser, ErrNoUser)
		return
	}

	var admin, al = GetAdmin(c, arg.CID)
	if admin != user && al&ALgame == 0 {
		Ret403(c, SEC_prop_join_noaccess, ErrNoAccess)
		return
	}

	var alias = util.ToID(arg.Alias)
	var maker, has = links.GameFactory[alias]
	if !has {
		Ret400(c, SEC_game_join_noalias, ErrNoAliase)
		return
	}
	var slotgame = maker(club.GainRTP)
	if slotgame == nil {
		Ret400(c, SEC_game_join_noreels, ErrNoReels)
		return
	}

	var scene = &Scene{
		Story: Story{
			Alias: alias,
			CID:   arg.CID,
			UID:   arg.UID,
			Flow:  true,
		},
		Game: slotgame.(game.SlotGame),
		Scrn: slotgame.(game.SlotGame).NewScreen(),
	}
	// make game screen object
	scene.Game.Spin(scene.Scrn)

	if err = SafeTransaction(cfg.XormStorage, func(session *Session) (err error) {
		if _, err = session.Insert(&scene.Story); err != nil {
			Ret500(c, SEC_game_join_open, err)
			return
		}

		// ensure that wallet record is exist
		if !user.props.Has(arg.CID) {
			var props = &Props{
				CID: arg.CID,
				UID: arg.UID,
			}
			if _, err = session.Insert(props); err != nil {
				Ret500(c, SEC_game_join_props, err)
				return
			}

			user.InsertProps(props)
		}

		return
	}); err != nil {
		return
	}

	Scenes.Set(scene.GID, scene)
	user.games.Set(scene.GID, scene)

	ret.GID = scene.GID
	ret.Game = scene.Game
	ret.Scrn = scene.Scrn
	ret.Wallet = user.GetWallet(arg.CID)

	RetOk(c, ret)
}

// Removes instance of opened game.
func SpiGamePart(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_game_part_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_game_part_nogid, ErrNoGID)
		return
	}

	var scene *Scene
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, SEC_game_part_notopened, ErrNotOpened)
		return
	}

	var user *User
	if user, ok = Users.Get(scene.UID); !ok {
		Ret500(c, SEC_game_part_nouser, ErrNoUser)
		return
	}

	var admin, al = GetAdmin(c, scene.CID)
	if admin != user && al&ALgame == 0 {
		Ret403(c, SEC_prop_part_noaccess, ErrNoAccess)
		return
	}

	scene.Flow = false
	if _, err = cfg.XormStorage.ID(arg.GID).Cols("flow").Update(&scene.Story); err != nil {
		Ret500(c, SEC_prop_part_update, err)
		return
	}

	Scenes.Delete(arg.GID)
	user.games.Delete(arg.GID)

	c.Status(http.StatusOK)
}

// Returns full info of game scene with given GID, and balance on wallet.
func SpiGameInfo(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid"`
	}
	var ret struct {
		XMLName xml.Name      `json:"-" yaml:"-" xml:"ret"`
		GID     uint64        `json:"gid" yaml:"gid" xml:"gid,attr"`
		Alias   string        `json:"alias" yaml:"alias" xml:"alias"`
		CID     uint64        `json:"cid" yaml:"cid" xml:"cid,attr"`
		UID     uint64        `json:"uid" yaml:"uid" xml:"uid,attr"`
		SID     uint64        `json:"sid" yaml:"sid" xml:"sid,attr"`
		Game    game.SlotGame `json:"game" yaml:"game" xml:"game"`
		Scrn    game.Screen   `json:"scrn" yaml:"scrn" xml:"scrn"`
		Wins    game.Wins     `json:"wins,omitempty" yaml:"wins,omitempty" xml:"wins,omitempty"`
		Wallet  float64       `json:"wallet" yaml:"wallet" xml:"wallet"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_game_info_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_game_info_nogid, ErrNoGID)
		return
	}

	var scene *Scene
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, SEC_game_info_notopened, ErrNotOpened)
		return
	}

	var user *User
	if user, ok = Users.Get(scene.UID); !ok {
		Ret500(c, SEC_game_info_nouser, ErrNoUser)
		return
	}

	var admin, al = GetAdmin(c, scene.CID)
	if admin != user && al&ALgame == 0 {
		Ret403(c, SEC_prop_state_noaccess, ErrNoAccess)
		return
	}

	var props *Props
	if props, ok = user.props.Get(scene.CID); !ok {
		Ret500(c, SEC_game_info_noprops, ErrNoWallet)
		return
	}

	ret.GID = arg.GID
	ret.Alias = scene.Alias
	ret.CID = scene.CID
	ret.UID = scene.UID
	ret.SID = scene.SID
	ret.Game = scene.Game
	ret.Scrn = scene.Scrn
	ret.Wins = scene.Wins
	ret.Wallet = props.Wallet

	RetOk(c, ret)
}

// Returns bet value.
func SpiGameBetGet(c *gin.Context) {
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
		Ret400(c, SEC_game_betget_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_game_betget_nogid, ErrNoGID)
		return
	}

	var scene *Scene
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, SEC_game_betget_notopened, ErrNotOpened)
		return
	}

	var admin, al = GetAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALgame == 0 {
		Ret403(c, SEC_prop_betget_noaccess, ErrNoAccess)
		return
	}

	ret.Bet = scene.Game.GetBet()

	RetOk(c, ret)
}

// Set bet value.
func SpiGameBetSet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr"`
		Bet     float64  `json:"bet" yaml:"bet" xml:"bet" binding:"required"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_game_betset_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_game_betset_nogid, ErrNoGID)
		return
	}

	var scene *Scene
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, SEC_game_betset_notopened, ErrNotOpened)
		return
	}

	var admin, al = GetAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALgame == 0 {
		Ret403(c, SEC_prop_betset_noaccess, ErrNoAccess)
		return
	}

	if err = scene.Game.SetBet(arg.Bet); err != nil {
		Ret403(c, SEC_game_betset_badbet, err)
		return
	}

	c.Status(http.StatusOK)
}

// Returns selected bet lines bitset.
func SpiGameSblGet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid"`
	}
	var ret struct {
		XMLName xml.Name    `json:"-" yaml:"-" xml:"ret"`
		SBL     game.Bitset `json:"sbl" yaml:"sbl" xml:"sbl"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_game_sblget_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_game_sblget_nogid, ErrNoGID)
		return
	}

	var scene *Scene
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, SEC_game_sblget_notopened, ErrNotOpened)
		return
	}

	var admin, al = GetAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALgame == 0 {
		Ret403(c, SEC_prop_sblget_noaccess, ErrNoAccess)
		return
	}

	ret.SBL = scene.Game.GetLines()

	RetOk(c, ret)
}

// Set selected bet lines bitset.
func SpiGameSblSet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name    `json:"-" yaml:"-" xml:"arg"`
		GID     uint64      `json:"gid" yaml:"gid" xml:"gid,attr"`
		SBL     game.Bitset `json:"sbl" yaml:"sbl" xml:"sbl" binding:"required"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_game_sblset_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_game_sblset_nogid, ErrNoGID)
		return
	}

	var scene *Scene
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, SEC_game_sblset_notopened, ErrNotOpened)
		return
	}

	var admin, al = GetAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALgame == 0 {
		Ret403(c, SEC_prop_sblset_noaccess, ErrNoAccess)
		return
	}

	if err = scene.Game.SetLines(arg.SBL); err != nil {
		Ret403(c, SEC_game_sblset_badlines, err)
		return
	}

	c.Status(http.StatusOK)
}

// Returns reels descriptor for given GID.
func SpiGameReelsGet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		RTP     float64  `json:"rtp" yaml:"rtp" xml:"rtp"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_game_rdget_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_game_rdget_nogid, ErrNoGID)
		return
	}

	var scene *Scene
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, SEC_game_rdget_notopened, ErrNotOpened)
		return
	}

	var admin, al = GetAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALgame == 0 {
		Ret403(c, SEC_prop_rdget_noaccess, ErrNoAccess)
		return
	}

	ret.RTP = scene.Game.GetRTP()

	RetOk(c, ret)
}

// Set reels descriptor for given GID. Only game admin can change reels.
func SpiGameReelsSet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr"`
		RTP     float64  `json:"rtp" yaml:"rtp" xml:"rtp" binding:"required"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_game_rdset_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_game_rdset_nogid, ErrNoGID)
		return
	}

	var scene *Scene
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, SEC_game_rdset_notopened, ErrNotOpened)
		return
	}

	// only game admin can change reels
	var _, al = GetAdmin(c, scene.CID)
	if al&ALgame == 0 {
		Ret403(c, SEC_prop_rdset_noaccess, ErrNoAccess)
		return
	}

	if err = scene.Game.SetRTP(arg.RTP); err != nil {
		Ret403(c, SEC_game_rdset_badreels, err)
		return
	}

	c.Status(http.StatusOK)
}

// Make a spin.
func SpiGameSpin(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid"`
	}
	var ret struct {
		XMLName xml.Name      `json:"-" yaml:"-" xml:"ret"`
		SID     uint64        `json:"sid" yaml:"sid" xml:"sid,attr"`
		Game    game.SlotGame `json:"game" yaml:"game" xml:"game"`
		Scrn    game.Screen   `json:"scrn" yaml:"scrn" xml:"scrn"`
		Wins    game.Wins     `json:"wins,omitempty" yaml:"wins,omitempty" xml:"wins,omitempty"`
		Wallet  float64       `json:"wallet" yaml:"wallet" xml:"wallet"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_game_spin_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_game_spin_nogid, ErrNoGID)
		return
	}

	var scene *Scene
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, SEC_game_spin_notopened, ErrNotOpened)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(scene.CID); !ok {
		Ret500(c, SEC_game_spin_noclub, ErrNoClub)
		return
	}

	var user *User
	if user, ok = Users.Get(scene.UID); !ok {
		Ret500(c, SEC_game_spin_nouser, ErrNoUser)
		return
	}

	var admin, al = GetAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALgame == 0 {
		Ret403(c, SEC_prop_spin_noaccess, ErrNoAccess)
		return
	}

	var (
		fs       = scene.Game.FreeSpins()
		bet      = scene.Game.GetBet()
		sbl      = scene.Game.GetLines()
		totalbet float64
		banksum  float64
	)
	if fs == 0 {
		totalbet = bet * float64(sbl.Num())
	}

	var props *Props
	if props, ok = user.props.Get(scene.CID); !ok {
		Ret500(c, SEC_game_spin_noprops, ErrNoWallet)
		return
	}
	if props.Wallet < totalbet {
		Ret403(c, SEC_game_spin_nomoney, ErrNoMoney)
		return
	}

	// spin until gain less than bank value
	club.mux.RLock()
	var bank = club.Bank
	club.mux.RUnlock()
	var wins game.Wins
	var n = 0
	for {
		scene.Game.Spin(scene.Scrn)
		scene.Game.Scanner(scene.Scrn, &wins)
		scene.Game.Spawn(scene.Scrn, wins)
		banksum = totalbet - wins.Gain()
		if bank+banksum >= 0 || (bank < 0 && banksum > 0) {
			break
		}
		wins.Reset()
		if n >= cfg.Cfg.MaxSpinAttempts {
			Ret500(c, SEC_game_spin_badbank, ErrBadBank)
			return
		}
		n++
	}

	// write gain and total bet as transaction
	if Cfg.BankBufferSize > 1 {
		go BankBat[scene.CID].Put(cfg.XormStorage, scene.UID, banksum)
	} else if err = BankBat[scene.CID].Put(cfg.XormStorage, scene.UID, banksum); err != nil {
		Ret500(c, SEC_game_spin_sqlbank, err)
		return
	}

	// make changes to memory data
	club.mux.Lock()
	club.Bank += banksum
	club.mux.Unlock()

	props.Wallet -= banksum
	scene.Game.Apply(scene.Scrn, wins)
	scene.Wins.Reset() // throw old wins
	scene.Wins = wins

	// write spin result to log and get spin ID
	var sid = atomic.AddUint64(&SpinCounter, 1)
	scene.SID = sid
	go func() {
		var rec = Spinlog{
			SID:    sid,
			GID:    arg.GID,
			Gain:   scene.Game.GetGain(),
			Wallet: props.Wallet,
		}
		_ = rec.MarshalState(scene)
		if err = SpinBuf.Put(cfg.XormSpinlog, rec); err != nil {
			log.Printf("can not write to spin log: %s", err.Error())
		}
	}()

	// prepare result
	ret.SID = sid
	ret.Game = scene.Game
	ret.Scrn = scene.Scrn
	ret.Wins = scene.Wins
	ret.Wallet = props.Wallet

	RetOk(c, ret)
}

// Double up gamble on last gain.
func SpiGameDoubleup(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid"`
		Mult    int      `json:"mult" yaml:"mult" xml:"mult" form:"mult"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		ID      uint64   `json:"id" yaml:"id" xml:"id,attr"`
		Gain    float64  `json:"gain" yaml:"gain" xml:"gain"`
		Wallet  float64  `json:"wallet" yaml:"wallet" xml:"wallet"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_game_doubleup_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_game_doubleup_nogid, ErrNoGID)
		return
	}
	if arg.Mult < 2 {
		Ret400(c, SEC_game_doubleup_nomult, ErrNoMult)
		return
	}
	if arg.Mult > 10 {
		Ret400(c, SEC_game_doubleup_bigmult, ErrBigMult)
		return
	}

	var scene *Scene
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, SEC_game_doubleup_notopened, ErrNotOpened)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(scene.CID); !ok {
		Ret500(c, SEC_game_doubleup_noclub, ErrNoClub)
		return
	}

	var user *User
	if user, ok = Users.Get(scene.UID); !ok {
		Ret500(c, SEC_game_doubleup_nouser, ErrNoUser)
		return
	}

	var admin, al = GetAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALgame == 0 {
		Ret403(c, SEC_prop_doubleup_noaccess, ErrNoAccess)
		return
	}

	var props *Props
	if props, ok = user.props.Get(scene.CID); !ok {
		Ret500(c, SEC_game_doubleup_noprops, ErrNoWallet)
		return
	}

	var risk = scene.Game.GetGain()
	if risk == 0 {
		Ret403(c, SEC_game_doubleup_nomoney, ErrNoMoney)
		return
	}

	club.mux.RLock()
	var bank = club.Bank
	var rtp = club.GainRTP
	club.mux.RUnlock()

	var multgain float64 // new multiplied gain
	if bank >= risk*float64(arg.Mult) {
		var r = rand.Float64()
		var side = 1 / float64(arg.Mult) * rtp / 100
		if r < side {
			multgain = risk * float64(arg.Mult)
		}
	}
	var banksum = risk - multgain

	// write gain and total bet as transaction
	if Cfg.BankBufferSize > 1 {
		go BankBat[scene.CID].Put(cfg.XormStorage, scene.UID, banksum)
	} else if err = BankBat[scene.CID].Put(cfg.XormStorage, scene.UID, banksum); err != nil {
		Ret500(c, SEC_game_doubleup_sqlbank, err)
		return
	}

	// make changes to memory data
	club.mux.Lock()
	club.Bank += banksum
	club.mux.Unlock()

	props.Wallet -= banksum

	scene.Game.SetGain(multgain)
	scene.Wins.Reset()

	// write doubleup result to log and get spin ID
	var id = atomic.AddUint64(&MultCounter, 1)
	go func() {
		var rec = Multlog{
			ID:     id,
			GID:    arg.GID,
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

func SpiGameCollect(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_game_collect_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_game_collect_nogid, ErrNoGID)
		return
	}

	var scene *Scene
	if scene, ok = Scenes.Get(arg.GID); !ok {
		Ret404(c, SEC_game_collect_notopened, ErrNotOpened)
		return
	}

	var admin, al = GetAdmin(c, scene.CID)
	if admin.UID != scene.UID && al&ALgame == 0 {
		Ret403(c, SEC_prop_collect_noaccess, ErrNoAccess)
		return
	}

	if err = scene.Game.SetGain(0); err != nil {
		Ret403(c, SEC_prop_collect_denied, err)
		return
	}

	c.Status(http.StatusOK)
}
