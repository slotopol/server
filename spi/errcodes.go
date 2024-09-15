package spi

import "errors"

const (
	SECnull = iota

	// authorization

	SEC_auth_absent
	SEC_auth_scheme
	SEC_basic_decode
	SEC_basic_nouser
	SEC_basic_deny
	SEC_token_nouser
	SEC_token_malform
	SEC_token_notsign
	SEC_token_badclaims
	SEC_token_expired
	SEC_token_notyet
	SEC_token_issuer
	SEC_token_error

	// 404

	SEC_nourl

	// 405

	SEC_nomethod

	// GET /signis

	SEC_signis_nobind
	SEC_signis_noemail

	// GET /sendcode

	SEC_sendcode_nobind
	SEC_sendcode_noemail
	SEC_sendcode_nouser
	SEC_sendcode_update
	SEC_sendcode_code

	// GET /activate

	SEC_activate_nobind
	SEC_activate_noemail
	SEC_activate_nouser
	SEC_activate_oldcode
	SEC_activate_badcode
	SEC_activate_update

	// POST /signup

	SEC_signup_nobind
	SEC_signup_smallsec
	SEC_signup_code
	SEC_signup_sql

	// POST /signin

	SEC_signin_nobind
	SEC_signin_noemail
	SEC_signin_nosecret
	SEC_signin_smallsec
	SEC_signin_nouser
	SEC_signin_activate
	SEC_signin_oldcode
	SEC_signin_badcode
	SEC_signin_denypass
	SEC_signin_sigtime
	SEC_signin_timeout
	SEC_signin_hs256
	SEC_signin_denyhash

	// POST /game/join

	SEC_game_join_nobind
	SEC_game_join_norid
	SEC_game_join_nouid
	SEC_game_join_noclub
	SEC_game_join_nouser
	SEC_game_join_noaccess
	SEC_game_join_noalias
	SEC_game_join_sql

	// POST /game/part

	SEC_game_part_nobind
	SEC_game_part_nogid
	SEC_game_part_notopened
	SEC_game_part_nouser
	SEC_game_part_noaccess
	SEC_game_part_sql

	// POST /slot/info

	SEC_game_info_nobind
	SEC_game_info_nogid
	SEC_game_info_notopened
	SEC_game_info_nouser
	SEC_game_info_noaccess
	SEC_game_info_noprops

	// POST /game/rtp/get

	SEC_game_rdget_nobind
	SEC_game_rdget_nogid
	SEC_game_rdget_notopened
	SEC_game_rdget_noclub
	SEC_game_rdget_nouser
	SEC_game_rdget_noaccess

	// POST /slot/bet/get

	SEC_slot_betget_nobind
	SEC_slot_betget_nogid
	SEC_slot_betget_notopened
	SEC_slot_betget_notslot
	SEC_slot_betget_noaccess

	// POST /slot/bet/set

	SEC_slot_betset_nobind
	SEC_slot_betset_nogid
	SEC_slot_betset_notopened
	SEC_slot_betset_notslot
	SEC_slot_betset_noaccess
	SEC_slot_betset_badbet

	// POST /slot/sel/get

	SEC_slot_selget_nobind
	SEC_slot_selget_nogid
	SEC_slot_selget_notopened
	SEC_slot_selget_notslot
	SEC_slot_selget_noaccess

	// POST /slot/sel/set

	SEC_slot_selset_nobind
	SEC_slot_selset_nogid
	SEC_slot_selset_notopened
	SEC_slot_selset_notslot
	SEC_slot_selset_noaccess
	SEC_slot_selset_badlines

	// POST /slot/spin

	SEC_slot_spin_nobind
	SEC_slot_spin_nogid
	SEC_slot_spin_notopened
	SEC_slot_spin_notslot
	SEC_slot_spin_noclub
	SEC_slot_spin_nouser
	SEC_slot_spin_noaccess
	SEC_slot_spin_noprops
	SEC_slot_spin_nomoney
	SEC_slot_spin_badbank
	SEC_slot_spin_sqlbank

	// POST /slot/doubleup

	SEC_slot_doubleup_nobind
	SEC_slot_doubleup_nogid
	SEC_slot_doubleup_nomult
	SEC_slot_doubleup_bigmult
	SEC_slot_doubleup_notopened
	SEC_slot_doubleup_notslot
	SEC_slot_doubleup_noclub
	SEC_slot_doubleup_nouser
	SEC_slot_doubleup_noaccess
	SEC_slot_doubleup_noprops
	SEC_slot_doubleup_nomoney
	SEC_slot_doubleup_sqlbank

	// POST /slot/collect

	SEC_slot_collect_nobind
	SEC_slot_collect_nogid
	SEC_slot_collect_notopened
	SEC_slot_collect_notslot
	SEC_slot_collect_noaccess
	SEC_slot_collect_denied

	// POST /keno/bet/get

	SEC_keno_betget_nobind
	SEC_keno_betget_nogid
	SEC_keno_betget_notopened
	SEC_keno_betget_notslot
	SEC_keno_betget_noaccess

	// POST /keno/bet/set

	SEC_keno_betset_nobind
	SEC_keno_betset_nogid
	SEC_keno_betset_notopened
	SEC_keno_betset_notslot
	SEC_keno_betset_noaccess
	SEC_keno_betset_badbet

	// POST /keno/sel/get

	SEC_keno_selget_nobind
	SEC_keno_selget_nogid
	SEC_keno_selget_notopened
	SEC_keno_selget_notslot
	SEC_keno_selget_noaccess

	// POST /keno/sel/set

	SEC_keno_selset_nobind
	SEC_keno_selset_nogid
	SEC_keno_selset_notopened
	SEC_keno_selset_notslot
	SEC_keno_selset_noaccess
	SEC_keno_selset_badlines

	// POST /prop/wallet/get

	SEC_prop_walletget_nobind
	SEC_prop_walletget_norid
	SEC_prop_walletget_nouid
	SEC_prop_walletget_noclub
	SEC_prop_walletget_nouser
	SEC_prop_walletget_noaccess

	// POST /prop/wallet/add

	SEC_prop_walletadd_nobind
	SEC_prop_walletadd_norid
	SEC_prop_walletadd_nouid
	SEC_prop_walletadd_noadd
	SEC_prop_walletadd_limit
	SEC_prop_walletadd_noclub
	SEC_prop_walletadd_nouser
	SEC_prop_walletadd_noaccess
	SEC_prop_walletadd_noprops
	SEC_prop_walletadd_nomoney
	SEC_prop_walletadd_sql

	// POST /prop/al/get

	SEC_prop_alget_nobind
	SEC_prop_alget_norid
	SEC_prop_alget_nouid
	SEC_prop_alget_noclub
	SEC_prop_alget_nouser
	SEC_prop_alget_noaccess

	// POST /prop/al/set

	SEC_prop_alset_nobind
	SEC_prop_alset_norid
	SEC_prop_alset_nouid
	SEC_prop_alset_noclub
	SEC_prop_alset_nouser
	SEC_prop_alset_noaccess
	SEC_prop_alset_noprops
	SEC_prop_alset_nolevel
	SEC_prop_alset_sql

	// POST /prop/rtp/get

	SEC_prop_rtpget_nobind
	SEC_prop_rtpget_norid
	SEC_prop_rtpget_nouid
	SEC_prop_rtpget_noclub
	SEC_prop_rtpget_nouser
	SEC_prop_rtpget_noaccess

	// POST /prop/rtp/set

	SEC_prop_rtpset_nobind
	SEC_prop_rtpset_norid
	SEC_prop_rtpset_nouid
	SEC_prop_rtpset_noclub
	SEC_prop_rtpset_nouser
	SEC_prop_rtpset_noaccess
	SEC_prop_rtpset_noprops
	SEC_prop_rtpset_sql

	// POST /user/rename

	SEC_user_rename_nobind
	SEC_user_rename_nouid
	SEC_user_rename_nouser
	SEC_user_rename_noaccess
	SEC_user_rename_update

	// POST /user/secret

	SEC_user_secret_nobind
	SEC_user_secret_nouid
	SEC_user_secret_smallsec
	SEC_user_secret_nouser
	SEC_user_secret_noaccess
	SEC_user_secret_nosecret
	SEC_user_secret_update

	// POST /user/delete

	SEC_user_delete_nobind
	SEC_user_delete_nouid
	SEC_user_delete_nouser
	SEC_user_delete_noaccess
	SEC_user_delete_nosecret
	SEC_user_delete_sqluser
	SEC_user_delete_sqllock
	SEC_user_delete_sqlprops
	SEC_user_delete_sqlgames

	// POST /club/is

	SEC_club_is_nobind
	SEC_club_is_nouid

	// POST /club/info

	SEC_club_info_nobind
	SEC_club_info_nouid
	SEC_club_info_noclub
	SEC_club_info_noaccess

	// POST /club/rename

	SEC_club_rename_nobind
	SEC_club_rename_nouid
	SEC_club_rename_noclub
	SEC_club_rename_noaccess
	SEC_club_rename_update

	// POST /club/cashin

	SEC_club_cashin_nobind
	SEC_club_cashin_nouid
	SEC_club_cashin_nosum
	SEC_club_cashin_noclub
	SEC_club_cashin_noaccess
	SEC_club_cashin_bankout
	SEC_club_cashin_fundout
	SEC_club_cashin_lockout
	SEC_club_cashin_sqlbank
	SEC_club_cashin_sqllog
)

