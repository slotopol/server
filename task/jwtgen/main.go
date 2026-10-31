package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	flag "github.com/spf13/pflag"
)

const jwtIssuer = "slotopol"

type Claims struct {
	jwt.RegisteredClaims
	SID uint64 `json:"sid,omitempty"`
}

var (
	key string
	sid uint64
	sub string
	nbf time.Time
	ttl time.Duration
)

func main() {
	var f = flag.CommandLine
	f.StringVarP(&key, "key", "k", "skJgM4NsbP3fs4k7vh0gfdkgGl8dJTszdLxZ1sQ9ksFnxbgvw2RsGH8xxddUV479", "key to sign HS-256 JWT-token")
	f.Uint64VarP(&sid, "sid", "s", 0, "session identifier")
	f.StringVarP(&sub, "user", "u", "3", "user identifier or user email")
	f.TimeVarP(&nbf, "nbf", "n", time.Now(), []string{time.RFC3339, time.DateTime, time.TimeOnly}, "'not before' time value to set to token")
	f.DurationVarP(&ttl, "ttl", "t", 30*time.Minute, "token time-to-live")
	f.SortFlags = false
	flag.Parse()

	var claims = &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			NotBefore: jwt.NewNumericDate(nbf),
			ExpiresAt: jwt.NewNumericDate(nbf.Add(ttl)),
			Subject:   sub,
		},
		SID: sid,
	}
	var token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(key))
	if err != nil {
		panic(err)
	}
	fmt.Println(token)
}
