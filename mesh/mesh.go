// Package for dealing with SU2 meshes
// For mesh definition, see: http://adl-public.stanford.edu/docs/display/SUSQUARED/Mesh+files
// For similar Python code, see https://github.com/su2code/SU2/blob/master/SU2_PY/SU2/mesh/tools.py

package mesh

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type PointID int
type ElementID int
type VTKType int

const (
	Quadrilateral VTKType = 9
)

type SU2 struct {
	Elements []*Element
	Points   []*Point
	Markers  []*Marker
	Dim      int
}

type Marker struct {
	Tag string
	//TypeID   VTKType
	Elements []Element
}

type Element struct {
	Id        ElementID
	Type      VTKType
	VertexIds []PointID
}

type Point struct {
	Id        PointID
	Location  []float64
	Neighbors map[PointID]*Point
}

// ReadFrom reads the SU2 mesh from an io.Reader creating the mesh
func (s *SU2) ReadFrom(r io.Reader) (n int64, err error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		str := scanner.Text()
		str = strings.TrimSpace(str)
		if len(str) == 0 {
			continue
		}
		switch {
		case strings.HasPrefix(str, "%"):
			// ignore because comment
		case strings.HasPrefix(str, "NDIME="):
			// Parse the number of dimensions
			strs := strings.Split(str, "=")
			if len(strs) != 2 {
				return n, errors.New("more than one equals sign in NDIME line")
			}
			str = strings.TrimSpace(strs[1])
			s.Dim, err = strconv.Atoi(str)
			if err != nil {
				return n, errors.New("error parsing NDIME: " + err.Error())
			}
		case strings.HasPrefix(str, "NELEM="):
			read, err := s.parseElements(scanner, str)
			n += read
			if err != nil {
				return n, err
			}
		case strings.HasPrefix(str, "NPOIN="):
			read, err := s.parsePoints(scanner, str)
			n += read
			if err != nil {
				return n, err
			}
		case strings.HasPrefix(str, "NMARK="):
			read, err := s.parseMarkers(scanner, str)
			n += read
			if err != nil {
				return n, err
			}
		}
	}
	err = s.initialize()
	if err != nil {
		return n, err
	}
	return n, nil
}

func (s *SU2) parseElements(scanner *bufio.Scanner, str string) (n int64, err error) {
	// first, get the number of elements
	strs := strings.Split(str, "=")
	if len(strs) != 2 {
		return n, errors.New("more than one equals sign in NELEM line")
	}
	str = strings.TrimSpace(strs[1])
	nelem, err := strconv.Atoi(str)
	if err != nil {
		return n, errors.New("error parsing NELEM: " + err.Error())
	}
	elements := make([]*Element, nelem)
	for i := 0; i < nelem; i++ {
		scanner.Scan()
		if scanner.Err() != nil {
			return n, fmt.Errorf("error scanning element %d: " + scanner.Err().Error())
		}
		elem := &Element{}
		str := scanner.Text()
		str = strings.TrimSpace(str)
		strs := strings.Fields(str)
		if len(strs) == 0 {
			return n, fmt.Errorf("error scanning element %d: unexpected blank line", i)
		}
		if len(strs) < 2 {
			return n, fmt.Errorf("error scanning element %d: should be at least two numbers (typeID elemID)", i)
		}
		t, err := strconv.Atoi(strs[0])
		elem.Type = VTKType(t)
		if err != nil {
			return n, fmt.Errorf("error scanning element %d: cannot parse typeID: %v", i, err)
		}
		t, err = strconv.Atoi(strs[len(strs)-1])
		elem.Id = ElementID(t)
		if err != nil {
			return n, fmt.Errorf("error scanning element %d: cannot parse nodeID: %v", i, err)
		}
		if elem.Id != ElementID(i) {
			return n, fmt.Errorf("error scanning element %d: bad nodeID: %v", i, err)
		}
		strs = strs[1 : len(strs)-1]
		elem.VertexIds = make([]PointID, len(strs))
		for j, str := range strs {
			t, err := strconv.Atoi(str)
			elem.VertexIds[j] = PointID(t)
			if err != nil {
				return n, fmt.Errorf("error scanning element %d: cannot parse node number: %v", i, err)
			}
		}
		elements[i] = elem
	}
	s.Elements = elements
	return n, nil
}

