package spi

import (
	"encoding/json"
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

func SpiGameJoin(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		RID     uint64   `json:"rid" yaml:"rid" xml:"rid,attr" form:"rid"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr" form:"uid"`
		Alias   string   `json:"alias" yaml:"alias" xml:"alias" form:"alias"`
	}
	var ret struct {
		XMLName xml.Name    `json:"-" yaml:"-" xml:"ret"`
		GID     uint64      `json:"gid" yaml:"gid" xml:"gid,attr"`
		Screen  game.Screen `json:"screen" yaml:"screen" xml:"screen"`
		Wallet  int         `json:"wallet" yaml:"wallet" xml:"wallet"`
	}

	if err = c.Bind(&arg); err != nil {
		Ret400(c, SEC_game_join_nobind, err)
		return
	}
	if arg.RID == 0 {
		Ret400(c, SEC_game_join_norid, ErrNoRID)
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

	var room *Room
	if room, ok = Rooms.Get(arg.RID); !ok {
		Ret404(c, SEC_game_join_noroom, ErrNoRoom)
		return
	}
	_ = room

	var user *User
	if user, ok = Users.Get(arg.UID); !ok {
		Ret404(c, SEC_game_join_nouser, ErrNoUser)
		return
	}

	var alias = util.ToID(arg.Alias)
	var gname string
	if gname, ok = cfg.GameAliases[alias]; !ok {
		Ret400(c, SEC_game_join_noalias, ErrNoAliase)
		return
	}

	var maker = cfg.GameFactory[gname]
	var slotgame = maker("96")
	if slotgame == nil {
		Ret400(c, SEC_game_join_noreels, ErrNoReels)
		return
	}

	var og = OpenGame{
		RID:   arg.RID,
		UID:   arg.UID,
		Alias: alias,
		game:  slotgame.(game.SlotGame),
	}
	if _, err = cfg.XormStorage.Transaction(func(session *xorm.Session) (ret interface{}, err error) {
		if _, err = session.Insert(&og); err != nil {
			Ret500(c, SEC_game_join_open, err)
			return
		}

		// ensure that wallet record is exist
		if !user.props.Has(arg.RID) {
			var props = &Props{
				RID: arg.RID,
				UID: arg.UID,
			}
			if _, err = session.Insert(props); err != nil {
				Ret500(c, SEC_game_join_props, err)
				return
			}

			user.props.Set(arg.RID, props)
		}

		return
	}); err != nil {
		return
	}

	OpenGames.Set(og.GID, og)
	user.games.Set(og.GID, og)

	var scrn = og.game.NewScreen()
	og.game.Spin(scrn)

	ret.GID = og.GID
	ret.Screen = scrn
	ret.Wallet = user.GetWallet(arg.RID)

	RetOk(c, ret)
}

func SpiGamePart(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid"`
	}

	if err = c.Bind(&arg); err != nil {
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

	OpenGames.Delete(arg.GID)
	user.games.Delete(arg.GID)

	c.Status(http.StatusOK)
}

func SpiGameState(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid"`
	}
	var ret struct {
		XMLName xml.Name      `json:"-" yaml:"-" xml:"ret"`
		Game    game.SlotGame `json:"game" yaml:"game" xml:"game"`
		Wallet  int           `json:"wallet" yaml:"wallet" xml:"wallet"`
	}

	if err = c.Bind(&arg); err != nil {
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

	var props *Props
	if props, ok = user.props.Get(og.RID); !ok {
		Ret500(c, SEC_game_state_noprops, ErrNoWallet)
		return
	}

	ret.Game = og.game
	ret.Wallet = props.Wallet

	RetOk(c, ret)
}

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

	if err = c.Bind(&arg); err != nil {
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

	ret.Bet = og.game.GetBet()

	RetOk(c, ret)
}

func SpiGameBetSet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr"`
		Bet     int      `json:"bet" yaml:"bet" xml:"bet"`
	}

	if err = c.Bind(&arg); err != nil {
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

	if err = og.game.SetBet(arg.Bet); err != nil {
		Ret403(c, SEC_game_betset_badbet, err)
		return
	}

	c.Status(http.StatusOK)
}

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

	if err = c.Bind(&arg); err != nil {
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

	ret.SBL = og.game.GetLines()

	RetOk(c, ret)
}

