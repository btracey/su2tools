package config

import (
	"os"
	"testing"
)

func TestWriteDefault(t *testing.T) {
	d := NewOptions()
	f, err := os.Create("testwriter.txt")
	if err != nil {
		panic(err)
	}
	d.WriteConfig(f, PrintAll)
}
