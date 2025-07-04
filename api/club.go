package api

import (
	"encoding/xml"

	"github.com/gin-gonic/gin"
	cfg "github.com/slotopol/server/config"
)

const (
	sqlclub = "UPDATE club SET `bank`=`bank`+?, `fund`=`fund`+?, `lock`=`lock`+?, `utime`=CURRENT_TIMESTAMP WHERE `cid`=?"
)

func ApiClubList(c *gin.Context) {
	type item struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"club"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr"`
		Name    string   `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		List    []item   `json:"list" yaml:"list" xml:"list>club" form:"list"`
	}

	ret.List = make([]item, 0, Clubs.Len())
	for cid, club := range Clubs.Items() {
		ret.List = append(ret.List, item{CID: cid, Name: club.Name()})
	}

	RetOk(c, ret)
}

// Returns current club state.
func ApiClubIs(c *gin.Context) {
	var err error
	var ok bool
	type item struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"user"`
		CID     uint64   `json:"cid,omitempty" yaml:"cid,omitempty" xml:"cid,attr,omitempty"`
		Name    string   `json:"name,omitempty" yaml:"name,omitempty" xml:"name,attr,omitempty"`
	}
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		List    []item   `json:"list" yaml:"list" xml:"list>user" form:"list" binding:"required"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		List    []item   `json:"list" yaml:"list" xml:"list>user"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_club_is_nobind, err)
		return
	}

	ret.List = make([]item, len(arg.List))
	for i, ai := range arg.List {
		var ri item
		if ai.CID != 0 {
			var club *Club
			if club, ok = Clubs.Get(ai.CID); ok {
				ri.CID = club.CID()
				ri.Name = club.Name()
			}
		} else {
			for _, club := range Clubs.Items() {
				if name := club.Name(); name == ai.Name {
					ri.CID = club.CID()
					ri.Name = name
					break
				}
			}
		}
		ret.List[i] = ri
	}

	RetOk(c, ret)
}

// Returns current club state.
func ApiClubInfo(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid" binding:"required"`
	}
	var ret struct {
		XMLName  xml.Name `json:"-" yaml:"-" xml:"ret"`
		ClubData `yaml:",inline"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_club_info_nobind, err)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(arg.CID); !ok {
		Ret500(c, AEC_club_info_noclub, ErrNoClub)
		return
	}

	var _, al = MustAdmin(c, arg.CID)
	if al&ALmaster == 0 {
		Ret403(c, AEC_club_info_noaccess, ErrNoAccess)
		return
	}

	ret.ClubData = club.Get()

	RetOk(c, ret)
}

// Returns full jackpot value.
func ApiClubJpfund(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid" binding:"required"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		JpFund  float64  `json:"jpfund" yaml:"jpfund" xml:"jpfund"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_club_jpfund_nobind, err)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(arg.CID); !ok {
		Ret500(c, AEC_club_jpfund_noclub, ErrNoClub)
		return
	}

	ret.JpFund = club.Fund()

	RetOk(c, ret)
}

// Rename the club.
func ApiClubRename(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid" binding:"required"`
		Name    string   `json:"name" yaml:"name" xml:"name" form:"name" binding:"required"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_club_rename_nobind, err)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(arg.CID); !ok {
		Ret500(c, AEC_club_rename_noclub, ErrNoClub)
		return
	}

	var _, al = MustAdmin(c, arg.CID)
	if al&ALmaster == 0 {
		Ret403(c, AEC_club_rename_noaccess, ErrNoAccess)
		return
	}

	if _, err = cfg.XormStorage.ID(club.CID()).Cols("name").Update(&ClubData{Name: arg.Name}); err != nil {
		Ret500(c, AEC_club_rename_update, err)
		return
	}

	club.SetName(arg.Name)

	Ret204(c)
}

// Adding or withdrawing coins from club bank, fund and lock balances.
func ApiClubCashin(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid" binding:"required"`
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
		Ret400(c, AEC_club_cashin_nobind, err)
		return
	}
	if arg.BankSum == 0 && arg.FundSum == 0 && arg.LockSum == 0 {
		Ret400(c, AEC_club_cashin_nosum, ErrNoAddSum)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(arg.CID); !ok {
		Ret500(c, AEC_club_cashin_noclub, ErrNoClub)
		return
	}

	var _, al = MustAdmin(c, arg.CID)
	if al&ALmaster == 0 {
		Ret403(c, AEC_club_cashin_noaccess, ErrNoAccess)
		return
	}

	var bank, fund, lock = club.GetCash()
	bank += arg.BankSum
	fund += arg.FundSum
	lock += arg.LockSum

	if bank < 0 {
		Ret403(c, AEC_club_cashin_bankout, ErrBankOut)
		return
	}
	if fund < 0 {
		Ret403(c, AEC_club_cashin_fundout, ErrFundOut)
		return
	}
	if lock < 0 {
		Ret403(c, AEC_club_cashin_lockout, ErrLockOut)
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
		if _, err = session.Exec(sqlclub, arg.BankSum, arg.FundSum, arg.LockSum, club.CID()); err != nil {
			Ret500(c, AEC_club_cashin_sqlbank, err)
			return
		}
		if _, err = session.Insert(&rec); err != nil {
			Ret500(c, AEC_club_cashin_sqllog, err)
			return
		}
		return
	}); err != nil {
		return
	}

	// make changes to memory data
	club.AddCash(arg.BankSum, arg.FundSum, arg.LockSum)

	ret.BID = rec.ID
	ret.Bank = bank
	ret.Fund = fund
	ret.Lock = lock

	RetOk(c, ret)
}
