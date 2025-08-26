package imaging

import (
	"fmt"
	"image"
	"os"

	c "image/color"
	_ "image/jpeg" // register JPEG format
	_ "image/png"  // register PNG format

	"github.com/ramadoka/penguin-logic/pkg/color"
	"github.com/ramadoka/penguin-logic/pkg/euclidean"
	"github.com/ramadoka/penguin-logic/pkg/memoize"
)

type image_ struct {
	i        image.Image
	memoizer memoize.IStore
}

type Image interface {
	TopLeft() euclidean.Point
	BottomRight() euclidean.Point
	Width() euclidean.W
	Height() euclidean.H
	Top() euclidean.Y
	Bottom() euclidean.Y
	Left() euclidean.X
	Right() euclidean.X
	Red(point euclidean.Point) color.Red
	Green(point euclidean.Point) color.Green
	Blue(point euclidean.Point) color.Blue
	Integral() (IntegralImage, error)
	Extract(channel color.Channel) image.Image
}

func Load(path string) (Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return New(img), nil
}

func New(i image.Image) Image {
	store := memoize.Store()
	return &image_{i, store}
}

func (i *image_) Extract(channel color.Channel) image.Image {
	newImage := image.NewRGBA(i.i.Bounds())
	for y := 0; y < newImage.Bounds().Dy(); y++ {
		for x := 0; x < newImage.Bounds().Dx(); x++ {
			coord := euclidean.P2(euclidean.X(x), euclidean.Y(y))
			newColor := c.RGBA{}
			// other := 0
			switch channel {
			case color.ChannelRed:
				v := uint8(i.Red(coord) >> 8)
				// newColor = c.RGBA{R: 255 - v, G: 255, B: 255, A: 255}
				newColor = c.RGBA{R: 255 - v, A: 255}
				// newColor = c.RGBA{R: v, G: v, B: v, A: 255}
			case color.ChannelGreen:
				v := uint8(i.Green(coord) >> 8)
				// newColor = c.RGBA{G: 255 - v, R: 255, B: 255, A: 255}
				newColor = c.RGBA{G: 255 - v, A: 255}
				// newColor = c.RGBA{R: v, G: v, B: v, A: 255}
			case color.ChannelBlue:
				v := uint8(i.Blue(coord) >> 8)
				// newColor = c.RGBA{B: 255 - v, G: 255, R: 255, A: 255}
				newColor = c.RGBA{B: 255 - v, A: 255}
				// newColor = c.RGBA{R: v, G: v, B: v, A: 255}
			}
			if (y+x)%100 == 0 {
				fmt.Printf("x: %d, y: %d => %d\n", x, y, newColor)
			}
			newImage.Set(x, y, newColor)
		}
	}
	return newImage
}

func (i *image_) Integral() (IntegralImage, error) {
	return Integral(*i)
}

func (img *image_) TopLeft() euclidean.Point {
	return memoize.Memoize("topLeft", img.memoizer, img._topLeft)
}

func (img *image_) BottomRight() euclidean.Point {
	return memoize.Memoize("bottomRight", img.memoizer, img._bottomRight)
}

func (img *image_) Width() euclidean.W {
	return memoize.Memoize("width", img.memoizer, img._width)
}

func (img *image_) Height() euclidean.H {
	return memoize.Memoize("height", img.memoizer, img._height)
}

func (img *image_) Top() euclidean.Y {
	return img.TopLeft().Y
}

func (img *image_) Bottom() euclidean.Y {
	return img.BottomRight().Y
}

func (img *image_) Left() euclidean.X {
	return img.TopLeft().X
}

func (img *image_) Right() euclidean.X {
	return img.BottomRight().X
}

func (img *image_) Red(point euclidean.Point) color.Red {
	index := point.ToIndex(img.TopLeft(), img.Width())
	return img.reds()[index]
}

func (img *image_) Green(point euclidean.Point) color.Green {
	index := point.ToIndex(img.TopLeft(), img.Width())
	return img.greens()[index]
}

func (img *image_) Blue(point euclidean.Point) color.Blue {
	index := point.ToIndex(img.TopLeft(), img.Width())
	return img.blues()[index]
}
