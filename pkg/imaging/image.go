package imaging

import (
	"image"
	"os"

	c "image/color"
	_ "image/jpeg" // register JPEG format
	"image/png"

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
	Invert() Image
	Crop(b euclidean.IBound) Image
	Save(path string) error
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

func (i *image_) Crop(bound euclidean.IBound) Image {
	cropped := image.NewRGBA(image.Rect(0, 0, int(bound.Width()), int(bound.Height())))
	inners := bound.InnerCoords()
	for _, coord := range inners {
		shifted := coord.ShiftNeg(bound.TopLeft())
		cropped.Set(int(shifted.X), int(shifted.Y), i.i.At(int(coord.X), int(coord.Y)))
	}
	return New(cropped)
}

func (i *image_) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, i.i)
}

func (i *image_) Invert() Image {
	newImage := image.NewRGBA(i.i.Bounds())
	for y := 0; y < newImage.Bounds().Dy(); y++ {
		for x := 0; x < newImage.Bounds().Dx(); x++ {
			// coord := euclidean.P2(euclidean.X(x), euclidean.Y(y))
			oldColor := i.i.At(x, y)
			r, g, b, a := oldColor.RGBA()
			newColor := c.RGBA{R: 255 - uint8(r), G: 255 - uint8(g), B: 255 - uint8(b), A: uint8(a)}
			newImage.Set(x, y, newColor)
		}
	}
	return &image_{
		newImage,
		memoize.Store(),
	}
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
				newColor = c.RGBA{R: v, A: 255}
			case color.ChannelGreen:
				v := uint8(i.Green(coord) >> 8)
				newColor = c.RGBA{G: v, A: 255}
			case color.ChannelBlue:
				v := uint8(i.Blue(coord) >> 8)
				newColor = c.RGBA{B: v, A: 255}
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
