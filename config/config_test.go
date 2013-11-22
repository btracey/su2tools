package config

import (
	"os"
	"reflect"
	"testing"
)

func TestWriteDefault(t *testing.T) {
	d := NewOptions()
	f, err := os.Create("testwriter.txt")
	if err != nil {
		panic(err)
	}
	d.WriteConfig(f, PrintAll)

	f.Close()

	f2, err := os.Open("testwriter.txt")
	if err != nil {
		panic(err)
	}
	options, err := ReadConfig(f2)
	if err != nil {
		t.Errorf(err.Error())
	}
	f2.Close()

	// Set some options
	options.CflNumber = 12.0
	err = options.SetEnum("PhysicalProblem", "RANS")
	if err != nil {
		t.Errorf("Error setting enum")
	}
	options.CflRamp = []float64{1.0, 30.0, 45.0}

	// Write
	f3, err := os.Create("modified_config.txt")
	if err != nil {
		panic(err)
	}
	options.WriteConfig(f3, PrintAll)
	f3.Close()

	f4, err := os.Open("modified_config.txt")
	if err != nil {
		panic(err)
	}
	options2, err := ReadConfig(f4)
	if err != nil {
		t.Errorf("error reading: ", err.Error())
	}
	if !reflect.DeepEqual(options, options2) {
		t.Errorf("options do not match")
	}
}
