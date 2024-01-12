package spi

import "errors"

const (
	SECnull = iota

	// game/join
	SEC_game_join_nobind
	SEC_game_join_nodata
	SEC_game_join_noalias
	SEC_game_join_noreels

	// game/part
	SEC_game_part_nobind
	SEC_game_part_nodata
	SEC_game_part_notopened
)

var (
	ErrNoData    = errors.New("data does not provided or empty")
	ErrNoAliase  = errors.New("no game aliase")
	ErrNoReels   = errors.New("no reels for given game with selected percentage")
	ErrNotOpened = errors.New("game with given ID is not opened")
)
