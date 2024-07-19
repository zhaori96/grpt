package grpt

type Positioned struct {
	Offset   Offset
	Position Alignment
	Child    Element
}

func (p Positioned) GetSize() Size {
	return NewSize(0, 0)
}

func (p *Positioned) Measure(boundries Size, renderer *DocumentRenderer) {
	p.Child.Measure(boundries, renderer)
}

func (p *Positioned) Render(renderer *DocumentRenderer) error {
	defer renderer.SetOffset(renderer.GetCurrentOffset())

	renderer.SetOffset(p.Offset)
	return p.Child.Render(renderer)
}
