package main

import (
	"flag"
	"log"
	"math"
	"os"
	"sort"
	"strconv"

	"github.com/btracey/numcsv"
	"github.com/btracey/su2tools/mesh/elements"
	"github.com/gonum/matrix/mat64"
)

var (
	dataname    string
	outname     string
	comma       string
	xcolname    string
	ycolname    string
	xIdxColname string
	yIdxColname string
)

func init() {
	flag.StringVar(&dataname, "data", "data.txt", "Filename containing the structured mesh data")
	flag.StringVar(&outname, "out", "mesh.su2", "Output mesh filename")
	flag.StringVar(&comma, "comma", ",", "column delimiter")
	flag.StringVar(&xcolname, "x", "X", "column name for the x locations")
	flag.StringVar(&ycolname, "y", "Y", "column name for the y locations")
	flag.StringVar(&xIdxColname, "xidx", "XIdx", "column name for the x index")
	flag.StringVar(&yIdxColname, "yidx", "YIdx", "column name for the y index")
}

func main() {
	flag.Parse()
	f, err := os.Open(dataname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := numcsv.NewReader(f)
	r.Comma = comma
	headings, err := r.ReadHeading()
	if err != nil {
		log.Fatalf("error reading heading: %v", err)
	}
	data, err := r.ReadAll()
	if err != nil {
		log.Fatalf("error loading data: %v", err)
	}

	// Look for the headings
	XIdx := findString(headings, xIdxColname)
	if XIdx == -1 {
		log.Fatalf("no 'XIdx' column, headings are: %v", headings)
	}
	YIdx := findString(headings, yIdxColname)
	if YIdx == -1 {
		log.Fatalf("no 'YIdx' column")
	}
	X := findString(headings, xcolname)
	if X == -1 {
		log.Fatalf("no 'X' column")
	}
	Y := findString(headings, ycolname)
	if Y == -1 {
		log.Fatalf("no 'Y' column")
	}

	rows, _ := data.Dims()
	if rows == 0 {
		log.Fatal("no data")
	}

	// Put the data into a matrix
	grid := mat64.NewDense(rows, 4, nil)
	for i := 0; i < rows; i++ {
		grid.Set(i, 0, data.At(i, XIdx))
		grid.Set(i, 1, data.At(i, YIdx))
		grid.Set(i, 2, data.At(i, X))
		grid.Set(i, 3, data.At(i, Y))
	}

	// Sort the data so that the xs are ordered and then the Ys
	sort.Sort(gridSorter{grid})

	// Check that the x and y indexes are all integers
	for i := 0; i < rows; i++ {
		xidx := grid.At(i, 0)
		if math.Floor(xidx) != xidx {
			log.Fatalf("x index at row %v is not an integer. Value is %v.", i, xidx)
		}
		yidx := grid.At(i, 1)
		if math.Floor(yidx) != yidx {
			log.Fatalf("x index at row %v is not an integer. Value is %v.", i, yidx)
		}
	}

	firstX := int(grid.At(0, 0))
	firstY := int(grid.At(0, 1))
	lastX := int(grid.At(rows-1, 0))
	lastY := int(grid.At(rows-1, 1))

	nX := lastX - firstX + 1
	nY := lastY - firstY + 1

	if rows != nX*nY {
		log.Fatalf("nX = %v, nY = %v, but rows = %v", nX, nY, rows)
	}

	// Check that we have a full grid and bring the locs out to it
	xlocs := mat64.NewDense(nX, nY, nil)
	ylocs := mat64.NewDense(nX, nY, nil)
	for i := 0; i < rows; i++ {
		xidx := int(grid.At(i, 0))
		yidx := int(grid.At(i, 1))
		col := i % nY
		row := i / nY
		if xidx != row+firstX {
			log.Fatalf("xidx not correct. Want %v, Got %v.", row+firstX, xidx)
		}
		if yidx != col+firstY {
			log.Fatalf("yidx not correct. Want %v, Got %v.", col+firstY, yidx)
		}
		x := grid.At(i, 2)
		y := grid.At(i, 3)
		xlocs.Set(row, col, x)
		ylocs.Set(row, col, y)
	}

	out, err := os.Create(outname)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	out.WriteString("NDIME= 2\n")
	// Print out all of the vertices
	out.WriteString("NPOIN= " + strconv.Itoa(nX*nY) + "\n")
	for i := 0; i < nX; i++ {
		for j := 0; j < nY; j++ {
			v := &elements.Vertex{
				Idx:      i*nY + j,
				Location: []float64{xlocs.At(i, j), ylocs.At(i, j)},
			}
			b, err := v.MarshalSU2()
			if err != nil {
				log.Fatal(err)
			}
			out.Write(b)
			out.WriteString("\n")
		}
	}

	// Print out the elements
	out.WriteString("NELEM= " + strconv.Itoa((nX-1)*(nY-1)) + "\n")
	for i := 0; i < nX-1; i++ {
		for j := 0; j < nY-1; j++ {
			q := &elements.Quadrilateral{
				Idx:      i*nY + j,
				Vertices: [4]int{i*nY + j, (i+1)*nY + j, (i+1)*nY + j + 1, i*nY + j + 1},
			}
			b, err := q.MarshalSU2()
			if err != nil {
				log.Fatal(err)
			}
			out.Write(b)
			out.WriteString("\n")
		}
	}

	// Print out the edges and markers
	out.WriteString("NMARK= 4\n")
	out.WriteString("MARKER_TAG= lower\n")
	out.WriteString("MARKER_ELEMS= " + strconv.Itoa(nX-1) + "\n")
	for i := 0; i < nX-1; i++ {
		l := &elements.Line{
			Vertices: [2]int{i * nY, (i + 1) * nY},
		}
		b, err := l.MarshalSU2()
		if err != nil {
			log.Fatal(err)
		}
		out.Write(b)
		out.WriteString("\n")
	}
	out.WriteString("MARKER_TAG= upper\n")
	out.WriteString("MARKER_ELEMS= " + strconv.Itoa(nX-1) + "\n")
	for i := 0; i < nX-1; i++ {
		l := &elements.Line{
			Vertices: [2]int{i*nY + nY - 1, (i+1)*nY + nY - 1},
		}
		b, err := l.MarshalSU2()
		if err != nil {
			log.Fatal(err)
		}
		out.Write(b)
		out.WriteString("\n")
	}
	out.WriteString("MARKER_TAG= left\n")
	out.WriteString("MARKER_ELEMS= " + strconv.Itoa(nY-1) + "\n")
	for j := 0; j < nY-1; j++ {
		l := &elements.Line{
			Vertices: [2]int{j, j + 1},
		}
		b, err := l.MarshalSU2()
		if err != nil {
			log.Fatal(err)
		}
		out.Write(b)
		out.WriteString("\n")
	}
	out.WriteString("MARKER_TAG= right\n")
	out.WriteString("MARKER_ELEMS= " + strconv.Itoa(nY-1) + "\n")
	for j := 0; j < nY-1; j++ {
		l := &elements.Line{
			Vertices: [2]int{(nX-1)*nY + j, (nX-1)*nY + j + 1},
		}
		b, err := l.MarshalSU2()
		if err != nil {
			log.Fatal(err)
		}
		out.Write(b)
		out.WriteString("\n")
	}
}

type gridSorter struct {
	data *mat64.Dense
}

func (g gridSorter) Len() int {
	r, _ := g.data.Dims()
	return r
}

func (g gridSorter) Less(i, j int) bool {
	xi := g.data.At(i, 0)
	xj := g.data.At(j, 0)
	if xi < xj {
		return true
	}
	if xi > xj {
		return false
	}
	// The x indices are the same, sort by y indices
	yi := g.data.At(i, 1)
	yj := g.data.At(j, 1)
	return yi < yj
}

func (g gridSorter) Swap(i, j int) {
	rowi := g.data.Row(nil, i)
	rowj := g.data.Row(nil, j)
	g.data.SetRow(i, rowj)
	g.data.SetRow(j, rowi)
}

func findString(strs []string, str string) int {
	for i, s := range strs {
		if str == s {
			return i
		}
	}
	return -1
}