func SpiGameSblSet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr"`
		SBL     game.SBL `json:"sbl" yaml:"sbl" xml:"sbl"`
	}

	if err = c.Bind(&arg); err != nil {
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

	if err = og.game.SetLines(arg.SBL); err != nil {
		Ret403(c, SEC_game_sblset_badlines, err)
		return
	}

	c.Status(http.StatusOK)
}

func SpiGameSpin(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid"`
	}
	var ret struct {
		XMLName xml.Name       `json:"-" yaml:"-" xml:"ret"`
		SID     uint64         `json:"sid" yaml:"sid" xml:"sid,attr" form:"sid"`
		Screen  game.Screen    `json:"screen" yaml:"screen" xml:"screen"`
		Wins    []game.WinItem `json:"wins,omitempty" yaml:"wins,omitempty" xml:"wins,omitempty"`
		FS      int            `json:"fs,omitempty" yaml:"fs,omitempty" xml:"fs,omitempty"`
		Gain    int            `json:"gain" yaml:"gain" xml:"gain"`
		Wallet  int            `json:"wallet" yaml:"wallet" xml:"wallet"`
	}

	if err = c.Bind(&arg); err != nil {
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

	var room *Room
	if room, ok = Rooms.Get(og.RID); !ok {
		Ret500(c, SEC_game_spin_noroom, ErrNoRoom)
		return
	}

	var user *User
	if user, ok = Users.Get(og.UID); !ok {
		Ret500(c, SEC_game_spin_nouser, ErrNoUser)
		return
	}

	var (
		fs       = og.game.FreeSpins()
		bet      = og.game.GetBet()
		sbl      = og.game.GetLines()
		totalbet int
		totalwin int
	)
	if fs == 0 {
		totalbet = bet * sbl.Num()
	}

	var props *Props
	if props, ok = user.props.Get(og.RID); !ok {
		Ret500(c, SEC_game_spin_noprops, ErrNoWallet)
		return
	}
	if props.Wallet < totalbet {
		Ret403(c, SEC_game_spin_nomoney, ErrNoMoney)
		return
	}

	// get game screen object
	var scrn = og.game.NewScreen()

	// spin until gain less than bank value
	room.mux.RLock()
	var bank = room.Bank
	room.mux.RUnlock()
	var ws game.WinScan
	var n = 0
	for {
		og.game.Spin(scrn)
		og.game.Scanner(scrn, &ws)
		og.game.Spawn(scrn, &ws)
		totalwin = ws.Gain()
		if bank+float64(totalbet-totalwin) >= 0 || (bank < 0 && totalbet > totalwin) {
			break
		}
		if n >= cfg.Cfg.MaxSpinAttempts {
			Ret500(c, SEC_game_spin_badbank, ErrBadBank)
			return
		}
		n++
	}

	// write gain and total bet as transaction
	if _, err = cfg.XormStorage.Transaction(func(session *xorm.Session) (ret interface{}, err error) {
		defer func() {
			if err != nil {
				session.Rollback()
			}
		}()

		const sql1 = `UPDATE room SET bank=bank+? WHERE rid=?`
		if ret, err = session.Exec(sql1, totalbet-totalwin, room.RID); err != nil {
			Ret500(c, SEC_game_spin_sqlbank, err)
			return
		}

		const sql2 = `UPDATE props SET wallet=wallet+? WHERE uid=? AND rid=?`
		if ret, err = session.Exec(sql2, totalwin-totalbet, props.UID, props.RID); err != nil {
			Ret500(c, SEC_game_spin_sqlupdate, err)
			return
		}

		return
	}); err != nil {
		return
	}

	// make changes to memory data
	room.mux.Lock()
	room.Bank += float64(totalbet - totalwin)
	room.mux.Unlock()

	props.Wallet += totalwin - totalbet

	og.game.SetGain(totalwin)
	og.game.Apply(scrn, &ws)

	// write spin result to log and get spin ID
	var sl = Spinlog{
		GID:    arg.GID,
		Gain:   totalwin,
		Wallet: props.Wallet,
	}
	var b []byte
	b, _ = json.Marshal(scrn)
	sl.Screen = util.B2S(b)
	b, _ = json.Marshal(og.game)
	sl.Game = util.B2S(b)
	b, _ = json.Marshal(ws.Wins)
	sl.Wins = util.B2S(b)
	if _, err = cfg.XormSpinlog.Insert(&sl); err != nil {
		log.Printf("can not write to spin log: %s", err.Error())
	}

	// prepare result
	ret.SID = sl.SID
	ret.Screen = scrn
	ret.Wins = ws.Wins
	ret.FS = fs
	ret.Gain = totalwin
	ret.Wallet = props.Wallet

	RetOk(c, ret)
}

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

	if err = c.Bind(&arg); err != nil {
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

	var room *Room
	if room, ok = Rooms.Get(og.RID); !ok {
		Ret500(c, SEC_game_doubleup_noroom, ErrNoRoom)
		return
	}

	var user *User
	if user, ok = Users.Get(og.UID); !ok {
		Ret500(c, SEC_game_doubleup_nouser, ErrNoUser)
		return
	}

	var props *Props
	if props, ok = user.props.Get(og.RID); !ok {
		Ret500(c, SEC_game_doubleup_noprops, ErrNoWallet)
		return
	}

	var gain = og.game.GetGain()
	if gain == 0 {
		Ret403(c, SEC_game_doubleup_nomoney, ErrNoMoney)
		return
	}

	room.mux.RLock()
	var bank = room.Bank
	var rtp = room.GainRTP
	room.mux.RUnlock()

	var multgain int // new multiplied gain
	if bank >= float64(gain*arg.Mult) {
		var r = rand.Float64()
		var side = 1 / float64(arg.Mult) * rtp / 100
		if r < side {
			multgain = gain * arg.Mult
		}
	}

	// write gain and total bet as transaction
	if _, err = cfg.XormStorage.Transaction(func(session *xorm.Session) (ret interface{}, err error) {
		defer func() {
			if err != nil {
				session.Rollback()
			}
		}()

		const sql1 = `UPDATE room SET bank=bank-? WHERE rid=?`
		if ret, err = session.Exec(sql1, multgain-gain, room.RID); err != nil {
			Ret500(c, SEC_game_spin_sqlbank, err)
			return
		}

		const sql2 = `UPDATE props SET wallet=wallet+? WHERE uid=? AND rid=?`
		if ret, err = session.Exec(sql2, multgain-gain, props.UID, props.RID); err != nil {
			Ret500(c, SEC_game_spin_sqlupdate, err)
			return
		}

		return
	}); err != nil {
		return
	}

	// make changes to memory data
	room.mux.Lock()
	room.Bank -= float64(multgain - gain)
	room.mux.Unlock()

	props.Wallet += multgain - gain

	og.game.SetGain(multgain)

	// write doubleup result to log and get spin ID
	var sl = Spinlog{
		GID:    arg.GID,
		Gain:   multgain,
		Wallet: props.Wallet,
	}
	var b []byte
	b, _ = json.Marshal(og.game)
	sl.Game = util.B2S(b)
	if _, err = cfg.XormSpinlog.Insert(&sl); err != nil {
		log.Printf("can not write to spin log: %s", err.Error())
	}

	// prepare result
	ret.Gain = multgain
	ret.Wallet = props.Wallet

	RetOk(c, ret)
}
