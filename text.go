package grpt

import (
	"fmt"
)

type TextStyle struct {
	Font      *Font
	Alignment Alignment
	Borders   []Border
	Padding   EdgeInsets
	WordWrap  bool
	Multiline bool
	Overflow  string
}

func (t TextStyle) Merge(other TextStyle) TextStyle {
	t.Alignment = t.Alignment | other.Alignment

	if t.Font == nil {
		t.Font = other.Font
	}

	if t.Padding.IsZero() {
		t.Padding = other.Padding
	}

	if len(t.Borders) == 0 {
		t.Borders = other.Borders
	}

	if len(t.Overflow) == 0 {
		t.Overflow = other.Overflow
	}

	if !t.WordWrap {
		t.WordWrap = other.WordWrap
	}

	if !t.Multiline {
		t.Multiline = other.Multiline
	}

	return t
}

type Text struct {
	Value string
	Size  Size
	Style TextStyle

	wasMeasuredAtLeastOnce bool
	originalSize           Size
}

func (t Text) GetSize() Size {
	return t.Size
}

func (l *Text) Measure(boundries Size, renderer *DocumentRenderer) {
	if l.wasMeasuredAtLeastOnce {
		l.Size = l.originalSize
	} else {
		l.originalSize = l.Size
	}
	l.wasMeasuredAtLeastOnce = true

	l.Size = l.Size.Merge(boundries)
	if l.Size.HasZeroValue() {
		size, _ := renderer.MeasureText(l.Value, &l.Style, l.Size.Merge(boundries))
		l.Size = l.Size.Merge(size)
	}
}

func (t *Text) Render(renderer *DocumentRenderer) error {
	if t.Size.HasZeroValue() {
		panic(fmt.Errorf(
			"Text size cannot have width or height of 0: {w:%v, h:%v}: %w",
			t.Size.Width,
			t.Size.Height,
			ErrInvalidSize,
		))
	}

	return renderer.DrawText(t.Value, t.Size, &t.Style)
}
