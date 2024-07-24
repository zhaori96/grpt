package grpt

type Visible struct {
	Visible           bool
	AlwaysOccupySpace bool
	Child             Element

	size Size
}

func (v Visible) GetSize() Size {
	return v.size
}

func (v *Visible) Measure(boundries Size, renderer *DocumentRenderer) {
	if !v.Visible && !v.AlwaysOccupySpace {
		v.size = Size{}
		return
	}

	v.Child.Measure(boundries, renderer)
	v.size = v.Child.GetSize()
}

func (v *Visible) Render(renderer *DocumentRenderer) error {
	if v.Visible {
		return v.Child.Render(renderer)
	}
	return nil
}
