package spi

import "errors"

const (
	SECnull = iota

	// POST game/join
	SEC_game_join_nobind
	SEC_game_join_nodata
	SEC_game_join_nouser
	SEC_game_join_noalias
	SEC_game_join_noreels

	// POST game/part
	SEC_game_part_nobind
	SEC_game_part_nodata
	SEC_game_part_notopened
	SEC_game_part_nouser

	// GET game/bet
	SEC_game_getbet_nobind
	SEC_game_getbet_nodata
	SEC_game_getbet_notopened

	// PUT game/bet
	SEC_game_setbet_nobind
	SEC_game_setbet_nodata
	SEC_game_setbet_notopened
	SEC_game_setbet_badbet

	// GET game/sbl
	SEC_game_getsbl_nobind
	SEC_game_getsbl_nodata
	SEC_game_getsbl_notopened

	// PUT game/sbl
	SEC_game_setsbl_nobind
	SEC_game_setsbl_nodata
	SEC_game_setsbl_notopened
	SEC_game_setsbl_badlines

	// POST game/spin
	SEC_game_spin_nobind
	SEC_game_spin_nodata
	SEC_game_spin_notopened
	SEC_game_spin_nouser
	SEC_game_spin_nomoney
	SEC_game_spin_noroom
	SEC_game_spin_sqlbank
	SEC_game_spin_sqlbalance
)

var (
	ErrNoData    = errors.New("data does not provided or empty")
	ErrNoRoom    = errors.New("room with given ID does not found")
	ErrNoUser    = errors.New("user with given ID does not found")
	ErrNoMoney   = errors.New("not enough money")
	ErrNoAliase  = errors.New("no game alias")
	ErrNoReels   = errors.New("no reels for given game with selected percentage")
	ErrNotOpened = errors.New("game with given ID is not opened")
)