var (
	Err404       = errors.New("page not found")
	Err405       = errors.New("method not allowed")
	ErrNoUID     = errors.New("user ID does not provided")
	ErrNoCID     = errors.New("club ID does not provided")
	ErrNoGID     = errors.New("game ID does not provided")
	ErrNoClub    = errors.New("club with given ID does not found")
	ErrNoUser    = errors.New("user with given ID does not found")
	ErrNoProps   = errors.New("properties for given user and club does not found")
	ErrNoAddSum  = errors.New("no sum to change balance of bank or fund or deposit")
	ErrNoMoney   = errors.New("not enough money on balance")
	ErrBankOut   = errors.New("not enough money at bank")
	ErrFundOut   = errors.New("not enough money at jackpot fund")
	ErrLockOut   = errors.New("not enough money at deposit")
	ErrNotSlot   = errors.New("specified GID refers to non-slot game")
	ErrNotKeno   = errors.New("specified GID refers to non-keno game")
	ErrNoAccess  = errors.New("no access rights for this feature")
	ErrNoLevel   = errors.New("admin have no privilege to modify specified access level to user")
	ErrNotConf   = errors.New("password confirmation does not pass")
	ErrNoMult    = errors.New("gamble multiplier not given")
	ErrBigMult   = errors.New("gamble multiplier too big")
	ErrZero      = errors.New("given value is zero")
	ErrTooBig    = errors.New("given value exceeds the limit")
	ErrNoAliase  = errors.New("no game alias")
	ErrNotOpened = errors.New("game with given ID is not opened")
	ErrBadBank   = errors.New("can not generate spin with current bank balance")
)
