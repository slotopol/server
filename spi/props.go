package spi

import (
	"encoding/xml"

	"github.com/gin-gonic/gin"
	cfg "github.com/slotopol/server/config"
	"xorm.io/xorm"
)

func SpiPropsWalletGet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		RID     uint64   `json:"rid" yaml:"rid" xml:"rid,attr" form:"rid"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr" form:"uid"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		Wallet  int      `json:"wallet" yaml:"wallet" xml:"wallet"`
	}

	if err = c.Bind(&arg); err != nil {
		Ret400(c, SEC_prop_walletget_nobind, err)
		return
	}
	if arg.RID == 0 {
		Ret400(c, SEC_prop_walletget_norid, ErrNoRID)
		return
	}
	if arg.UID == 0 {
		Ret400(c, SEC_prop_walletget_nouid, ErrNoUID)
		return
	}

	var room *Room
	if room, ok = Rooms.Get(arg.RID); !ok {
		Ret404(c, SEC_prop_walletget_noroom, ErrNoRoom)
		return
	}
	_ = room

	var user *User
	if user, ok = Users.Get(arg.UID); !ok {
		Ret404(c, SEC_prop_walletget_nouser, ErrNoUser)
		return
	}

	ret.Wallet = user.GetWallet(arg.RID)

	RetOk(c, ret)
}

func SpiPropsWalletAdd(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		RID     uint64   `json:"rid" yaml:"rid" xml:"rid,attr" form:"rid"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr" form:"uid"`
		Addend  int      `json:"addend" yaml:"addend" xml:"addend"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		Wallet  int      `json:"wallet" yaml:"wallet" xml:"wallet"`
	}

	if err = c.Bind(&arg); err != nil {
		Ret400(c, SEC_prop_walletadd_nobind, err)
		return
	}
	if arg.RID == 0 {
		Ret400(c, SEC_prop_walletadd_norid, ErrNoRID)
		return
	}
	if arg.UID == 0 {
		Ret400(c, SEC_prop_walletadd_nouid, ErrNoUID)
		return
	}
	if arg.Addend == 0 {
		Ret400(c, SEC_prop_walletadd_noadd, ErrZero)
		return
	}
	if arg.Addend > cfg.Cfg.AdjunctLimit {
		Ret400(c, SEC_prop_walletadd_limit, ErrTooBig)
		return
	}

	var room *Room
	if room, ok = Rooms.Get(arg.RID); !ok {
		Ret404(c, SEC_prop_walletadd_noroom, ErrNoRoom)
		return
	}
	_ = room

	var user *User
	if user, ok = Users.Get(arg.UID); !ok {
		Ret404(c, SEC_prop_walletadd_nouser, ErrNoUser)
		return
	}

	var props *Props
	var hasprops bool
	if props, hasprops = user.props.Get(arg.RID); !hasprops {
		props = &Props{
			RID: arg.RID,
			UID: arg.UID,
		}
	}
	if props.Wallet+arg.Addend < 0 {
		Ret403(c, SEC_prop_walletadd_nomoney, ErrNoMoney)
		return
	}

	// update wallet as transaction
	if _, err = cfg.XormStorage.Transaction(func(session *xorm.Session) (ret interface{}, err error) {
		defer func() {
			if err != nil {
				session.Rollback()
			}
		}()

		var wl = Walletlog{
			RID:    arg.RID,
			UID:    arg.UID,
			AdmID:  arg.UID,
			Wallet: props.Wallet + arg.Addend,
			Addend: arg.Addend,
		}

		if hasprops {
			const sql = `UPDATE props SET wallet=wallet+? WHERE uid=? AND rid=?`
			if ret, err = session.Exec(sql, arg.Addend, props.UID, props.RID); err != nil {
				Ret500(c, SEC_prop_walletadd_sqlupdate, err)
				return
			}
		} else {
			props.Wallet = arg.Addend
			if _, err = session.Insert(props); err != nil {
				Ret500(c, SEC_prop_walletadd_sqlinsert, err)
				return
			}
		}

		if _, err = session.Insert(&wl); err != nil {
			Ret500(c, SEC_prop_walletadd_sqllog, err)
			return
		}

		return
	}); err != nil {
		return
	}

	// make changes to memory data
	if hasprops {
		props.Wallet += arg.Addend
	} else {
		user.props.Set(props.RID, props)
	}

	ret.Wallet = props.Wallet

	RetOk(c, ret)
}
