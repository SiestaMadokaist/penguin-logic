package euclidean

import (
	"fmt"
	"image"
)

type Point struct {
	X X
	Y Y
}

func (p Point) ToString() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

func P2(x X, y Y) Point {
	return Point{X: x, Y: y}
}

func P(p image.Point) Point {
	return Point{X: X(p.X), Y: Y(p.Y)}
}

func (p Point) Up() Point {
	return Point{X: p.X, Y: p.Y - 1}
}

func (p Point) AddX(x X) Point {
	return Point{X: p.X + x, Y: p.Y}
}

func (p Point) AddY(y Y) Point {
	return Point{X: p.X, Y: p.Y + y}
}

func (p Point) Add(other Point) Point {
	return Point{X: p.X + other.X, Y: p.Y + other.Y}
}

func (p Point) Down() Point {
	return Point{X: p.X, Y: p.Y + 1}
}

func (p Point) Left() Point {
	return Point{X: p.X - 1, Y: p.Y}
}

func (p Point) Right() Point {
	return Point{X: p.X + 1, Y: p.Y}
}

func (p Point) ToIndex0(width W) int {
	return p.ToIndex(Point{X: 0, Y: 0}, width)
}

func (p Point) ToIndex(topLeft Point, width W) int {
	yDiff := p.Y.Dist(topLeft.Y)
	xDiff := p.X.Dist(topLeft.X)
	return int(int(yDiff)*int(width) + int(xDiff))
}

func (p Point) Scale(factor float64, center Point) Point {
	return Point{
		X: X(float64(p.X-center.X)*factor) + center.X,
		Y: Y(float64(p.Y-center.Y)*factor) + center.Y,
	}
}

func FromIndex(width W, height H, index int) Point {
	x := index % int(width)
	y := index / int(width)
	return Point{X: X(x), Y: Y(y)}
}
