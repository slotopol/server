// This code prints current time in ISO 8601 format.
package main

import (
	"fmt"
	"time"
)

// See https://tc39.es/ecma262/#sec-date-time-string-format
// time format acceptable for Date constructors.
const ISO8601 = "2006-01-02T15:04:05.999Z07:00"

func main() {
	fmt.Println(time.Now().UTC().Format(ISO8601))
}
