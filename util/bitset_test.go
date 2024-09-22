package util_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/slotopol/server/util"
)

func ExampleMakeBitset64() {
	fmt.Printf("bitset(1, 2, 3, 5, 14, 8) is 0b%016b", util.MakeBitset64(1, 2, 3, 5, 14, 8))
	// Output:
	// bitset(1, 2, 3, 5, 14, 8) is 0b0100000100101110
}

func ExampleMakeBitset128() {
	var bs = util.MakeBitset128(1, 2, 3, 5, 14, 8, 63, 64, 65, 125, 127)
	fmt.Printf("bitset(1, 2, 3, 5, 14, 8, 63, 64, 65, 125, 127) is 0x%016x%016x", bs[1], bs[0])
	// Output:
	// bitset(1, 2, 3, 5, 14, 8, 63, 64, 65, 125, 127) is 0xa000000000000003800000000000412e
}

func ExampleMakeBitNum64() {
	fmt.Printf("bitnum(7, 3) is %b", util.MakeBitNum64(7, 3))
	// Output:
	// bitnum(7, 3) is 1111111000
}

func ExampleMakeBitNum128() {
	var bs = util.MakeBitNum128(77, 3)
	fmt.Printf("bitnum(77, 3) is %b%b", bs[1], bs[0])
	// Output:
	// bitnum(77, 3) is 11111111111111111111111111111111111111111111111111111111111111111111111111111000
}

func TestBitset64_Num(t *testing.T) {
	var test = []struct {
		bs util.Bitset64
		c  int
	}{
		{0x3, 2},                   // test #1
		{0x4F28, 7},                // test #2
		{0xc000_0000_0000_0001, 3}, // test #3
	}
	for tn, v := range test {
		if c := v.bs.Num(); c != v.c {
			t.Errorf("test #%d, wrong counts for bitset %x, expected %d, get %d", tn+1, v.bs, v.c, c)
		}
	}
}

func TestBitset128_Num(t *testing.T) {
	var test = []struct {
		bs util.Bitset128
		c  int
	}{
		{ // test #1
			util.Bitset128{0x3, 0}, 2,
		},
		{ // test #2
			util.Bitset128{0x4F28, 0}, 7,
		},
		{ // test #3
			util.Bitset128{0xc000_0000_0000_0001, 0}, 3,
		},
		{ // test #4
			util.Bitset128{0x4F28, 0xAD16}, 15,
		},
		{ // test #5
			util.Bitset128{0xc000_0000_0000_0001, 0xa000_0000_0400_0001}, 7,
		},
	}
	for tn, v := range test {
		if c := v.bs.Num(); c != v.c {
			t.Errorf("test #%d, wrong counts for bitset %v, expected %d, get %d", tn+1, v.bs, v.c, c)
		}
	}
}

func TestBitset64_Next(t *testing.T) {
	var test = []struct {
		bs util.Bitset64
		a  []int
	}{
		{0x3, []int{0, 1}}, // test #1
		{0b0100_1111_0010_1000, []int{3, 5, 8, 9, 10, 11, 14}}, // test #2
		{0xc000_0000_0000_0001, []int{0, 62, 63}},              // test #3
	}
	for tn, v := range test {
		var i int
		for n := v.bs.Next(-1); n != -1; n = v.bs.Next(n) {
			if n != v.a[i] {
				t.Errorf("test #%d, expected %d number, get %d at position %d", tn+1, v.a[i], n, i)
			}
			i++
		}
		if i != len(v.a) {
			t.Errorf("test #%d, expected %d iterations, gets %d", tn+1, len(v.a), i)
		}
	}
}

