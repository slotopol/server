// This code prints current time in ISO 8601 format.
package main

import (
	"fmt"
	"time"
)

// See https://tc39.es/ecma262/#sec-date-time-string-format
// time format acceptable for Date constructors.
func main() {
	fmt.Println(time.Now().UTC().Format(time.RFC3339Nano))
}
