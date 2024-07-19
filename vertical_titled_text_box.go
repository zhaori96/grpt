package grpt

type VerticalTitledTextBox struct {
	Size             Size
	DefaultChildSize Size
	Title            *Text
	Text             *Text
	Style            TextStyle
	Padding          EdgeInsets
	Justfy           JustifyContent
	Separator        Element
	Border           Border
	Borders          []Border

	wasMeasuredAtLeastOnce bool
	originalSize           Size
}

func (t VerticalTitledTextBox) GetSize() Size {
	return t.Size
}

func (t *VerticalTitledTextBox) Measure(
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
	if t.Size.Width == 0 {
		maxChildWidth := MaxWitdh(children)
		if maxChildWidth > 0 && maxChildWidth != MaxSize {
			t.Size.Width = maxChildWidth
		} else {
			t.Size.Width = boundries.Width
		}
	}

	if t.DefaultChildSize.Width == 0 || t.DefaultChildSize.Width == MaxSize {
		t.DefaultChildSize.Width = t.Size.Width
	}

	if t.DefaultChildSize.Height == 0 || t.DefaultChildSize.Height == MaxSize {
		size := CalculateUnsizedElementSize(
			children,
			t.Size.Merge(boundries),
			VerticalAxis,
		)
		t.DefaultChildSize.Height = size.Height
	}

	MeasureAll(children, t.DefaultChildSize.Merge(t.Size).Merge(boundries), renderer)

	if t.Size.Height == 0 {
		t.Size.Height = TotalHeight(children)
		if t.Separator != nil {
			separatorCount := float64(len(children) - 1)
			separatorHeight := t.Separator.GetSize().Height * separatorCount
			t.Size.Height += separatorHeight
		}
	}
}

func (t *VerticalTitledTextBox) Render(renderer *DocumentRenderer) error {
	t.Title.Style = t.Title.Style.Merge(t.Style)
	t.Text.Style = t.Text.Style.Merge(t.Style)
	element := &Container{
		Size:    t.Size,
		Border:  t.Border,
		Borders: t.Borders,
		Child: &Column{
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
