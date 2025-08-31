package imaging

import (
	"github.com/ramadoka/penguin-logic/pkg/color"
	"github.com/ramadoka/penguin-logic/pkg/euclidean"
	"github.com/ramadoka/penguin-logic/pkg/memoize"
)

func (img *image_) _topLeft() euclidean.Point {
	topLeft := euclidean.P(img.i.Bounds().Min)
	return topLeft
}
func (img *image_) _bottomRight() euclidean.Point {
	bottomRight := euclidean.P(img.i.Bounds().Max)
	return bottomRight
}

func (img *image_) _width() euclidean.W {
	return img.BottomRight().X.Dist(img.TopLeft().X)
}

func (img *image_) _height() euclidean.H {
	return img.BottomRight().Y.Dist(img.TopLeft().Y)
}

func (img *image_) colors() [5][]uint32 {
	return memoize.Memoize("colors", img.memoizer, img._colors)
}

func Colors(img *image_) [5][]uint32 {
	return img._colors()
}

func (img *image_) _colors() [5][]uint32 {
	colors := [5][]uint32{}
	for y := img.Top(); y < img.Bottom(); y++ {
		for x := img.Left(); x < img.Right(); x++ {
			r, g, b, a := img.i.At(int(x), int(y)).RGBA()
			colors[0] = append(colors[0], r)
			colors[1] = append(colors[1], g)
			colors[2] = append(colors[2], b)
			colors[3] = append(colors[3], a)
			gray := (r + g + b) / 3
			colors[4] = append(colors[4], gray)
		}
	}
	return colors
}

func (i *image_) reds() []color.Red {
	return memoize.Memoize("reds", i.memoizer, i._reds)
}

func (i *image_) greens() []color.Green {
	return memoize.Memoize("greens", i.memoizer, i._greens)
}

func (i *image_) blues() []color.Blue {
	return memoize.Memoize("blues", i.memoizer, i._blues)
}

func (i *image_) grays() []color.Gray {
	return memoize.Memoize("grays", i.memoizer, i._grays)
}

func (i *image_) _reds() []color.Red {
	return wrap(i, 0, color.R)
}

func (i *image_) _greens() []color.Green {
	return wrap(i, 1, color.G)
}

func (i *image_) _blues() []color.Blue {
	return wrap(i, 2, color.B)
}

func (i *image_) _grays() []color.Gray {
	return wrap(i, 4, color.BW)
}

func wrap[T ~uint32](img *image_, colorIdx int, wrapper func(uint32) T) []T {
	values := img.colors()[colorIdx]
	wrapped := make([]T, len(values))
	for i, v := range values {
		wrapped[i] = wrapper(v)
	}
	return wrapped
}
