package textwidth

import (
	"testing"

	"github.com/mattn/go-runewidth"
)

// test cases copied from https://github.com/mattn/go-runewidth/raw/master/runewidth_test.go

var stringwidthtests = []struct {
	in    string
	out   int
	eaout int
}{
	{"■㈱の世界①", 10, 12},
	{"スター☆", 7, 8},
	{"つのだ☆HIRO", 11, 12},
}

func BenchmarkStringWidth(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WidthOfString(stringwidthtests[i%len(stringwidthtests)].in)
	}
}

func BenchmarkStringWidthOriginal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runewidth.StringWidth(stringwidthtests[i%len(stringwidthtests)].in)
	}
}

func TestStringWidth(t *testing.T) {
	for _, tt := range stringwidthtests {
		if out := WidthOfString(tt.in); out != tt.out {
			t.Errorf("WidthOfString(%q) = %d, want %d", tt.in, out, tt.out)
		}
	}
	//c := runewidth.NewCondition()
	//c.EastAsianWidth = false
	//for _, tt := range stringwidthtests {
	//	if out := c.StringWidth(tt.in); out != tt.out {
	//		t.Errorf("StringWidth(%q) = %d, want %d", tt.in, out, tt.out)
	//	}
	//}
	//c.EastAsianWidth = true
	//for _, tt := range stringwidthtests {
	//	if out := c.StringWidth(tt.in); out != tt.eaout {
	//		t.Errorf("StringWidth(%q) = %d, want %d (EA)", tt.in, out, tt.eaout)
	//	}
	//}
}

var slicewidthtests = []struct {
	in    []byte
	out   int
	eaout int
}{
	{[]byte("■㈱の世界①"), 10, 12},
	{[]byte("スター☆"), 7, 8},
	{[]byte("つのだ☆HIRO"), 11, 12},
}

func TestSliceWidth(t *testing.T) {
	for _, tt := range slicewidthtests {
		if out := Width(tt.in); out != tt.out {
			t.Errorf("Width(%q) = %d, want %d", tt.in, out, tt.out)
		}
	}
}

var runewidthtests = []struct {
	in     rune
	out    int
	eaout  int
	nseout int
}{
	{'世', 2, 2, 2},
	{'界', 2, 2, 2},
	{'ｾ', 1, 1, 1},
	{'ｶ', 1, 1, 1},
	{'ｲ', 1, 1, 1},
	{'☆', 1, 2, 2}, // double width in ambiguous
	{'☺', 1, 1, 2},
	{'☻', 1, 1, 2},
	{'♥', 1, 2, 2},
	{'♦', 1, 1, 2},
	{'♣', 1, 2, 2},
	{'♠', 1, 2, 2},
	{'♂', 1, 2, 2},
	{'♀', 1, 2, 2},
	{'♪', 1, 2, 2},
	{'♫', 1, 1, 2},
	{'☼', 1, 1, 2},
	{'↕', 1, 2, 2},
	{'‼', 1, 1, 2},
	{'↔', 1, 2, 2},
	{'\x00', 0, 0, 0},
	{'\x01', 0, 0, 0},
	{'\u0300', 0, 0, 0},
	{'\u2028', 0, 0, 0},
	{'\u2029', 0, 0, 0},
	{'a', 1, 1, 1}, // ASCII classified as "na" (narrow)
	{'⟦', 1, 1, 1}, // non-ASCII classified as "na" (narrow)
	{'👁', 1, 1, 2},
}

func BenchmarkRuneWidth(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WidthOfRune(runewidthtests[i%len(runewidthtests)].in)
	}
}

func BenchmarkRuneWidthOriginal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runewidth.RuneWidth(runewidthtests[i%len(runewidthtests)].in)
	}
}

func TestRuneWidth(t *testing.T) {
	for i, tt := range runewidthtests {
		if out := WidthOfRune(tt.in); out != tt.out {
			t.Errorf("case %d: WidthOfRune(%q) = %d, want %d", i, tt.in, out, tt.out)
		}
	}
	//c := runewidth.NewCondition()
	//c.EastAsianWidth = false
	//for _, tt := range runewidthtests {
	//	if out := c.RuneWidth(tt.in); out != tt.out {
	//		t.Errorf("RuneWidth(%q) = %d, want %d (EastAsianWidth=false)", tt.in, out, tt.out)
	//	}
	//}
	//c.EastAsianWidth = true
	//for _, tt := range runewidthtests {
	//	if out := c.RuneWidth(tt.in); out != tt.eaout {
	//		t.Errorf("RuneWidth(%q) = %d, want %d (EastAsianWidth=true)", tt.in, out, tt.eaout)
	//	}
	//}
}
