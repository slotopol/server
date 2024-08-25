package spi

import (
	"encoding/xml"
	"net/http"

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

	var props, hasprops = user.props.Get(arg.CID)
	if !hasprops {
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
	if Cfg.WalletlogBufferSize > 1 {
		go BankBat[arg.CID].Add(cfg.XormStorage, arg.UID, admin.UID, props.Wallet+arg.Sum, arg.Sum, !hasprops)
	} else if err = BankBat[arg.CID].Add(cfg.XormStorage, arg.UID, admin.UID, props.Wallet+arg.Sum, arg.Sum, !hasprops); err != nil {
		Ret500(c, SEC_prop_walletadd_sql, err)
		return
	}

	// make changes to memory data
	props.Wallet += arg.Sum
	if !hasprops {
		user.InsertProps(props)
	}

	ret.Wallet = props.Wallet

	RetOk(c, ret)
}

// Returns master RTP for pointed user at pointed club.
// This RTP if it set have more priority then club RTP.
func SpiPropsRtpGet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr" form:"uid"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		MRTP    float64  `json:"mrtp" yaml:"mrtp" xml:"mrtp"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_prop_rtpget_nobind, err)
		return
	}
	if arg.CID == 0 {
		Ret400(c, SEC_prop_rtpget_norid, ErrNoCID)
		return
	}
	if arg.UID == 0 {
		Ret400(c, SEC_prop_rtpget_nouid, ErrNoUID)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(arg.CID); !ok {
		Ret404(c, SEC_prop_rtpget_noclub, ErrNoClub)
		return
	}
	_ = club

	var user *User
	if user, ok = Users.Get(arg.UID); !ok {
		Ret404(c, SEC_prop_rtpget_nouser, ErrNoUser)
		return
	}

	var admin, al = GetAdmin(c, arg.CID)
	if admin != user && al&ALgame == 0 {
		Ret403(c, SEC_prop_rtpget_noaccess, ErrNoAccess)
		return
	}

	ret.MRTP = user.GetRTP(arg.CID)

	RetOk(c, ret)
}

// Set personal master RTP for given user at given club.
func SpiPropsRtpSet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr" form:"uid"`
		MRTP    float64  `json:"mrtp" yaml:"mrtp" xml:"mrtp"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_prop_rtpset_nobind, err)
		return
	}
	if arg.CID == 0 {
		Ret400(c, SEC_prop_rtpset_norid, ErrNoCID)
		return
	}
	if arg.UID == 0 {
		Ret400(c, SEC_prop_rtpset_nouid, ErrNoUID)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(arg.CID); !ok {
		Ret404(c, SEC_prop_rtpset_noclub, ErrNoClub)
		return
	}
	_ = club

	var user *User
	if user, ok = Users.Get(arg.UID); !ok {
		Ret404(c, SEC_prop_rtpset_nouser, ErrNoUser)
		return
	}

	var props, hasprops = user.props.Get(arg.CID)
	if !hasprops {
		props = &Props{
			CID: arg.CID,
			UID: arg.UID,
		}
	}

	var admin, al = GetAdmin(c, arg.CID)
	if al&ALgame == 0 {
		Ret403(c, SEC_prop_rtpset_noaccess, ErrNoAccess)
		return
	}
	_ = admin

	// update wallet as transaction
	if Cfg.WalletlogBufferSize > 1 {
		go BankBat[arg.CID].MRTP(cfg.XormStorage, arg.UID, arg.MRTP, !hasprops)
	} else if err = BankBat[arg.CID].MRTP(cfg.XormStorage, arg.UID, arg.MRTP, !hasprops); err != nil {
		Ret500(c, SEC_prop_rtpset_sql, err)
		return
	}

	// make changes to memory data
	props.MRTP = arg.MRTP
	if !hasprops {
		user.InsertProps(props)
	}

	c.Status(http.StatusOK)
}
