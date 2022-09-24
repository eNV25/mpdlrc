package urunewidth

import (
	"reflect"
	"testing"

	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
)

func TestGraphemeWidth(t *testing.T) {
	s := " H̡̫̤̤̣͉̤ͭ̓̓̇͗̎̀ơ̯̗̱̘̮͒̄̀̈ͤ̀͡w͓̲͙͖̥͉̹͋ͬ̊ͦ̂̀̚ ͎͉͖̌ͯͅͅd̳̘̿̃̔̏ͣ͂̉̕ŏ̖̙͋ͤ̊͗̓͟͜e͈͕̯̮̙̣͓͌ͭ̍̐̃͒s͙͔̺͇̗̱̿̊̇͞ ̸̤͓̞̱̫ͩͩ͑̋̀ͮͥͦ̊Z̆̊͊҉҉̠̱̦̩͕ą̟̹͈̺̹̋̅ͯĺ̡̘̹̻̩̩͋͘g̪͚͗ͬ͒o̢̖͇̬͍͇͓̔͋͊̓ ̢͈͙͂ͣ̏̿͐͂ͯ͠t̛͓̖̻̲ͤ̈ͣ͝e͋̄ͬ̽͜҉͚̭͇ͅx͎̬̠͇̌ͤ̓̂̓͐͐́͋͡ț̗̹̝̄̌̀ͧͩ̕͢ ̮̗̩̳̱̾w͎̭̤͍͇̰̄͗ͭ̃͗ͮ̐o̢̯̻̰̼͕̾ͣͬ̽̔̍͟ͅr̢̪͙͍̠̀ͅǩ̵̶̗̮̮ͪ́?̙͉̥̬͙̟̮͕ͤ̌͗ͩ̕͡ "
	sw := runewidth.StringWidth(s)

	gw := 0
	g := uniseg.NewGraphemes(s)
	for g.Next() {
		gw += GraphemeWidth(g.Runes())
	}

	if !reflect.DeepEqual(sw, gw) {
		t.Errorf("StringWidth(%q) = %v, sum of GraphemeWidth = %v", s, sw, gw)
	}
	t.Logf("StringWidth(%q) = %v", s, sw)
}
