package grpt

import (
	"fmt"
	"slices"
)

type OverflowMode int

const (
	OverflowModeTruncate OverflowMode = iota
	OverflowModeMustFitPage
	OverflowModeContinueOnNextPage
)

type Element interface {
	GetSize() Size
	Measure(boundries Size, renderer *DocumentRenderer)
	Render(renderer *DocumentRenderer) error
}

type ElementsSummary struct {
	TotalSize     Size
	MaxSize       Size
	MinSize       Size
	UnsizedWidth  int
	UnsizedHeight int
}

func (es ElementsSummary) UnsizedAxis(axis Axis) int {
	switch axis {
	case HorizontalAxis:
		return es.UnsizedWidth
	case VerticalAxis:
		return es.UnsizedHeight
	default:
		panic(fmt.Errorf(
			"ElementsSummary.UnsizedAxis: %w",
			ErrInvalidAxis,
		))
	}
}

type Elements []Element

func TotalWidth(elements []Element) float64 {
	var total float64
	for _, element := range elements {
		total += element.GetSize().Width
	}
	return total
}

func MaxWitdh(elements []Element) float64 {
	if len(elements) == 0 {
		return 0
	}
	return slices.MaxFunc(elements, func(a, b Element) int {
		if b.GetSize().Width == MaxSize {
			return 1
		}
		return int(a.GetSize().Width - b.GetSize().Width)
	}).GetSize().Width
}

func MaxHeight(elements []Element) float64 {
	if len(elements) == 0 {
		return 0
	}
	return slices.MaxFunc(elements, func(a, b Element) int {
		if b.GetSize().Height == MaxSize {
			return 1
		}
		return int(a.GetSize().Height - b.GetSize().Height)
	}).GetSize().Height
}

func MaxAxisSize(elements []Element, axis Axis) float64 {
	switch axis {
	case HorizontalAxis:
		return MaxWitdh(elements)
	case VerticalAxis:
		return MaxHeight(elements)
	default:
		panic("Elements.MaxFromAxis: invalid axis")
	}
}

func TotalHeight(elements []Element) float64 {
	var total float64
	for _, element := range elements {
		total += element.GetSize().Height
	}
	return total
}

func TotalAxisSize(elements []Element, axis Axis) float64 {
	switch axis {
	case HorizontalAxis:
		return TotalWidth(elements)
	case VerticalAxis:
		return TotalHeight(elements)
	default:
		panic("Elements.TotalFromAxis: invalid axis")
	}
}

func TotalSize(elements []Element) Size {
	total := Size{}
	for _, element := range elements {
		boundries := element.GetSize()
		total.Width += boundries.Width
		total.Height += boundries.Height
	}
	return total
}

func FitsParent(elements []Element, parent Size, axis Axis) bool {
	switch axis {
	case HorizontalAxis:
		return parent.Width >= TotalWidth(elements) && parent.Height >= MaxHeight(elements)
	case VerticalAxis:
		return parent.Height >= TotalHeight(elements) && parent.Width >= MaxWitdh(elements)
	default:
		panic("List.Render: invalid Axis")
	}
}

func FitsParentWithPadding(
	elements []Element,
	parent Size,
	padding EdgeInsets,
	axis Axis,
) bool {
	return FitsParent(elements, parent.WithPadding(padding), axis)
}

func SpaceBetween(elements []Element, parent Size, axis Axis) Offset {
	if parent.IsZero() {
		return NewZeroOffset()
	}

	elementCount := len(elements)
	if elementCount == 0 {
		return NewZeroOffset()
	}

	switch axis {
	case HorizontalAxis:
		totalSize := TotalWidth(elements)
		return NewOffsetX((parent.Width - totalSize) / float64(elementCount-1))
	case VerticalAxis:
		totalSize := TotalHeight(elements)
		return NewOffsetY((parent.Height - totalSize) / float64(elementCount-1))
	default:
		panic(fmt.Errorf("Elements.SpaceBetween: %w", ErrInvalidAxis))
	}
}

