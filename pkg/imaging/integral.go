package imaging

import (
	"fmt"

	"github.com/ramadoka/penguin-logic/pkg/color"
	"github.com/ramadoka/penguin-logic/pkg/euclidean"
)

type integral struct {
	reds   []color.Red
	greens []color.Green
	blues  []color.Blue
	w      euclidean.W
	h      euclidean.H
}

type IntegralImage interface {
	ApplyFeat(channel color.Channel, bound euclidean.IBound, pattern IPattern) int32
	ExtractFeat(channel color.Channel, bound euclidean.IBound, pattern IPattern) []Feature
	Calculate(channel color.Channel, bound euclidean.IBound) int32
}

func (i integral) ApplyFeat(channel color.Channel, bound euclidean.IBound, pattern IPattern) int32 {
	feats := i.ExtractFeat(channel, bound, pattern)
	sum := int32(0)
	for _, feat := range feats {
		value := i.Calculate(channel, feat.Bound)
		normalized := int32(feat.Multiplier) * value
		sum += normalized
	}
	return sum
}

func (i integral) ExtractFeat(channel color.Channel, bound euclidean.IBound, pattern IPattern) []Feature {
	split := Split(bound, pattern)
	return split
}

func (i integral) Calculate(channel color.Channel, bound euclidean.IBound) int32 {
	topLeft := bound.TopLeft()
	bottomRight := bound.BottomRight()
	switch channel {
	case color.ChannelRed:
		return int32(i.SumRed(topLeft, bottomRight))
	case color.ChannelGreen:
		return int32(i.SumGreen(topLeft, bottomRight))
	case color.ChannelBlue:
		return int32(i.SumBlue(topLeft, bottomRight))
	}
	return 0
}

func (i integral) SumRed(topLeft, bottomRight euclidean.Point) color.Red {
	return sumColor(i.w, i.h, i.reds, topLeft, bottomRight)
}

func (i integral) SumGreen(topLeft, bottomRight euclidean.Point) color.Green {
	return sumColor(i.w, i.h, i.greens, topLeft, bottomRight)
}

func (i integral) SumBlue(topLeft, bottomRight euclidean.Point) color.Blue {
	return sumColor(i.w, i.h, i.blues, topLeft, bottomRight)
}

func Integral(image image_) (integral, error) {
	zero := integral{}
	reds, err := Integrate(image.Width(), image.Height(), image.reds())
	if err != nil {
		return zero, err
	}
	greens, err := Integrate(image.Width(), image.Height(), image.greens())
	if err != nil {
		return zero, err
	}
	blues, err := Integrate(image.Width(), image.Height(), image.blues())
	if err != nil {
		return zero, err
	}
	return integral{reds: reds, greens: greens, blues: blues, w: image.Width(), h: image.Height()}, nil
}

func sumColor[T ~uint32](w euclidean.W, h euclidean.H, colors []T, topLeft, bottomRight euclidean.Point) T {
	A := getColorAt(w, h, colors, bottomRight)
	B := getColorAt(w, h, colors, euclidean.Point{X: topLeft.X - 1, Y: bottomRight.Y})
	C := getColorAt(w, h, colors, euclidean.Point{X: bottomRight.X, Y: topLeft.Y - 1})
	D := getColorAt(w, h, colors, euclidean.Point{X: topLeft.X - 1, Y: topLeft.Y - 1})
	return A - B - C + D
}

func getColorAt[T ~uint32](w euclidean.W, h euclidean.H, colors []T, coord euclidean.Point) T {
	zero := euclidean.Point{X: 0, Y: 0}
	if coord.X < 0 {
		return 0
	}
	if coord.Y < 0 {
		return 0
	}
	if coord.X >= zero.X.Add(w) {
		return 0
	}
	if coord.Y >= zero.Y.Add(h) {
		return 0
	}
	index := coord.ToIndex0(w)
	return colors[index]
}

func Integrate[T ~uint32](width euclidean.W, height euclidean.H, colors []T) ([]T, error) {
	if (len(colors) != int(width.Mul(height))) || width == 0 || height == 0 {
		return nil, fmt.Errorf("invalid dimensions: %d x %d for %d colors", width, height, len(colors))
	}
	size := width.Mul(height)
	result := make([]T, size)
	for i := range int(size) {
		coord := euclidean.FromIndex(width, height, i)
		up := coord.Up()
		colorUp := getColorAt(width, height, result, up)
		// fmt.Printf("up: %v, colorUp: %d\n", up, colorUp)

		left := coord.Left()
		colorLeft := getColorAt(width, height, result, left)
		// fmt.Printf("left: %v, colorLeft: %d\n", left, colorLeft)

		diag := coord.Left().Up()
		colorDiag := getColorAt(width, height, result, diag)
		// fmt.Printf("diag: %v, colorDiag: %d\n", diag, colorDiag)
		result[i] = colorUp + colorLeft - colorDiag + colors[i]
	}
	return result, nil
}
