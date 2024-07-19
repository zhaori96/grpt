package grpt

type LineStyle string

const (
	LineStyleDotted LineStyle = "dotted"
	LineStyleDashed LineStyle = "dashed"
	LineStyleSolid  LineStyle = "solid"
)

func (l LineStyle) IsValid() bool {
	switch l {
	case LineStyleDashed, LineStyleDotted, LineStyleSolid:
		return true
	default:
		return false
	}
}

type LineOptions struct {
	StrokeWidth float64
	Style       LineStyle
	Color       *Color
}

func NewLineOptions(stroke float64, style LineStyle, color *Color) LineOptions {
	return LineOptions{
		StrokeWidth: stroke,
		Style:       style,
		Color:       color,
	}
}