func TestBitset128_Next(t *testing.T) {
	var test = []struct {
		bs util.Bitset128
		a  []int
	}{
		{ // test #1
			util.Bitset128{0x3, 0},
			[]int{0, 1},
		},
		{ // test #2
			util.Bitset128{0b0100_1111_0010_1000, 0},
			[]int{3, 5, 8, 9, 10, 11, 14},
		},
		{ // test #3
			util.Bitset128{0xc000_0000_0000_0001, 0},
			[]int{0, 62, 63},
		},
		{ // test #4
			util.Bitset128{0x4F28, 0xAD16},
			[]int{
				3, 5, 8, 9, 10, 11, 14,
				64 + 1, 64 + 2, 64 + 4, 64 + 8, 64 + 10, 64 + 11, 64 + 13, 64 + 15,
			},
		},
		{ // test #5
			util.Bitset128{0xc000_0000_0000_0001, 0xa000_0000_0400_0001},
			[]int{0, 62, 63, 64, 90, 125, 127},
		},
	}
	for tn, v := range test {
		var i int
		for n := v.bs.Next(-1); n != -1; n = v.bs.Next(n) {
			if n != v.a[i] {
				t.Errorf("test #%d, expected %d number, get %d at position %d", tn+1, v.a[i], n, i)
			}
			i++
		}
		if i != len(v.a) {
			t.Errorf("test #%d, expected %d iterations, gets %d", tn+1, len(v.a), i)
		}
	}
}

func TestBitset64_Bits(t *testing.T) {
	var test = []struct {
		bs util.Bitset64
		a  []int
	}{
		{0x3, []int{0, 1}}, // test #1
		{0b0100_1111_0010_1000, []int{3, 5, 8, 9, 10, 11, 14}}, // test #2
		{0xc000_0000_0000_0001, []int{0, 62, 63}},              // test #3
	}
	for tn, v := range test {
		var i int
		for n := range v.bs.Bits() {
			if n != v.a[i] {
				t.Errorf("test #%d, expected %d number, get %d at position %d", tn+1, v.a[i], n, i)
			}
			i++
		}
		if i != len(v.a) {
			t.Errorf("test #%d, expected %d iterations, gets %d", tn+1, len(v.a), i)
		}
	}
}

func TestBitset128_Bits(t *testing.T) {
	var test = []struct {
		bs util.Bitset128
		a  []int
	}{
		{ // test #1
			util.Bitset128{0x3, 0},
			[]int{0, 1},
		},
		{ // test #2
			util.Bitset128{0b0100_1111_0010_1000, 0},
			[]int{3, 5, 8, 9, 10, 11, 14},
		},
		{ // test #3
			util.Bitset128{0xc000_0000_0000_0001, 0},
			[]int{0, 62, 63},
		},
		{ // test #4
			util.Bitset128{0x4F28, 0xAD16},
			[]int{
				3, 5, 8, 9, 10, 11, 14,
				64 + 1, 64 + 2, 64 + 4, 64 + 8, 64 + 10, 64 + 11, 64 + 13, 64 + 15,
			},
		},
		{ // test #5
			util.Bitset128{0xc000_0000_0000_0001, 0xa000_0000_0400_0001},
			[]int{0, 62, 63, 64, 90, 125, 127},
		},
	}
	for tn, v := range test {
		var i int
		for n := range v.bs.Bits() {
			if n != v.a[i] {
				t.Errorf("test #%d, expected %d number, get %d at position %d", tn+1, v.a[i], n, i)
			}
			i++
		}
		if i != len(v.a) {
			t.Errorf("test #%d, expected %d iterations, gets %d", tn+1, len(v.a), i)
		}
	}
}

func TestBitset64_Is(t *testing.T) {
	var bs = util.Bitset64(0b1010_0101)
	if bs.Is(1) {
		t.Error("bit #1 should be detected as 0")
	}
	if !bs.Is(2) {
		t.Error("bit #2 should be detected as 1")
	}
	if bs.Is(6) {
		t.Error("bit #6 should be detected as 0")
	}
}

func TestBitset128_Is(t *testing.T) {
	var bs = util.Bitset128{0b1010_0101, 0b0101_1010}
	if bs.Is(1) {
		t.Error("bit #1 should be detected as 0")
	}
	if bs.Is(2 + 64) {
		t.Error("bit #66 should be detected as 0")
	}
	if !bs.Is(6 + 64) {
		t.Error("bit #70 should be detected as 1")
	}
}

