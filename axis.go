package grpt

type Axis int

const (
	HorizontalAxis = iota
	VerticalAxis
)

func (a Axis) IsValid() bool {
	return a >= 0 && a <= 1
}

func (a Axis) Cross() Axis {
	if a == HorizontalAxis {
		return VerticalAxis
	}
	return HorizontalAxis
}
