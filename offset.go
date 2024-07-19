package grpt

type Offset struct {
	X float64
	Y float64
}

func NewOffset(x float64, y float64) Offset {
	return Offset{
		X: x,
		Y: y,
	}
}

func NewOffsetX(value float64) Offset {
	return NewOffset(value, 0)
}

func NewOffsetY(value float64) Offset {
	return NewOffset(0, value)
}

func NewZeroOffset() Offset {
	return Offset{}
}

func (o Offset) IsZero() bool {
	return o.X == 0 && o.Y == 0
}

func (o Offset) IsValid() bool {
	return o.X >= 0 && o.Y >= 0
}

func (o Offset) FromAxis(axis Axis) float64 {
	switch axis {
	case HorizontalAxis:
		return o.X
	case VerticalAxis:
		return o.Y
	default:
		panic("Offset.FromAxis: invalid axis")
	}
}
