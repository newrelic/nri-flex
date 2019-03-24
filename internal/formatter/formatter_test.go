package formatter

import (
	"fmt"
	"testing"
)

func TestValueParse(t *testing.T) {
	testVal := "helloThereHowAreYou"
	testRegex := "There"
	str := ValueParse(testVal, testRegex)
	if str != "There" {
		t.Errorf("want: There got %v", str)
	}

	testRegex2 := "There2"
	str2 := ValueParse(testVal, testRegex2)
	if str2 != "helloThereHowAreYou" {
		t.Errorf("want: helloThereHowAreYou got %v", str2)
	}
}
func TestSplitKey(t *testing.T) {
	a, b, split := SplitKey("name=bob", "=")
	if a != "name" {
		t.Errorf("key want: name, received: %v", a)
	}
	if b != "bob" {
		t.Errorf("value want: bob, received: %v", b)
	}
	if split == false {
		t.Errorf("split %v failed", split)
	}
	a, b, split = SplitKey("namebob", "=")
	if a != "" {
		t.Errorf("key want: '', received: %v", a)
	}
	if b != "" {
		t.Errorf("value want: '', received: %v", b)
	}
	if split == true {
		t.Errorf("should not have split")
	}

}

func TestPercToDecimal(t *testing.T) {
	var value interface{} = "10.5%"
	PercToDecimal(&value)
	if fmt.Sprintf("%v", value) != "10.5" {
		t.Errorf("want 10.5, received %v", fmt.Sprintf("%v", value))
	}
}

func TestSnakeToCamel(t *testing.T) {
	key := "hello_there_batman"
	SnakeCaseToCamelCase(&key)
	if key != "helloThereBatman" {
		t.Errorf("want helloThereBatman, received %v", key)
	}
}

func TestRegSplit(t *testing.T) {
	expect := []string{"hello", "there", "batman"}
	strings := RegSplit("hello  there  batman", `\s{1,}`)
	for i := range expect {
		if expect[i] != strings[i] {
			t.Errorf("does not match %v : %v", expect[i], strings[i])
		}
	}
}

func TestRegMatch(t *testing.T) {
	expect := []string{"hello", "there", "batman"}
	strings := RegMatch("hello  there  batman", `(\w+)\s+(\w+)\s+(\w+)`)
	for i := range expect {
		if expect[i] != strings[i] {
			t.Errorf("does not match %v : %v", expect[i], strings[i])
		}
	}

	strings = RegMatch("hello  there  batman", `blah`)
	if strings != nil {
		t.Errorf("nothing should have matched!")
	}
}

func TestKvFinder(t *testing.T) {
	found := KvFinder("prefix", "batman", "bat")
	if !found {
		t.Errorf("not found")
	}
	found = KvFinder("suffix", "batman", "man")
	if !found {
		t.Errorf("not found")
	}
	found = KvFinder("contains", "batman", "atm")
	if !found {
		t.Errorf("not found")
	}
	found = KvFinder("regex", "batman", "man$")
	if !found {
		t.Errorf("not found")
	}
	found = KvFinder("contains", "batman", "cat")
	if found {
		t.Errorf("should not have been found")
	}
}
