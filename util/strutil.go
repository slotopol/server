package util

import (
	"os"
	"unsafe"
)

// B2S converts bytes slice to string without memory allocation.
func B2S(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// S2B converts string to bytes slice without memory allocation.
func S2B(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// ToSlash brings filenames to true slashes
// without superfluous allocations if it possible.
func ToSlash(s string) string {
	var b = S2B(s)
	var bc = b
	var c bool
	for i, v := range b {
		if v == '\\' {
			if !c {
				bc, c = []byte(s), true
			}
			bc[i] = '/'
		}
	}
	return B2S(bc)
}

// ToLower is high performance function to bring filenames to lower case in ASCII
// without superfluous allocations if it possible.
func ToLower(s string) string {
	var b = S2B(s)
	var bc = b
	var c bool
	for i, v := range b {
		if v >= 'A' && v <= 'Z' {
			if !c {
				bc, c = []byte(s), true
			}
			bc[i] |= 0x20
		}
	}
	return B2S(bc)
}

// ToUpper is high performance function to bring filenames to upper case in ASCII
// without superfluous allocations if it possible.
func ToUpper(s string) string {
	var b = S2B(s)
	var bc = b
	var c bool
	for i, v := range b {
		if v >= 'a' && v <= 'z' {
			if !c {
				bc, c = []byte(s), true
			}
			bc[i] &= 0xdf
		}
	}
	return B2S(bc)
}

// ToKey is high performance function to bring filenames to lower case in ASCII
// and true slashes at once without superfluous allocations if it possible.
func ToKey(s string) string {
	var b = S2B(s)
	var bc = b
	var c bool
	for i, v := range b {
		if v >= 'A' && v <= 'Z' {
			if !c {
				bc, c = []byte(s), true
			}
			bc[i] |= 0x20
		} else if v == '\\' {
			if !c {
				bc, c = []byte(s), true
			}
			bc[i] = '/'
		}
	}
	return B2S(bc)
}

// ToID is high performance function to bring filenames to lower case identifier
// with only letters, digits and '_', without superfluous allocations if it possible.
func ToID(s string) string {
	var b = S2B(s)
	var bc = b
	var n int
	for _, v := range b {
		if VarChar[v] {
			n++
		}
	}
	var c bool
	if n != len(b) {
		bc = make([]byte, n)
		c = true
	}
	var i int
	for _, v := range b {
		if VarChar[v] {
			if v >= 'A' && v <= 'Z' {
				if !c {
					bc, c = []byte(s), true
				}
				bc[i] = v | 0x20
			} else if c {
				bc[i] = v
			}
			i++
		}
	}
	return B2S(bc)
}

// JoinPath performs fast join of two UNIX-like path chunks.
func JoinPath(dir, base string) string {
	if dir == "" || dir == "." {
		return base
	}
	if base == "" || base == "." {
		return dir
	}
	if dir[len(dir)-1] == '/' {
		if base[0] == '/' {
			return dir + base[1:]
		} else {
			return dir + base
		}
	}
	if base[0] == '/' {
		return dir + base
	}
	return dir + "/" + base
}

// OS-specific path separator string
const PathSeparator = string(os.PathSeparator)

// JoinFilePath performs fast join of two file path chunks.
// In some cases concatenates with OS-specific separator.
func JoinFilePath(dir, base string) string {
	if dir == "" || dir == "." {
		return base
	}
	if base == "" || base == "." {
		return dir
	}
	var isd = os.IsPathSeparator(dir[len(dir)-1])
	var isb = os.IsPathSeparator(base[0])
	if isd {
		if isb {
			return dir + base[1:]
		} else {
			return dir + base
		}
	}
	if isb {
		return dir + base
	}
	return dir + PathSeparator + base
}

// PathName returns name of file in given file path without extension.
func PathName(fpath string) string {
	var j = len(fpath)
	if j == 0 {
		return ""
	}
	var i = j - 1
	for {
		if fpath[i] == '\\' || fpath[i] == '/' {
			i++
			break
		}
		if fpath[i] == '.' {
			j = i
		}
		if i == 0 {
			break
		}
		i--
	}
	return fpath[i:j]
}

// VarCharFirst is table for fast check that ASCII code is acceptable first symbol of variable.
var VarCharFirst [256]bool = func() (a [256]bool) {
	a['_'] = true
	for c := 'A'; c <= 'Z'; c++ {
		a[c] = true
		a[c+32] = true
	}
	return
}()

// VarChar is table for fast check that ASCII code is acceptable symbol of variable.
var VarChar [256]bool = func() (a [256]bool) {
	a['_'] = true
	for c := 'A'; c <= 'Z'; c++ {
		a[c] = true
		a[c+32] = true
	}
	for c := '0'; c <= '9'; c++ {
		a[c] = true
	}
	return
}()

// Envfmt replaces environment variables entries in file path to there values.
// Environment variables must be followed by those 3 patterns: $VAR, ${VAR}, %VAR%.
// Environment variables are looked at first in 'envmap' if it given, and then by os-call.
// This function works by two string passes, without superfluous memory allocations.
func Envfmt(fpath string, envmap map[string]string) string {
	var (
		mode int          // pattern mode, 1 - $VAR, 2 - ${VAR}, 3 - %VAR%
		i    int          // iterator in incoming path
		j    int          // iterator in expanded path
		p1   int          // VAR start position
		p2   int          // VAR end position
		pc   int          // position for copy start
		pe   int          // position for copy end
		n    int          // number of entries
		fl   = len(fpath) // len of incoming path
		el   = len(fpath) // len of expanded path
	)

	// 1st passage, string length calculation
	var b = S2B(fpath)
	i = 0
	for i < fl {
		if mode == 0 {
			if b[i] == '$' {
				if i < fl-2 && VarCharFirst[b[i+1]] {
					mode, p1 = 1, i+1
					i += 2
				} else if i < fl-3 && b[i+1] == '{' && VarCharFirst[b[i+2]] {
					mode, p1 = 2, i+2
					i += 3
				} else {
					i++
				}
			} else if b[i] == '%' && i < fl-2 && VarCharFirst[b[i+1]] {
				mode, p1 = 3, i+1
				i += 2
			} else {
				i++
			}
			if mode > 0 { // starts some mode
				p2 = 0
			}
		} else { // mode > 0
			if VarChar[b[i]] {
				i++
				continue
			}
			switch mode {
			case 1:
				p2 = i
			case 2:
				if b[i] == '}' {
					p2 = i
					i++
				} else {
					mode, p1 = 0, 0 // error in pattern, skip VAR
				}
			case 3:
				if b[i] == '%' {
					p2 = i
					i++
				} else {
					mode, p1 = 0, 0 // error in pattern, skip VAR
				}
			}
			if p1 > 0 && p2 > 0 { // pattern ready
				var v = B2S(b[p1:p2])
				var env string
				var ok bool
				if envmap != nil {
					env, ok = envmap[v]
				}
				if !ok {
					env, ok = os.LookupEnv(v)
				}
				if ok {
					el += len(env) - len(v)
					switch mode {
					case 1:
						el -= 1 // $VAR
					case 2:
						el -= 3 // ${VAR}
					case 3:
						el -= 2 // %VAR%
					}
					n++
				}
				mode = 0
			}
		}
	}
	if mode == 1 { // special case: $VAR can be at the end of fpath
		p2 = i
		var v = b[p1:p2]
		if env, ok := os.LookupEnv(B2S(v)); ok {
			el += len(env) - len(v) - 1
			n++
		}
		mode = 0
	}

	if n == 0 { // there is no changes
		return fpath
	}

	// 2nd passage, fill out result
	var e = make([]byte, el)
	i = 0
	for i < fl {
		if mode == 0 {
			if b[i] == '$' {
				if i < fl-2 && VarCharFirst[b[i+1]] {
					pe = i
					mode, p1 = 1, i+1
					i += 2
				} else if i < fl-3 && b[i+1] == '{' && VarCharFirst[b[i+2]] {
					pe = i
					mode, p1 = 2, i+2
					i += 3
				} else {
					i++
				}
			} else if b[i] == '%' && i < fl-2 && VarCharFirst[b[i+1]] {
				pe = i
				mode, p1 = 3, i+1
				i += 2
			} else {
				i++
			}
			if mode > 0 { // starts some mode
				copy(e[j:], b[pc:pe])
				j += pe - pc
				pc = pe
				p2 = 0
			}
		} else { // mode > 0
			if VarChar[b[i]] {
				i++
				continue
			}
			switch mode {
			case 1:
				p2 = i
			case 2:
				if b[i] == '}' {
					p2 = i
					i++
				} else {
					mode, p1 = 0, 0 // error in pattern, skip VAR
				}
			case 3:
				if b[i] == '%' {
					p2 = i
					i++
				} else {
					mode, p1 = 0, 0 // error in pattern, skip VAR
				}
			}
			if p1 > 0 && p2 > 0 { // pattern ready
				var v = B2S(b[p1:p2])
				var env string
				var ok bool
				if envmap != nil {
					env, ok = envmap[v]
				}
				if !ok {
					env, ok = os.LookupEnv(v)
				}
				if ok {
					copy(e[j:], env)
					j += len(env)
					pc = i
				}
				mode = 0
			}
		}
	}
	if mode == 1 { // special case: $VAR can be at the end of fpath
		p2 = i
		var v = b[p1:p2]
		if env, ok := os.LookupEnv(B2S(v)); ok {
			copy(e[j:], env)
		}
		mode = 0
	} else { // copy remainder
		copy(e[j:], b[pc:])
	}
	return B2S(e)
}
