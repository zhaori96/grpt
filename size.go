package grpt

import (
	"fmt"
	"math"

	"github.com/signintech/gopdf"
)

const MaxSize float64 = math.MaxFloat32

type Size struct {
	Width  float64
	Height float64
}

func NewSize(width, height float64) Size {
	return Size{
		Width:  width,
		Height: height,
	}
}

func NewMaxSize() Size {
	return Size{Width: MaxSize, Height: MaxSize}
}

func NewSizeWithPadding(width, height float64, padding EdgeInsets) Size {
	size := Size{
		Width:  width,
		Height: height,
	}

	return size.WithPadding(padding)
}

func NewSizeFromAxis(mainAxis float64, crossAxis float64, axis Axis) Size {
	switch axis {
	case HorizontalAxis:
		return Size{Width: mainAxis, Height: crossAxis}

	case VerticalAxis:
		return Size{Width: crossAxis, Height: mainAxis}

	default:
		panic(fmt.Errorf("NewSizeFromAxis: %w", ErrInvalidAxis))
	}
}

func NewWidth(width float64) Size {
	return Size{
		Width:  width,
		Height: 0,
	}
}

func NewMaxWidth() Size {
	return Size{
		Width:  MaxSize,
		Height: 0,
	}
}

func NewWidthWithPadding(width float64, padding EdgeInsets) Size {
	size := Size{Width: width}.WithPadding(padding)
	size.Height = 0
	return size
}

func NewHeight(height float64) Size {
	return Size{
		Width:  0,
		Height: height,
	}
}

func NewMaxHeight() Size {
	return Size{
		Width:  0,
		Height: MaxSize,
	}
}

func NewHeightWithPadding(height float64, padding EdgeInsets) Size {
	size := Size{Height: height}.WithPadding(padding)
	size.Width = 0
	return size
}

func NewSquareSize(value float64) Size {
	return Size{Width: value, Height: value}
}

func NewSquareSizeWithPadding(value float64, padding EdgeInsets) Size {
	return NewSizeWithPadding(value, value, padding)
}

func (s Size) IsZero() bool {
	return s.Width == 0 && s.Height == 0
}

func (s Size) HasZeroValue() bool {
	return s.Width == 0 || s.Height == 0
}

func (s Size) IsValid() bool {
	return s.Width >= 0 && s.Height >= 0
}

func (s *Size) ToRect() *gopdf.Rect {
	return &gopdf.Rect{
		W: s.Width,
		H: s.Height,
	}
}

func (s Size) ToOffset() Offset {
	return Offset{X: s.Width, Y: s.Height}
}

func (s Size) WithPadding(padding EdgeInsets) Size {
	return Size{
		Width:  s.Width - padding.Left - padding.Right,
		Height: s.Height - padding.Top - padding.Bottom,
	}
}

func (s Size) WithoutPadding(padding EdgeInsets) Size {
	return Size{
		Width:  s.Width + padding.Left + padding.Right,
		Height: s.Height + padding.Top + padding.Bottom,
	}
}

func (s Size) WithHorizontalPadding(left float64, right float64) Size {
	return Size{
		Width:  s.Width - left - right,
		Height: s.Height,
	}
}

func (s Size) WithVerticalPadding(top float64, bottom float64) Size {
	return Size{
		Width:  s.Width,
		Height: s.Height - top - bottom,
	}
}

func (s Size) GetAxis(axis Axis) float64 {
	switch axis {
	case HorizontalAxis:
		return s.Width
	case VerticalAxis:
		return s.Height
	default:
		panic(fmt.Errorf("Size.FromAxis: %w", ErrInvalidAxis))
	}
}

func (s Size) SetMainAxis(value float64, axis Axis) Size {
	switch axis {
	case HorizontalAxis:
		s.Width = value
	case VerticalAxis:
		s.Height = value
	default:
		panic(fmt.Errorf("Size.SetAxis: %w", ErrInvalidAxis))
	}
	return s
}

func (s Size) SetCrossAxis(value float64, axis Axis) Size {
	switch axis {
	case HorizontalAxis:
		s.Height = value
	case VerticalAxis:
		s.Width = value
	default:
		panic(fmt.Errorf("Size.SetAxis: %w", ErrInvalidAxis))
	}
	return s
}

func (s Size) FitsContainer(container Size) bool {
	return s.Width <= container.Width && s.Height <= container.Height
}

func (s Size) FitsContainerAxis(container Size, axis Axis) bool {
	switch axis {
	case HorizontalAxis:
		return s.Width <= container.Width
	case VerticalAxis:
		return s.Height <= container.Height
	default:
		panic(fmt.Errorf("Size.FitsFromAxis: %w", ErrInvalidAxis))
	}
}

func (s Size) Merge(size Size) Size {
	if s.Width == 0 || s.Width == MaxSize {
		s.Width = size.Width
	}

	if s.Height == 0 || s.Height == MaxSize {
		s.Height = size.Height
	}

	return s
}

func (s Size) Reverse() Size {
	return Size{
		Width:  s.Height,
		Height: s.Width,
	}
}

func (s Size) Difference(other Size) Size {
	return Size{
		Width:  s.Width - other.Width,
		Height: s.Height - other.Height,
	}
}
