package imaging_test

import (
	"fmt"
	"image/png"
	"os"
	"testing"

	"github.com/ramadoka/penguin-logic/pkg/color"
	"github.com/ramadoka/penguin-logic/pkg/euclidean"
	"github.com/ramadoka/penguin-logic/pkg/imaging"
)

func TestIntegrate(t *testing.T) {
	reds := make([]color.Red, 24)
	w := euclidean.W(6)
	h := euclidean.H(4)
	size := w.Mul(h)
	for i := range size {
		reds[i] = color.Red(i)
	}
	integratedRed, err := imaging.Integrate(w, h, reds)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v", integratedRed)
	t.Errorf("---")
}

type tcBound struct {
	x0 euclidean.X
	y0 euclidean.Y
	x1 euclidean.X
	y1 euclidean.Y
}

func (b tcBound) ToBound() euclidean.IBound {
	return euclidean.Bound(euclidean.P2(b.x0, b.y0), euclidean.P2(b.x1, b.y1))
}

var Ys = [][]euclidean.Y{{31, 164}, {216, 345}, {391, 515}}
var Xs = [][]euclidean.X{
	{46, 179},
	{130, 250},
	{192, 318},
	{343, 473},
	{483, 626},
	{634, 771},
	{781, 919},
	{924, 1069},
	{1078, 1203},
}

func tcs() []euclidean.IBound {
	acc := []euclidean.IBound{}
	for _, yPair := range Ys {
		for _, xYpair := range Xs {
			b := tcBound{y0: yPair[0], y1: yPair[1], x0: xYpair[0], x1: xYpair[1]}
			acc = append(acc, b.ToBound())
		}
	}
	return acc
}

func TestExtract(t *testing.T) {
	inputDir := "/home/ramadoka/development/personal/penguin-logic/integrations/inputs/"
	i, err := imaging.Load(fmt.Sprintf("%s/%s", inputDir, "1.png"))
	if err != nil {
		t.Fatalf("failed to load image: %v", err)
	}
	outputDir := "/home/ramadoka/development/personal/penguin-logic/integrations/outputs/"
	colors := []color.Channel{color.ChannelRed, color.ChannelGreen, color.ChannelBlue}
	for n, color := range colors {
		outPath := fmt.Sprintf("%s/%d.png", outputDir, n)
		outF, _ := os.Create(outPath)
		extracted := i.Extract(color)
		png.Encode(outF, extracted)
		outF.Close()
	}
	t.Error("Extraction test not implemented")
}

func TestIntegral(t *testing.T) {
	i, err := imaging.Load("/home/ramadoka/development/personal/penguin-logic/integrations/inputs/1.png")
	if err != nil {
		t.Fatalf("failed to load image: %v", err)
	}
	integral, _ := i.Integral()
	colors := []color.Channel{color.ChannelRed, color.ChannelGreen, color.ChannelBlue}
	patterns := []imaging.IPattern{imaging.FeatHorizontal(), imaging.FeatVertical(), imaging.FeatDiagonal(), imaging.FeatInner()}
	for tc, bound := range tcs() {
		for _, channel := range colors {
			for n, p := range patterns {
				feat := integral.ApplyFeat(channel, bound, p)
				fmt.Printf("tc #%d. (%s) (Feat#%d) channel %s: %d\n", tc, bound, n, channel, feat)
			}
			fmt.Println()
		}
		fmt.Println("")
	}
	// for _, feat := range feats {
	// 	c0 := integral.Calculate(color.ChannelBlue, feat.Bound)
	// 	fmt.Printf("%s: %d\n", feat, c0)
	// }
	// fmt.Printf("%s", feats)
	t.Errorf("---")

}
