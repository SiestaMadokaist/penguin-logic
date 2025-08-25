package playground_test

import (
	"fmt"
	"testing"

	"github.com/ramadoka/penguin-logic/pkg/euclidean"
	"github.com/ramadoka/penguin-logic/pkg/imaging"
)

func TestMain(t *testing.T) {
	filePath := "../inputs/1.png"
	image, err := imaging.Load(filePath)
	if err != nil {
		t.Fatalf("failed to open image: %v", err)
	}
	b := image.Blue(euclidean.Point{X: 0, Y: 0})
	fmt.Println(b)
	// integral, err := image.Integral()
	// if err != nil {
	// 	t.Fatalf("failed to compute integral image: %v", err)
	// }
	// topLeft := euclidean.Point{X: 0, Y: 0}
	// bottomRight := euclidean.Point{X: 10, Y: 10}
	// blue := integral.SumBlue(topLeft, bottomRight)
	// fmt.Println(blue)
}
