package euclidean

type Y int
type X int
type W int
type H int
type Scale int
type Area int

func (w W) Mul(h H) Area {
	return Area(int(w) * int(h))
}

func (h H) Mul(w W) Area {
	return w.Mul(h)
}

func Neg[X int | Y | W | H | Scale](a X) X {
	return -a
}

func (y Y) Add(dy H) Y {
	return y + Y(dy)
}

func (y Y) Sub(dy H) Y {
	return y - Y(dy)
}

// not abs dist
func (y Y) Dist(y1 Y) H {
	return H(y - y1)
}

func (x X) Add(dx W) X {
	return x + X(dx)
}

func (x X) Sub(dx W) X {
	return x - X(dx)
}

// not abs dist
func (x X) Dist(x1 X) W {
	return W(x - x1)
}