func TestBitset64_Set(t *testing.T) {
	var bs = util.Bitset64(0b1010_0101)
	if bs.Set(1); !bs.Is(1) {
		t.Error("bit #1 should be detected as 1")
	}
	if bs.Set(2); !bs.Is(2) {
		t.Error("bit #2 should be detected as 1")
	}
}

func TestBitset128_Set(t *testing.T) {
	var bs = util.Bitset128{0b1010_0101, 0b0101_1010}
	if bs.Set(1); !bs.Is(1) {
		t.Error("bit #1 should be detected as 1")
	}
	if bs.Set(64 + 2); !bs.Is(64 + 2) {
		t.Error("bit #66 should be detected as 1")
	}
	if bs.Set(64 + 3); !bs.Is(64 + 3) {
		t.Error("bit #67 should be detected as 1")
	}
}

func TestBitset64_Res(t *testing.T) {
	var bs = util.Bitset64(0b1010_0101)
	if bs.Res(0); bs.Is(0) {
		t.Error("bit #0 should be detected as 0")
	}
	if bs.Res(4); bs.Is(4) {
		t.Error("bit #4 should be detected as 4")
	}
}

func TestBitset128_Res(t *testing.T) {
	var bs = util.Bitset128{0b1010_0101, 0b0101_1010}
	if bs.Res(0); bs.Is(0) {
		t.Error("bit #0 should be detected as 0")
	}
	if bs.Res(64 + 1); bs.Is(64 + 1) {
		t.Error("bit #65 should be detected as 0")
	}
	if bs.Res(64 + 2); bs.Is(64 + 2) {
		t.Error("bit #66 should be detected as 0")
	}
}

func TestBitset64_Toggle(t *testing.T) {
	var bs = util.Bitset64(0b1010_0101)
	if bs.Toggle(1); !bs.Is(1) {
		t.Error("bit #1 should be detected as 1")
	}
	if bs.Toggle(2); bs.Is(2) {
		t.Error("bit #2 should be detected as 0")
	}
}

func TestBitset128_Toggle(t *testing.T) {
	var bs = util.Bitset128{0b1010_0101, 0b0101_1010}
	if bs.Toggle(1); !bs.Is(1) {
		t.Error("bit #1 should be detected as 1")
	}
	if bs.Toggle(64 + 2); !bs.Is(64 + 2) {
		t.Error("bit #66 should be detected as 1")
	}
	if bs.Toggle(64 + 3); bs.Is(64 + 3) {
		t.Error("bit #67 should be detected as 0")
	}
}

func TestBitset128_LShift(t *testing.T) {
	var test = []struct {
		bs1, bs2 util.Bitset128
		shift    int
	}{
		{ // test #1
			util.Bitset128{0x0000_0000_0000_ffff, 0},
			util.Bitset128{0x0000_0000_001f_ffe0, 0},
			5,
		},
		{ // test #2
			util.Bitset128{0xfa5f_0000_0000_b00d, 0},
			util.Bitset128{0x5f00_0000_00b0_0d00, 0x00fa},
			8,
		},
		{ // test #3
			util.Bitset128{0xcccc_3333_5555_aaaa, 0xbbbb_dddd_6666_9999},
			util.Bitset128{0x5555_aaaa_0000_0000, 0x6666_9999_cccc_3333},
			32,
		},
	}
	for tn, v := range test {
		var bs = v.bs1
		bs.LShift(v.shift)
		if bs != v.bs2 {
			t.Errorf("test #%d, bitset does not equal to pattern", tn+1)
		}
	}
}