func SpaceEvenly(elements []Element, parent Size, axis Axis) Offset {
	if parent.IsZero() {
		return NewZeroOffset()
	}

	elementCount := len(elements)
	if elementCount == 0 {
		return NewZeroOffset()
	}

	switch axis {
	case HorizontalAxis:
		totalSize := TotalWidth(elements)
		return NewOffsetX((parent.Width - totalSize) / float64(elementCount+1))
	case VerticalAxis:
		totalSize := TotalHeight(elements)
		return NewOffsetY((parent.Height - totalSize) / float64(elementCount+1))
	default:
		panic("Elements.SpaceEvenly: invalid axis")
	}
}

func SpaceAround(elements []Element, parent Size, axis Axis) (Offset, Offset) {
	if parent.IsZero() {
		return NewZeroOffset(), NewZeroOffset()
	}

	elementCount := len(elements)
	if elementCount == 0 {
		return NewZeroOffset(), NewZeroOffset()
	}

	switch axis {
	case HorizontalAxis:
		totalSize := TotalWidth(elements)
		extraSpace := parent.Width - totalSize
		spacing := extraSpace / float64(elementCount*2)
		return NewOffsetX(spacing), NewOffsetX(spacing * 2)
	case VerticalAxis:
		totalSize := TotalHeight(elements)
		extraSpace := parent.Height - totalSize
		spacing := extraSpace / float64(elementCount*2)
		return NewOffsetY(spacing), NewOffsetY(spacing * 2)

	default:
		panic("Elements.SpaceAround: invalid axis")
	}
}

func MeasureAll(elements []Element, boundries Size, renderer *DocumentRenderer) {
	for _, element := range elements {
		element.Measure(boundries, renderer)
	}
}

func CalculateSpacing(
	elements []Element,
	parent Size,
	justify JustifyContent,
	axis Axis,
) (Offset, Offset) {
	var edgeGap, betweenGap Offset
	switch justify {
	case JustifyContentSpaceAround:
		edgeGap, betweenGap = SpaceAround(elements, parent, axis)

	case JustifyContentSpaceBetween:
		betweenGap = SpaceBetween(elements, parent, axis)

	case JustifyContentSpaceEvenly:
		betweenGap = SpaceEvenly(elements, parent, axis)
		edgeGap = betweenGap
	}
	return edgeGap, betweenGap
}

func CalculateSpacingWithPadding(
	elements []Element,
	parent Size,
	padding EdgeInsets,
	justify JustifyContent,
	axis Axis,
) (Offset, Offset) {
	return CalculateSpacing(elements, parent.WithPadding(padding), justify, axis)
}

func CalculateUnsizedElementSize(
	elements []Element,
	availableSpace Size,
	axis Axis,
) Size {
	unsizedElements := 0
	totalAxisSize := 0.0
	for _, child := range elements {
		axisSize := child.GetSize().GetAxis(axis)
		if axisSize == 0 || axisSize == MaxSize {
			unsizedElements++
			continue
		}
		totalAxisSize += axisSize
	}

	if unsizedElements == 0 {
		return NewSize(0, 0)
	}

	totalSpace := availableSpace.GetAxis(axis) - totalAxisSize
	size := totalSpace / float64(unsizedElements)
	chesque := NewSizeFromAxis(size, availableSpace.GetAxis(axis.Cross()), axis)

	return chesque
}

func CalculateSummary(elements []Element) ElementsSummary {
	summary := ElementsSummary{}
	for _, element := range elements {
		size := element.GetSize()
		summary.TotalSize.Width += size.Width
		summary.TotalSize.Height += size.Height

		if size.Height > 0 && size.Height < summary.MinSize.Height {
			summary.MinSize.Height = size.Height
		}

		if size.Width == 0 {
			summary.UnsizedWidth++
		} else {
			if size.Width > summary.MaxSize.Width {
				summary.MaxSize.Width = size.Width
			} else if size.Width < summary.MinSize.Width {
				summary.MinSize.Width = size.Width
			}
		}

		if size.Height == 0 {
			summary.UnsizedHeight++
		} else {
			if size.Height > summary.MaxSize.Height {
				summary.MaxSize.Height = size.Height
			} else if size.Height < summary.MinSize.Height {
				summary.MinSize.Height = size.Height
			}
		}
	}

	return summary
}
