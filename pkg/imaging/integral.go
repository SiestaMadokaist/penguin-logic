package imaging

import (
	"fmt"
	"math"

	"github.com/ramadoka/penguin-logic/pkg/color"
	"github.com/ramadoka/penguin-logic/pkg/euclidean"
)

type integral struct {
	reds   []color.Red
	greens []color.Green
	blues  []color.Blue
	grays  []color.Gray
	w      euclidean.W
	h      euclidean.H
}

type IntegralImage interface {
	ApplyFeat(channel color.Channel, bound euclidean.IBound, pattern IPattern) int64
	ExtractFeat(channel color.Channel, bound euclidean.IBound, pattern IPattern) []Feature
	Guess(bound euclidean.IBound) ([]int64, bool)
	Calculate(channel color.Channel, bound euclidean.IBound) int64
	CenterOfMass(channel color.Channel, bound euclidean.IBound) euclidean.Point
	BoundRecenter(channel color.Channel, bound euclidean.IBound, maxIter int) euclidean.IBound
}

// 1234567890 => 1.234.567.890
func prettyPrint(n int) string {
	s := fmt.Sprintf("%d", n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "." + s[i:]
	}
	return s
}

func (i integral) BoundRecenter(channel color.Channel, bound euclidean.IBound, maxIter int) euclidean.IBound {
	centerOfBound := bound.Center()
	if maxIter <= 0 {
		return bound
	}
	centerOfMass := i.CenterOfMass(channel, bound)
	shift := centerOfMass.ShiftNeg(centerOfBound)
	fmt.Printf("shift: %v\n", shift)
	fmt.Printf("%d. enter of Bound: %v, Center of Mass: %v\n", maxIter, centerOfBound, centerOfMass)
	dx := math.Abs(float64(shift.X))
	dy := math.Abs(float64(shift.Y))
	if dx < 5 && dy < 5 {
		return bound
	}
	newBound := bound.ShiftPos(shift)
	fmt.Printf("Center of Bound: %v, Center of Mass: %v\n", centerOfBound, centerOfMass)
	return i.BoundRecenter(channel, newBound, maxIter-1)
}

func (i integral) CenterOfMass(channel color.Channel, bound euclidean.IBound) euclidean.Point {
	total := i.Calculate(channel, bound)
	half := total / 2
	limY := [2]euclidean.Y{bound.Top(), bound.Bottom()}
	limX := [2]euclidean.X{bound.Left(), bound.Right()}
	for range 100 {
		if limX[0] >= limX[1] {
			break
		}
		newX := (limX[0] + limX[1]) / 2
		bottomRight := euclidean.P2(euclidean.X(newX), bound.Bottom())
		partialBound := euclidean.Bound(bound.TopLeft(), bottomRight)
		newSum := i.Calculate(channel, partialBound)
		if newSum > half {
			limX[1] = newX
		} else {
			limX[0] = newX + 1
		}
	}

	for range 100 {
		if limY[0] >= limY[1] {
			break
		}
		newY := (limY[0] + limY[1]) / 2
		bottomRight := euclidean.P2(bound.Right(), euclidean.Y(newY))
		partialBound := euclidean.Bound(bound.TopLeft(), bottomRight)
		newSum := i.Calculate(channel, partialBound)
		if newSum > half {
			limY[1] = newY
		} else {
			limY[0] = newY + 1
		}
	}

	return euclidean.P2(euclidean.X(limX[0]-1), euclidean.Y(limY[0]-1))
}

func (i integral) ApplyFeat(channel color.Channel, bound euclidean.IBound, pattern IPattern) int64 {
	feats := i.ExtractFeat(channel, bound, pattern)
	sum := int64(0)
	for _, feat := range feats {
		value := i.Calculate(channel, feat.Bound)
		normalized := int64(feat.Multiplier) * int64(value)
		sum += normalized
		// fmt.Printf("%d. Feature %v: %d x %d = %d => %d\n", idx, feat.Bound, feat.Multiplier, value, normalized, sum)
	}
	// fmt.Printf("Sum: %d\n", sum)
	// fmt.Println("---")
	return sum
}

func (i integral) ExtractFeat(channel color.Channel, bound euclidean.IBound, pattern IPattern) []Feature {
	split := Split(bound, pattern)
	return split
}

func (i integral) Guess(bound euclidean.IBound) ([]int64, bool) {
	feat := FeatInner5()
	outs := []int64{0, 0, 0}
	offense := 0
	for idx, channel := range color.Channels() {
		v := i.ApplyFeat(channel, bound, feat)
		outs[idx] = v
		if v < 0 {
			offense++
		}
	}

	return outs, offense <= 1
}

func (i integral) Calculate(channel color.Channel, bound euclidean.IBound) int64 {
	topLeft := bound.TopLeft()
	bottomRight := bound.BottomRight()
	switch channel {
	case color.ChannelRed:
		return i.SumRed(topLeft, bottomRight)
	case color.ChannelGreen:
		return i.SumGreen(topLeft, bottomRight)
	case color.ChannelBlue:
		return i.SumBlue(topLeft, bottomRight)
	case color.ChannelGray:
		return i.SumGray(topLeft, bottomRight)
	}
	return 0
}

func (i integral) SumRed(topLeft, bottomRight euclidean.Point) int64 {
	return sumColor(i.w, i.h, i.reds, topLeft, bottomRight)
}

func (i integral) SumGray(topLeft, bottomRight euclidean.Point) int64 {
	return sumColor(i.w, i.h, i.grays, topLeft, bottomRight)
}

func (i integral) SumGreen(topLeft, bottomRight euclidean.Point) int64 {
	return sumColor(i.w, i.h, i.greens, topLeft, bottomRight)
}

func (i integral) SumBlue(topLeft, bottomRight euclidean.Point) int64 {
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
	grays, err := Integrate(image.Width(), image.Height(), image.grays())
	if err != nil {
		return zero, err
	}
	return integral{
		reds:   reds,
		greens: greens,
		blues:  blues,
		grays:  grays,
		w:      image.Width(),
		h:      image.Height(),
	}, nil
}

func sumColor[T ~uint32](w euclidean.W, h euclidean.H, colors []T, topLeft, bottomRight euclidean.Point) int64 {
	A := getColorAt(w, h, colors, bottomRight)
	B := getColorAt(w, h, colors, euclidean.Point{X: topLeft.X - 1, Y: bottomRight.Y})
	C := getColorAt(w, h, colors, euclidean.Point{X: bottomRight.X, Y: topLeft.Y - 1})
	D := getColorAt(w, h, colors, euclidean.Point{X: topLeft.X - 1, Y: topLeft.Y - 1})
	return int64(A - B - C + D)
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
