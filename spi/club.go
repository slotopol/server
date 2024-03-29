package spi

import (
	"encoding/xml"
	"net/http"

	"github.com/gin-gonic/gin"
	cfg "github.com/slotopol/server/config"
)

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

	var _, al = GetAdmin(c, arg.CID)
	if al&ALclub == 0 {
		Ret403(c, SEC_club_rename_noaccess, ErrNoAccess)
		return
	}

	if _, err = cfg.XormStorage.Cols("name").Update(&Club{CID: arg.CID, Name: arg.Name}); err != nil {
		Ret500(c, SEC_club_rename_update, err)
		return
	}
	club.Name = arg.Name

	c.Status(http.StatusOK)
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

	var _, al = GetAdmin(c, arg.CID)
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
	if _, err = cfg.XormStorage.Transaction(func(session *Session) (_ interface{}, err error) {
		defer func() {
			if err != nil {
				session.Rollback()
			}
		}()

		const sql1 = `UPDATE club SET bank=bank+?, fund=fund+?, lock=lock+? WHERE cid=?`
		if _, err = session.Exec(sql1, arg.BankSum, arg.FundSum, arg.LockSum, club.CID); err != nil {
			Ret500(c, SEC_game_cashin_sqlbank, err)
			return
		}

		if _, err = session.Insert(&rec); err != nil {
			Ret500(c, SEC_game_cashin_sqllog, err)
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
