package grpt

type JustifyContent int

const (
	JustfiyContentNone JustifyContent = iota
	JustifyContentSpaceBetween
	JustifyContentSpaceAround
	JustifyContentSpaceEvenly
)

func (j JustifyContent) IsValid() bool {
	return j >= 0 && j <= 3
}
