package color

type Channel int

const (
	ChannelRed Channel = iota
	ChannelGreen
	ChannelBlue
)

func (c Channel) String() string {
	switch c {
	case ChannelRed:
		return "#red"
	case ChannelGreen:
		return "#green"
	case ChannelBlue:
		return "#blue"
	default:
		return "(unknown)"
	}
}
