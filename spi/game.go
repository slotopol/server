package spi

import (
	"encoding/xml"
	"log"
	"math/rand"
	"net/http"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/util"
	"xorm.io/xorm"

	"github.com/gin-gonic/gin"
	"github.com/slotopol/server/game"
)

// Joins to game and creates new instance of game.
func SpiGameJoin(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr" form:"uid"`
		Alias   string   `json:"alias" yaml:"alias" xml:"alias" form:"alias"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr"`
		State   `yaml:",inline"`
		Wallet  int `json:"wallet" yaml:"wallet" xml:"wallet"`
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
	if arg.Alias == "" {
		Ret400(c, SEC_game_join_nodata, ErrNoData)
		return
	}

	var alias = util.ToID(arg.Alias)
	var gname string
	if gname, ok = cfg.GameAliases[alias]; !ok {
		Ret400(c, SEC_game_join_noalias, ErrNoAliase)
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

	var maker = cfg.GameFactory[gname]
	var slotgame = maker("96")
	if slotgame == nil {
		Ret400(c, SEC_game_join_noreels, ErrNoReels)
		return
	}

	var og = OpenGame{
		CID:   arg.CID,
		UID:   arg.UID,
		Alias: alias,
		State: State{
			Game: slotgame.(game.SlotGame),
			Scrn: slotgame.(game.SlotGame).NewScreen(),
		},
	}
	// make game screen object
	og.Game.Spin(og.Scrn)

	if _, err = cfg.XormStorage.Transaction(func(session *xorm.Session) (_ interface{}, err error) {
		if _, err = session.Insert(&og); err != nil {
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

			user.props.Set(arg.CID, props)
		}

		return
	}); err != nil {
		return
	}

	OpenGames.Set(og.GID, og)
	user.games.Set(og.GID, og)

	ret.GID = og.GID
	ret.State = og.State
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

	var og OpenGame
	if og, ok = OpenGames.Get(arg.GID); !ok {
		Ret404(c, SEC_game_part_notopened, ErrNotOpened)
		return
	}

	var user *User
	if user, ok = Users.Get(og.UID); !ok {
		Ret500(c, SEC_game_part_nouser, ErrNoUser)
		return
	}

	var admin, al = GetAdmin(c, og.CID)
	if admin != user && al&ALgame == 0 {
		Ret403(c, SEC_prop_part_noaccess, ErrNoAccess)
		return
	}

	OpenGames.Delete(arg.GID)
	user.games.Delete(arg.GID)

	c.Status(http.StatusOK)
}

// Returns full state of game with given GID, and balance on wallet.
func SpiGameState(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		State   `yaml:",inline"`
		Wallet  int `json:"wallet" yaml:"wallet" xml:"wallet"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_game_state_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_game_state_nogid, ErrNoGID)
		return
	}

	var og OpenGame
	if og, ok = OpenGames.Get(arg.GID); !ok {
		Ret404(c, SEC_game_state_notopened, ErrNotOpened)
		return
	}

	var user *User
	if user, ok = Users.Get(og.UID); !ok {
		Ret500(c, SEC_game_state_nouser, ErrNoUser)
		return
	}

	var admin, al = GetAdmin(c, og.CID)
	if admin != user && al&ALgame == 0 {
		Ret403(c, SEC_prop_state_noaccess, ErrNoAccess)
		return
	}

	var props *Props
	if props, ok = user.props.Get(og.CID); !ok {
		Ret500(c, SEC_game_state_noprops, ErrNoWallet)
		return
	}

	ret.State = og.State
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
		Bet     int      `json:"bet" yaml:"bet" xml:"bet"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_game_betget_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_game_betget_nogid, ErrNoGID)
		return
	}

	var og OpenGame
	if og, ok = OpenGames.Get(arg.GID); !ok {
		Ret404(c, SEC_game_betget_notopened, ErrNotOpened)
		return
	}

	var admin, al = GetAdmin(c, og.CID)
	if admin.UID != og.UID && al&ALgame == 0 {
		Ret403(c, SEC_prop_betget_noaccess, ErrNoAccess)
		return
	}

	ret.Bet = og.Game.GetBet()

	RetOk(c, ret)
}

// Set bet value.
func SpiGameBetSet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr"`
		Bet     int      `json:"bet" yaml:"bet" xml:"bet"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_game_betset_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_game_betset_nogid, ErrNoGID)
		return
	}
	if arg.Bet == 0 {
		Ret400(c, SEC_game_betset_nodata, ErrNoData)
		return
	}

	var og OpenGame
	if og, ok = OpenGames.Get(arg.GID); !ok {
		Ret404(c, SEC_game_betset_notopened, ErrNotOpened)
		return
	}

	var admin, al = GetAdmin(c, og.CID)
	if admin.UID != og.UID && al&ALgame == 0 {
		Ret403(c, SEC_prop_betset_noaccess, ErrNoAccess)
		return
	}

	if err = og.Game.SetBet(arg.Bet); err != nil {
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
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		SBL     game.SBL `json:"sbl" yaml:"sbl" xml:"sbl"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_game_sblget_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_game_sblget_nogid, ErrNoGID)
		return
	}

	var og OpenGame
	if og, ok = OpenGames.Get(arg.GID); !ok {
		Ret404(c, SEC_game_sblget_notopened, ErrNotOpened)
		return
	}

	var admin, al = GetAdmin(c, og.CID)
	if admin.UID != og.UID && al&ALgame == 0 {
		Ret403(c, SEC_prop_sblget_noaccess, ErrNoAccess)
		return
	}

	ret.SBL = og.Game.GetLines()

	RetOk(c, ret)
}

