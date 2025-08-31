package imaging_test

import (
	"fmt"
	"image/png"
	"math/rand"
	"os"
	"testing"
	"time"

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
	// integratedRed, err := imaging.Integrate(w, h, reds)
	// if err != nil {
	// 	t.Error(err)
	// }
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
	{46 + 30, 179 + 30},
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
	patterns := []imaging.IPattern{imaging.FeatHorizontal(), imaging.FeatVertical(), imaging.FeatDiagonal(), imaging.FeatInner3()}
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
	t.Errorf("---")
}

func randomBound(w euclidean.W, h euclidean.H) euclidean.IBound {
	x0 := euclidean.X(rand.Intn(int(1200 - w)))
	y0 := euclidean.Y(rand.Intn(int(500 - h)))
	x1 := x0.Add(w)
	y1 := y0.Add(h)
	return euclidean.Bound(euclidean.P2(x0, y0), euclidean.P2(x1, y1))
}

/** @todo */
const WORK_DIR = "/home/ramadoka/development/personal/penguin-logic"

// const OUTDIR =
// this should return true for all test cases
func TestFixed1(t *testing.T) {
	original, err := imaging.Load(fmt.Sprintf("%s/integrations/inputs/1.png", WORK_DIR))
	if err != nil {
		t.Error(err)
	}
	i := original.Invert()
	integral, _ := i.Integral()
	now := time.Now()
	outDir := fmt.Sprintf("%s/integrations/outputs/guesses/%d", WORK_DIR, now.Unix())
	_ = os.MkdirAll(outDir, os.ModePerm)
	tcsAll := []euclidean.IBound{
		euclidean.Bound(euclidean.P2(522, 169), euclidean.P2(652, 299)),
		// euclidean.Bound(euclidean.P2(555, 143), euclidean.P2(685, 273)),
	}
	// tcsAll = append(tcsAll, tcs()...)
	for idx, tc := range tcsAll {
		values, _ := integral.Guess(tc)
		outName := fmt.Sprintf("%d %d.%d.%d", idx, values[0], values[1], values[2])
		outPath := fmt.Sprintf("%s/%s", outDir, outName)
		cropped := i.Crop(tc)
		// err := cropped.Save(fmt.Sprintf("%s.png", outPath))
		// if err != nil {
		// 	t.Error(err)
		// }
		for _, channel := range color.Channels() {
			extracted := cropped.Extract(channel)
			outF, _ := os.Create(fmt.Sprintf("%s_%s.png", outPath, channel))
			png.Encode(outF, extracted)
		}
		// t.Errorf("%d: Guessed Right for %s", idx, tc)
		t.Errorf("%d. Guess for %s (%+v)", idx, tc, values)
	}
	t.Error()
}

func TestFixed(t *testing.T) {
	original, err := imaging.Load(fmt.Sprintf("%s/integrations/inputs/1.png", WORK_DIR))
	if err != nil {
		t.Error(err)
	}
	i := original.Invert()
	integral, _ := i.Integral()
	now := time.Now()
	outDir := fmt.Sprintf("%s/integrations/outputs/guesses/%d", WORK_DIR, now.Unix())
	_ = os.MkdirAll(outDir, os.ModePerm)
	// tcsAll := tcs()
	tcsAll := []euclidean.IBound{tcs()[0]}
	// TopLeft (555, 143), BottomRight: (685, 273)
	// tcsAll = append(tcsAll, tcs()...)
	for idx, tc := range tcsAll {
		// values, guess := integral.Guess(tc)
		recenteredBound := integral.BoundRecenter(color.ChannelGray, tc, 5)
		fmt.Printf("%d. Center of Mass for %s\n", idx, tc)
		fmt.Printf("OriginalCenter: %v - CenterOfMass: %v\n", tc.Center(), recenteredBound.Center())
		outName := fmt.Sprintf("%d", idx)
		outPath := fmt.Sprintf("%s/%s", outDir, outName)
		cropped := i.Crop(recenteredBound)
		// if !guess {
		// 	t.Errorf("%d failed to guess for %+v => %+v", idx, tc, values)
		// }
		for _, channel := range color.Channels() {
			extracted := cropped.Extract(channel)
			outF, _ := os.Create(fmt.Sprintf("%s_%s.png", outPath, channel))
			png.Encode(outF, extracted)
		}
		// // t.Errorf("%d: Guessed Right for %s", idx, tc)
		// t.Errorf("%d. Guess for %s (%+v)", idx, tc, values)
	}
	t.Error()
}

func TestRandom(t *testing.T) {
	original, err := imaging.Load(fmt.Sprintf("%s/integrations/inputs/1.png", WORK_DIR))
	if err != nil {
		t.Error(err)
	}
	i := original.Invert()
	integral, _ := i.Integral()
	now := time.Now()
	outDir := fmt.Sprintf("%s/integrations/outputs/guesses/%d", WORK_DIR, now.Unix())
	_ = os.MkdirAll(outDir, os.ModePerm)
	for idx := range 50 {
		b := randomBound(130, 130)
		_, yesno := integral.Guess(b)
		if yesno {
			rebound := integral.BoundRecenter(color.ChannelGray, b, 5)
			outName := fmt.Sprintf("YES-%d-%v", idx, rebound.Center())
			outPath := fmt.Sprintf("%s/%s", outDir, outName)
			cropped := i.Invert().Crop(rebound)
			cropped.Save(fmt.Sprintf("%s.png", outPath))
			t.Errorf("%d: Guessed Right for %s, Center: %v", idx, b, rebound.Center())
		} else {
			// continue
			outName := fmt.Sprintf("NO-%d-%v", idx, b.Center())
			outPath := fmt.Sprintf("%s/%s", outDir, outName)
			cropped := i.Crop(b)
			cropped.Save(fmt.Sprintf("%s.png", outPath))
			t.Errorf("%d: Guessed Wrong for %s", idx, b)
		}
	}
}
