package imaging

import (
	"errors"
)

type INamedPattern func() IPattern

type IPattern interface {
	validate() (bool, error)
	Height() int
	Width() int
	Inverse() [][]int
	XY(x int, y int) (int, error)
}

type pattern [][]int

func (p pattern) Inverse() [][]int {
	return p
}

func (p pattern) XY(x int, y int) (int, error) {
	if y < 0 || y >= len(p) {
		return 0, errors.New("y out of range")
	}
	if x < 0 || x >= len(p[0]) {
		return 0, errors.New("x out of range")
	}
	return p[y][x], nil
}

func (p pattern) validate() (bool, error) {
	if len(p) == 0 || len(p[0]) == 0 {
		return false, errors.New("invalid pattern")
	}
	width := len(p[0])
	for _, row := range p {
		w := len(row)
		if w != width {
			return false, errors.New("inconsistent row width")
		}
	}
	return true, nil
}

func (p pattern) Height() int {
	return len(p)
}

func (p pattern) Width() int {
	return len(p[0])
}

func InitFeat(data [][]int) (IPattern, error) {
	p := pattern(data)
	valid, err := p.validate()
	if valid {
		return p, nil
	}
	return nil, err
}

func FeatVertical() IPattern {
	p, _ := InitFeat([][]int{{1}, {-1}})
	return p
}

func FeatHorizontal() IPattern {
	p, _ := InitFeat([][]int{{1, -1}})
	return p
}

func FeatDiagonal() IPattern {
	p, _ := InitFeat([][]int{{1, -1}, {-1, 1}})
	return p
}

func FeatInner3() IPattern {
	p, _ := InitFeat([][]int{{-1, -1, -1}, {-1, 8, -1}, {-1, -1, -1}})
	return p
}

func FeatInner4() IPattern {
	p, _ := InitFeat([][]int{{-1, -1, -1, -1}, {-1, 3, 3, -1}, {-1, 3, 3, -1}, {-1, -1, -1, -1}})
	return p
}

func FeatDynamicHorizontal(size int) IPattern {
	items := make([][]int, 1)
	items[0] = make([]int, size)
	for i := 0; i < size; i++ {
		items[0][i] = 1
	}
	p, _ := InitFeat(items)
	return p
}

func FeatDynamicVertical(size int) IPattern {
	items := make([][]int, size)
	for i := 0; i < size; i++ {
		items[i] = make([]int, 1)
		items[i][0] = 1
	}
	p, _ := InitFeat(items)
	return p
}

func FeatInner5() IPattern {
	neg := -1
	pos := 0  // 8 * 1 = 8
	pos3 := 2 // 4 * 3 * 1 = 12
	pos4 := 4 // 1 * 4 * 1 = 4
	p, _ := InitFeat([][]int{
		{neg, neg, pos, neg, neg},
		{neg, pos, pos3, pos, neg},
		{pos, pos3, pos4, pos3, pos},
		{neg, pos, pos3, pos, neg},
		{neg, neg, pos, neg, neg},
	})
	return p
}