// Set selected bet lines bitset.
func SpiGameSblSet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr"`
		SBL     game.SBL `json:"sbl" yaml:"sbl" xml:"sbl"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_game_sblset_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_game_sblset_nogid, ErrNoGID)
		return
	}
	if arg.SBL == 0 {
		Ret400(c, SEC_game_sblset_nodata, ErrNoData)
		return
	}

	var og OpenGame
	if og, ok = OpenGames.Get(arg.GID); !ok {
		Ret404(c, SEC_game_sblset_notopened, ErrNotOpened)
		return
	}

	var admin, al = GetAdmin(c, og.CID)
	if admin.UID != og.UID && al&ALgame == 0 {
		Ret403(c, SEC_prop_sblset_noaccess, ErrNoAccess)
		return
	}

	if err = og.Game.SetLines(arg.SBL); err != nil {
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
		RD      string   `json:"rd" yaml:"rd" xml:"rd"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_game_rdget_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_game_rdget_nogid, ErrNoGID)
		return
	}

	var og OpenGame
	if og, ok = OpenGames.Get(arg.GID); !ok {
		Ret404(c, SEC_game_rdget_notopened, ErrNotOpened)
		return
	}

	var admin, al = GetAdmin(c, og.CID)
	if admin.UID != og.UID && al&ALgame == 0 {
		Ret403(c, SEC_prop_rdget_noaccess, ErrNoAccess)
		return
	}

	ret.RD = og.Game.GetReels()

	RetOk(c, ret)
}

