package grpt

import (
	"database/sql/driver"
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
	Value          any
	Formatter      TextFormatter
	FormatterFunc  Formatter
	SkipFormatting bool
	Size           Size
	Style          TextStyle

	text                   string
	wasMeasuredAtLeastOnce bool
	originalSize           Size
}

func (t Text) GetSize() Size {
	return t.Size
}

func (t *Text) Measure(boundries Size, renderer *DocumentRenderer) {
	if t.wasMeasuredAtLeastOnce {
		t.Size = t.originalSize
	} else {
		t.originalSize = t.Size
	}
	t.wasMeasuredAtLeastOnce = true

	t.text = t.parseValue()

	t.Size = t.Size.Merge(boundries)
	if t.Size.HasZeroValue() {
		size, _ := renderer.MeasureText(t.text, &t.Style, t.Size.Merge(boundries))
		t.Size = t.Size.Merge(size)
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

	if !t.wasMeasuredAtLeastOnce {
		t.text = t.parseValue()
	}

	return renderer.DrawText(t.text, t.Size, &t.Style)
}

func (t *Text) parseValue() string {
	if !t.SkipFormatting {
		if t.Formatter != nil {
			return t.Formatter.Format(t.Value)
		}
		if t.FormatterFunc != nil {
			return t.FormatterFunc(t.Value)
		}
		if value, ok := t.Value.(FormattedText); ok {
			return value.Formatted()
		}
	}

	if t.Value == nil {
		return ""
	}

	var text string
	switch raw := t.Value.(type) {
	case driver.Valuer:
		if value, _ := raw.Value(); value != nil {
			text = fmt.Sprint(value)
		}
	default:
		text = fmt.Sprint(raw)
	}

	return text
}