func TestBitset64_SetNum(t *testing.T) {
	var test = []struct {
		count, from int
		bs          util.Bitset64
	}{
		{ // test #1
			3, 0, 0b111,
		},
		{ // test #2
			7, 3, 0b1111111000,
		},
		{ // test #3
			60, 4, 0xffff_ffff_ffff_fff0,
		},
		{ // test #4
			60, 24, 0xffff_ffff_ff00_0000,
		},
	}
	for tn, v := range test {
		var bs util.Bitset64
		bs.SetNum(v.count, v.from)
		if bs != v.bs {
			t.Errorf("test #%d, bitset does not equal to pattern", tn+1)
		}
	}
}

func TestBitset128_SetNum(t *testing.T) {
	var test = []struct {
		count, from int
		bs          util.Bitset128
	}{
		{ // test #1
			3, 0, util.Bitset128{0b111, 0},
		},
		{ // test #2
			7, 3, util.Bitset128{0b1111111000, 0},
		},
		{ // test #3
			60, 4, util.Bitset128{0xffff_ffff_ffff_fff0, 0},
		},
		{ // test #4
			60, 24, util.Bitset128{0xffff_ffff_ff00_0000, 0x000f_ffff},
		},
		{ // test #5
			80, 40, util.Bitset128{0xffff_ff00_0000_0000, 0x00ff_ffff_ffff_ffff},
		},
		{ // test #6
			120, 24, util.Bitset128{0xffff_ffff_ff00_0000, 0xffff_ffff_ffff_ffff},
		},
	}
	for tn, v := range test {
		var bs util.Bitset128
		bs.SetNum(v.count, v.from)
		if bs != v.bs {
			t.Errorf("test #%d, bitset does not equal to pattern", tn+1)
		}
	}
}

var packtest64 = []struct {
	idx []int
	bs  util.Bitset64
}{
	{[]int{0, 1, 2}, 0b111},                           // test #1
	{[]int{3, 4, 5, 6, 7, 8, 9}, 0b1111111000},        // test #2
	{[]int{1, 2, 3, 5, 14, 8}, 0b0100_0001_0010_1110}, // test #3
}

var packtest128 = []struct {
	idx []int
	bs  util.Bitset128
}{
	{ // test #1
		[]int{0, 1, 2},
		util.Bitset128{0b111, 0},
	},
	{ // test #2
		[]int{3, 4, 5, 6, 7, 8, 9},
		util.Bitset128{0b1111111000, 0},
	},
	{ // test #3
		[]int{1, 2, 3, 5, 14, 8},
		util.Bitset128{0b0100_0001_0010_1110, 0},
	},
	{ // test #4
		[]int{1, 2, 3, 5, 14, 8, 63, 64, 65, 125, 127},
		util.Bitset128{0x800000000000412e, 0xa000_0000_0000_0003},
	},
}

func TestBitset64_Pack(t *testing.T) {
	for tn, v := range packtest64 {
		var bs util.Bitset64
		bs.Pack(v.idx)
		if bs != v.bs {
			t.Errorf("test #%d, bitset does not equal to pattern", tn+1)
		}
	}
}

func TestBitset128_Pack(t *testing.T) {
	for tn, v := range packtest128 {
		var bs util.Bitset128
		bs.Pack(v.idx)
		if bs != v.bs {
			t.Errorf("test #%d, bitset does not equal to pattern", tn+1)
		}
	}
}

func TestBitset64_Expand(t *testing.T) {
	for tn, v := range packtest64 {
		var idx = v.bs.Expand()
		var vidx = append([]int{}, v.idx...)
		slices.Sort(vidx)
		if !slices.Equal(idx, vidx) {
			t.Errorf("test #%d, bitset does not equal to pattern", tn+1)
		}
	}
}

func TestBitset128_Expand(t *testing.T) {
	for tn, v := range packtest128 {
		var idx = v.bs.Expand()
		var vidx = append([]int{}, v.idx...)
		slices.Sort(vidx)
		if !slices.Equal(idx, vidx) {
			t.Errorf("test #%d, bitset does not equal to pattern", tn+1)
		}
	}
}

