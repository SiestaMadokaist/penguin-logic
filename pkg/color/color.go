package color

type Red uint32
type Green uint32
type Blue uint32
type Alpha uint32

func R(u uint32) Red {
	return Red(u)
}

func G(u uint32) Green {
	return Green(u)
}

func B(u uint32) Blue {
	return Blue(u)
}

func A(u uint32) Alpha {
	return Alpha(u)
}
