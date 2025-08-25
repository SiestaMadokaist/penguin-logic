package imaging_test

import (
	"fmt"
	"testing"

	"github.com/ramadoka/penguin-logic/pkg/color"
	"github.com/ramadoka/penguin-logic/pkg/euclidean"
	"github.com/ramadoka/penguin-logic/pkg/imaging"
)

func TestColor(t *testing.T) {
	i, err := imaging.Load("/home/ramadoka/development/personal/penguin-logic/integrations/inputs/1.png")
	if err != nil {
		t.Fatalf("failed to load image: %v", err)
	}
	start := euclidean.P2(0, 0)
	windowSize := euclidean.P2(64, 64)
	bound := euclidean.Bound(start, windowSize)

	integral, _ := i.Integral()
	blue := integral.Calculate(color.ChannelBlue, bound)
	fmt.Println(blue)
}
