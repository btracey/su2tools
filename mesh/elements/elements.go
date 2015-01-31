package elements

import "strconv"

type Quadrilateral struct {
	Idx      int
	Vertices [4]int
}

func (q *Quadrilateral) MarshalSU2() ([]byte, error) {
	str := "9"
	for i := range q.Vertices {
		str += " " + strconv.Itoa(q.Vertices[i])
	}
	str += " " + strconv.Itoa(q.Idx)
	return []byte(str), nil
}

type Vertex struct {
	Idx      int
	Location []float64
}

func (v *Vertex) MarshalSU2() ([]byte, error) {
	str := strconv.FormatFloat(v.Location[0], 'g', 16, 64)
	for _, val := range v.Location[1:] {
		str += " " + strconv.FormatFloat(val, 'g', 16, 64)
	}
	str += " " + strconv.Itoa(v.Idx)
	return []byte(str), nil
}

type Line struct {
	Vertices [2]int
}

func (l *Line) MarshalSU2() ([]byte, error) {
	return []byte("3" + " " + strconv.Itoa(l.Vertices[0]) + " " + strconv.Itoa(l.Vertices[1])), nil
}