func TestBitset64_And(t *testing.T) {
	var op1 = util.Bitset64(0b1010_0101)
	var op2 = util.Bitset64(0b1001_0110)
	var res = util.Bitset64(0b1000_0100)
	if op1.And(op2); op1 != res {
		t.Error("bitset result does not equal to pattern")
	}
}

func TestBitset128_And(t *testing.T) {
	var op1 = util.Bitset128{0b1010_0101, 0b0101_1010}
	var op2 = util.Bitset128{0b1001_0110, 0b1001_0110}
	var res = util.Bitset128{0b1000_0100, 0b0001_0010}
	if op1.And(op2); op1 != res {
		t.Error("bitset result does not equal to pattern")
	}
}

func TestBitset64_Or(t *testing.T) {
	var op1 = util.Bitset64(0b1010_0101)
	var op2 = util.Bitset64(0b1001_0110)
	var res = util.Bitset64(0b1011_0111)
	if op1.Or(op2); op1 != res {
		t.Error("bitset result does not equal to pattern")
	}
}

func TestBitset128_Or(t *testing.T) {
	var op1 = util.Bitset128{0b1010_0101, 0b0101_1010}
	var op2 = util.Bitset128{0b1001_0110, 0b1001_0110}
	var res = util.Bitset128{0b1011_0111, 0b1101_1110}
	if op1.Or(op2); op1 != res {
		t.Error("bitset result does not equal to pattern")
	}
}

func TestBitset64_AndNot(t *testing.T) {
	var op1 = util.Bitset64(0b1010_0101)
	var op2 = util.Bitset64(0b1001_0110)
	var res = util.Bitset64(0b0010_0001)
	if op1.AndNot(op2); op1 != res {
		t.Error("bitset result does not equal to pattern")
	}
}

func TestBitset128_AndNot(t *testing.T) {
	var op1 = util.Bitset128{0b1010_0101, 0b0101_1010}
	var op2 = util.Bitset128{0b1001_0110, 0b1001_0110}
	var res = util.Bitset128{0b0010_0001, 0b0100_1000}
	if op1.AndNot(op2); op1 != res {
		t.Error("bitset result does not equal to pattern")
	}
}

func TestBitset64_Xor(t *testing.T) {
	var op1 = util.Bitset64(0b1010_0101)
	var op2 = util.Bitset64(0b1001_0110)
	var res = util.Bitset64(0b0011_0011)
	if op1.Xor(op2); op1 != res {
		t.Error("bitset result does not equal to pattern")
	}
}

func TestBitset128_Xor(t *testing.T) {
	var op1 = util.Bitset128{0b1010_0101, 0b0101_1010}
	var op2 = util.Bitset128{0b1001_0110, 0b1001_0110}
	var res = util.Bitset128{0b0011_0011, 0b1100_1100}
	if op1.Xor(op2); op1 != res {
		t.Error("bitset result does not equal to pattern")
	}
}

func TestBitset64_IsZero(t *testing.T) {
	var test = []struct {
		bs util.Bitset64
		is bool
	}{
		{0x3, false},    // test #1
		{0x0F28, false}, // test #2
		{0, true},       // test #3
	}
	for tn, v := range test {
		if is := v.bs.IsZero(); is != v.is {
			t.Errorf("test #%d, wrong zero detection, expected %t, get %t", tn+1, v.is, is)
		}
	}
}

func TestBitset128_IsZero(t *testing.T) {
	var test = []struct {
		bs util.Bitset128
		is bool
	}{
		{util.Bitset128{0x3, 0}, false},         // test #1
		{util.Bitset128{0x4F28, 0}, false},      // test #2
		{util.Bitset128{0, 0}, true},            // test #3
		{util.Bitset128{0x4F28, 0xAD16}, false}, // test #4
		{util.Bitset128{0, 0xAD16}, false},      // test #5
	}
	for tn, v := range test {
		if is := v.bs.IsZero(); is != v.is {
			t.Errorf("test #%d, wrong zero detection, expected %t, get %t", tn+1, v.is, is)
		}
	}
}
