package spi

type Room struct {
	RID  uint64
	Bank float64 // users win/lost balance, in coins
	Fund float64 // jackpot fund, in coins
}

type User struct {
	UID     uint64
	Balance int // in coins
}
