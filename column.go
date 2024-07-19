package grpt

type Column struct {
	Name             string
	Size             Size
	Separator        Element
	Justify          JustifyContent
	DefaultChildSize Size
	OverflowMode     OverflowMode
	Children         Elements

	originalPage           int
	currentChildIndex      int
	wasMeasuredAtLeastOnce bool
	originalSize           Size
}

func (c Column) GetSize() Size {
	return c.Size
}

func (c *Column) Measure(boundries Size, renderer *DocumentRenderer) {
	if c.wasMeasuredAtLeastOnce {
		c.Size = c.originalSize
	} else {
		c.originalSize = c.Size
		c.originalPage = renderer.GetCurrentPage()
	}
	c.wasMeasuredAtLeastOnce = true

	if c.Size.Width == MaxSize {
		c.Size.Width = boundries.Width
	}

	if c.Size.Height == MaxSize {
		c.Size.Height = boundries.Height
	}

	if c.Size.Width == 0 {
		maxChildWidth := MaxWitdh(c.Children)
		if maxChildWidth > 0 && maxChildWidth != MaxSize {
			c.Size.Width = maxChildWidth
		} else {
			c.Size.Width = boundries.Width
		}
	}

	if c.DefaultChildSize.Width == 0 || c.DefaultChildSize.Width == MaxSize {
		c.DefaultChildSize.Width = c.Size.Width
	}

	if c.DefaultChildSize.Height == 0 || c.DefaultChildSize.Height == MaxSize {
		size := CalculateUnsizedElementSize(
			c.Children,
			c.Size.Merge(boundries),
			VerticalAxis,
		)
		c.DefaultChildSize.Height = size.Height
	}

	MeasureAll(c.Children, c.DefaultChildSize.Merge(c.Size), renderer)

	if c.Size.Height == 0 {
		c.Size.Height = TotalHeight(c.Children[c.currentChildIndex:])
		if c.Separator != nil {
			separatorCount := float64(len(c.Children) - 1)
			separatorHeight := c.Separator.GetSize().Height * separatorCount
			c.Size.Height += separatorHeight
		}
	}

	if c.Size.Height > boundries.Height {
		c.Size.Height = boundries.Height
	}
}

func (c *Column) Render(renderer *DocumentRenderer) error {
	defer renderer.SetOffset(renderer.GetCurrentOffset())

	var edgeGap, gap Offset
	if c.Separator == nil {
		edgeGap, gap = CalculateSpacing(
			c.Children,
			c.Size,
			c.Justify,
			VerticalAxis,
		)
	}

	renderer.AddOffset(edgeGap)

	if c.OverflowMode == OverflowModeMustFitPage {
		if !renderer.FitsCurrentPage(c.Size.Height) {
			renderer.AddPage()
		}
	}

	position := renderer.GetCurrentOffset()
	for index, child := range c.Children {
		c.currentChildIndex = index
		if len(c.Children) == 150 {
			println("X")
		}
		if renderer.GetCurrentPage() == 2 && len(c.Children) == 150 {
			println("Y")
		}
		if c.OverflowMode == OverflowModeContinueOnNextPage {
			if !renderer.FitsIn(child.GetSize(), position, c.Size) {
				renderer.AddPage()
				position = renderer.GetCurrentOffset()
			}
		}

		if err := child.Render(renderer); err != nil {
			return err
		}

		renderer.AddOffsetFromAxis(child.GetSize().ToOffset(), VerticalAxis)
		if index < len(c.Children)-1 {
			if c.Separator != nil {
				c.Separator.Render(renderer)
				renderer.AddOffsetFromAxis(
					c.Separator.GetSize().ToOffset(),
					VerticalAxis,
				)
			} else {
				renderer.AddOffset(gap)
			}
		}
	}

	return nil
}
