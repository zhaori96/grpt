package grpt

type Alignment int

const (
	LeftAlignment             Alignment = 8
	RightAlignment            Alignment = 2
	HorizontalCenterAlignment Alignment = 16
	VerticalCenterAlignment   Alignment = 32
	TopAlignment              Alignment = 4
	BottomAlignment           Alignment = 1

	CenterAlignment Alignment = HorizontalCenterAlignment | VerticalCenterAlignment
)

func (a Alignment) IsValid() bool {
	switch a {
	case LeftAlignment, RightAlignment,
		HorizontalCenterAlignment, VerticalCenterAlignment, CenterAlignment,
		TopAlignment, BottomAlignment:
		if a&LeftAlignment != 0 &&
			(a&RightAlignment != 0 || a&HorizontalCenterAlignment != 0) {
			return false
		}

		if a&TopAlignment != 0 &&
			(a&BottomAlignment != 0 || a&VerticalCenterAlignment != 0) {
			return false
		}

		return true
	default:
		return false
	}
}
