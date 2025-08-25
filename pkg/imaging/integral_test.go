package imaging_test

import (
	"fmt"
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

func TestIntegral(t *testing.T) {
	i, err := imaging.Load("/home/ramadoka/development/personal/penguin-logic/integrations/inputs/1.png")
	if err != nil {
		t.Fatalf("failed to load image: %v", err)
	}
	start := euclidean.P2(0, 0)
	windowSize := euclidean.P2(64, 64)
	bound := euclidean.Bound(start, windowSize)
	integral, _ := i.Integral()
	feats := integral.FeatExtract(color.ChannelBlue, bound, imaging.FeatDiagonal())
	for _, feat := range feats {
		c0 := integral.Calculate(color.ChannelBlue, feat.Bound)
		fmt.Printf("%s: %d\n", feat, c0)
	}
	// fmt.Printf("%s", feats)
	t.Errorf("---")

}
