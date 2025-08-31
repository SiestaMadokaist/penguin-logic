package imaging

type recognition struct {
	integral IntegralImage
	pattern  IPattern
}

type IRecognition interface {
}

func Recognition(i *IntegralImage, p IPattern) *recognition {
	return &recognition{
		integral: *i,
		pattern:  p,
	}
}
