package grpt

type HorizontalTitledTextBox struct {
	Size             Size
	DefaultChildSize Size
	Title            *Text
	Text             *Text
	Padding          EdgeInsets
	Justfy           JustifyContent
	Separator        Element
	Border           Border
	Borders          []Border

	wasMeasuredAtLeastOnce bool
	originalSize           Size
}

func (t HorizontalTitledTextBox) GetSize() Size {
	return t.Size
}

func (t *HorizontalTitledTextBox) Measure(
	boundries Size,
	renderer *DocumentRenderer,
) {
	if t.wasMeasuredAtLeastOnce {
		t.Size = t.originalSize
	} else {
		t.originalSize = t.Size
	}
	t.wasMeasuredAtLeastOnce = true

	if t.Size.Width == MaxSize {
		t.Size.Width = boundries.Width
	}

	if t.Size.Height == MaxSize {
		t.Size.Height = boundries.Height
	}

	children := Elements{t.Title, t.Text}
	if t.Size.Height == 0 {
		maxChildHeight := MaxHeight(children)
		if maxChildHeight > 0 && maxChildHeight != MaxSize {
			t.Size.Height = maxChildHeight
		} else {
			t.Size.Height = boundries.Height
		}
	}

	if t.DefaultChildSize.Height == 0 || t.DefaultChildSize.Height == MaxSize {
		t.DefaultChildSize.Height = t.Size.Height
	}

	if t.DefaultChildSize.Width == 0 || t.DefaultChildSize.Width == MaxSize {
		size := CalculateUnsizedElementSize(
			children,
			t.Size.Merge(boundries),
			HorizontalAxis,
		)
		t.DefaultChildSize.Width = size.Width
	}

	MeasureAll(children, t.DefaultChildSize, renderer)

	if t.Size.Width == 0 {
		t.Size.Width = TotalWidth(children)
		if t.Separator != nil {
			separatorCount := float64(len(children) - 1)
			separatorWidth := t.Separator.GetSize().Width * separatorCount
			t.Size.Width += separatorWidth
		}
	}
}

func (t *HorizontalTitledTextBox) Render(renderer *DocumentRenderer) error {
	element := &Container{
		Size:    t.Size,
		Border:  t.Border,
		Borders: t.Borders,
		Padding: t.Padding,
		Child: &Row{
			Justify:          t.Justfy,
			Separator:        t.Separator,
			DefaultChildSize: t.DefaultChildSize,
			Children: Elements{
				t.Title,
				t.Text,
			},
		},
	}

	return element.Render(renderer)
}
