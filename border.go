package grpt

type BorderSide int

const (
	BorderLeft BorderSide = 1 << iota
	BorderRigh
	BorderTop
	BorderBottom

	BorderAll = BorderLeft | BorderRigh | BorderTop | BorderBottom
)

func (b BorderSide) IsValid() bool {
	return b&BorderLeft != 0 || b&BorderRigh != 0 ||
		b&BorderTop != 0 || b&BorderBottom != 0
}

type Border struct {
	Side    BorderSide
	Options LineOptions
}

func NewBorder(side BorderSide) Border {
	return Border{
		Side: side,
	}
}

func NewBorderWithOptions(side BorderSide, options LineOptions) Border {
	return Border{
		Side:    side,
		Options: options,
	}
}

func NewBorderAll() Border {
	return Border{
		Side: BorderAll,
	}
}

func NewBorderAllWithOptions(options LineOptions) Border {
	return Border{
		Side:    BorderAll,
		Options: options,
	}
}

func NewBorderLeft() Border {
	return Border{
		Side:    BorderLeft,
		Options: LineOptions{},
	}
}

func NewBorderLeftWithOptions(options LineOptions) Border {
	return Border{
		Side:    BorderLeft,
		Options: options,
	}
}

func NewBorderRight() Border {
	return Border{
		Side:    BorderRigh,
		Options: LineOptions{},
	}
}

func NewBorderRightWithOptions(options LineOptions) Border {
	return Border{
		Side:    BorderRigh,
		Options: options,
	}
}

func NewBorderTop() Border {
	return Border{
		Side:    BorderTop,
		Options: LineOptions{},
	}
}

func NewBorderTopWithOptions(options LineOptions) Border {
	return Border{
		Side:    BorderTop,
		Options: options,
	}
}

func NewBorderBottom() Border {
	return Border{
		Side:    BorderBottom,
		Options: LineOptions{},
	}
}

func NewBorderBottomWithOptions(options LineOptions) Border {
	return Border{
		Side:    BorderBottom,
		Options: options,
	}
}
