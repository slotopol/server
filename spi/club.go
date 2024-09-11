package spi

import (
	"encoding/xml"

	"github.com/gin-gonic/gin"
	cfg "github.com/slotopol/server/config"
)

const (
	sqlclub = `UPDATE club SET bank=bank+?, fund=fund+?, lock=lock+?, utime=CURRENT_TIMESTAMP WHERE cid=?`
)

// Returns current club state.
func SpiClubIs(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid"`
		Name    string   `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty" form:"name"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr"`
		Name    string   `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_club_is_nobind, err)
		return
	}
	if arg.CID == 0 {
		Ret400(c, SEC_club_is_nouid, ErrNoCID)
		return
	}

	if arg.CID != 0 {
		var club *Club
		if club, ok = Clubs.Get(arg.CID); ok {
			ret.CID = club.CID
			ret.Name = club.Name
		}
	} else {
		for _, club := range Clubs.Items() {
			if club.Name == arg.Name {
				ret.CID = club.CID
				ret.Name = club.Name
				break
			}
		}
	}

	RetOk(c, ret)
}

// Returns current club state.
func SpiClubInfo(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		Name    string   `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
		Bank    float64  `json:"bank" yaml:"bank" xml:"bank"` // users win/lost balance, in coins
		Fund    float64  `json:"fund" yaml:"fund" xml:"fund"` // jackpot fund, in coins
		Lock    float64  `json:"lock" yaml:"lock" xml:"lock"` // not changed deposit within games

		JptRate float64 `json:"jptrate" yaml:"jptrate" xml:"jptrate"`
		MRTP    float64 `json:"mrtp" yaml:"mrtp" xml:"mrtp"` // master RTP
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_club_info_nobind, err)
		return
	}
	if arg.CID == 0 {
		Ret400(c, SEC_club_info_nouid, ErrNoCID)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(arg.CID); !ok {
		Ret500(c, SEC_club_info_noclub, ErrNoClub)
		return
	}

	var _, al = MustAdmin(c, arg.CID)
	if al&ALclub == 0 {
		Ret403(c, SEC_club_info_noaccess, ErrNoAccess)
		return
	}

	club.mux.Lock()
	ret.Name = club.Name
	ret.Bank = club.Bank
	ret.Fund = club.Fund
	ret.Lock = club.Lock
	ret.JptRate = club.JptRate
	ret.MRTP = club.MRTP
	club.mux.Unlock()

	RetOk(c, ret)
}

// Rename the club.
func SpiClubRename(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid"`
		Name    string   `json:"name" yaml:"name" xml:"name" form:"name" binding:"required"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_club_rename_nobind, err)
		return
	}
	if arg.CID == 0 {
		Ret400(c, SEC_club_rename_nouid, ErrNoCID)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(arg.CID); !ok {
		Ret500(c, SEC_club_rename_noclub, ErrNoClub)
		return
	}

	var _, al = MustAdmin(c, arg.CID)
	if al&ALclub == 0 {
		Ret403(c, SEC_club_rename_noaccess, ErrNoAccess)
		return
	}

	if _, err = cfg.XormStorage.ID(club.CID).Cols("name").Update(&Club{Name: arg.Name}); err != nil {
		Ret500(c, SEC_club_rename_update, err)
		return
	}

	club.mux.Lock()
	club.Name = arg.Name
	club.mux.Unlock()

	Ret204(c)
}

// Adding or withdrawing coins from club bank, fund and lock balances.
func SpiClubCashin(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid"`
		BankSum float64  `json:"banksum" yaml:"banksum" xml:"banksum" form:"banksum"`
		FundSum float64  `json:"fundsum" yaml:"fundsum" xml:"fundsum" form:"fundsum"`
		LockSum float64  `json:"locksum" yaml:"locksum" xml:"locksum" form:"locksum"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		BID     uint64   `json:"bid" yaml:"bid" xml:"bid,attr"`
		Bank    float64  `json:"bank" yaml:"bank" xml:"bank"`
		Fund    float64  `json:"fund" yaml:"fund" xml:"fund"`
		Lock    float64  `json:"lock" yaml:"lock" xml:"lock"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_club_cashin_nobind, err)
		return
	}
	if arg.CID == 0 {
		Ret400(c, SEC_club_cashin_nouid, ErrNoCID)
		return
	}
	if arg.BankSum == 0 && arg.FundSum == 0 && arg.LockSum == 0 {
		Ret400(c, SEC_club_cashin_nosum, ErrNoAddSum)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(arg.CID); !ok {
		Ret500(c, SEC_club_cashin_noclub, ErrNoClub)
		return
	}

	var _, al = MustAdmin(c, arg.CID)
	if al&ALclub == 0 {
		Ret403(c, SEC_club_cashin_noaccess, ErrNoAccess)
		return
	}

	club.mux.Lock()
	var bank = club.Bank + arg.BankSum
	var fund = club.Fund + arg.FundSum
	var lock = club.Lock + arg.LockSum
	club.mux.Unlock()

	if bank < 0 {
		Ret403(c, SEC_club_cashin_bankout, ErrBankOut)
		return
	}
	if fund < 0 {
		Ret403(c, SEC_club_cashin_fundout, ErrFundOut)
		return
	}
	if lock < 0 {
		Ret403(c, SEC_club_cashin_lockout, ErrLockOut)
		return
	}

	var rec = Banklog{
		Bank:    bank,
		Fund:    fund,
		Lock:    lock,
		BankSum: arg.BankSum,
		FundSum: arg.FundSum,
		LockSum: arg.LockSum,
	}
	if err = SafeTransaction(cfg.XormStorage, func(session *Session) (err error) {
		if _, err = session.Exec(sqlclub, arg.BankSum, arg.FundSum, arg.LockSum, club.CID); err != nil {
			Ret500(c, SEC_club_cashin_sqlbank, err)
			return
		}
		if _, err = session.Insert(&rec); err != nil {
			Ret500(c, SEC_club_cashin_sqllog, err)
			return
		}
		return
	}); err != nil {
		return
	}

	// make changes to memory data
	club.mux.Lock()
	club.Bank += arg.BankSum
	club.Fund += arg.FundSum
	club.Lock += arg.LockSum
	club.mux.Unlock()

	ret.BID = rec.ID
	ret.Bank = bank
	ret.Fund = fund
	ret.Lock = lock

	RetOk(c, ret)
}
