package imaging

import (
	"fmt"

	"github.com/ramadoka/penguin-logic/pkg/euclidean"
)

type Feature struct {
	Bound      euclidean.IBound
	Multiplier int
}

func (f Feature) String() string {
	return fmt.Sprintf("%d@%s", f.Multiplier, f.Bound.ToString())
}

func Split(zone euclidean.IBound, p IPattern) []Feature {
	pattern := p.Inverse()
	rows := len(pattern)
	if rows == 0 {
		return nil
	}
	cols := len(pattern[0])
	for i := 1; i < rows; i++ {
		if len(pattern[i]) != cols {
			panic("pattern must be a rectangular matrix")
		}
	}

	topLeft := zone.TopLeft()
	// bottomRight := zone.BottomRight()
	w, h := zone.Width(), zone.Height()
	if w <= 0 || h <= 0 {
		return nil
	}

	xCuts := splitEdges(int(w), cols)
	yCuts := splitEdges(int(h), rows)

	out := make([]Feature, 0, rows*cols)
	for r := 0; r < rows; r++ {
		y0 := topLeft.Y.Add(euclidean.H(yCuts[r]))
		y1 := topLeft.Y.Add(euclidean.H(yCuts[r+1]))
		for c := 0; c < cols; c++ {
			x0 := topLeft.X.Add(euclidean.W(xCuts[c]))
			x1 := topLeft.X.Add(euclidean.W(xCuts[c+1]))
			topLeft := euclidean.Point{X: x0, Y: y0}
			bottomRight := euclidean.Point{X: x1, Y: y1}
			bound := euclidean.Bound(topLeft, bottomRight)
			out = append(out, Feature{
				Bound:      bound,
				Multiplier: pattern[r][c],
			})
		}
	}
	return out
}

func splitEdges(size, n int) []int {
	if n <= 0 {
		panic("n must be > 0")
	}
	edges := make([]int, n+1)
	q, r := size/n, size%n
	acc := 0
	edges[0] = 0
	for i := 0; i < n; i++ {
		step := q
		if i < r {
			step++
		}
		acc += step
		edges[i+1] = acc
	}
	return edges
}