// Set reels descriptor for given GID. Only game admin can change reels.
func SpiGameReelsSet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr"`
		RD      string   `json:"rd" yaml:"rd" xml:"rd"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_game_rdset_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_game_rdset_nogid, ErrNoGID)
		return
	}
	if arg.RD == "" {
		Ret400(c, SEC_game_rdset_nodata, ErrNoData)
		return
	}

	var og OpenGame
	if og, ok = OpenGames.Get(arg.GID); !ok {
		Ret404(c, SEC_game_rdset_notopened, ErrNotOpened)
		return
	}

	// only game admin can change reels
	var _, al = GetAdmin(c, og.CID)
	if al&ALgame == 0 {
		Ret403(c, SEC_prop_rdset_noaccess, ErrNoAccess)
		return
	}

	if err = og.Game.SetReels(arg.RD); err != nil {
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
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		SID     uint64   `json:"sid" yaml:"sid" xml:"sid,attr" form:"sid"`
		State   `yaml:",inline"`
		Wallet  int `json:"wallet" yaml:"wallet" xml:"wallet"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_game_spin_nobind, err)
		return
	}
	if arg.GID == 0 {
		Ret400(c, SEC_game_spin_nogid, ErrNoGID)
		return
	}

	var og OpenGame
	if og, ok = OpenGames.Get(arg.GID); !ok {
		Ret404(c, SEC_game_spin_notopened, ErrNotOpened)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(og.CID); !ok {
		Ret500(c, SEC_game_spin_noclub, ErrNoClub)
		return
	}

	var user *User
	if user, ok = Users.Get(og.UID); !ok {
		Ret500(c, SEC_game_spin_nouser, ErrNoUser)
		return
	}

	var admin, al = GetAdmin(c, og.CID)
	if admin.UID != og.UID && al&ALgame == 0 {
		Ret403(c, SEC_prop_spin_noaccess, ErrNoAccess)
		return
	}

	var (
		fs       = og.Game.FreeSpins()
		bet      = og.Game.GetBet()
		sbl      = og.Game.GetLines()
		totalbet int
		totalwin int
	)
	if fs == 0 {
		totalbet = bet * sbl.Num()
	}

	var props *Props
	if props, ok = user.props.Get(og.CID); !ok {
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
	var ws game.WinScan
	var n = 0
	for {
		og.Game.Spin(og.Scrn)
		og.Game.Scanner(og.Scrn, &ws)
		og.Game.Spawn(og.Scrn, &ws)
		totalwin = ws.Gain()
		if bank+float64(totalbet-totalwin) >= 0 || (bank < 0 && totalbet > totalwin) {
			break
		}
		if n >= cfg.Cfg.MaxSpinAttempts {
			Ret500(c, SEC_game_spin_badbank, ErrBadBank)
			return
		}
		ws.Reset()
		n++
	}

	// write gain and total bet as transaction
	if _, err = cfg.XormStorage.Transaction(func(session *xorm.Session) (_ interface{}, err error) {
		defer func() {
			if err != nil {
				session.Rollback()
			}
		}()

		const sql1 = `UPDATE club SET bank=bank+? WHERE cid=?`
		if _, err = session.Exec(sql1, totalbet-totalwin, club.CID); err != nil {
			Ret500(c, SEC_game_spin_sqlbank, err)
			return
		}

		const sql2 = `UPDATE props SET wallet=wallet+? WHERE uid=? AND cid=?`
		if _, err = session.Exec(sql2, totalwin-totalbet, props.UID, props.CID); err != nil {
			Ret500(c, SEC_game_spin_sqlupdate, err)
			return
		}

		return
	}); err != nil {
		return
	}

	// make changes to memory data
	club.mux.Lock()
	club.Bank += float64(totalbet - totalwin)
	club.mux.Unlock()

	props.Wallet += totalwin - totalbet
	og.Game.Apply(og.Scrn, &ws)
	og.WinScan.Reset() // throw old wins
	og.WinScan = ws

	// write spin result to log and get spin ID
	var sl = Spinlog{
		GID:    arg.GID,
		Gain:   og.Game.GetGain(),
		Wallet: props.Wallet,
	}
	_ = sl.MarshalState(&og.State)
	if _, err = cfg.XormSpinlog.Insert(&sl); err != nil {
		log.Printf("can not write to spin log: %s", err.Error())
	}

	// prepare result
	ret.SID = sl.SID
	ret.State = og.State
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
		SID     uint64   `json:"sid" yaml:"sid" xml:"sid,attr" form:"sid"`
		Gain    int      `json:"gain" yaml:"gain" xml:"gain"`
		Wallet  int      `json:"wallet" yaml:"wallet" xml:"wallet"`
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

	var og OpenGame
	if og, ok = OpenGames.Get(arg.GID); !ok {
		Ret404(c, SEC_game_doubleup_notopened, ErrNotOpened)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(og.CID); !ok {
		Ret500(c, SEC_game_doubleup_noclub, ErrNoClub)
		return
	}

	var user *User
	if user, ok = Users.Get(og.UID); !ok {
		Ret500(c, SEC_game_doubleup_nouser, ErrNoUser)
		return
	}

	var admin, al = GetAdmin(c, og.CID)
	if admin.UID != og.UID && al&ALgame == 0 {
		Ret403(c, SEC_prop_doubleup_noaccess, ErrNoAccess)
		return
	}

	var props *Props
	if props, ok = user.props.Get(og.CID); !ok {
		Ret500(c, SEC_game_doubleup_noprops, ErrNoWallet)
		return
	}

	var gain = og.Game.GetGain()
	if gain == 0 {
		Ret403(c, SEC_game_doubleup_nomoney, ErrNoMoney)
		return
	}

	club.mux.RLock()
	var bank = club.Bank
	var rtp = club.GainRTP
	club.mux.RUnlock()

	var multgain int // new multiplied gain
	if bank >= float64(gain*arg.Mult) {
		var r = rand.Float64()
		var side = 1 / float64(arg.Mult) * rtp / 100
		if r < side {
			multgain = gain * arg.Mult
		}
	}

	// write gain and total bet as transaction
	if _, err = cfg.XormStorage.Transaction(func(session *xorm.Session) (_ interface{}, err error) {
		defer func() {
			if err != nil {
				session.Rollback()
			}
		}()

		const sql1 = `UPDATE club SET bank=bank-? WHERE cid=?`
		if _, err = session.Exec(sql1, multgain-gain, club.CID); err != nil {
			Ret500(c, SEC_game_spin_sqlbank, err)
			return
		}

		const sql2 = `UPDATE props SET wallet=wallet+? WHERE uid=? AND cid=?`
		if _, err = session.Exec(sql2, multgain-gain, props.UID, props.CID); err != nil {
			Ret500(c, SEC_game_spin_sqlupdate, err)
			return
		}

		return
	}); err != nil {
		return
	}

	// make changes to memory data
	club.mux.Lock()
	club.Bank -= float64(multgain - gain)
	club.mux.Unlock()

	props.Wallet += multgain - gain

	og.Game.SetGain(multgain)
	og.WinScan.Reset()

	// write doubleup result to log and get spin ID
	var sl = Spinlog{
		GID:    arg.GID,
		Gain:   multgain,
		Wallet: props.Wallet,
	}
	_ = sl.MarshalState(&og.State)
	if _, err = cfg.XormSpinlog.Insert(&sl); err != nil {
		log.Printf("can not write to spin log: %s", err.Error())
	}

	// prepare result
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

	var og OpenGame
	if og, ok = OpenGames.Get(arg.GID); !ok {
		Ret404(c, SEC_game_collect_notopened, ErrNotOpened)
		return
	}

	var admin, al = GetAdmin(c, og.CID)
	if admin.UID != og.UID && al&ALgame == 0 {
		Ret403(c, SEC_prop_collect_noaccess, ErrNoAccess)
		return
	}

	if err = og.Game.SetGain(0); err != nil {
		Ret403(c, SEC_prop_collect_denied, err)
		return
	}

	c.Status(http.StatusOK)
}
