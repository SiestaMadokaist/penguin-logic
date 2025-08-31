package euclidean

import (
	"fmt"
)

type bound struct {
	topLeft     Point
	bottomRight Point
}

type IBound interface {
	Contains(coord Point) bool
	Width() W
	Height() H
	Top() Y
	Bottom() Y
	Left() X
	Right() X
	TopLeft() Point
	BottomRight() Point
	ToString() string
	InnerCoords() []Point
	Center() Point
	ShiftNeg(p Point) IBound
	ShiftPos(p Point) IBound
}

func Bound(topLeft, bottomRight Point) IBound {
	return &bound{
		topLeft:     topLeft,
		bottomRight: bottomRight,
	}
}

func (b bound) ShiftNeg(p Point) IBound {
	newTopLeft := b.topLeft.ShiftNeg(p)
	newBottomRight := b.bottomRight.ShiftNeg(p)
	return Bound(newTopLeft, newBottomRight)
}

func (b bound) ShiftPos(p Point) IBound {
	newTopLeft := b.topLeft.ShiftPos(p)
	newBottomRight := b.bottomRight.ShiftPos(p)
	return Bound(newTopLeft, newBottomRight)
}

func (b bound) Center() Point {
	centerX := (b.Left() + b.Right()) / 2
	centerY := (b.Top() + b.Bottom()) / 2
	return P2(centerX, centerY)
}

func (b bound) InnerCoords() []Point {
	accumulator := []Point{}
	for y := b.Top() + 1; y < b.Bottom(); y++ {
		for x := b.Left() + 1; x < b.Right(); x++ {
			accumulator = append(accumulator, P2(x, y))
		}
	}
	return accumulator
}

func (b bound) ToString() string {
	return fmt.Sprintf("TopLeft: %s, BottomRight: %s", b.topLeft.ToString(), b.bottomRight.ToString())
}

func (b bound) String() string {
	return b.ToString()
}

func (b bound) TopLeft() Point {
	return b.topLeft
}

func (b bound) BottomRight() Point {
	return b.bottomRight
}

func (b bound) Width() W {
	return b.bottomRight.X.Dist(b.topLeft.X)
}

func (b bound) Height() H {
	return b.bottomRight.Y.Dist(b.topLeft.Y)
}

func (b bound) Left() X {
	return b.topLeft.X
}

func (b bound) Right() X {
	return b.bottomRight.X
}

func (b bound) Top() Y {
	return b.topLeft.Y
}

func (b bound) Bottom() Y {
	return b.bottomRight.Y
}

func (b bound) Contains(coord Point) bool {
	if coord.X < b.topLeft.X {
		return false
	}
	if coord.Y < b.topLeft.Y {
		return false
	}
	if coord.X > b.bottomRight.X {
		return false
	}
	if coord.Y > b.bottomRight.Y {
		return false
	}
	return true
}

// func (b bound) SplitBy(p pattern.IPattern) []IBound {
// 	patternHeight := p.Height()
// 	patternWidth := p.Width()
// 	boundWidth := int(b.Width())
// 	boundHeight := int(b.Height())
// 	var bounds []IBound
// 	for y := 0; y < boundHeight; y += patternHeight {
// 		for x := 0; x < boundWidth; x += patternWidth {
// 			subTopLeft := Point{
// 				X: b.topLeft.X + X(x),
// 				Y: b.topLeft.Y + Y(y),
// 			}
// 			subBottomRight := Point{
// 				X: b.topLeft.X + X(min(x+patternWidth, boundWidth)-1),
// 				Y: b.topLeft.Y + Y(min(y+patternHeight, boundHeight)-1),
// 			}
// 			bounds = append(bounds, Bound(subTopLeft, subBottomRight))
// 		}
// 	}
// 	return bounds
// }