func (s *SU2) parsePoints(scanner *bufio.Scanner, str string) (n int64, err error) {
	strs := strings.Split(str, "=")
	if len(strs) != 2 {
		return n, errors.New("more than one equals sign in NPOIN line")
	}
	str = strings.TrimSpace(strs[1])
	npoints, err := strconv.Atoi(str)
	if err != nil {
		return n, errors.New("error parsing NPOIN: " + err.Error())
	}
	points := make([]*Point, npoints)
	for i := 0; i < npoints; i++ {
		scanner.Scan()
		if scanner.Err() != nil {
			return n, fmt.Errorf("error scanning point %d: " + scanner.Err().Error())
		}
		str := scanner.Text()
		str = strings.TrimSpace(str)
		strs := strings.Fields(str)
		if len(strs) != s.Dim+1 {
			return n, fmt.Errorf("point %d has wrong number of entries (should have nDim + 1)", i)
		}
		point := &Point{}
		point.Location = make([]float64, s.Dim)
		for j := 0; j < s.Dim; j++ {
			point.Location[j], err = strconv.ParseFloat(strs[j], 64)
			if err != nil {
				return n, fmt.Errorf("error parsing location of point %d: %v", i, err)
			}
		}
		t, err := strconv.Atoi(strs[len(strs)-1])
		point.Id = PointID(t)
		if err != nil {
			return n, fmt.Errorf("error parsing index of point %d: %v", i, err)
		}
		if point.Id != PointID(i) {
			return n, fmt.Errorf("bad point id %d: %v", i, err)
		}
		points[i] = point
	}
	s.Points = points
	return n, nil
}

func (s *SU2) parseMarkers(scanner *bufio.Scanner, str string) (n int64, err error) {
	strs := strings.Split(str, "=")
	if len(strs) != 2 {
		return n, errors.New("more than one equals sign in NMARK line")
	}
	str = strings.TrimSpace(strs[1])
	nMarkers, err := strconv.Atoi(str)
	if err != nil {
		return n, errors.New("error parsing NMARK: " + err.Error())
	}
	markers := make([]*Marker, nMarkers)
	for i := 0; i < nMarkers; i++ {
		scanner.Scan()
		if scanner.Err() != nil {
			return n, fmt.Errorf("error scanning marker %d: " + scanner.Err().Error())
		}
		str := scanner.Text()
		if !strings.HasPrefix(str, "MARKER_TAG") {
			return n, fmt.Errorf("no MARKER_TAG for marker %d", i)
		}
		strs := strings.Split(str, "=")
		if len(strs) != 2 {
			return n, fmt.Errorf("marker %d: MARKER_TAG doesn't have exactly one equals sign ", i)
		}
		marker := &Marker{}
		marker.Tag = strs[1]
		scanner.Scan()
		if scanner.Err() != nil {
			return n, fmt.Errorf("error scanning MARKER_ELEMS %d: " + scanner.Err().Error())
		}
		str = scanner.Text()
		if !strings.HasPrefix(str, "MARKER_ELEMS") {
			return n, fmt.Errorf("MARKER_ELEMS prefix not found for marker %d", i)
		}
		strs = strings.Split(str, "=")
		if len(strs) != 2 {
			return n, fmt.Errorf("marker %d: MARKER_ELEMS doesn't have exactly one equals sign", i)
		}

		str = strings.TrimSpace(strs[1])
		nMarkerElems, err := strconv.Atoi(str)
		if err != nil {
			return n, fmt.Errorf("marker %d:  MARKER_ELEMS parsing error: ", err.Error())
		}
		marker.Elements = make([]Element, nMarkerElems)
		for j := 0; j < nMarkerElems; j++ {
			scanner.Scan()
			str := scanner.Text()
			strs := strings.Fields(str)
			// Parse the type ID
			t, err := strconv.Atoi(strs[0])
			marker.Elements[j].Type = VTKType(t)
			if err != nil {
				return n, fmt.Errorf("marker %d, element %d: bad type id: %v", i, j, err)
			}
			marker.Elements[j].VertexIds = make([]PointID, len(strs)-1)
			for k := 0; k < len(strs)-1; k++ {
				str = strs[k+1]
				str = strings.TrimSpace(str)
				t, err := strconv.Atoi(str)
				marker.Elements[j].VertexIds[k] = PointID(t)
				if err != nil {
					return n, fmt.Errorf("marked %d, element %d, point %d: bad point %v", i, j, k, err)
				}
			}
			marker.Elements[j].Id = -1 // marker elements don't have an ID
		}
		markers[i] = marker
	}
	s.Markers = markers
	return n, nil
}

func (s *SU2) initialize() error {
	// Create all of the neighbor maps
	for _, point := range s.Points {
		point.Neighbors = make(map[PointID]*Point)
	}
	// Add all of the neighboring points
	for _, elem := range s.Elements {
		switch elem.Type {
		case Quadrilateral:
			l := len(elem.VertexIds)
			// The neighbors are the adjacent nodes in the list
			for j, id := range elem.VertexIds {
				neighbor := elem.VertexIds[(j+1)%l]
				s.Points[id].Neighbors[neighbor] = s.Points[neighbor]
				neighbor = elem.VertexIds[(j+l-1)%l]
				s.Points[id].Neighbors[neighbor] = s.Points[neighbor]
			}
		default:
			return fmt.Errorf("element type %d not implemented", elem.Type)
		}
	}
	return nil
}
