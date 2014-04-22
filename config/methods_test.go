package config

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/btracey/su2tools/config/enum"
	"github.com/btracey/su2tools/config/su2types"
)

var gopath string

var testname string

func init() {
	gopath = os.Getenv("GOPATH")
	testname = filepath.Join(gopath, "src", "github.com", "btracey", "su2tools", "config", "testconfig", "testprint.cfg")
}

func TestDefaultWriteTo(t *testing.T) {

	f, err := os.Create(testname)
	if err != nil {
		t.Errorf(err.Error())
	}
	defaultOptions.WriteTo(f, nil)
}

func TestReadAndWrite(t *testing.T) {

	o := NewOptions()

	// Change a bunch of values from their defaults
	o.ExtraOutput = true                  // Boolean
	o.DampNacelleInflow = 8.1234123412    // Double
	o.MeshFilename = "swamp"              // String
	o.TimeInstances = 12                  // Uint16
	o.UnstRestartIter = 51                //int
	o.MathProblem = enum.AdjointProblem   // enum
	o.CflRamp = [3]float64{6, 7, 8}       // array
	o.DvParam = &su2types.DVParam{"waka"} // Custom
	o.RegimeType = enum.Incompressible

	o.MarkerOut1d = []string{"paul", "y", "esther"}                                         // slice
	o.GridMovementKind = []enum.Gridmovement{enum.Aeroelastic, enum.AeroelasticRigidMotion} // Enum list

	// 8:18pm 4/18/14

	b := make([]byte, 0)
	buf := bytes.NewBuffer(b)
	_, err := o.WriteTo(buf, ForceAll)
	if err != nil {
		t.Errorf("Error writing: %v", err)
	}

	b2 := buf.Bytes()

	reader := bytes.NewReader(b2)

	newOpt, _, err := Read(reader)
	if err != nil {
		t.Errorf("Error reading: %v", err)
	}

	diff := Diff(o, newOpt)
	if diff != nil {
		t.Errorf("Option structs not equal after marshal and unmarshal\n%v", diff)
	}

}
