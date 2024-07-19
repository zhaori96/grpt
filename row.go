package grpt

type Row struct {
	Name             string
	Size             Size
	Separator        Element
	Justify          JustifyContent
	DefaultChildSize Size
	Children         Elements

	wasMeasuredAtLeastOnce bool
	originalSize           Size
}

func (r Row) GetSize() Size {
	return r.Size
}

func (r *Row) Measure(boundries Size, renderer *DocumentRenderer) {
	if r.wasMeasuredAtLeastOnce {
		r.Size = r.originalSize
	} else {
		r.originalSize = r.Size
	}
	r.wasMeasuredAtLeastOnce = true

	if r.Size.Width == MaxSize {
		r.Size.Width = boundries.Width
	}

	if r.Size.Height == MaxSize {
		r.Size.Height = boundries.Height
	}

	if r.Size.Height == 0 {
		maxChildHeight := MaxHeight(r.Children)
		if maxChildHeight > 0 && maxChildHeight != MaxSize {
			r.Size.Height = maxChildHeight
		} else {
			r.Size.Height = boundries.Height
		}
	}

	if r.DefaultChildSize.Height == 0 || r.DefaultChildSize.Height == MaxSize {
		r.DefaultChildSize.Height = r.Size.Height
	}

	if r.DefaultChildSize.Width == 0 || r.DefaultChildSize.Width == MaxSize {
		size := CalculateUnsizedElementSize(
			r.Children,
			r.Size.Merge(boundries),
			HorizontalAxis,
		)
		r.DefaultChildSize.Width = size.Width
	}

	MeasureAll(r.Children, r.DefaultChildSize, renderer)

	if r.Size.Width == 0 {
		r.Size.Width = TotalWidth(r.Children)
		if r.Separator != nil {
			separatorCount := float64(len(r.Children) - 1)
			separatorWidth := r.Separator.GetSize().Width * separatorCount
			r.Size.Width += separatorWidth
		}
	}
}

func (r *Row) Render(renderer *DocumentRenderer) error {
	defer renderer.SetOffset(renderer.GetCurrentOffset())

	var edgeGap, gap Offset
	if r.Separator == nil {
		edgeGap, gap = CalculateSpacing(
			r.Children,
			r.Size,
			r.Justify,
			HorizontalAxis,
		)
	}

	renderer.AddOffset(edgeGap)
	for index, child := range r.Children {
		if err := child.Render(renderer); err != nil {
			return err
		}

		renderer.AddOffsetFromAxis(child.GetSize().ToOffset(), HorizontalAxis)
		if index < len(r.Children)-1 {
			if r.Separator != nil {
				r.Separator.Render(renderer)
				renderer.AddOffsetFromAxis(
					r.Separator.GetSize().ToOffset(),
					HorizontalAxis,
				)
			} else {
				renderer.AddOffset(gap)
			}
		}
	}

	return nil
}
