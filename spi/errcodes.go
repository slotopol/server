package spi

import "errors"

const (
	SECnull = iota

	// game/join
	SEC_game_join_nobind
	SEC_game_join_nodata
	SEC_game_join_nouser
	SEC_game_join_noalias
	SEC_game_join_noreels

	// game/part
	SEC_game_part_nobind
	SEC_game_part_nodata
	SEC_game_part_notopened
	SEC_game_part_nouser

	// game/bet
	SEC_game_bet_nobind
	SEC_game_bet_nodata
	SEC_game_bet_notopened
	SEC_game_bet_badbet

	// game/line
	SEC_game_line_nobind
	SEC_game_line_nodata
	SEC_game_line_notopened
	SEC_game_line_badlines
)

var (
	ErrNoData    = errors.New("data does not provided or empty")
	ErrNoUser    = errors.New("user with given ID does not found")
	ErrNoAliase  = errors.New("no game alias")
	ErrNoReels   = errors.New("no reels for given game with selected percentage")
	ErrNotOpened = errors.New("game with given ID is not opened")
)
