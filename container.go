package grpt

type Container struct {
	Size             Size
	Padding          EdgeInsets
	Border           Border
	Borders          []Border
	ContentAlignment Alignment
	Child            Element

	wasMeasuredAtLeastOnce bool
	originalSize           Size
}

func (c Container) GetSize() Size {
	return c.Size
}

func (c *Container) Measure(boundries Size, renderer *DocumentRenderer) {
	if c.wasMeasuredAtLeastOnce {
		c.Size = c.originalSize
	} else {
		c.originalSize = c.Size
	}
	c.wasMeasuredAtLeastOnce = true

	c.Size = c.Size.Merge(boundries)
	if c.Size.HasZeroValue() {
		c.Child.Measure(c.Size.WithPadding(c.Padding), renderer)
		c.Size = c.Size.Merge(c.Child.GetSize().WithoutPadding(c.Padding))
	} else {
		c.Child.Measure(c.Size.WithPadding(c.Padding), renderer)
	}
}

func (c *Container) Render(document *DocumentRenderer) error {
	defer document.SetOffset(document.GetCurrentOffset())
	if c.Size.HasZeroValue() {
		panic(ErrInvalidSize)
	}

	if len(c.Borders) > 0 {
		document.DrawBoxWithBorders(c.Size, c.Borders...)
	} else if c.Border.Side != 0 {
		document.DrawBoxWithBorders(c.Size, c.Border)
	}

	paddedSize := c.Size.WithPadding(c.Padding)
	document.AddXY(c.Padding.Left, c.Padding.Top)

	emptySpace := paddedSize.Difference(c.Child.GetSize())
	if emptySpace.IsValid() && !emptySpace.IsZero() && c.ContentAlignment.IsValid() {
		if c.ContentAlignment&RightAlignment != 0 {
			document.AddX(emptySpace.Width)
		}

		if c.ContentAlignment&BottomAlignment != 0 {
			document.AddY(emptySpace.Height)
		}

		if c.ContentAlignment&HorizontalCenterAlignment != 0 {
			document.AddX(emptySpace.Width / 2)
		}

		if c.ContentAlignment&VerticalCenterAlignment != 0 {
			document.AddY(emptySpace.Height / 2)
		}
	}

	return c.Child.Render(document)
}
