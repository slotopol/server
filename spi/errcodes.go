package spi

import "errors"

const (
	SECnull = iota

	// POST game/join
	SEC_game_join_nobind
	SEC_game_join_norid
	SEC_game_join_nouid
	SEC_game_join_nodata
	SEC_game_join_noroom
	SEC_game_join_nouser
	SEC_game_join_noalias
	SEC_game_join_noreels
	SEC_game_join_insert

	// POST game/part
	SEC_game_part_nobind
	SEC_game_part_nogid
	SEC_game_part_notopened
	SEC_game_part_nouser

	// GET game/bet
	SEC_game_betget_nobind
	SEC_game_betget_nogid
	SEC_game_betget_notopened

	// PUT game/bet
	SEC_game_betset_nobind
	SEC_game_betset_nogid
	SEC_game_betset_nodata
	SEC_game_betset_notopened
	SEC_game_betset_badbet

	// GET game/sbl
	SEC_game_sblget_nobind
	SEC_game_sblget_nogid
	SEC_game_sblget_notopened

	// PUT game/sbl
	SEC_game_sblset_nobind
	SEC_game_sblset_nogid
	SEC_game_sblset_nodata
	SEC_game_sblset_notopened
	SEC_game_sblset_badlines

	// POST game/spin
	SEC_game_spin_nobind
	SEC_game_spin_nogid
	SEC_game_spin_notopened
	SEC_game_spin_noroom
	SEC_game_spin_nouser
	SEC_game_spin_nomoney
	SEC_game_spin_badbank
	SEC_game_spin_sqlbank
	SEC_game_spin_sqlupdate
	SEC_game_spin_sqlinsert

	// GET prop/wallet
	SEC_prop_walletget_nobind
	SEC_prop_walletget_norid
	SEC_prop_walletget_nouid
	SEC_prop_walletget_noroom
	SEC_prop_walletget_nouser

	// PUT prop/wallet
	SEC_prop_walletadd_nobind
	SEC_prop_walletadd_norid
	SEC_prop_walletadd_nouid
	SEC_prop_walletadd_noadd
	SEC_prop_walletadd_limit
	SEC_prop_walletadd_noroom
	SEC_prop_walletadd_nouser
	SEC_prop_walletadd_nomoney
	SEC_prop_walletadd_sqlupdate
	SEC_prop_walletadd_sqlinsert
	SEC_prop_walletadd_sqllog
)

var (
	ErrNoUID     = errors.New("user ID does not provided")
	ErrNoRID     = errors.New("room ID does not provided")
	ErrNoGID     = errors.New("game ID does not provided")
	ErrNoData    = errors.New("data does not provided or empty")
	ErrNoRoom    = errors.New("room with given ID does not found")
	ErrNoUser    = errors.New("user with given ID does not found")
	ErrNoMoney   = errors.New("not enough money")
	ErrZero      = errors.New("given value is zero")
	ErrTooBig    = errors.New("given value exceeds the limit")
	ErrNoAliase  = errors.New("no game alias")
	ErrNoReels   = errors.New("no reels for given game with selected percentage")
	ErrNotOpened = errors.New("game with given ID is not opened")
	ErrBadBank   = errors.New("can not generate spin with current bank balance")
)
