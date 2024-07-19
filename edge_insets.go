package grpt

type EdgeInsets struct {
	Left   float64
	Right  float64
	Top    float64
	Bottom float64
}

func (e EdgeInsets) IsZero() bool {
	return e.Left == 0 && e.Right == 0 && e.Top == 0 && e.Bottom == 0
}

func (e EdgeInsets) HasZeroValue() bool {
	return e.Left == 0 || e.Right == 0 || e.Top == 0 || e.Bottom == 0
}

func (e EdgeInsets) Merge(other EdgeInsets) EdgeInsets {
	if e.Left == 0 {
		e.Left = other.Left
	}

	if e.Right == 0 {
		e.Right = other.Right
	}

	if e.Top == 0 {
		e.Top = other.Top
	}

	if e.Bottom == 0 {
		e.Bottom = other.Bottom
	}

	return e
}

func NewEdgeInsets(left, right, top, bottom float64) EdgeInsets {
	return EdgeInsets{
		Left:   left,
		Right:  right,
		Top:    top,
		Bottom: bottom,
	}
}

func NewSquareEdgeInsets(value float64) EdgeInsets {
	return EdgeInsets{
		Left:   value,
		Right:  value,
		Top:    value,
		Bottom: value,
	}
}

func NewHorizontalEdgeInsets(left, right float64) EdgeInsets {
	return EdgeInsets{
		Left:  left,
		Right: right,
	}
}

func NewVerticalEdgeInsets(top, bottom float64) EdgeInsets {
	return EdgeInsets{
		Top:    top,
		Bottom: bottom,
	}
}
