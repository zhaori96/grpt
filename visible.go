package grpt

type Visible struct {
	Size              Size
	Visible           bool
	AlwaysOccupySpace bool
	Child             Element

	wasMeasuredAtLeastOnce bool
	originalSize           Size
}

func (v *Visible) GetSize() Size {
	if v.Size.HasZeroValue() {
		return v.Size.Merge(v.Child.GetSize())
	}
	return v.Size
}

func (v *Visible) Measure(boundries Size, renderer *DocumentRenderer) {
	if v.wasMeasuredAtLeastOnce {
		v.Size = v.originalSize
	} else {
		v.originalSize = v.Size
	}
	v.wasMeasuredAtLeastOnce = true

	if !v.Visible && !v.AlwaysOccupySpace {
		v.Size = Size{}
		return
	}

	v.Size = v.Size.Merge(boundries)
	if v.Size.HasZeroValue() {
		v.Child.Measure(v.Size, renderer)
		v.Size = v.Size.Merge(v.Child.GetSize())
	} else {
		v.Child.Measure(v.Size, renderer)
	}
}

func (v *Visible) Render(renderer *DocumentRenderer) error {
	if v.Visible {
		return v.Child.Render(renderer)
	}
	return nil
}
