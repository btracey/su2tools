package mesh

import (
	"os"
	"testing"
)

func TestReadFrom(t *testing.T) {
	filename := "mesh_flatplate_turb_137x97.su2"
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf(err.Error())
	}
	s := &SU2{}
	_, err = s.ReadFrom(f)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if s.Dim != 2 {
		t.Errorf("Dimension mismatch. Expected %v, found %v", 2, s.Dim)
	}
}
