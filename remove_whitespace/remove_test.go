package remove_whitespace

import (
	"bytes"
	"fmt"
	"testing"
)

var TestFile string = "Title1 Title2 \n23	54	23		\n15	23 16"
var Nowhitespace string = "Title1 Title2\n23	54	23\n15	23 16"

func TestMain(t *testing.T) {
	err := mainargs([]string{"name"})
	if err == nil {
		t.Errorf("Main runs with no arguments")
	}
	err = mainargs([]string{"name", "one", "two", "three"})
	if err == nil {
		t.Errorf("main runs with three arguments")
	}
}

func TestRemove(t *testing.T) {
	b := []byte(TestFile)
	writer := new(bytes.Buffer)
	err := RemoveTrailingWhitespace(b, writer)
	if err != nil {
		t.Errorf(err.Error())
	}
	if string(writer.Bytes()) != Nowhitespace {
		t.Errorf("writing mismatch. \nExpecting " + Nowhitespace + "\nFound " + string(writer.Bytes()))
		fmt.Println([]byte(Nowhitespace))
		fmt.Println(writer.Bytes())
	}
}

// Need to add something where actually create the files to test main
