package util_test

import (
	"fmt"
	"os"
	"testing"
	"unsafe"

	"github.com/slotopol/server/util"
)

func TestS2B(t *testing.T) {
	var s = "some string"
	var ps = unsafe.Pointer(unsafe.StringData(s))
	/*var b = []byte(s)
	var pb = unsafe.Pointer(unsafe.SliceData(b))
	if ps == pb {
		t.Error("string pointer is equal to pointer on new allocated bytes slice")
	}*/
	var b = util.S2B(s)
	var pb = unsafe.Pointer(unsafe.SliceData(b))
	if ps != pb {
		t.Error("string pointer is not equal to pointer on same bytes slice")
	}
}

func TestB2S(t *testing.T) {
	var b = []byte("some string")
	var pb = unsafe.Pointer(unsafe.SliceData(b))
	var s = string(b)
	var ps = unsafe.Pointer(unsafe.StringData(s))
	if pb == ps {
		t.Error("bytes slice pointer is equal to pointer on new allocated string")
	}
	s = util.B2S(b)
	ps = unsafe.Pointer(unsafe.StringData(s))
	if pb != ps {
		t.Error("bytes slice pointer is not equal to pointer on same string")
	}
}

func ExampleToSlash() {
	fmt.Println(util.ToSlash("C:\\Windows\\Temp"))
	// Output: C:/Windows/Temp
}

func ExampleToLower() {
	fmt.Println(util.ToLower("C:\\Windows\\Temp"))
	// Output: c:\windows\temp
}

func ExampleToUpper() {
	fmt.Println(util.ToUpper("C:\\Windows\\Temp"))
	// Output: C:\WINDOWS\TEMP
}

func ExampleToKey() {
	fmt.Println(util.ToKey("C:\\Windows\\Temp"))
	// Output: c:/windows/temp
}

func ExampleToID() {
	fmt.Println(util.ToID("Joker Dolphin"))
	fmt.Println(util.ToID("Lucky Lady's Charm"))
	fmt.Println(util.ToID("BetSoft/2 Million B.C."))
	// Output:
	// jokerdolphin
	// luckyladyscharm
	// betsoft/2millionbc
}

func ExampleJoinPath() {
	fmt.Println(util.JoinPath("dir", "base.ext"))
	fmt.Println(util.JoinPath("dir/", "base.ext"))
	fmt.Println(util.JoinPath("dir", "/base.ext"))
	fmt.Println(util.JoinPath("dir/", "/base.ext"))
	// Output:
	// dir/base.ext
	// dir/base.ext
	// dir/base.ext
	// dir/base.ext
}

func ExampleJoinFilePath() {
	fmt.Println(util.JoinFilePath("dir/", "base.ext"))
	fmt.Println(util.JoinFilePath("dir", "/base.ext"))
	fmt.Println(util.JoinFilePath("dir/", "/base.ext"))
	// Output:
	// dir/base.ext
	// dir/base.ext
	// dir/base.ext
}

func ExamplePathName() {
	fmt.Println(util.PathName("C:\\Windows\\system.ini"))
	fmt.Println(util.PathName("/go/bin/wpkbuild_win_x64.exe"))
	fmt.Println(util.PathName("wpkbuild_win_x64.exe"))
	fmt.Println(util.PathName("/go/bin/wpkbuild_linux_x64"))
	fmt.Println(util.PathName("wpkbuild_linux_x64"))
	fmt.Printf("'%s'\n", util.PathName("/go/bin/"))
	// Output:
	// system
	// wpkbuild_win_x64
	// wpkbuild_win_x64
	// wpkbuild_linux_x64
	// wpkbuild_linux_x64
	// ''
}

func ExampleEnvfmt() {
	os.Setenv("VAR", "/go")
	// successful patterns
	fmt.Println(util.Envfmt("$VAR/bin/", nil))
	fmt.Println(util.Envfmt("${VAR}/bin/", nil))
	fmt.Println(util.Envfmt("%VAR%/bin/", nil))
	fmt.Println(util.Envfmt("/home$VAR", nil))
	fmt.Println(util.Envfmt("/home%VAR%", map[string]string{"VAR": "/any/path"}))
	fmt.Println(util.Envfmt("$VAR%VAR%${VAR}", nil))
	// patterns with unknown variable
	fmt.Println(util.Envfmt("$VYR/bin/", nil))
	fmt.Println(util.Envfmt("${VAR}/${_foo_}", nil))
	// patterns with errors
	fmt.Println(util.Envfmt("$VAR$/bin/", nil))
	fmt.Println(util.Envfmt("${VAR/bin/", nil))
	fmt.Println(util.Envfmt("%VAR/bin/", nil))
	fmt.Println(util.Envfmt("/home${VAR", nil))
	// Output:
	// /go/bin/
	// /go/bin/
	// /go/bin/
	// /home/go
	// /home/any/path
	// /go/go/go
	// $VYR/bin/
	// /go/${_foo_}
	// /go$/bin/
	// ${VAR/bin/
	// %VAR/bin/
	// /home${VAR
}
