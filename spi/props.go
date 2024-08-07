package spi

import (
	"encoding/xml"

	"github.com/gin-gonic/gin"
	cfg "github.com/slotopol/server/config"
)

// Returns balance at wallet for pointed user at pointed club.
func SpiPropsWalletGet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr" form:"uid"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		Wallet  float64  `json:"wallet" yaml:"wallet" xml:"wallet"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_prop_walletget_nobind, err)
		return
	}
	if arg.CID == 0 {
		Ret400(c, SEC_prop_walletget_norid, ErrNoCID)
		return
	}
	if arg.UID == 0 {
		Ret400(c, SEC_prop_walletget_nouid, ErrNoUID)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(arg.CID); !ok {
		Ret404(c, SEC_prop_walletget_noclub, ErrNoClub)
		return
	}
	_ = club

	var user *User
	if user, ok = Users.Get(arg.UID); !ok {
		Ret404(c, SEC_prop_walletget_nouser, ErrNoUser)
		return
	}

	var admin, al = GetAdmin(c, arg.CID)
	if admin != user && al&ALuser == 0 {
		Ret403(c, SEC_prop_walletget_noaccess, ErrNoAccess)
		return
	}

	ret.Wallet = user.GetWallet(arg.CID)

	RetOk(c, ret)
}

// Adds some coins to user wallet. Sum can be < 0 to remove some coins.
func SpiPropsWalletAdd(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr" form:"uid"`
		Sum     float64  `json:"sum" yaml:"sum" xml:"sum"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		Wallet  float64  `json:"wallet" yaml:"wallet" xml:"wallet"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_prop_walletadd_nobind, err)
		return
	}
	if arg.CID == 0 {
		Ret400(c, SEC_prop_walletadd_norid, ErrNoCID)
		return
	}
	if arg.UID == 0 {
		Ret400(c, SEC_prop_walletadd_nouid, ErrNoUID)
		return
	}
	if arg.Sum == 0 {
		Ret400(c, SEC_prop_walletadd_noadd, ErrZero)
		return
	}
	if arg.Sum > cfg.Cfg.AdjunctLimit || arg.Sum < -cfg.Cfg.AdjunctLimit {
		Ret400(c, SEC_prop_walletadd_limit, ErrTooBig)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(arg.CID); !ok {
		Ret404(c, SEC_prop_walletadd_noclub, ErrNoClub)
		return
	}
	_ = club

	var user *User
	if user, ok = Users.Get(arg.UID); !ok {
		Ret404(c, SEC_prop_walletadd_nouser, ErrNoUser)
		return
	}

	var props *Props
	var hasprops bool
	if props, hasprops = user.props.Get(arg.CID); !hasprops {
		props = &Props{
			CID: arg.CID,
			UID: arg.UID,
		}
	}
	if props.Wallet+arg.Sum < 0 {
		Ret403(c, SEC_prop_walletadd_nomoney, ErrNoMoney)
		return
	}

	var admin, al = GetAdmin(c, arg.CID)
	if al&ALuser == 0 {
		Ret403(c, SEC_prop_walletadd_noaccess, ErrNoAccess)
		return
	}

	// update wallet as transaction
	if err = SafeTransaction(cfg.XormStorage, func(session *Session) (err error) {
		var rec = Walletlog{
			CID:    arg.CID,
			UID:    arg.UID,
			AdmID:  admin.UID,
			Wallet: props.Wallet + arg.Sum,
			Sum:    arg.Sum,
		}

		if hasprops {
			const sql = `UPDATE props SET wallet=wallet+? WHERE uid=? AND cid=?`
			if _, err = session.Exec(sql, arg.Sum, props.UID, props.CID); err != nil {
				Ret500(c, SEC_prop_walletadd_sqlupdate, err)
				return
			}
		} else {
			props.Wallet = arg.Sum
			if _, err = session.Insert(props); err != nil {
				Ret500(c, SEC_prop_walletadd_sqlinsert, err)
				return
			}
		}

		if _, err = session.Insert(&rec); err != nil {
			Ret500(c, SEC_prop_walletadd_sqllog, err)
			return
		}

		return
	}); err != nil {
		return
	}

	// make changes to memory data
	if hasprops {
		props.Wallet += arg.Sum
	} else {
		user.InsertProps(props)
	}

	ret.Wallet = props.Wallet

	RetOk(c, ret)
}
